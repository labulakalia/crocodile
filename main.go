package main

import (
	"fmt"
	"time"
)

var getsql = `SELECT 
			name,
			taskid,
			starttime,
			endtime,
			totalruntime,
			status,
			trigger,
			errcode,
			errmsg,
			errtasktype,
			errtaskid,
			errtask
		FROM 
			crocodile_log
		WHERE 
			name=?`

func main() {

	now := time.Now().UnixNano() / 1e6

	fmt.Println("v1.1.8" < "v1.2.1")

	fmt.Println(now, int64(time.Second/1e6))
}
