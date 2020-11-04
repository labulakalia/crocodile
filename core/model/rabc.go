package model

import (
	"context"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/config"
	"go.uber.org/zap"
)

var (
	enforcer *casbin.Enforcer
)

// InitRabc init rabc
func InitRabc() {
	modeltext := `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)
`
	dbcfg := config.CoreConf.Server.DB
	m, err := model.NewModelFromString(modeltext)
	if err != nil {
		log.Panic("NewModelFromString Err", zap.Error(err))
	}
	a, err := gormadapter.NewAdapter(dbcfg.Drivename, dbcfg.Dsn, true)
	if err != nil {
		log.Panic("NewAdapter Err", zap.Error(err))
	}

	enforcer, err = casbin.NewEnforcer(m, a)
	if err != nil {
		log.Fatal("InitRabc failed", zap.Error(err))
	}
	

}

// initRabcData init rabc data
func initRabcData(ctx context.Context) error {
	casbinroles := []CasbinRule{
		{"p", "Admin", "/api/v1/hostgroup*", "(GET)|(POST)|(DELETE)|(PUT)", "", "", ""},
		{"p", "Normal", "/api/v1/hostgroup*", "(GET)|(POST)|(DELETE)|(PUT)", "", "", ""},
		{"p", "Guest", "/api/v1/hostgroup*", "(GET)", "", "", ""},
		{"p", "Admin", "/api/v1/task*", "(GET)|(POST)|(DELETE)|(PUT)", "", "", ""},
		{"p", "Normal", "/api/v1/task*", "(GET)|(POST)|(DELETE)|(PUT)", "", "", ""},
		{"p", "Guest", "/api/v1/task*", "(GET)", "", "", ""},
		{"p", "Admin", "/api/v1/host*", "(GET)|(POST)|(DELETE)|(PUT)", "", "", ""},
		{"p", "Normal", "/api/v1/host*", "(GET)|(POST)|(DELETE)|(PUT)", "", "", ""},
		{"p", "Guest", "/api/v1/host*", "(GET)", "", "", ""},
		{"p", "Admin", "/api/v1/user/info", "(GET)|(PUT)", "", "", ""},
		{"p", "Normal", "/api/v1/user/info", "(GET)|(PUT)", "", "", ""},
		{"p", "Guest", "/api/v1/user/info", "(GET)", "", "", ""},
		{"p", "Admin", "/api/v1/user/select", "(GET)", "", "", ""},
		{"p", "Normal", "/api/v1/user/select", "(GET)", "", "", ""},
		{"p", "Guest", "/api/v1/user/select", "(GET)", "", "", ""},
		{"p", "Admin", "/api/v1/user/registry", "(POST)", "", "", ""},
		{"p", "Admin", "/api/v1/user/all", "(GET)", "", "", ""},
		{"p", "Admin", "/api/v1/user/admin", "(PUT)|(DELETE)", "", "", ""},
		{"p", "Admin", "/api/v1/user/alarmstatus", "(GET)", "", "", ""},
		{"p", "Normal", "/api/v1/user/alarmstatus", "(GET)", "", "", ""},
		{"p", "Guest", "/api/v1/user/alarmstatus", "(GET)", "", "", ""},
		{"p", "Admin", "/api/v1/user/operate", "(GET)", "", "", ""},
		{"p", "Admin", "/api/v1/notify", "(GET)|(PUT)", "", "", ""},
		{"p", "Normal", "/api/v1/notify", "(GET)|(PUT)", "", "", ""},
		{"p", "Guest", "/api/v1/notify", "(GET)", "", "", ""},
	}
	err := gormdb.WithContext(ctx).Create(casbinroles).Error
	return err
}

// GetEnforcer get casbin auth
func GetEnforcer() *casbin.Enforcer {
	return enforcer
}
