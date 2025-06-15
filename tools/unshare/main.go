package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
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

func main() {
	//as soon as you add the third command, this no longer functions
	unshare("/home/lordralex/testchroot", "pwd")                        //are we in the right place
	unshare("/home/lordralex/testchroot", "ls", "-l")                   //do we see anything
	unshare("/home/lordralex/testchroot", "whoami")                     //are we the correct user
	unshare("/home/lordralex/testchroot", "touch", "test")              //can we write and it persist
	unshare("/home/lordralex/testchroot", "curl", "1.1.1.1")            //can we access an IP
	unshare("/home/lordralex/testchroot", "curl", "google.com")         //does DNS work
	unshare("/home/lordralex/testchroot", "curl", "https://google.com") //does SSL work
}

func unshare(dir, cmd string, args ...string) {
	var err error

	c := exec.Command("bash", "-c",
		strings.Join([]string{"mkdir -p {dev,bin,usr,lib,lib64,etc,tmp,run/systemd/resolve," + strings.TrimPrefix(dir, "/") + "}",
			"mount -t tmpfs -o size=100m tmpfs tmp",
			"mount --bind /bin bin",
			"mount --bind /lib lib",
			"mount --bind /lib64 lib64",
			"mount --rbind /etc etc",
			"mount --rbind /dev dev",
			"mount --rbind /run/systemd/resolve run/systemd/resolve",
			"mount --bind " + dir + " " + strings.TrimPrefix(dir, "/"),
			"mount --rbind / .",
			fmt.Sprintf("unshare -UR . -w %s --map-user=%d --map-group=%d %s %s", dir, os.Getuid(), os.Getgid(), cmd, strings.Join(args, " ")),
		}, " && "))
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
