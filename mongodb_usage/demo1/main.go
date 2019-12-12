package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func main() {
	//client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://@localhost:27017"))
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	//defer cancel()
	//err = client.Connect(ctx)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//建立链接
	var (
		client     *mongo.Client
		err        error
		database   *mongo.Database
		collection *mongo.Collection
	)
	if client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://@127.0.0.1:27017"), options.Client().SetConnectTimeout(5*time.Second)); err != nil {
		fmt.Println("client failed err :", err)
		return
	}
	//选择数据库my_db
	database = client.Database("my_db")
	// 选择表my_collection
	collection = database.Collection("my_collection")

	collection = collection
}
