package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/servers/tty"
	"github.com/pufferpanel/pufferpanel/v3/utils"
)

// This is a test package to see how unshare/chroot works in the Golang work
// if this works well, this wil be merged into the main codebase

//things we need to do:
//configure apparmor to permit us to create user namespaces
/* /etc/apparmor.d/unshare
abi <abi/4.0>,
include <tunables/global>

profile lordralex /usr/bin/unshare flags=(unconfined) {
  userns,
}
*/ //sudo apparmor_parser -r /etc/apparmor.d/unshare

//after more testing, AppArmor can shove it. They are causing our issues. SOMEHOW EVEN THOUGH THEY DONT GIVE A DAMN INDICATION
// - /etc/default/grub -> GRUB_CMDLINE_LINUX="apparmor=0"
// - update-grub
// - reboot

//ensure uidmap is installed (?) - removing it works so far

/*
mkdir -p {tmp,proc,run/systemd/resolve}
unshare --mount-proc=proc --map-users 1000:1000:1 -muipfCr bash -c 'mkdir -p {pufferpanel,dev,bin,usr,lib,lib64,etc}; mount --bind /bin bin; mount --bind /usr usr; mount --bind /lib lib; mount --bind /lib64 lib64; mount --rbind /etc etc; mount --rbind /dev dev; mount --rbind /run/systemd/resolve run/systemd/resolve; mount -t tmpfs -o size=100m tmpfs tmp; mount --rbind / .; unshare -UR . -w pufferpanel bash'
*/

///etc/apparmor.d/pufferpanel
//# Last Modified: Sat Jun 14 14:12:25 2025
//abi <abi/3.0>,
//
//include <tunables/global>
//
///home/lordralex/.cache/JetBrains/GoLand2025.1/tmp/GoLand/___1go_build_github_com_pufferpanel_pufferpanel_v3_tools_unshare {
//  include <abstractions/base>
//
//  capability setgid,
//  capability setuid,
//  capability sys_admin,
//
//  mount options=(rprivate, rw) -> /,
//
//  /home/lordralex/.cache/JetBrains/GoLand2025.1/tmp/GoLand/___1go_build_github_com_pufferpanel_pufferpanel_v3_tools_unshare mr,
//  /usr/bin/* ux,
//  owner /proc/*/gid_map w,
//  owner /proc/*/setgroups w,
//  owner /proc/*/uid_map w,
//}

//what we've ended up with is a nightmare. But, a working one
//the current thing works if we remove new PID (this causes issues with the forking because Go can't do the -f flag)
//we also had to disable the net namespace because we could not actually use the network, unsure why (but I bet... go)
//we store our working folder into /tmp so that it's cleaner where we actually persist files
//this defers back to our "regular" username so we should be okay there
//move the dir we're "chrooted" in to the arguments so we can test it elsewhere
//long term, we might need to look at persistence, and see if we can "persist" these so we don't need to keep remaking it
//this required several rounds of coffee, mountain dew, and midori

func main() {
	dir, _ := os.Getwd()
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	for _, unshare := range []bool{true, false} {
		_ = config.SecurityDisableUnshare.Set(unshare, false)

		fmt.Printf("Unshare status: %t", config.SecurityDisableUnshare.Value())

		//as soon as you add the third command, this no longer functions
		unshareImplementation(dir, "pwd")                        //are we in the right place
		unshareImplementation(dir, "ls", "-l")                   //do we see anything
		unshareImplementation(dir, "whoami")                     //are we the correct user
		unshareImplementation(dir, "touch", "test")              //can we write and it persist
		unshareImplementation(dir, "curl", "1.1.1.1")            //can we access an IP
		unshareImplementation(dir, "curl", "google.com")         //does DNS work
		unshareImplementation(dir, "curl", "https://google.com") //does SSL work
		unshareImplementation(dir, "ls", "-l", "/")              //is root clean

		//now let's test the servers!
		//unshare(dir, "ls", "-l", "/usr/lib/jvm/java-21-openjdk-amd64/lib")
		//unshare(dir, "env")
		//unshareImplementation(dir, "java", "-Xmx4G", "-jar", "server.jar", "nogui")

		unshareImplementation(dir, "pwd && pwd")
	}
}

var cmdList = []string{
	"mkdir -p {dev,bin,usr,lib,lib64,etc,tmp,proc}",
	"mount -t tmpfs -o size=50m tmpfs tmp",
	"mount --bind /bin bin",
	"mount --bind /lib lib",
	"mount --bind /lib64 lib64",
	"mount --rbind /usr usr",
	"mount --rbind /etc etc",
	"mount --rbind /dev dev",
	"mount --rbind /proc proc",
}

func unshareImplementation(dir, cmd string, args ...string) {
	factory := tty.EnvironmentFactory{}

	env := pufferpanel.Environment{
		RootDirectory: dir, Wait: &sync.Mutex{}, ConsoleBuffer: pufferpanel.CreateCache(), ConsoleTracker: pufferpanel.CreateTracker(),
		StatsTracker: pufferpanel.CreateTracker(), StatusTracker: pufferpanel.CreateTracker(),
		Wrapper: os.Stdout,
	}
	env.Implementation = factory.Create()

	c := cmd
	if len(args) > 0 {
		c += " " + utils.MergeArguments(args)
	}

	err := env.Execute(pufferpanel.ExecutionData{
		Command:      c,
		StdInConfig:  pufferpanel.StdinConsoleConfiguration{},
		DisableQuery: true,
		DisableStats: true,
	})
	if err != nil {
		panic(err)
	}
}

func unshareDirect(dir, cmd string, args ...string) {
	var err error

	c := exec.Command("bash", "-c",
		strings.Join(append(cmdList,
			"mkdir -p "+strings.TrimPrefix(dir, "/"),
			"mount --bind "+dir+" "+strings.TrimPrefix(dir, "/"),
			"mount --rbind / .",
			fmt.Sprintf("unshare -UR . -w %s --map-user=%d --map-group=%d %s %s", dir, os.Getuid(), os.Getgid(), cmd, strings.Join(args, " ")),
		), " && "))
	c.Dir, err = os.MkdirTemp("", "unshare-pp-")
	defer func() {
		_ = os.RemoveAll(c.Dir)
	}()
	if err != nil {
		panic(err)
	}
	c.SysProcAttr = &syscall.SysProcAttr{
		Unshareflags: syscall.CLONE_NEWUSER |
			syscall.CLONE_NEWNS |
			syscall.CLONE_FILES |
			syscall.CLONE_NEWCGROUP |
			syscall.CLONE_NEWIPC |
			//syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUTS, //|
		//syscall.CLONE_NEWPID,
		Credential: &syscall.Credential{Uid: 0, Gid: 0, NoSetGroups: true},
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}
	fmt.Println(c.String())
	output, err := c.CombinedOutput()
	fmt.Printf("%s\n", output)
	if err != nil {
		panic(err)
	}
}
