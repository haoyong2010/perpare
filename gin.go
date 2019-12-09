package main

import "os/exec"

func main() {
	exec.Command("/bin/bash", "-c", "echo 1;echo 2;")
}
