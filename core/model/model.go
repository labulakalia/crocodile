package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/labulaka521/crocodile/common/db"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// InitDb init db
func InitDb() error {
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
		return err
	}
	log.Debug("InitDb Success", zap.String("drive", dbcfg.Drivename), zap.String("DSN", dbcfg.Dsn))
	return nil

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
	// NameCreateByUID check name's createByUID
	NameCreateByUID
	// HostGroupID check hostgroup is used by tasks
	HostGroupID
	// CreateByID check use is used by hostgroup or tasks
	CreateByID
	// UserName check exist user name 
	UserName
)

const (
	// TBUser select ccrocodile_user
	TBUser string = "crocodile_user"
	// TBHostgroup select ccrocodile_user
	TBHostgroup string = "crocodile_hostgroup"
	// TBTask select crocodile_task
	TBTask string = "crocodile_task"
	// TBHost select crocodile_host
	TBHost string = "crocodile_host"
	// TBLog log table
	TBLog string = "crocodile_log"
	// TBNotify notify table
	TBNotify string = "crocodile_notify"
	// TBOperate operate table
	TBOperate string = "crocodile_operate"
	// TBCasbin casbin table
	TBCasbin string = "casbin_rule"
)

// Check check some msg is valid
func Check(ctx context.Context, table string, checkType checkType, args ...interface{}) (bool, error) {
	check := fmt.Sprintf("select COUNT(id) FROM %s WHERE ", table)
	switch checkType {
	case Email:
		check += "email=?"
	case Name:
		check += "name=?"
	case ID:
		check += "id=?"
	case IDCreateByUID:
		// 检查ID的createByUID字段是否位当前登陆用户
		// 如果当前用户为Admin 则世界返回true
		check += "id=? AND createByID=?"
	case NameCreateByUID:
		// 检查ID的createByUID字段是否位当前登陆用户
		// 如果当前用户为Admin 则世界返回true
		check += "name=? AND createByID=?"
	case UID:
		// 检查UID状态是否正常
		check += "id=? AND forbid=false"
	case HostGroupID:
		check += "hostGroupID=?"
	case CreateByID:
		check += "createByID=?"
	case UserName:
		// 修改用户名 检查新的用户名不与除自已外其他的用户名重复
		check += "name=? AND id!=?"
	default:
		return false, errors.New("reqType unSupport")
	}
	conn, err := db.GetConn(ctx)
	if err != nil {
		return false, fmt.Errorf("sqlDb.GetConn failed: %w", err)
	}
	defer conn.Close()
	log.Debug("check sql", zap.String("sql", check),zap.Any("args", args))
	stmt, err := conn.PrepareContext(ctx, check)
	if err != nil {
		return false, fmt.Errorf("conn.PrepareContext failed: %w", err)
	}
	defer stmt.Close()

	res := 0
	err = stmt.QueryRowContext(ctx, args...).Scan(&res)
	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("stmt.QueryRowContext failed: %w", err)
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
		return 0, fmt.Errorf("sqlDb.GetConn failed: %w", err)
	}
	defer conn.Close()
	var role define.Role

	rolesql := `SELECT role FROM crocodile_user WHERE id=?`
	stmt, err := conn.PrepareContext(ctx, rolesql)
	if err != nil {
		return 0, fmt.Errorf("conn.PrepareContext failed: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, uid).Scan(&role)
	if err != nil {
		return 0, fmt.Errorf("stmt.QueryRowContext failed: %w", err)
	}
	return role, nil
}

// GetNameID get return name,id
func GetNameID(ctx context.Context, t string) ([]define.KlOption, error) {
	getsql := `SELECT id,name FROM ` + string(t)
	if t == TBHost {
		getsql = `SELECT id,addr,lastUpdateTimeUnix FROM ` + string(t)
	}
	conn, err := db.GetConn(ctx)
	if err != nil {
		return nil, fmt.Errorf("sqlDb.GetConn failed: %w", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, getsql)
	if err != nil {
		return nil, fmt.Errorf("conn.PrepareContext failed: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("stmt.QueryContext failed: %w", err)
	}
	kloptions := []define.KlOption{}
	for rows.Next() {
		var (
			id, name           string
			lastUpdateTimeUnix int64
		)
		kloption := define.KlOption{}
		if t == TBHost {
			err = rows.Scan(&id, &name, &lastUpdateTimeUnix)
			if err != nil {
				log.Error("rows.Scan failed", zap.Error(err))
				continue
			}

			if lastUpdateTimeUnix+maxWorkerTTL > time.Now().Unix() {
				kloption.Online = 1
			} else {
				kloption.Online = -1
			}
		} else {
			err = rows.Scan(&id, &name)
			if err != nil {
				log.Error("rows.Scan failed", zap.Error(err))
				continue
			}
		}
		kloption.Label = name
		kloption.Value = id
		kloptions = append(kloptions, kloption)
	}
	return kloptions, nil
}

func countColums(ctx context.Context, querysql string, args ...interface{}) (int, error) {
	querysql2 := gencountsql(querysql)
	conn, err := db.GetConn(ctx)
	if err != nil {
		return 0, fmt.Errorf("db.GetConn failed: %w", err)
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, querysql2)
	if err != nil {
		return 0, fmt.Errorf("conn.PrepareContext failed: %w", err)
	}
	defer stmt.Close()
	var count int
	err = stmt.QueryRowContext(ctx, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("stmt.QueryRowContext Scan failed: %w", err)
	}
	return count, nil
}

func gencountsql(querysql string) string {
	to := strings.Index(querysql, "FROM")
	from := 6
	return strings.Replace(querysql, querysql[from:to], " count(*) ", -1)

}
