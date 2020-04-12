package main

import (
	"fmt"

	"github.com/go-redis/redis/v7"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	// client.RPush("test1112", 111)
	// client.RPush("test1112", 222)
	// client.RPush("test1112", 333)
	// client.RPush("test1112", 444)
	// client.RPush("test1112", 555)
	// client.RPush("test1112", 666)
	client.SMembers(key string)
}
