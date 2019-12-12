package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

//任务执行时间点
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

func main() {
	var (
		client     *mongo.Client
		err        error
		database   *mongo.Database
		collection *mongo.Collection
		record     *LogRecord
		result     *mongo.InsertOneResult
		docId      primitive.ObjectID
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
	//插入记录
	record = &LogRecord{
		JobName: "job10",
		Command: "echo hello",
		Err:     "",
		Content: "hello",
		TimePoint: TimePoint{
			StartTime: time.Now().Unix(),
			EndTime:   time.Now().Unix() + 10,
		},
	}
	if result, err = collection.InsertOne(context.TODO(), record); err != nil {
		fmt.Println("insertOne failed err:", err)
		return
	}
	//_id:默认生成一个全局唯一ID，ObjectID：12字节的二进制
	docId = result.InsertedID.(primitive.ObjectID)
	fmt.Println("自增ID：", docId.Hex())
}
