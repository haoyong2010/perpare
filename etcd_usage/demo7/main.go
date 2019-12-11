package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"time"
)

func main() {
	var (
		config clientv3.Config
		client *clientv3.Client
		err    error
		kv     clientv3.KV
		//putResp *clientv3.PutResponse
		//delResp *clientv3.DeleteResponse
		getResp        *clientv3.GetResponse
		watchStartResp int64
		watcher        clientv3.Watcher
		watchRespChan  <-chan clientv3.WatchResponse
		watchResp      clientv3.WatchResponse
		event          *clientv3.Event
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
	//链接数据操作kv
	kv = clientv3.NewKV(client)
	//模拟etcd中kv的变化
	go func() {
		for {
			//if _, err = kv.Put(context.TODO(), "/cron/jobs/job7", "I am job7"); err != nil {
			//	fmt.Println("put crontab err:", err)
			//}

			//if _, err = kv.Delete(context.TODO(), "cron/jobs/jbo7"); err != nil {
			//	fmt.Println("del crontab err:", err)
			//}
			kv.Put(context.TODO(), "/cron/jobs/job7", "I am job7")
			kv.Delete(context.TODO(), "/cron/jobs/job7")
			time.Sleep(1 * time.Second)
		}

	}()
	//time.Sleep(time.Second * 100)
	//先get当前的值，并监听后续的变化
	if getResp, err = kv.Get(context.TODO(), "/cron/jobs/job7"); err != nil {
		fmt.Println("get crontab err:", err)
		return
	}
	//如果key有值
	if len(getResp.Kvs) != 0 {
		fmt.Println("当前值为：", string(getResp.Kvs[0].Value))
	}
	//当前etcd集群事物ID，单调递增的
	watchStartResp = getResp.Header.Revision + 1
	//创建一个监听器watch
	watcher = clientv3.NewWatcher(client)
	//启动监听
	fmt.Println("从该版本向后监听：", watchStartResp)
	ctx, canceFun := context.WithCancel(context.TODO())
	time.AfterFunc(5*time.Second, func() {
		canceFun()
	})
	watchRespChan = watcher.Watch(ctx, "/cron/jobs/job7", clientv3.WithRev(watchStartResp))
	//处理kv变化事件
	for watchResp = range watchRespChan {
		for _, event = range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("修改为：", string(event.Kv.Value), "Revision:", event.Kv.CreateRevision, event.Kv.ModRevision)
			case mvccpb.DELETE:
				fmt.Println("删除了,Revision:", event.Kv.ModRevision)
			}
		}
	}
}
