package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		err     error
		kv      clientv3.KV
		delResp *clientv3.DeleteResponse
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
	//数据操作
	kv = clientv3.NewKV(client)
	//删除操作
	if delResp, err = kv.Delete(context.TODO(), "/cron/jobs/", clientv3.WithPrefix()); err != nil {
		fmt.Println("delete failed err:", err)
		return
	}
	//被删除之前的value是什么
	if len(delResp.PrevKvs) != 0 {
		for _, kvpair := range delResp.PrevKvs {
			fmt.Println("删除了：", string(kvpair.Key), string(kvpair.Value))
		}
	}

}
