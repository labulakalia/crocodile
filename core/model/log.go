package model

import (
	"context"
	"encoding/json"

	"github.com/labulaka521/crocodile/common/db"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/common/utils"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// SaveLog  save task reps log
func SaveLog(ctx context.Context, l *define.Log) error {
	log.Info("start savelog", zap.Any("tasklog", l))
	savesql := `INSERT INTO crocodile_log
				(name,
				taskid,
				starttime,
				endtime,
				totalruntime,
				status,
				taskresps,
				errcode,
				errmsg,
				errtasktype,
				errtaskid,
				errtask
			)
			VALUES
			(?,?,?,?,?,?,?,?,?,?,?,?)`
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
	_, err = stmt.ExecContext(ctx, l.Name, l.RunByTaskID,
		l.StartTime, l.EndTime, l.TotalRunTime,
		l.Status, taskresps, l.ErrCode, l.ErrMsg,
		l.ErrTasktype, l.ErrTaskID, l.ErrTask)
	if err != nil {
		return errors.Wrap(err, "stmt.ExecContext")
	}
	return nil
}

// GetLog get task resp log by taskid
func GetLog(ctx context.Context, taskid string, offset, limit int) ([]*define.Log, error) {
	logs := []*define.Log{}
	getsql := `SELECT 
					name,
					taskid,
					starttime,
					endtime,
					totalruntime,
					status,
					taskresps,
					errcode,
					errmsg,
					errtasktype,
					errtaskid,
					errtask
				FROM 
					crocodile_log
			   	WHERE 
					taskid=?
				ORDER BY id DESC
				LIMIT ? OFFSET ?`

	args := []interface{}{taskid, limit, offset}

	conn, err := db.GetConn(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "db.GetConn")
	}
	defer conn.Close()
	// fmt.Println(getsql, args)
	stmt, err := conn.PrepareContext(ctx, getsql)
	if err != nil {
		return nil, errors.Wrap(err, "conn.PrepareContext")
	}

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, errors.Wrap(err, "stmt.QueryContext")
	}
	for rows.Next() {
		getlog := define.Log{}
		taskrepos := []*define.TaskResp{}
		var taskreposbyte []byte
		err = rows.Scan(
			&getlog.Name,
			&getlog.RunByTaskID,
			&getlog.StartTime,
			&getlog.EndTime,
			&getlog.TotalRunTime,
			&getlog.Status,
			&taskreposbyte,
			&getlog.ErrCode,
			&getlog.ErrMsg,
			&getlog.ErrTasktype,
			&getlog.ErrTaskID,
			&getlog.ErrTask,
		)
		if err != nil {
			log.Error("rows.Scan failed", zap.Error(err))
			continue
		}
		err = json.Unmarshal(taskreposbyte, &taskrepos)
		if err != nil {
			log.Error("json.Unmarshal failed", zap.Error(err))
			continue
		}
		getlog.ErrTaskTypeStr = getlog.ErrTasktype.String()
		getlog.TaskResps = taskrepos
		getlog.StartTimeStr = utils.UnixToStr(getlog.StartTime / 1e3)
		getlog.EndTimeStr = utils.UnixToStr(getlog.EndTime / 1e3)
		logs = append(logs, &getlog)
	}
	return logs, nil
}
