package schedule

import (
	"context"
	"fmt"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/model"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"strings"
	"sync"
	"time"
)

var (
	_Parse *cron.Parser
	Cron   *cacheSchedule
)

func init() {
	_Parse = getParse()
	Cron = &cacheSchedule{
		sch: make(map[string]*task),
	}
}

type task struct {
	cronexpr  string
	close     chan struct{}
	failCount int  // run failed count
	runTotal  int  // totoal run count
	running   bool // task is running
}

type cacheSchedule struct {
	sync.RWMutex
	sch map[string]*task
}

// start run already exists task from db
func Init() {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	eps, err := model.GetTasks(ctx)
	if err != nil {
		log.Error("GetTasks failed", zap.Error(err))
	}
	for _, t := range eps {
		Cron.Add(t.Id, t.CronExpr)
	}
	log.Info("Init already exist cron task", zap.Int("Total", len(eps)))
}

// add task to schedule
func (s *cacheSchedule) Add(taskId string, cronExpr string) {

	cronexpr := fullcronExpr(cronExpr)
	s.Del(taskId)
	t := task{
		cronexpr: cronexpr,
		close:    make(chan struct{}),
	}
	s.Lock()
	s.sch[taskId] = &t
	s.Unlock()
	log.Info("Add Task success", zap.String("taskid", taskId))
	go s.runSchedule(taskId)
}

// start run cronexpr schedule
func (s *cacheSchedule) runSchedule(taskid string) {
	t, exist := s.gettask(taskid)
	if !exist {
		return
	}
	log.Info("Start Run Cronexpr", zap.String("cronexpr", t.cronexpr), zap.String("id", taskid))
	sch, _ := _Parse.Parse(t.cronexpr)
	for {
		now := time.Now()
		next := sch.Next(now)
		select {
		case <-t.close:
			log.Info("Close Schedule", zap.String("taskID", taskid), zap.Any("task", t))
			return
		case <-time.After(next.Sub(now)):
			go func() {
				err := s.RunTask(taskid)
				if err != nil {
					log.Error("ExecPlan failed", zap.Error(err))
				}
			}()
		}
	}
}

// del schedule
func (s *cacheSchedule) Del(id string) {
	t, ok := s.gettask(id)
	if ok {
		s.Lock()
		close(t.close)
		delete(s.sch, id)
		s.Unlock()
		log.Info("Del task success", zap.String("taskid", id))
		return
	}
}

func (s *cacheSchedule) gettask(id string) (*task, bool) {
	s.Lock()
	defer s.Unlock()

	t, ok := s.sch[id]
	return t, ok
}

// start run task by execplanid
func (s *cacheSchedule) RunTask(id string) error {

	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	task, err := model.GetTaskByID(ctx, id)
	if err != nil {
		return err
	}
	if len(task.ParentTaskIds) != 0 {
		log.Info("Start Run ParentTasks", zap.Strings("taskids", task.ParentTaskIds))
		s.runMultiTasks(ctx, task.ParentRunParallel, task.ParentTaskIds...)
	}
	log.Info("Start Run Task", zap.String("taskid", task.Id))
	s.runTask(ctx, task.Id)
	if len(task.ChildTaskIds) != 0 {
		log.Info("Start Run ChildTasks", zap.Strings("taskids", task.ChildTaskIds))
		s.runMultiTasks(ctx, task.ChildRunParallel, task.ChildTaskIds...)
	}

	return nil
}
func (s *cacheSchedule) runTask(ctx context.Context, id string) error {
	task, err := model.GetTaskByID(ctx, id)
	if err != nil {
		return err
	}
	task = task

	return nil
}

func (s *cacheSchedule) runMultiTasks(ctx context.Context, RunParallel int, taskids ...string) {
	if RunParallel == 1 {
		var wg sync.WaitGroup
		wg.Add(len(taskids))
		for _, id := range taskids {
			go func(id string) {
				s.runTask(ctx, id)
				wg.Done()
			}(id)
		}
		wg.Wait()
	} else {
		for _, id := range taskids {
			s.runTask(ctx, id)

		}
	}

}

func fullcronExpr(cronexpr string) string {
	if len(strings.Fields(cronexpr)) == 5 {
		cronexpr = fmt.Sprintf("* %s", cronexpr)
	}
	return cronexpr
}

func getParse() *cron.Parser {
	p := cron.NewParser(
		cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)
	return &p
}
