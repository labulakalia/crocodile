package subscriber

import (
	"context"
	"crocodile/common/event"
	"crocodile/service/executor/execute"
	pbexecutor "crocodile/service/executor/proto/executor"
	"fmt"
	"github.com/labulaka521/logging"
	"github.com/micro/go-micro/broker"
	"github.com/micro/util/go/lib/addr"
	"strings"
)

type Executor struct {
	PubSub broker.Broker
}

// 接收运行 或者强杀任务的消息事件
func (e *Executor) ExecEvent(ctx context.Context, exMsg *pbexecutor.ExecuteMsg) (err error) {
	var (
		port string
		run  bool
	)
	// 判断任务执行的主机IP
	logging.Debugf("Recevice Event %+v", exMsg)

	switch exMsg.Event {
	case event.Run_Task:
		port = strings.Split(e.PubSub.Options().Addrs[0], ":")[1]
		for _, ip := range addr.IPs() {
			address := fmt.Sprintf("%s:%s", ip, port)
			if address == exMsg.Runhost {
				run = true
			}
		}

		if !run {
			return
		}
		// 运行任务的事件
		if err = execute.RunTask(ctx, exMsg); err != nil {
			logging.Errorf("RunTask Err: %v", err)
		}
	case event.Kill_Task:
		// 强杀任务的事件
		if err = execute.KillTask(exMsg); err != nil {
			logging.Errorf("KillTask Err: %v", err)
		}
	default:
		logging.Warn("Unsupport exMsg Event ID", exMsg.Event)
	}
	return
}
