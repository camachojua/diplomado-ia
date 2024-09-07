package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

func run() {
	detect_machine_type()

	fmt.Printf("parent running as PID %d\n", os.Getpid())

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// set UTS, PID for child
	cmd.SysProcAttr = &syscall.SysProcAttr{Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID}
	must(cmd.Run())
}

func child(current_dir string) {
	linux_flavor := ""
	command := ""
	args := []string{}

	// We need to create the proper folders
	ensure_directory(current_dir, "alpine")
	ensure_directory(current_dir, "ubuntu")

	// We need a linux base filesystem for chrooting as an example we can use:
	// Ubuntu base => http://cdimage.ubuntu.com/ubuntu-base/releases/24.04/release/
	// Alpine mini root filesystem => https://www.alpinelinux.org/downloads/
	// For this we need to download the corresponding .tar.gz file for our target architecture
	// After that we need to create two directories: mkdir alpine ubuntu
	// And we need to extract the filesystem to those directories:
	// tar xzf ~/Downloads/ubuntu-*.tar.gz -C ubuntu/
	// tar xzf ~/Downloads/alpine-minirootfs-*.tar.gz -C alpine/
	switch os.Args[2] {
	case "alpine":
		linux_flavor = "alpine"
		get_base_fs(linux_flavor, detect_machine_type())
		fmt.Printf("child running %v as PID %d\n", os.Args[3:], os.Getpid())
		command = os.Args[3]
		args = os.Args[4:]
	case "ubuntu":
		linux_flavor = "ubuntu"
		get_base_fs(linux_flavor, detect_machine_type())
		fmt.Printf("child running %v as PID %d\n", os.Args[3:], os.Getpid())
		command = os.Args[3]
		args = os.Args[4:]
	default:
		linux_flavor = "ubuntu"
		get_base_fs(linux_flavor, detect_machine_type())
		fmt.Printf("child running %v as PID %d\n", os.Args[2:], os.Getpid())
		command = os.Args[2]
		args = os.Args[3:]
	}

	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	syscall.Sethostname(([]byte("gcv6_" + linux_flavor)))
	fmt.Printf("Chrooting the system at: %v\n", current_dir+"/"+linux_flavor+"/")
	must(syscall.Chroot(current_dir + "/" + linux_flavor + "/"))
	must(os.Chdir("/"))
	must(syscall.Mount("proc", "proc", "proc", 0, ""))
	defer syscall.Unmount("proc", 0)

	must(cmd.Run())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func usage(current_file string) {
	fmt.Print("Usage: " + current_file + " run|child [ubuntu|alpine] command\n" +
		"When using 'child' you can specify the flavor of the base linux system (ubuntu|alpine)\n")
	os.Exit(3)
}

func decompress_file(linux_flavor string) {
	// Before decompressing we must be shure that the base filesystem doesn't exist in the target directory
	if _, err := os.Stat(linux_flavor + "/bin/ls"); err == nil {
		fmt.Printf("=> File system found\n")
		return
	}

	fmt.Printf("Starting file decompression...")
	cmd := exec.Command("tar", "xzf", linux_flavor+".tar.gz", "-C", "./"+linux_flavor+"/")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	must(cmd.Run())
	fmt.Printf("[ok]\nDecompression completed\n")
}

func download_file(linux_flavor string, url string) {
	out, _ := os.Create(linux_flavor + ".tar.gz")
	defer out.Close()

	req, _ := http.NewRequest("GET", url, nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	buf := make([]byte, 32*1024)
	var downloaded int64

	for {
		n, err := resp.Body.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error while downloading: %v", err)
			panic("")
		}
		if n > 0 {
			out.Write(buf[:n])
			downloaded += int64(n)
			fmt.Printf("\rDownloading %v ... %.2f%%", linux_flavor, float64(downloaded)/float64(resp.ContentLength)*100)
		}
	}
}

func get_base_fs(linux_flavor, machine_type string) {
	// If the base filesystem exists we do not download it form the internet
	if _, err := os.Stat(linux_flavor + ".tar.gz"); err == nil {
		decompress_file(linux_flavor)
		return
	}

	url := ""

	if strings.Contains(machine_type, "x86_64") {
		url = "http://cdimage.ubuntu.com/ubuntu-base/releases/24.04/release/ubuntu-base-24.04.1-base-amd64.tar.gz"
		if strings.Contains(linux_flavor, "alpine") {
			url = "https://dl-cdn.alpinelinux.org/alpine/v3.20/releases/x86_64/alpine-minirootfs-3.20.3-x86_64.tar.gz"
		}
		download_file(linux_flavor, url)
	}

	if strings.Contains(machine_type, "aarch64") || strings.Contains(machine_type, "arm64") || strings.Contains(machine_type, "armv7") {
		url = "http://cdimage.ubuntu.com/ubuntu-base/releases/24.04/release/ubuntu-base-24.04.1-base-arm64.tar.gz"
		if strings.Contains(linux_flavor, "alpine") {
			url = "https://dl-cdn.alpinelinux.org/alpine/v3.20/releases/aarch64/alpine-minirootfs-3.20.3-aarch64.tar.gz"
		}
		download_file(linux_flavor, url)
	}

	if url == "" {
		panic("CPU unsupported.")
	}

	fmt.Printf("Download finished!!!!!!!!!!!!!!!!\n")
	decompress_file(linux_flavor)
}

func ensure_directory(current_dir, dirname string) {
	directory := filepath.Join(current_dir, dirname)
	must(os.MkdirAll(directory, os.ModePerm))
}

func detect_machine_type() string {
	var utsname syscall.Utsname
	syscall.Uname(&utsname)
	machine := (*[len(utsname.Machine)]byte)(unsafe.Pointer(&utsname.Machine))
	machine_type := string(machine[:])
	fmt.Printf("Our current machine type is: %v\n", machine_type)
	return machine_type
}

func main() {
	current_file, err := os.Executable()

	if len(os.Args) < 2 {
		usage(current_file)
	}

	if err != nil {
		panic(err)
	}

	current_dir := filepath.Dir(current_file)
	fmt.Printf("Our current working dir is: %v\n", current_dir)

	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child(current_dir)
	default:
		usage(current_file)
	}
}
