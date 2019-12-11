package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	var (
		config         clientv3.Config
		client         *clientv3.Client
		err            error
		lease          clientv3.Lease
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId        clientv3.LeaseID
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
		keepResp       *clientv3.LeaseKeepAliveResponse
		ctx            context.Context
		cancleFunc     context.CancelFunc
		kv             clientv3.KV
		txn            clientv3.Txn
		txnResp        *clientv3.TxnResponse
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
	//1￿￿·上锁（创建租约，自动续租，拿着租约去抢占一个key)
	lease = clientv3.NewLease(client)
	//申请一个5秒的租约
	if leaseGrantResp, err = lease.Grant(context.TODO(), 5); err != nil {
		fmt.Println("grant lease failed err:", err)
		return
	}
	//取得租约id
	leaseId = leaseGrantResp.ID
	//准备一个用于取消自动续租的context
	ctx, cancleFunc = context.WithCancel(context.TODO())
	defer cancleFunc()
	//立即销毁租约
	defer lease.Revoke(context.TODO(), leaseId)
	//取消自动续约
	if keepRespChan, err = lease.KeepAlive(ctx, leaseId); err != nil {
		fmt.Println("lease keepAlive failed err:", err)
		return
	}
	//处理续租应答的协程
	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepRespChan == nil {
					fmt.Println("租约已经失效")
					goto END
				} else { //每秒会自动续约一次，所以就会收到一次应答
					fmt.Println("收到自动续租应答：", keepResp.ID)
				}
			}
		}
	END:
	}()

	//如果不存在key，那么设置它，否则强锁失败
	kv = clientv3.NewKV(client)
	//创建事务
	txn = kv.Txn(context.TODO())
	//定义事务
	//如果key不存在
	txn.If(clientv3.Compare(clientv3.CreateRevision("/cron/jobs/job9"), "=", 0)).Then(clientv3.OpPut("/cron/jobs/job9", "xxx", clientv3.WithLease(leaseId))).Else(clientv3.OpGet("/cron/jobs/job9")) //抢锁失败
	//提交事务
	if txnResp, err = txn.Commit(); err != nil {
		fmt.Println("commit txn failed err:", err)
		return
	}
	//判断是否抢到了锁
	if !txnResp.Succeeded {
		fmt.Println("锁被占用：", string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
		return
	}
	//2·处理事务
	fmt.Println("处理任务")
	time.Sleep(time.Second * 5)
	//在锁内
	//3·释放锁（取消自动续租，释放租约）
	//defer 会把租约释放掉，关联的kv就被删除了
}
