package registry

import (
	"crocodile/common/cfg"
	"crocodile/common/log"
	"testing"
)

func TestGetEtcdListServices(t *testing.T) {
	log.Init()
	cfg.Init()
	GetEtcdListServices(cfg.EtcdConfig.Endpoints...)
}
