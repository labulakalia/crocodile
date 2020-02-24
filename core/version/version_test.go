package version

import (
	"testing"

	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/model"
	mylog "github.com/labulaka521/crocodile/core/utils/log"
)

func Test_checkverson(t *testing.T) {
	config.Init("/Users/labulakalia/workerspace/golang/crocodile/core/config/core.toml")
	mylog.Init()
	model.InitDb()
	checkverson()
}
