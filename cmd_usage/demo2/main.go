package main

import (
	"fmt"
	"os/exec"
)

func main() {
	var (
		cmd    *exec.Cmd
		output []byte
		err    error
	)
	//生成Cmd
	cmd = exec.Command("/bin/bash", "-c", "sleep 1;ls -l")
	//执行命令，捕获进程的输出pipe
	if output, err = cmd.CombinedOutput(); err != nil {
		fmt.Println("err=", err)
		return
	}
	//打印子进程的输出
	fmt.Println(string(output))
}
