package main

import (
	"fmt"
	"time"
)

func main() {
	// client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})

	// fmt.Println(client.Set("11", int(define.TsWait), 0).Err())

	// res, _ := client.Get("11").Int()
	// fmt.Println(define.TaskStatus(res))

	// res := "2020-04-15 21:51:00: Task Run Finished,Return Code:    0\n"
	// codebyte := res[len(res)-5:]
	// fmt.Println(strings.TrimSpace(string(codebyte)))
	// code, err := strconv.Atoi(strings.TrimLeft(string(codebyte), " "))
	// fmt.Println(code, err)
	fmt.Println(time.Now().UnixNano() / 1e6)
	fmt.Println(time.Now().Unix())
}
