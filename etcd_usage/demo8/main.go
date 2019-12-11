package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	var (
		config clientv3.Config
		client *clientv3.Client
		err    error
		kv     clientv3.KV
		putOp  clientv3.Op
		opResp clientv3.OpResponse
		getOp  clientv3.Op
	)
	//初始化配置
	config = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}
	//建立链接
	if client, err = clientv3.New(config); err != nil {
		fmt.Println("client failed err:", err)
		return
	}
	//kv
	kv = clientv3.NewKV(client)
	//创建PutOp:operation
	putOp = clientv3.OpPut("/cron/jobs/job8", "222")
	//执行putOp
	if opResp, err = kv.Do(context.TODO(), putOp); err != nil {
		fmt.Println("put op failed err:", err)
		return
	}
	fmt.Println("写入Revision:", opResp.Put().Header.Revision)
	//创建GetOp
	getOp = clientv3.OpGet("/cron/jobs/job8")
	//执行getOp
	if opResp, err = kv.Do(context.TODO(), getOp); err != nil {
		fmt.Println("get op failed err:", err)
		return
	}
	//打印数据
	fmt.Println("获取Resvision:", opResp.Get().Kvs[0].ModRevision)
	fmt.Println("数据value:", string(opResp.Get().Kvs[0].Value))
}
