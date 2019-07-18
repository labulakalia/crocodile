package main

import (
	"crocodile/common/cfg"
	"crocodile/common/log"
	"crocodile/common/registry"
	"crocodile/common/wrapper"
	"crocodile/service/executor/execute"
	"crocodile/service/executor/subscriber"
	"github.com/labulaka521/logging"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-plugins/wrapper/trace/opentracing"
	goopentracing "github.com/opentracing/opentracing-go"
	"time"
)

func main() {
	var (
		err error
	)
	cfg.Init()
	log.Init()
	t, io, err := wrapper.NewTracer("crocodile.srv.executor", "")
	if err != nil {
		logging.Fatal(err)
	}
	defer io.Close()
	goopentracing.SetGlobalTracer(t)

	service := micro.NewService(
		micro.Name("crocodile.srv.executor"),
		micro.Version("latest"),
		micro.Registry(registry.Etcd(cfg.EtcdConfig.Endpoints...)),
		micro.RegisterInterval(15*time.Second),
		micro.RegisterTTL(30*time.Second),
		micro.Broker(
			broker.NewBroker(
				broker.Addrs(cfg.ExecuteConfig.Address),
				broker.Registry(registry.Etcd(cfg.EtcdConfig.Endpoints...)),
			),
		),
		micro.WrapSubscriber(opentracing.NewSubscriberWrapper()),
	)

	// Initialise service
	service.Init(
		micro.Action(func(c *cli.Context) {
			execute.Init(service.Client())
		}),
	)

	bk := service.Server().Options().Broker

	if err = bk.Connect(); err != nil {
		logging.Fatalf("broker Connect Err:%v", err)
	}

	executor := &subscriber.Executor{
		PubSub: bk,
	}

	//Register Struct as Subscriber
	err = micro.RegisterSubscriber("crocodile.srv.executor", service.Server(), executor)
	if err != nil {
		logging.Fatalf("Register Subscriber Err: %v", err)
	}

	// Run service
	if err := service.Run(); err != nil {
		logging.Fatal(err)
	}
}
