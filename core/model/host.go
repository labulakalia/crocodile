package model

import (
	"context"
	"fmt"
	"github.com/labulaka521/crocodile/common/db"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/common/utils"
	pb "github.com/labulaka521/crocodile/core/proto"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

const (
	maxWorkerTTL int64 = 60
)

func RegistryNewHost(ctx context.Context, req *pb.RegistryReq) (string, error) {
	hostsql := `INSERT INTO crocodile_host 
					(id,hostname,
					addr,
					version,
					lastUpdateTimeUnix)
 			  	VALUES
					(?,?,?,?,?)`
	addr := fmt.Sprintf("%s:%d", req.Ip, req.Port)
	hosts, err := getHosts(ctx, addr, "")
	if err != nil {
		return "", err
	}
	if len(hosts) == 1 {
		log.Info("Addr Already Registry", zap.String("addr", addr))
		return "", nil
	}
	conn, err := db.GetConn(ctx)
	if err != nil {
		return "", errors.Wrap(err, "db.GetConn")
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, hostsql)
	if err != nil {
		return "", errors.Wrap(err, "conn.PrepareContext")
	}
	defer stmt.Close()
	id := utils.GetId()
	_, err = stmt.ExecContext(ctx,
		id,
		req.Hostname,
		addr,
		req.Version,
		time.Now().Unix())
	if err != nil {
		return "", errors.Wrap(err, "stmt.ExecContext")
	}
	log.Info("New Client Registry ", zap.String("addr", addr))
	return id, nil
}

func UpdateRunningTask(ctx context.Context, hbreq *pb.HeartbeatReq) error {
	updatesql := `UPDATE crocodile_host set lastUpdateTimeUnix=? WHERE addr=?`
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
	_, err = stmt.ExecContext(ctx,
		time.Now().Unix(),
		fmt.Sprintf("%s:%d", hbreq.Ip, hbreq.Port))
	if err != nil {
		return errors.Wrap(err, "stmt.ExecContext")
	}
	return nil
}

func getHosts(ctx context.Context, addr, id string) ([]*define.Host, error) {
	getsql := "SELECT id,addr,hostname,runingTasks,version,lastUpdateTimeUnix FROM crocodile_host"
	args := []interface{}{}
	if addr != "" {
		getsql += " WHERE addr=?"
		args = append(args, addr)
	}
	if id != "" {
		getsql += " WHERE id=?"
		args = append(args, id)
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

	hosts := []*define.Host{}
	for rows.Next() {
		var (
			h     define.Host
			rtask string
		)

		err := rows.Scan(&h.Id, &h.Addr, &h.HostName, &rtask, &h.Version, &h.LastUpdateTimeUnix)
		if err != nil {
			log.Error("Scan failed", zap.Error(err))
		}
		if h.LastUpdateTimeUnix+maxWorkerTTL > time.Now().Unix() {
			h.Online = 1
		}
		h.LastUpdateTime = utils.UnixToStr(h.LastUpdateTimeUnix)
		hosts = append(hosts, &h)
	}
	return hosts, nil
}

func GetHost(ctx context.Context) ([]*define.Host, error) {
	return getHosts(ctx, "", "")
}

func GetHostByAddr(ctx context.Context, addr string) (*define.Host, error) {
	hosts, err := getHosts(ctx, addr, "")
	if err != nil {
		return nil, err
	}
	if len(hosts) != 1 {
		return nil, errors.New("can not find hostid")
	}
	return hosts[0], nil
}

func ExistAddr(ctx context.Context, addr string) (*define.Host, bool, error) {
	hosts, err := getHosts(ctx, addr, "")
	if err != nil {
		return nil, false, err
	}
	if len(hosts) != 1 {
		return nil, false, nil
	}
	return hosts[0], true, nil
}

func GetHostById(ctx context.Context, id string) (*define.Host, error) {
	hosts, err := getHosts(ctx, "", id)
	if err != nil {
		return nil, err
	}
	if len(hosts) != 1 {
		return nil, errors.New("can not find hostid")
	}
	if hosts[0].Online == 0 {
		return nil, errors.New(fmt.Sprintf("host %s is not online", hosts[0].Addr))
	}
	return hosts[0], nil
}
