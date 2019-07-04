package scheduler

import (
	"context"
	"crocodile/common/util"
	"github.com/golang/protobuf/ptypes"
	"math"

	"crocodile/service/job/model/task"

	pbjob "crocodile/service/job/proto/job"

	"database/sql"
	"github.com/labulaka521/logging"
	"time"
)

// 调度器
// 检测任务运行的时机
// 获取全部的任务
// 比较任务的下一次执行时间

func Loop(exit chan int, db *sql.DB) {
	var (
		err             error
		taskService     task.Servicer
		t               *pbjob.Task
		tasks           []*pbjob.Task
		nextTime        time.Time
		now             time.Time
		nearTime        time.Duration
		tk              *time.Ticker
		defaultWaitTime time.Duration
	)

	defaultWaitTime = 10 * time.Second
	logging.Debug("Start Run Scheduler Loop...")

	taskService = &task.Service{
		DB: db,
	}
	for {
		now = time.Now()
		if tasks, err = taskService.GetJob(context.TODO(), ""); err != nil {
			logging.Errorf("Get JOb Err: %v", err)

		}

		// 获取全部的任务
		if len(tasks) == 0 {
			nearTime = defaultWaitTime
		}

		for _, t = range tasks {
			if nextTime, err = ptypes.Timestamp(t.Nexttime); err != nil {
				logging.Errorf("task %s Time format Err: %v", t.Taskname, err)
				continue
			}

			// 比较任务的下一次运行时间
			if nextTime.Before(now) || nextTime.Equal(now) {
				go func(t *pbjob.Task) {
					if t.Stop {
						logging.Warnf("Task %s is stop scheduler ", t.Taskname)
						return
					}
					if err = taskService.RunJob(context.TODO(), t); err != nil {
						logging.Errorf("Run Job %s Err: %v", t.Taskname, err)
					}
				}(t)
			}
			// 更新下一次的运行时间
			nextTime, _ = util.NextTime(t.Cronexpr, now)
			if err = taskService.UpdateNextTime(context.TODO(), t.Taskname, nextTime); err != nil {
				logging.Errorf("UpdatenexTime Task %s Err:%v", t.Taskname, err)
			}

			if nextTime.Sub(now) < nearTime || nearTime == 0 {
				nearTime = nextTime.Sub(now)
			}
		}

		tk = time.NewTicker(nearTime)
		select {
		case <-tk.C:
			nearTime = math.MaxInt64
		case <-exit:
			tk.Stop()
			logging.Info("Exiting Scheduler Loop...")
			return
		}
	}
}
