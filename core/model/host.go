package model

import (
	"context"
	"fmt"
	"time"

	"github.com/labulaka521/crocodile/common/db"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/common/utils"
	pb "github.com/labulaka521/crocodile/core/proto"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	maxWorkerTTL int64 = 60
)

// RegistryNewHost refistry new host
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
	id := utils.GetID()
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

// UpdateHostHearbeat update host last recv hearbeat time
func UpdateHostHearbeat(ctx context.Context, hbreq *pb.HeartbeatReq) error {
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

// get host by addr or id
func getHosts(ctx context.Context, addr, id string) ([]*define.Host, error) {
	getsql := "SELECT id,addr,hostname,runingTasks,stop,version,lastUpdateTimeUnix FROM crocodile_host"
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

		err := rows.Scan(&h.ID, &h.Addr, &h.HostName, &rtask, &h.Stop, &h.Version, &h.LastUpdateTimeUnix)
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

// GetHost get all hosts
func GetHost(ctx context.Context) ([]*define.Host, error) {
	return getHosts(ctx, "", "")
}

// GetHostByAddr get host by addr
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

// ExistAddr check already exist
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

// GetHostByID get host by hostid
func GetHostByID(ctx context.Context, id string) (*define.Host, error) {
	hosts, err := getHosts(ctx, "", id)
	if err != nil {
		return nil, err
	}
	if len(hosts) != 1 {
		return nil, errors.New("can not find hostid")
	}
	host := hosts[0]
	if host.Online == 0 {
		return nil, fmt.Errorf("host %s is not online", host.Addr)
	}
	if host.Stop == 0 {
		return nil, fmt.Errorf("host %s is stop", host.Addr)
	}

	return hosts[0], nil
}

// StopHost will stop run worker in hostid
func StopHost(ctx context.Context, hostid string, stop int) error {
	stopsql := `UPDATE crocodile_host SET stop=?`
	conn, err := db.GetConn(ctx)
	if err != nil {
		return errors.Wrap(err, "db.GetConn")
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, stopsql)
	if err != nil {
		return errors.Wrap(err, "conn.PrepareContext")
	}
	_, err = stmt.ExecContext(ctx, stop)
	if err != nil {
		return errors.Wrap(err, "stmt.ExecContext")
	}
	return nil
}

// DeleteHost will delete host
func DeleteHost(ctx context.Context, hostid string) error {
	err := StopHost(ctx, hostid, 0)
	if err != nil {
		return errors.Wrap(err, "StopHost")
	}
	go deleteHostFromHostGroup(hostid)
	return nil
}

func deleteHostFromHostGroup(hostid string) error {
	hostgroups, err := GetHostGroups(context.Background())
	if err != nil {
		return errors.Wrap(err, "GetHostGroups")
	}
	for _,hostgroup := range hostgroups {
		newhostid,ok := deletefromslice(hostid,hostgroup.HostsID)
		if !ok {
			continue
		}
		hostgroup.HostsID = newhostid
		err = ChangeHostGroup(context.Background(), &hostgroup)
		if err != nil {
			log.Error("CHangeHostGroup failed", zap.String("hostgroupid", hostgroup.ID))
		}
	}
	return nil
}

// delete from slice
func deletefromslice(deleteid string, ids []string) ([]string, bool) {
	var existid = -1
	for index, id := range ids {
		if id == deleteid {
			existid = index
			break
		}
	}
	if existid == -1 {
		// no found delete id
		return nil,false
	}
	ids = append(ids[:existid-1],ids[existid:]...)
	return ids, true
}

