package main

// containerV5 In this version, the container:
// (1) isolate/proctect the image
// (2) Set a new PID
// (3) Spawn a child process with

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("error")
	}
}

func run() {
	fmt.Printf("parent running as PID %d\n", os.Getpid())

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// set UTS, PID for child
	cmd.SysProcAttr = &syscall.SysProcAttr{Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID}
	must(cmd.Run())
}

func child() {
	fmt.Printf("child running %v as PID %d\n", os.Args[2:], os.Getpid())

	syscall.Sethostname(([]byte("gcv5")))

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(cmd.Run())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
