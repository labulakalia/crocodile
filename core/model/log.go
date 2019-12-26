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

// task log
func SaveLog(ctx context.Context, tasklog *define.Log) error {
	log.Info("start savelog", zap.Any("tasklog", tasklog))
	savesql := `INSERT INTO crocodile_log
    (taskid,starttime,endtime,totalruntime,taskresps)
VALUES
    (?,?,?,?,?)`
	conn, err := db.GetConn(ctx)
	if err != nil {
		return errors.Wrap(err, "db.GetConn")
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, savesql)
	if err != nil {
		return errors.Wrap(err, "conn.PrepareContext")
	}
	taskresps, err := json.Marshal(tasklog.TaskResps)
	if err != nil {
		return errors.Wrap(err, "json.Marshal")
	}
	_, err = stmt.ExecContext(ctx, tasklog.RunByTaskId,
		tasklog.StartTimne, tasklog.EndTime, tasklog.TotalRunTime, taskresps)
	if err != nil {
		return errors.Wrap(err, "stmt.ExecContext")
	}
	return nil
}
