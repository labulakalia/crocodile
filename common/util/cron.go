package util

import (
	"github.com/gorhill/cronexpr"
	"time"
)

// 获取下一个运行的时间点
func NextTime(cronLine string, fromtime time.Time) (nexttime time.Time, err error) {
	var (
		expr *cronexpr.Expression
	)
	if expr, err = cronexpr.Parse(cronLine); err != nil {
		return
	}
	nexttime = expr.Next(fromtime)

	return
}
