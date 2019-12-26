package config

import (
	"testing"
)

func TestInit(t *testing.T) {
	Init("/Users/labulakalia/workerspace/golang/crocodile/core/config/core.toml")
	t.Logf("%+v", CoreConf)
}
