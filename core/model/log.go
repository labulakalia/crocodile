package model

import (
	"context"
	"encoding/json"
	"github.com/labulaka521/crocodile/common/db"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// SaveLog  save task reps log
func SaveLog(ctx context.Context, l *define.Log) error {
	log.Info("start savelog", zap.Any("tasklog", l))
	savesql := `INSERT INTO crocodile_log
				(taskid,
				starttime,
				endtime,
				totalruntime,
				status,
				taskresps,
				errcode,
				errmsg,
				errtasktype,
				errtaskid)
			VALUES
			(?,?,?,?,?,?,?,?,?,?)`
	conn, err := db.GetConn(ctx)
	if err != nil {
		return errors.Wrap(err, "db.GetConn")
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, savesql)
	if err != nil {
		return errors.Wrap(err, "conn.PrepareContext")
	}
	taskresps, err := json.Marshal(l.TaskResps)
	if err != nil {
		return errors.Wrap(err, "json.Marshal")
	}
	_, err = stmt.ExecContext(ctx, l.RunByTaskID,
		l.StartTime, l.EndTime, l.TotalRunTime,
		l.Status, taskresps, l.ErrCode, l.ErrMsg,
		l.ErrTasktype, l.ErrTaskID)
	if err != nil {
		return errors.Wrap(err, "stmt.ExecContext")
	}
	return nil
}

// GetLog get task resp log by taskid
func GetLog(ctx context.Context, taskid string) ([]*define.Log, error) {
	logs := []*define.Log{}
	getsql := 	`SELECT 
					starttime,
					endtime,
					taskresps 
				FROM 
					crocodile_log
			   	WHERE 
			    	taskid=?`
	conn, err := db.GetConn(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "db.GetConn")
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, getsql)
	if err != nil {
		return nil, errors.Wrap(err, "conn.PrepareContext")
	}

	rows, err := stmt.QueryContext(ctx, taskid)
	if err != nil {
		return nil, errors.Wrap(err, "stmt.QueryContext")
	}
	for rows.Next() {
		getlog := define.Log{}
		taskrepos := []*define.TaskResp{}
		var taskreposbyte []byte
		err = rows.Scan(&getlog.StartTime, &getlog.EndTime, &taskreposbyte)
		if err != nil {
			log.Error("rows.Scan failed", zap.Error(err))
			continue
		}
		err = json.Unmarshal(taskreposbyte, &taskrepos)
		if err != nil {
			log.Error("json.Unmarshal failed", zap.Error(err))
			continue
		}
		getlog.TaskResps = taskrepos
		getlog.TotalRunTime = int(getlog.EndTime - getlog.StartTime)
		getlog.RunByTaskID = taskid
		logs = append(logs, &getlog)
	}
	return logs, nil
}
