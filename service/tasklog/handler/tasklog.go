package handler

import (
	"context"
	"crocodile/service/tasklog/model/tasklog"
	pbtasklog "crocodile/service/tasklog/proto/tasklog"
	"github.com/labulaka521/logging"
)

type TaskLog struct {
	Service tasklog.Servicer
}

func (tLog *TaskLog) CreateLog(ctx context.Context, log *pbtasklog.SimpleLog, resp *pbtasklog.Empty) (err error) {
	logging.Debugf("Create Task %s Log", log.Taskname)
	if err = tLog.Service.CreateLog(ctx, log); err != nil {
		logging.Errorf("Create Log Err:%v", err)
	}
	return
}

// 获取日志
func (tLog *TaskLog) GetLog(ctx context.Context, querylog *pbtasklog.QueryLog, resp *pbtasklog.RespLog) (err error) {
	var (
		rsp *pbtasklog.RespLog
	)
	if rsp, err = tLog.Service.GetLog(ctx, querylog); err != nil {
		logging.Errorf("Get Log Err:%v", err)
	}
	resp.Logs = rsp.Logs
	resp.Count = rsp.Count
	return
}
