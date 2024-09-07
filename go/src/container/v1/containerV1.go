package main

import (
	"fmt"
	"os"
)

func run() {
	fmt.Printf("running %v\n", os.Args[2:])
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

// go run containerV1.go run echo hello containerV1
func main() {
	switch os.Args[1] {
	case "run":
		run()
	default:
		panic("error: bad command")
	}
}
