package subscriber

import (
	"context"
	"github.com/micro/go-micro/util/log"

	actuator "crocodile/service/actuator/proto/actuator"
)

type Actuator struct{}

func (e *Actuator) Handle(ctx context.Context, msg *actuator.Message) error {
	log.Log("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *actuator.Message) error {
	log.Log("Function Received message: ", msg.Say)
	return nil
}
