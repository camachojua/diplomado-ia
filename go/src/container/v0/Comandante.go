package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func f1() {
	cmd := exec.Command("tr", "a-z", "A-Z")
	cmd.Stdin = strings.NewReader("some input")
	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("in all caps: %q\n", out.String())
}

func f2() {
	cmd := exec.Command("ls", "/usr/local/bin")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func f3() error {
	log, err := os.Create("lofas.log")
	if err != nil {
		return err
	}

	defer log.Close()
	cmd := exec.Command("ls", "/usr/local/bin")
	cmd.Stdout = log
	cmd.Stderr = log
	return cmd.Run()
}

func f4() error {
	r, w, err := os.Pipe()
	if err != nil {
		return err
	}

	defer r.Close()
	ls := exec.Command("ls", "/home/juan/Src/diplomado-ia/go/src/student_submissions/movielens/")
	ls.Stdout = w
	err = ls.Start()

	defer ls.Wait()
	w.Close()

	grep := exec.Command("grep", "go")
	grep.Stdin = r
	grep.Stdout = os.Stdout
	return grep.Run()
}

func main() {
	//f1()
	//f2()
	//f3()
	f4()
}
