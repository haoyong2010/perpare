package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

type result struct {
	err    error
	output []byte
}

func main() {
	var (
		resultChan chan *result
		res        *result
	)
	//创建结果队列
	resultChan = make(chan *result, 1000)
	//执行1个cmd,让它在一个协程执行，执行2秒，输出hello
	//1秒的时候杀死cmd
	ctx, cancleFunc := context.WithCancel(context.TODO())
	go func() {
		cmd := exec.CommandContext(ctx, "/bin/bash", "-c", "sleep 2;echo hello;")
		//执行任务，捕获输出
		output, err := cmd.CombinedOutput()
		//把任务输出结果，传给main协程
		resultChan <- &result{
			err:    err,
			output: output,
		}
	}()
	time.Sleep(time.Second)
	//取消上下文
	cancleFunc()

	//在main协程里，等待子协程的退出，并打印任务执行结果
	res = <-resultChan
	//打印任务执行结果
	fmt.Println(res.err, string(res.output))
}
