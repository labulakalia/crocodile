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
	"net"
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

	switch exMsg.Event {
	case event.Run_Task:
		_, port, _ = net.SplitHostPort(e.PubSub.Address())
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
		go func() {
			if err = execute.RunTask(ctx, exMsg); err != nil {
				logging.Errorf("RunTask Err: %v", err)
			}
		}()

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
