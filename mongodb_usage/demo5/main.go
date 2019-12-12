package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

//starttime 小于某时间
//	key :   val
//{"$lt":timestamp}
type TimeBeforecond struct {
	Before int64 `bson:"$lt"`
}

//key					: val
//{"timePoint.startTime":{"$lt":timestamp}}
type DeleteCont struct {
	beforeCond TimeBeforecond `bson:"timePoint.startTime"`
}

func main() {
	var (
		client     *mongo.Client
		err        error
		database   *mongo.Database
		collection *mongo.Collection
		delCond    *DeleteCont
		delResp    *mongo.DeleteResult
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
	//删除开始时间早于当前时间的所有日志
	//delete({"timePoint.startTime":{"$lt":当前时间}})
	delCond = &DeleteCont{beforeCond: TimeBeforecond{Before: time.Now().Unix()}}
	//执行删除
	if delResp, err = collection.DeleteMany(context.TODO(), delCond); err != nil {
		fmt.Println("delete failed err:", err)
		return
	}
	fmt.Println("删除的行数：", delResp.DeletedCount)
}
