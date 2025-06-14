package main

import (
	"fmt"
	"os"
	"os/exec"
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

func main() {
	unshare("/usr/bin/id")
}

func chroot(cmd string, args ...string) {
	attrs := &syscall.SysProcAttr{
		Chroot: ".",
	}

	c := exec.Command(cmd, args...)
	c.SysProcAttr = attrs
	c.Dir = "/"

	output, err := c.CombinedOutput()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", output)
}

func unshare(cmd string, args ...string) {
	c := exec.Command(cmd, args...)
	c.SysProcAttr = &syscall.SysProcAttr{
		Unshareflags: syscall.CLONE_NEWUSER |
			syscall.CLONE_NEWNS |
			syscall.CLONE_FILES |
			syscall.CLONE_NEWCGROUP |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID,
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
	//fmt.Println(c.String())
	output, err := c.CombinedOutput()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", output)
}
