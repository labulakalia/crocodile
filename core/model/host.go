package model

import (
	"context"
	"github.com/labulaka521/crocodile/common/db"
	pb "github.com/labulaka521/crocodile/core/proto"
	"github.com/pkg/errors"
	"strings"
	"time"
)

func UpdateHost(ctx context.Context, host *pb.HeartbeatReq) error {
	hostsql := `INSERT INTO crocodile_host 
					(hostname,
					ip,
					port,
					version,
					runingTasks,
					lastUpdateTime)
 			  	VALUES
					(?,?,?,?,?,?)`
	conn, err := db.GetConn(ctx)
	if err != nil {
		return errors.Wrap(err, "db.GetConn")
	}
	stmt, err := conn.PrepareContext(ctx, hostsql)
	if err != nil {
		return errors.Wrap(err, "conn.PrepareContext")
	}
	defer stmt.Close()
	runtask := strings.Join(host.RunningTask, ",")
	_, err = stmt.ExecContext(ctx,
		host.Hostname,
		host.Ip,
		host.Port,
		host.Version,
		runtask,
		time.Now().Unix())
	if err != nil {
		return errors.Wrap(err, "stmt.ExecContext")
	}
	return nil
}
