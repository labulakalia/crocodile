package model

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/labulaka521/crocodile/common/db"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/common/utils"
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/utils/asset"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)


var crcocodileTables = []string{
	TBHost,
	TBHostgroup,
	TBLog,
	TBNotify,
	TBOperate,
	TBTask,
	TBUser,
	TBCasbin,
}

// QueryIsInstall check table is create
func QueryIsInstall(ctx context.Context) (bool, error) {
	var querytable string
	needtables := []interface{}{}

	for _, tbname := range crcocodileTables {
		needtables = append(needtables, tbname)
	}
	var queryname string

	drivename := config.CoreConf.Server.DB.Drivename
	if drivename == "sqlite3" {
		querytable = `SELECT count() FROM sqlite_master WHERE type="table" AND (`
		queryname = "name"
	} else if drivename == "mysql" {
		querytable = `SELECT count() FROM information_schema.TABLES WHERE (`
		queryname = "table_name"
	} else {
		return false, fmt.Errorf("unsupport drive type %s, only support sqlite3 or mysql", drivename)
	}
	params := []string{}

	for i := 0; i < len(crcocodileTables); i++ {
		params = append(params, queryname+"=?")
	}
	querytable += strings.Join(params, " OR ")
	querytable += ")"
	var count int
	conn, err := db.GetConn(ctx)
	if err != nil {
		return false, errors.Wrap(err, "db.GetConn")
	}
	defer conn.Close()
	err = conn.QueryRowContext(ctx, querytable, needtables...).Scan(&count)
	if err != nil {
		log.Error("msg string", zap.Error(err))
		return false, errors.Wrap(err, "Scan")
	}

	if count != len(crcocodileTables) {
		return false, nil
	}
	return true, nil
}


// StartInstall start install system
func StartInstall(ctx context.Context, username, password string) error {
	// create table
	conn, err := db.GetConn(ctx)
	if err != nil {
		return errors.Wrap(err, "db.GetConn")
	}

	fs := &assetfs.AssetFS{
		Asset:     asset.Asset,
		AssetDir:  asset.AssetDir,
		AssetInfo: asset.AssetInfo,
	}

	defer conn.Close()
	for _, tbname := range crcocodileTables {
		// crocodile_host
		var name string
		if tbname != TBCasbin {
			name = tbname[10:]
		} else {
			name = tbname
		}
		sqlfilename := "sql/" + name + ".sql"
		file, err := fs.Open(sqlfilename)
		if err != nil {
			log.Error("fs.Open failed", zap.String("filename", sqlfilename), zap.Error(err))
			continue
		}

		content, err := ioutil.ReadAll(file)
		if err != nil {
			log.Error("ioutil.ReadAll failed", zap.Error(err))
			continue
		}
		_, err = conn.ExecContext(ctx, string(content))
		if err != nil {
			log.Error("conn.ExecContext failed", zap.Error(err), zap.String("tbname", tbname))
			return errors.Wrap(err, "conn.ExecContext")
		}
		// wait second
		time.Sleep(time.Second / 2)
	}
	log.Debug("Success Run All crocodile Sql")

	// create admin user
	hashpassword, err := utils.GenerateHashPass(password)
	if err != nil {
		return errors.Wrap(err, "utils.GenerateHashPass")
	}
	err = AddUser(ctx, username, hashpassword, define.AdminUser)
	if err != nil {
		return errors.Wrap(err, "AddUser")
	}
	err = enforcer.LoadPolicy()
	if err != nil {
		log.Error("enforcer.LoadPolicy failed", zap.Error(err))
		return errors.Wrap(err, "enforcer.LoadPolicy")
	}

	log.Debug("Success Install Crocodile")
	return nil
}
