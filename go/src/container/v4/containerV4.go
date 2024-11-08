package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func run() {
	fmt.Printf("running %v as PID %d\n", os.Args[2:], os.Getpid())

	cmd := exec.Command(os.Args[2], os.Args[3:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// set flags for UTS
	cmd.SysProcAttr = &syscall.SysProcAttr{Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID}

	must(cmd.Run())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	switch os.Args[1] {
	case "run":
		run()
	default:
		panic("error")
	}
}
