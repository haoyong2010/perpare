package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

//代表一个任务
type CronJob struct {
	expr     *cronexpr.Expression
	nextTime time.Time
}

func main() {
	//需要有1个调度协程，它定时检查所有的cron任务，谁过期了就执行谁
	var (
		cronJob       *CronJob
		expr          *cronexpr.Expression
		now           time.Time
		err           error
		scheduleTable map[string]*CronJob
	)
	scheduleTable = make(map[string]*CronJob)
	//当前时间
	now = time.Now()
	//定义2个cronJob
	expr, err = cronexpr.Parse("*/5 * * * * * *")
	if err != nil {
		fmt.Println("expr failed err :", err)
	}
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}
	//任务注册到调度表
	scheduleTable["job1"] = cronJob

	expr, err = cronexpr.Parse("*/5 * * * * * *")
	if err != nil {
		fmt.Println("expr failed err :", err)
	}
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}
	//任务注册到调度表
	scheduleTable["job2"] = cronJob

	//启动调度协程
	go func() {
		//定时检查任务调度表
		for {
			now = time.Now()
			for jobName, cronJob := range scheduleTable {
				//判断是否过期
				if cronJob.nextTime.Before(now) || cronJob.nextTime.Equal(now) {
					//启动协程执行这个任务
					go func(jobName string) {
						fmt.Println("执行：", jobName)
					}(jobName)

					//计算下一次调度时间

					cronJob.nextTime = cronJob.expr.Next(now)
					fmt.Println(jobName, "下次执行时间：", cronJob.nextTime)
				}
			}

			select {
			case <-time.NewTimer(100 * time.Millisecond).C: //将在100毫秒可读，返回
			}
		}
	}()
	time.Sleep(time.Second * 100)
}
