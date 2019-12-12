package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

//任务执行时间
type TimePoint struct {
	StartTime int64
	EndTime   int64
}

//一条日志
type LogRecord struct {
	JobName   string    `bson:"jobName"` //任务名
	Command   string    `bson:"command"` //shell命令
	Err       string    `err`            //脚本错误
	Content   string    `content`        //脚本输出
	TimePoint TimePoint `timePoint`      //执行时间点
}

//jobName过滤条件
type FindByJobName struct {
	JobName string `bson:"jobName"`
}

func main() {
	var (
		client     *mongo.Client
		err        error
		database   *mongo.Database
		collection *mongo.Collection
		cond       *FindByJobName
		cursor     *mongo.Cursor
		record     *LogRecord
	)
	//建立链接
	if client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://127.0.0.1:27017"), options.Client().SetConnectTimeout(5*time.Second)); err != nil {
		fmt.Println("client failed err:", err)
		return
	}
	//链接数据库
	database = client.Database("cron")
	//链接数据表
	collection = database.Collection("log")
	//按照jobName字段过滤，想找出jobname=job10,找出
	cond = &FindByJobName{JobName: "job10"}
	//查询数据
	opts := new(options.FindOptions)
	limit := int64(2)
	skip := int64((2 - 1) * 2)

	opts.Limit = &limit
	opts.Skip = &skip
	if cursor, err = collection.Find(context.TODO(), cond, opts); err != nil {
		fmt.Println("find failed err:", err)
		return
	}
	//释放游标
	defer cursor.Close(context.TODO())
	//遍历结果集
	for cursor.Next(context.TODO()) {
		//定义一个日志对象
		record = &LogRecord{}
		//反序列化bson到对象
		if err = cursor.Decode(record); err != nil {
			fmt.Println("decode failed err:", err)
			return
		}
		//把日志行打印出来
		fmt.Println(*record)
	}

}
