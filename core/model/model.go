package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/labulaka521/crocodile/common/db"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// InitDb init db
func InitDb() {
	var (
		err error
	)
	dbcfg := config.CoreConf.Server.DB
	err = db.NewDb(db.Drivename(dbcfg.Drivename),
		db.Dsn(dbcfg.Dsn),
		db.MaxIdleConnection(dbcfg.MaxIdle),
		db.MaxOpenConnection(dbcfg.MaxConn),
		db.MaxQueryTime(dbcfg.MaxQueryTime.Duration),
	)
	if err != nil {
		log.Fatal("InitDb failed", zap.Error(err))
	}
}

type checkType uint

const (
	// Email check email
	Email checkType = iota
	// Name check name
	Name
	// ID check id
	ID
	// IDCreateByUID check ID CreateByUID
	IDCreateByUID
	// UID check uid
	UID // uid正常
)

// Tb selcet table name
type Tb string

const (
	// TBUser select ccrocodile_user
	TBUser Tb = "crocodile_user"
	// TBHostgroup select ccrocodile_user
	TBHostgroup Tb = "crocodile_hostgroup"
    // TBTask select crocodile_task
	TBTask Tb = "crocodile_task"
)

// Check check some msg is valid
func Check(ctx context.Context, table Tb, checkType checkType, args ...interface{}) (bool, error) {
	check := fmt.Sprintf("select COUNT(id) FROM %s WHERE ", table)
	switch checkType {
	case Email:
		check += "email=?"
	case Name:
		check += "name=?"
	case ID:
		check += "id=?"
	case IDCreateByUID:
		// 检查ID的createBy字段是否位当前登陆用户
		// 如果当前用户为Admin 则世界返回true
		check += "id=? AND createByID=?"
	case UID:
		// 检查UID状态是否正常
		check += "id=? AND forbid=0"
	default:
		return false, errors.New("reqType Only Support email username")
	}
	conn, err := db.GetConn(ctx)
	if err != nil {
		return false, errors.Wrap(err, "sqlDb.GetConn")
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, check)
	if err != nil {
		return false, errors.Wrap(err, "conn.PrepareContext")
	}
	defer stmt.Close()

	res := 0
	err = stmt.QueryRowContext(ctx, args...).Scan(&res)
	if err != nil && err != sql.ErrNoRows {
		return false, errors.Wrap(err, "stmt.QueryRowContext")
	}
	if err == sql.ErrNoRows || res == 0 {
		return false, nil
	}
	return true, nil
}

// QueryUserRule query user rule by uid
func QueryUserRule(ctx context.Context, uid string) (define.Role, error) {
	conn, err := db.GetConn(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "sqlDb.GetConn")
	}
	defer conn.Close()
	var role define.Role

	rolesql := `SELECT role FROM crocodile_user WHERE id=?`
	stmt, err := conn.PrepareContext(ctx, rolesql)
	if err != nil {
		return 0, errors.Wrap(err, "conn.PrepareContext")
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, uid).Scan(&role)
	if err != nil {
		return 0, errors.Wrap(err, "stmt.QueryRowContext")
	}
	return role, nil
}
