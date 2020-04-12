package schedule

import (
	"fmt"
	"testing"

	"github.com/go-redis/redis"
)

func Test_task2_GetCronExpr(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	fmt.Println(client.Exists("testtask").Result())
	// v := map[string]interface{}{
	// 	"cronexpr1": 2,
	// }
	// err := client.HMSet("testtask", v).Err()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// res, err := client.HGet("testtask", "cronexpr1").Int()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(res)

}
