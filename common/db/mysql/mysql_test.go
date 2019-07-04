package mysql_test

import (
	"crocodile/common/db/mysql"
	"github.com/labulaka521/logging"
	"testing"
)

func TestNewClient(t *testing.T) {
	logging.SetLogLevel("FATAL")
	logging.Setup()
	dsn := "root:wang109097@tcp(127.0.0.1:3306)/crocodile?charset=utf8mb4&parseTime=true"
	_ = mysql.New(dsn, 1, 1)

	t.Log("test")
}
