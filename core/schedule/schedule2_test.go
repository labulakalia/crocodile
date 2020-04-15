package schedule

import (
	"testing"

	"github.com/go-redis/redis"
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/labulaka521/crocodile/core/utils/log"
)

func Test_task2_GetCronExpr(t *testing.T) {
	config.Init("/Users/labulakalia/workerspace/golang/crocodile/core.toml")
	log.Init()
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	task := task2{redis: client}
	err := task.SetData(define.ChildTask, "111", define.TsWait, taskstatus)
	if err != nil {
		t.Log(err)
	}
	// v := map[string]interface{}{
	// 	"cronexpr1": 2,
	// }
	// err := client.HMSet("testtask", v).Err()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// res, err := client.HGet("testtask", "cronexpr1").Int()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(res)

}

// func Test_task2_SetData(t *testing.T) {
// 	type fields struct {
// 		id          string
// 		name        string
// 		cronexpr    string
// 		cronsub     time.Duration
// 		close       chan struct{}
// 		ctxcancel   context.CancelFunc
// 		next        Next
// 		canrun      bool
// 		RWMutex     sync.RWMutex
// 		redis       *redis.Client
// 		once        sync.Once
// 		errTaskID   string
// 		errTask     string
// 		errCode     int
// 		errMsg      string
// 		errTasktype define.TaskRespType
// 	}
// 	type args struct {
// 		tasrunktype define.TaskRespType
// 		realid      string
// 		value       interface{}
// 		setdata     string
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t := &task2{
// 				id:          tt.fields.id,
// 				name:        tt.fields.name,
// 				cronexpr:    tt.fields.cronexpr,
// 				cronsub:     tt.fields.cronsub,
// 				close:       tt.fields.close,
// 				ctxcancel:   tt.fields.ctxcancel,
// 				next:        tt.fields.next,
// 				canrun:      tt.fields.canrun,
// 				RWMutex:     tt.fields.RWMutex,
// 				redis:       tt.fields.redis,
// 				once:        tt.fields.once,
// 				errTaskID:   tt.fields.errTaskID,
// 				errTask:     tt.fields.errTask,
// 				errCode:     tt.fields.errCode,
// 				errMsg:      tt.fields.errMsg,
// 				errTasktype: tt.fields.errTasktype,
// 			}
// 			if err := t.SetData(tt.args.tasrunktype, tt.args.realid, tt.args.value, tt.args.setdata); (err != nil) != tt.wantErr {
// 				t.Errorf("task2.SetData() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
