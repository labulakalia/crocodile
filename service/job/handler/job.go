package handler

import (
	"context"
	"crocodile/common/e"
	"crocodile/service/job/model/task"
	pbjob "crocodile/service/job/proto/job"
	"github.com/labulaka521/logging"
)

type Job struct {
	Service task.Servicer
}

func (job *Job) CreateJob(ctx context.Context, task *pbjob.Task, resp *pbjob.Response) (err error) {
	var (
		code int32
	)
	logging.Debugf("CreateJob %s", task.Taskname)
	code = e.SUCCESS

	if err = job.Service.CreateJob(ctx, task); err != nil {
		code = e.ERR_CREATE_JOB_FAIL
		return
	}

	resp.Code = code
	return
}

func (job *Job) DeleteJob(ctx context.Context, task *pbjob.Task, resp *pbjob.Response) (err error) {
	var (
		code int32
	)
	logging.Debugf("DeleteJob %s", task.Taskname)
	code = e.SUCCESS

	if err = job.Service.DeleteJob(ctx, task.Taskname); err != nil {
		code = e.ERR_DELETE_JOB_FAIL
		return
	}
	resp.Code = code
	return
}

func (job *Job) ChangeJob(ctx context.Context, task *pbjob.Task, resp *pbjob.Response) (err error) {
	var (
		code int32
	)
	logging.Debugf("ChangeJob %s", task.Taskname)
	code = e.SUCCESS

	if err = job.Service.ChangeJob(ctx, task); err != nil {
		code = e.ERR_CHANGE_JOB_FAIL
		return
	}
	resp.Code = code
	return
}

// taskname or
func (job *Job) GetJob(ctx context.Context, task *pbjob.Task, resp *pbjob.Response) (err error) {
	var (
		code  int32
		tasks []*pbjob.Task
	)
	logging.Debugf("GetJob %s", task.Taskname)
	code = e.SUCCESS

	if tasks, err = job.Service.GetJob(ctx, task.Taskname); err != nil {
		code = e.ERR_GET_JOB_FAIL
		return
	}
	resp.Code = code
	resp.Tasks = tasks
	return
}

// taskname
func (job *Job) KillJob(ctx context.Context, task *pbjob.Task, resp *pbjob.Response) (err error) {
	var (
		code int32
	)
	logging.Debugf("KillJob %s", task.Taskname)
	code = e.SUCCESS

	if err = job.Service.KillJob(ctx, task); err != nil {
		code = e.ERR_KILL_JOB_FAIL
		return
	}
	resp.Code = code
	return
}

// taskname
func (job *Job) RunJob(ctx context.Context, task *pbjob.Task, resp *pbjob.Response) (err error) {
	var (
		code    int32
		getTask []*pbjob.Task
	)
	logging.Debugf("RunJob %s", task.Taskname)
	if getTask, err = job.Service.GetJob(ctx, task.Taskname); err != nil || len(getTask) != 1 {
		code = e.ERR_GET_JOB_FAIL
		return
	}

	code = e.SUCCESS
	if err = job.Service.RunJob(ctx, getTask[0]); err != nil {
		code = e.ERR_RUN_JOB_FAIL
		return
	}
	resp.Code = code
	return
}
