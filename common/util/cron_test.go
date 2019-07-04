package util

import (
	"testing"
	"time"
)

func TestNextTime(t *testing.T) {
	var (
		next time.Time
		err  error
	)
	next = time.Now()
	t.Log(next)
	next, err = NextTime("*/5 * * * * * *", next)
	err = err
	t.Log(next)
	//for i:= 0;i<=10;i++ {
	//	now := time.Now()
	//
	//
	//	t.Log(next,err)
	//	tk:=time.NewTicker(next.Sub(now))
	//	select {
	//	case <-tk.C:
	//		t.Log("next")
	//	}
	//}
}
