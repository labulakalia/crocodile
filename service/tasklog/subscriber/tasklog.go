package subscriber

import (
	"context"
	"github.com/micro/go-micro/util/log"

	tasklog "crocodile/service/tasklog/proto/tasklog"
)

type Tasklog struct{}

func (e *Tasklog) Handle(ctx context.Context, msg *tasklog.Message) error {
	log.Log("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *tasklog.Message) error {
	log.Log("Function Received message: ", msg.Say)
	return nil
}
