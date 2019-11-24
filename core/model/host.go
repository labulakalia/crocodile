package model

import (
	"context"
	"fmt"
	"github.com/labulaka521/crocodile/common/db"
	"github.com/labulaka521/crocodile/common/log"
	pb "github.com/labulaka521/crocodile/core/proto"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"strings"
	"time"
)

func RegistryNewHost(ctx context.Context, host *pb.RegistryReq) error {
	hostsql := `INSERT INTO crocodile_host 
					(hostname,
					addr,
					version,
					lastUpdateTime)
 			  	VALUES
					(?,?,?,?)`
	addr := fmt.Sprintf("%s:%d", host.Ip, host.Port)
	hosts, err := getHosts(ctx, addr)
	if err != nil {
		return err
	}
	if len(hosts) == 1 {
		log.Info("Addr Already Registry", zap.String("addr", addr))
		return nil
	}
	conn, err := db.GetConn(ctx)
	if err != nil {
		return errors.Wrap(err, "db.GetConn")
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, hostsql)
	if err != nil {
		return errors.Wrap(err, "conn.PrepareContext")
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		host.Hostname,
		addr,
		host.Version,
		time.Now().Unix())
	if err != nil {
		return errors.Wrap(err, "stmt.ExecContext")
	}
	return nil
}

func UpdateRunningTask(ctx context.Context, host *pb.HeartbeatReq) error {
	updatesql := `UPDATE TABLE crocodile_host set lastUpdateTime=?,runingTasks=? WHERE addr=?`
	conn, err := db.GetConn(ctx)
	if err != nil {
		return errors.Wrap(err, "db.GetConn")
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, updatesql)
	if err != nil {
		return errors.Wrap(err, "conn.PrepareContext")
	}
	defer stmt.Close()
	runtask := strings.Join(host.RunningTask, ",")
	_, err = stmt.ExecContext(ctx,
		time.Now().Unix(),
		runtask,
		fmt.Sprintf("%s:%d", host.Ip, host.Port))
	if err != nil {
		return errors.Wrap(err, "stmt.ExecContext")
	}
	return nil
}

func getHosts(ctx context.Context, addr string) ([]define.Host, error) {
	getsql := "SELECT hostname,addr,runingTasks,version,lastUpdateTime FROM crocodile_host"
	args := []interface{}{}
	if addr != "" {
		getsql += " WHERE addr=?"
		args = append(args, addr)
	}

	conn, err := db.GetConn(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "db.GetConn")
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, getsql)
	if err != nil {
		return nil, errors.Wrap(err, "conn.PrepareContext")
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, errors.Wrap(err, "stmt.QueryContext")
	}

	hosts := []define.Host{}
	for rows.Next() {
		var (
			h     define.Host
			rtask string
		)

		err := rows.Scan(&h.HostName, &h.Addr, &rtask, &h.Version, &h.LastUpdateTime)
		if err != nil {
			log.Error("Scan failed", zap.Error(err))
		}
		hosts = append(hosts, h)
	}

	return hosts, nil
}

func GetHost(ctx context.Context) ([]define.Host, error) {
	return getHosts(ctx, "")
}

func GetHostByID(ctx context.Context, id string) (*define.Host, error) {
	hosts, err := getHosts(ctx, id)
	if err != nil {
		return nil, err
	}
	if len(hosts) != 1 {
		return nil, errors.New("can not find hostid")
	}
	return &hosts[0], nil
}
