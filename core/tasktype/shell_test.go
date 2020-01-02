package tasktype

import (
	"context"
	"testing"
	"time"
)

func TestDataShell_Run(t *testing.T) {
	var data TaskRuner = &DataShell{
		Name: "ls",
	}
	resp := data.Run(context.Background())
	t.Log(resp)
}

func TestDataApi_Run(t *testing.T) {
	//var data TaskRuner = &DataApi{
	//	Url:"http://httpbin.org",
	//	Method:"GET",
	//}
	//resp := data.Run(context.Background())
	//t.Log(resp)

	t.Log(time.Now().UnixNano() / 1e6)
	t.Log(time.Now().Unix())
}
