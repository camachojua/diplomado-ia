package main

import (
	"fmt"
	"os"
	"os/exec"
)

// go build container.go, then go run container.go cmd args.
// For example go run container.go run echo build youn own container in Go
// For root commands such as changing hostname use ... sudo ./container.go run /bin/bash

func main() {
	switch os.Args[1] {
	case "run":
		run()
	default:
		panic("error")
	}
}

func run() {
	fmt.Printf("running %v\n", os.Args[2:])

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
