package model

import (
	"context"
	"fmt"
	"strings"
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
	maxWorkerTTL int64 = 20 // defaultHearbeatInterval = 15
)

// RegistryToUpdateHost refistry new host
func RegistryToUpdateHost(ctx context.Context, req *pb.RegistryReq) error {
	updatesql := `UPDATE crocodile_host set weight=?,version=?,lastUpdateTimeUnix=?,remark=? WHERE addr=?`
	conn, err := db.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("db.GetConn failed: %w", err)
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, updatesql)
	if err != nil {
		return fmt.Errorf("conn.PrepareContext failed: %w", err)
	}
	defer stmt.Close()
	addr := fmt.Sprintf("%s:%d", req.Ip, req.Port)
	_, err = stmt.ExecContext(ctx, req.Weight, req.Version, time.Now().Unix(), req.Remark, addr)
	if err != nil {
		return fmt.Errorf("stmt.ExecContext faled: %w", err)
	}
	return nil
}

// RegistryNewHost refistry new host
func RegistryNewHost(ctx context.Context, req *pb.RegistryReq) (string, error) {
	hostsql := `INSERT INTO crocodile_host 
					(id,
					hostname,
					addr,
					weight,
					version,
					lastUpdateTimeUnix,
					remark
				)
 			  	VALUES
					(?,?,?,?,?,?,?)`
	addr := fmt.Sprintf("%s:%d", req.Ip, req.Port)
	hosts, _, err := getHosts(ctx, addr, nil, 0, 0)
	if err != nil {
		return "", err
	}
	if len(hosts) == 1 {
		log.Info("Addr Already Registry", zap.String("addr", addr))
		return "", nil
	}
	conn, err := db.GetConn(ctx)
	if err != nil {
		return "", fmt.Errorf("db.GetConn failed: %w", err)
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, hostsql)
	if err != nil {
		return "", fmt.Errorf("conn.PrepareContext failed: %w", err)
	}
	defer stmt.Close()
	id := utils.GetID()
	_, err = stmt.ExecContext(ctx,
		id,
		req.Hostname,
		addr,
		req.Weight,
		req.Version,
		time.Now().Unix(),
		req.Remark,
	)
	if err != nil {
		return "", fmt.Errorf("stmt.ExecContext failed: %w", err)
	}
	log.Info("New Client Registry ", zap.String("addr", addr))
	return id, nil
}

// UpdateHostHearbeat update host last recv hearbeat time
func UpdateHostHearbeat(ctx context.Context, ip string, port int32, runningtasks []string) error {
	updatesql := `UPDATE crocodile_host set lastUpdateTimeUnix=?,runningTasks=? WHERE addr=?`
	conn, err := db.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("db.GetConn failed: %w", err)
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, updatesql)
	if err != nil {
		return fmt.Errorf("conn.PrepareContext failed: %w", err)
	}
	defer stmt.Close()
	result, err := stmt.ExecContext(ctx,
		time.Now().Unix(),
		strings.Join(runningtasks, ","),
		fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return fmt.Errorf("stmt.ExecContext failed: %w", err)
	}
	line, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("stmt.ExecContext failed: %w", err)
	}
	if line <= 0 {
		return fmt.Errorf("host %s not registry, may be this host is delete", fmt.Sprintf("%s:%d", ip, port))
	}
	return nil
}

// get host by addr or id
func getHosts(ctx context.Context, addr string, ids []string, offset, limit int) ([]*define.Host, int, error) {
	getsql := `SELECT 
					id,
					addr,
					hostname,
					runningTasks,
					weight,
					stop,
					version,
					lastUpdateTimeUnix,
					remark
			   FROM 
					crocodile_host`
	var (
		count int
	)
	args := []interface{}{}
	// only use addr or ids query
	if addr != "" && len(ids) != 0 {
		return nil, 0, errors.New("only use addr or ids query")
	}
	if addr != "" {
		getsql += " WHERE addr=?"
		args = append(args, addr)
	}

	if len(ids) > 0 {
		var querys = []string{}
		for _, id := range ids {
			querys = append(querys, "id=?")
			args = append(args, id)
		}
		getsql += " WHERE " + strings.Join(querys, " OR ")

	}
	if limit > 0 {
		var err error
		count, err = countColums(ctx, getsql, args...)
		if err != nil {
			return nil, 0, fmt.Errorf("countColums failed: %w", err)
		}
		getsql += " LIMIT ? OFFSET ?"
		args = append(args, limit, offset)
	}

	conn, err := db.GetConn(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("db.GetConn failed: %w", err)
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, getsql)
	if err != nil {
		return nil, 0, fmt.Errorf("conn.PrepareContext failed: %w", err)
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("stmt.QueryContext failed: %w", err)
	}

	hosts := []*define.Host{}
	for rows.Next() {
		var (
			h     define.Host
			rtask string
		)
		err := rows.Scan(&h.ID,
			&h.Addr,
			&h.HostName,
			&rtask,
			&h.Weight,
			&h.Stop,
			&h.Version,
			&h.LastUpdateTimeUnix,
			&h.Remark)
		if err != nil {
			log.Error("Scan failed", zap.Error(err))
			continue
		}
		h.RunningTasks = []string{}
		if rtask != "" {
			h.RunningTasks = append(h.RunningTasks, strings.Split(rtask, ",")...)
		}
		if h.LastUpdateTimeUnix+maxWorkerTTL > time.Now().Unix() {
			h.Online = true
		}
		h.LastUpdateTime = utils.UnixToStr(h.LastUpdateTimeUnix)
		hosts = append(hosts, &h)
	}
	return hosts, count, nil
}

// GetHosts get all hosts
func GetHosts(ctx context.Context, offset, limit int) ([]*define.Host, int, error) {
	return getHosts(ctx, "", nil, offset, limit)
}

// GetHostByAddr get host by addr
func GetHostByAddr(ctx context.Context, addr string) (*define.Host, error) {
	hosts, _, err := getHosts(ctx, addr, nil, 0, 0)
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
	hosts, _, err := getHosts(ctx, addr, nil, 0, 0)
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
	hosts, _, err := getHosts(ctx, "", []string{id}, 0, 0)
	if err != nil {
		return nil, err
	}
	if len(hosts) != 1 {
		log.Warn("can not find hostid", zap.Error(err))
		err = define.ErrNotExist{Value: id}
		return nil, err
	}
	return hosts[0], nil
}

// GetHostByIDS get hosts by hostids
func GetHostByIDS(ctx context.Context, ids []string) ([]*define.Host, error) {
	hosts, _, err := getHosts(ctx, "", ids, 0, 0)
	if err != nil {
		return nil, err
	}
	return hosts, nil
}

// StopHost will stop run worker in hostid
func StopHost(ctx context.Context, hostid string, stop bool) error {
	stopsql := `UPDATE crocodile_host SET stop=? WHERE id=?`
	conn, err := db.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("db.GetConn failed: %w", err)
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, stopsql)
	if err != nil {
		return fmt.Errorf("conn.PrepareContext failed: %w", err)
	}
	_, err = stmt.ExecContext(ctx, stop, hostid)
	if err != nil {
		return fmt.Errorf("stmt.ExecContext failed: %w", err)
	}
	return nil
}

// DeleteHost will delete host
func DeleteHost(ctx context.Context, hostid string) error {
	err := StopHost(ctx, hostid, true)
	if err != nil {
		return fmt.Errorf("StopHost failed: %w", err)
	}
	deletehostsql := `DELETE FROM crocodile_host WHERE id=?`
	conn, err := db.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("db.GetConn failed: %w", err)
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, deletehostsql)
	if err != nil {
		return fmt.Errorf("conn.PrepareContext failed: %w", err)
	}
	_, err = stmt.ExecContext(ctx, hostid)
	if err != nil {
		return fmt.Errorf("stmt.ExecContext failed: %w", err)
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
		return ids, false
	}
	ids = append(ids[:existid], ids[existid+1:]...)
	return ids, true
}
