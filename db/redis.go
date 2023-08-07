package db

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var client *redis.Client

func init() {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis 服务器地址和端口
		DB:   0,                // 使用的 Redis 数据库索引，默认为 0
	})

	// 測試連線
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Failed to connect to Redis:", err)
	} else {
		fmt.Println("Redis connect success")
	}
}

func GetRedisClient() *redis.Client {
	return client
}

func CloseRedisClient() {
	err := client.Close()
	if err != nil {
		fmt.Println("Failed to close Redis:", err)
	} else {
		fmt.Println("Closed Redis connection.")
	}
}

func SetRedis(key string, value string) {
	//設置過期時間 1小時
	err := client.Set(context.Background(), key, value, time.Hour).Err()
	if err != nil {
		fmt.Println("Failed to set data:", err)
	} else {
		fmt.Println("Set data successfully.")
	}
}

func GetRedis(key string) string {
	val, err := client.Get(context.Background(), key).Result()
	if err != nil {
		fmt.Println("Failed to get data:", err)
		return ""
	} else {
		fmt.Println("Get data successfully.")
		return val
	}
}

func DelRedis(client *redis.Client, key string) {
	err := client.Del(context.Background(), key).Err()
	if err != nil {
		fmt.Println("Failed to delete data:", err)
	} else {
		fmt.Println("Delete data successfully.")
	}
}
