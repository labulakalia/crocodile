package schedule

import (
	"context"
	"encoding/json"
	"github.com/labulaka521/crocodile/core/utils/define"

	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/model"
	"go.uber.org/zap"
)

// TaskEvent task event
type TaskEvent uint8

const (
	// AddEvent recv add event
	AddEvent TaskEvent = iota + 1
	// ChangeEvent recv delete task
	ChangeEvent
	// DeleteEvent recv delete task
	DeleteEvent
	// RunEvent run a task
	RunEvent
	// KillEvent recv stop task
	KillEvent
)

const (
	pubsubChannel = "task.event"
)

// EventData sub data from redis
// 应用于调度节点集群，当添加任务、删除修改任务、终止任务时，
// 所有的集群调度节点都会接收到信息，然后进行相应的修改操作
type EventData struct {
	TaskID string    // task id
	TE     TaskEvent // task event: add change delete stop task
}

// RecvEvent recv task event
func RecvEvent() {
	sub := Cron2.redis.Subscribe(pubsubChannel)
	for msg := range sub.Channel() {
		go dealEvent([]byte(msg.Payload))
	}
}

func dealEvent(data []byte) {
	var subdata EventData
	err := json.Unmarshal(data, &subdata)
	if err != nil {
		log.Error("json.Unmarshal event data failed", zap.Error(err))
		return
	}
	switch subdata.TE {
	case AddEvent:
		fallthrough
	case ChangeEvent:
		ctx, cancel := context.WithTimeout(context.Background(), config.CoreConf.Server.DB.MaxQueryTime.Duration)
		defer cancel()
		task, err := model.GetTaskByID(ctx, subdata.TaskID)
		if err != nil {
			log.Error("model.GetTaskByID failed", zap.Error(err))
			return
		}
		Cron2.addtask(task.ID, task.Name, task.Cronexpr, GetRoutePolicy(task.HostGroupID, task.RoutePolicy),task.Run)
	case DeleteEvent:
		Cron2.deletetask(subdata.TaskID)
	case RunEvent:
		task,ok := Cron2.GetTask(subdata.TaskID)
		if !ok {
			log.Error("Can not get Task",zap.String("taskid",subdata.TaskID))
			return
		}
		go task.StartRun(define.Manual)

	case KillEvent:
		Cron2.killtask(subdata.TaskID)
	default:
		log.Warn("unsupport task event", zap.Any("data", subdata))
	}
}
