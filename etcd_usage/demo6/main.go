package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	var (
		config       clientv3.Config
		client       *clientv3.Client
		err          error
		lease        clientv3.Lease
		leaseGetResp *clientv3.LeaseGrantResponse
		leaseId      clientv3.LeaseID
		kv           clientv3.KV
		putResp      *clientv3.PutResponse
		getResp      *clientv3.GetResponse
		keepResp     *clientv3.LeaseKeepAliveResponse
		keepRespChan <-chan *clientv3.LeaseKeepAliveResponse
	)
	//初始化配置
	config = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}
	//链接
	if client, err = clientv3.New(config); err != nil {
		fmt.Println("client failed err:", err)
		return
	}
	//申请一个lease(租约)
	lease = clientv3.Lease(client)
	//申请一个10秒的租约
	if leaseGetResp, err = lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println("Grant failed err:", err)
		return
	}
	//获取租约的ID
	leaseId = leaseGetResp.ID
	//自动续约
	//ctx, _ := context.WithTimeout(context.TODO(), 5*time.Second)
	//续租5秒，停止续租，10秒的生命期，15秒的生命
	//5秒后会自动取消自动续约
	if keepRespChan, err = lease.KeepAlive(context.TODO(), leaseId); err != nil {
		fmt.Println("keepAlive failed err:", err)
		return
	}
	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepRespChan == nil {
					fmt.Println("租约已失效")
					goto END
				} else { //每秒会续约一次，所以就会收到一次应答
					fmt.Println("收到自动续约应答：", keepResp.ID)

				}
			}
		}
	END:
	}()

	//获取kv API子集
	kv = clientv3.NewKV(client)
	//PUT一个kv，与租约关联，从而实现10秒自动过期
	if putResp, err = kv.Put(context.TODO(), "/cron/lock/job1", "", clientv3.WithLease(leaseId)); err != nil {
		fmt.Println("put failed err:", err)
		return
	}
	fmt.Println("写入成功:", putResp.Header.Revision)
	for {
		if getResp, err = kv.Get(context.TODO(), "/cron/lock/job1"); err != nil {
			fmt.Println("get failed err:", err)
			return
		}
		if getResp.Count == 0 {
			fmt.Println("kv过期了")
			return
		}
		fmt.Println("还未过期：", getResp.Kvs)
		time.Sleep(time.Second * 2)
	}
}
