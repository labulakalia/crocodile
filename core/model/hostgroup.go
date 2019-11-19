package model

import (
	"context"
	"github.com/labulaka521/crocodile/common/db"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/common/utils"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"strings"
	"time"
)

// 主机组

func CreateHostgroup(ctx context.Context, hg *define.HostGroup) error {
	createsql := `INSERT INTO crocodile_hostgroup (id,name,remark,createByID,hosts,createTime,updateTime) VALUES(?,?,?,?,?,?,?)`
	conn, err := db.GetConn(ctx)
	if err != nil {
		return errors.Wrap(err, "db.Db.GetConn")
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, createsql)
	if err != nil {
		return errors.Wrap(err, "conn.PrepareContext")
	}

	createTime := time.Now().Unix()
	_, err = stmt.ExecContext(ctx,
		hg.Id,
		hg.Name,
		hg.Remark,
		hg.CreateByUId,
		strings.Join(hg.Hosts, ","),
		createTime,
		createTime)
	if err != nil {
		return errors.Wrap(err, "stmt.ExecContext")
	}
	return nil
}

func ChangeHostGroup(ctx context.Context, hg *define.HostGroup) error {
	changesql := `UPDATE crocodile_hostgroup SET hosts=?,remark=?,updateTime=? WHERE id=?`
	conn, err := db.GetConn(ctx)
	if err != nil {
		return errors.Wrap(err, "db.Db.GetConn")
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, changesql)
	if err != nil {
		return errors.Wrap(err, "conn.PrepareContext")
	}

	_, err = stmt.ExecContext(ctx, strings.Join(hg.Hosts, ","), hg.Remark, time.Now().Unix(), hg.Id)
	if err != nil {
		return errors.Wrap(err, "stmt.ExecContext")
	}
	return nil
}

func DeleteHostGroup(ctx context.Context, id string) error {
	sqldelete := `DELETE FROM crocodile_hostgroup WHERE id=?`
	conn, err := db.GetConn(ctx)
	if err != nil {
		return errors.Wrap(err, "db.Db.GetConn")
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, sqldelete)
	if err != nil {
		return errors.Wrap(err, "conn.PrepareContext")
	}

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return errors.Wrap(err, "stmt.ExecContext")
	}
	return nil
}

func GetHostGroups(ctx context.Context) ([]define.HostGroup, error) {
	hgs := []define.HostGroup{}

	sqlget := `SELECT 
					hg.id,
					hg.name,
					hg.remark,
					hg.hosts,
					hg.createByID,
					u.name,
					hg.createTime,
					hg.updateTime 
				FROM 
					crocodile_hostgroup as hg,crocodile_user as u
				WHERE
					hg.createByID == u.id`
	conn, err := db.GetConn(ctx)
	if err != nil {
		return hgs, errors.Wrap(err, "db.Db.GetConn")
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, sqlget)
	if err != nil {
		return hgs, errors.Wrap(err, "conn.PrepareContext")
	}
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return hgs, errors.Wrap(err, "stmt.QueryContext")
	}

	for rows.Next() {
		var (
			hg                     define.HostGroup
			hosts                  string
			createTime, updateTime int64
		)
		err := rows.Scan(&hg.Id, &hg.Name, &hg.Remark,
			&hosts, &hg.CreateByUId, &hg.CreateBy, &createTime, &updateTime)
		if err != nil {
			log.Info("Scan result failed", zap.Error(err))
		}
		hg.Hosts = []string{}
		if hosts != "" {
			hg.Hosts = strings.Split(hosts, "")
		}
		hg.CreateTime = utils.UnixToStr(createTime)
		hg.UpdateTime = utils.UnixToStr(updateTime)

		hgs = append(hgs, hg)
	}
	return hgs, nil
}
