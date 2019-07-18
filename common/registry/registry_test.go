package registry

import (
	"crocodile/common/cfg"
	"crocodile/common/log"
	"testing"
	"time"
)

func TestGetEtcdListServices(t *testing.T) {
	log.Init()
	cfg.Init()
	for i := 0; i < 100; i++ {
		GetEtcdListServicesIP("topic:crocodile.srv.executor", cfg.EtcdConfig.Endpoints...)
		time.Sleep(time.Second)
	}

}
