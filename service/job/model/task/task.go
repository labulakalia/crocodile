package task

import (
	"context"
	"crocodile/common/event"
	"crocodile/common/util"
	pbactator "crocodile/service/actuator/proto/actuator"
	pbexecutor "crocodile/service/executor/proto/executor"
	pbjob "crocodile/service/job/proto/job"
	"database/sql"
	"errors"
	"github.com/golang/protobuf/ptypes"
	"github.com/labulaka521/logging"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"math/rand"
	"time"
)

// 负责任务的增删该查和定时运行

type Servicer interface {
	CreateJob(ctx context.Context, task *pbjob.Task) (err error)
	DeleteJob(ctx context.Context, taskname string) (err error)
	ChangeJob(ctx context.Context, task *pbjob.Task) (err error)
	GetJob(ctx context.Context, taskname string) (resp []*pbjob.Task, err error)
	KillJob(ctx context.Context, task *pbjob.Task) (err error)
	RunJob(ctx context.Context, task *pbjob.Task) (err error)
	UpdateNextTime(ctx context.Context, taskname string, nexttime time.Time) (err error)
}

var _ Servicer = &Service{}

var (
	Pub           micro.Publisher
	ActatorCLient pbactator.ActuatorService
)

func Init(cient client.Client) {
	Pub = micro.NewPublisher("crocodile.srv.executor", cient)
	ActatorCLient = pbactator.NewActuatorService("crocodile.srv.actuator", cient)
}

type Service struct {
	DB            *sql.DB
	Pub           *micro.Publisher
	ActatorCLient *pbactator.ActuatorService
}

// 创建任务
// required
//   taskname
//	 command
//   cronexpr
//   createdby
//	 executors
// optional
//   retry
//   stop
//	 remark
//   timeout
func (s *Service) CreateJob(ctx context.Context, task *pbjob.Task) (err error) {
	var (
		createjob_sql string
		stmt          *sql.Stmt
		now           time.Time
		nexttime      time.Time
	)

	createjob_sql = `INSERT INTO crocodile_task 
					(taskname,command,cronexpr,createdby,remark,stop,timeout,nexttime,actuator)
					VALUE(?,?,?,?,?,?,?,?,?)
					`
	//

	if s.isExist(ctx, task.Taskname) {
		err = errors.New("Task %s Alreay Exist" + task.Taskname)
		return
	}

	// 获取程序下一次运行时间
	now = time.Now().Local()
	if nexttime, err = util.NextTime(task.Cronexpr, now); err != nil {
		return
	}
	logging.Info("change", nexttime)
	if task.Nexttime, err = ptypes.TimestampProto(nexttime); err != nil {
		return
	}

	if stmt, err = s.DB.PrepareContext(ctx, createjob_sql); err != nil {
		logging.Errorf("Prepare SQL %s Err: %v", createjob_sql, err)
		return
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, task.Taskname, task.Command, task.Cronexpr, task.Createdby, task.Remark, task.Stop,
		task.Timeout, nexttime, task.Actuator); err != nil {
		logging.Errorf("Exec Context Err: %v", err)
		return
	}

	logging.Debugf("Create Job success")
	return
}

// required
//   id
func (s *Service) DeleteJob(ctx context.Context, taskname string) (err error) {
	var (
		deletejob_sql string
		stmt          *sql.Stmt
	)

	deletejob_sql = "DELETE FROM crocodile_task WHERE taskname=?"
	if !s.isExist(ctx, taskname) {
		err = errors.New("Task Not Exist : " + taskname)
		return
	}
	if stmt, err = s.DB.PrepareContext(ctx, deletejob_sql); err != nil {
		logging.Errorf("Prepare SQL %s Err: %v", deletejob_sql, err)
		return
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, taskname); err != nil {
		logging.Errorf("Exec Context Err: %v", err)
		return
	}

	logging.Debugf("Delete Job success ")
	return
}

//
// required
//   id
//   taskname
//	 command
//   cronexpr
//   createdby
//	 executors
//   retry
//	 remark
//   timeout
//   stop
//   nexttime
//   actuator_id

func (s *Service) ChangeJob(ctx context.Context, task *pbjob.Task) (err error) {
	var (
		changejob_sql string
		stmt          *sql.Stmt
	)

	changejob_sql = `UPDATE crocodile_task 
					SET command=?,cronexpr=?,remark=?,stop=?,timeout=?,actuator=? 
					WHERE taskname=?`
	if !s.isExist(ctx, task.Taskname) {
		err = errors.New("Task Not Exist: " + task.Taskname)
		return
	}
	if stmt, err = s.DB.PrepareContext(ctx, changejob_sql); err != nil {
		logging.Errorf("Prepare SQL %s Err: %v", changejob_sql, err)
		return
	}
	defer stmt.Close()
	if _, err = stmt.ExecContext(ctx, task.Command, task.Cronexpr,
		task.Remark, task.Stop, task.Timeout, task.Actuator, task.Taskname); err != nil {
		logging.Errorf("Exec Context Err: %v", err)
		return
	}

	logging.Debugf("Change Job success")
	return
}

func (s *Service) UpdateNextTime(ctx context.Context, taskname string, nexttime time.Time) (err error) {
	var (
		updatejob_sql string
		stmt          *sql.Stmt
	)

	updatejob_sql = `UPDATE crocodile_task SET nexttime=? WHERE taskname=?`
	if stmt, err = s.DB.PrepareContext(ctx, updatejob_sql); err != nil {
		logging.Errorf("Prepare SQL %s Err: %v", updatejob_sql, err)
		return
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, nexttime, taskname); err != nil {
		logging.Errorf("Exec Context Err: %v", err)
		return
	}
	logging.Debugf("Update NextRun Time Success")
	return
}

// optional
//   id
func (s *Service) GetJob(ctx context.Context, taskname string) (resp []*pbjob.Task, err error) {
	var (
		getjob_sql string
		stmt       *sql.Stmt

		rows     *sql.Rows
		nextTime time.Time
	)
	resp = []*pbjob.Task{}
	if taskname == "" {
		taskname = "%"
	}
	getjob_sql = `SELECT id, taskname,command,cronexpr,createdby,remark,stop,timeout,nexttime,actuator
				FROM crocodile_task
				WHERE taskname LIKE ?`

	if stmt, err = s.DB.PrepareContext(ctx, getjob_sql); err != nil {
		logging.Errorf("Prepare SQL %s Err: %v", getjob_sql, err)
		return
	}
	defer stmt.Close()
	if rows, err = stmt.QueryContext(ctx, taskname); err != nil {

		logging.Errorf("SQL %s Query Err: %v", getjob_sql, err)
		return
	}

	for rows.Next() {
		task := pbjob.Task{}
		if err = rows.Scan(&task.Id, &task.Taskname, &task.Command, &task.Cronexpr,
			&task.Createdby, &task.Remark, &task.Stop, &task.Timeout, &nextTime, &task.Actuator); err != nil {
			logging.Errorf("Scan Task Err: %v\n", err)
			continue
		}
		task.Nexttime, _ = ptypes.TimestampProto(nextTime)
		resp = append(resp, &task)
	}

	logging.Debugf("GetJob success")
	return
}

// 强杀任务
// 发布一个强杀命令的消息
// executor 会接收到这个消息，然后查询自已是否拥有这个任务 并且是正在运行的时候 如果存在就杀死这个任务
// required
//   id
func (s *Service) KillJob(ctx context.Context, task *pbjob.Task) (err error) {

	exmsg := &pbexecutor.ExecuteMsg{
		Event: event.Kill_Task,
		Task:  task,
	}

	if err = Pub.Publish(ctx, exmsg); err != nil {
		logging.Errorf("Publish Err: %v", err)
	}
	return
}

// 发布一个执行的消息时 会附带这个任务的执行主机的IP 执行者收到这个订阅时 会先匹配自已的ip是否一致如果一致就运行
// required
// id
func (s *Service) RunJob(ctx context.Context, task *pbjob.Task) (err error) {
	var (
		resp   *pbactator.Response
		actuat *pbactator.Actuat
		addrs  []*pbactator.Addr
		addr   *pbactator.Addr
		exmsg  *pbexecutor.ExecuteMsg
	)
	// TODO 主机选择算法
	// 随机

	// 轮询
	// 权重
	// 最少负载

	rand.Seed(time.Now().Unix())
	actuat = &pbactator.Actuat{
		Name: task.Actuator,
	}

	if resp, err = ActatorCLient.GetActuator(ctx, actuat); err != nil {
		logging.Errorf("Get Actuator Err: %v", err)
		return
	}
	if len(resp.Actuators) != 1 {
		logging.Errorf("No Get  Actuator %s ", task.Taskname)
		return
	}
	// 执行器中包的的主机IP
	addrs = resp.Actuators[0].Address

	for i := 0; i < len(addrs)*2; i++ {
		// 随机选择一个在线的执行节点
		if addrs[rand.Intn(len(addrs))].Online {
			addr = addrs[rand.Intn(len(addrs))]
			break
		}
	}
	if addr == nil {
		logging.Warnf("No Found Online Host")
		err = errors.New("No Found Online Host")
		return
	}
	logging.Debugf("Select Run Host: %s", addr.Ip)

	exmsg = &pbexecutor.ExecuteMsg{
		Event:   event.Run_Task,
		Task:    task,
		Runhost: addr.Ip,
	}
	logging.Debugf("New Publish Job %v", exmsg)
	if err = Pub.Publish(ctx, exmsg); err != nil {
		logging.Errorf("Publish Err: %v", err)
	}
	return
}

// 任务是否存在
func (s *Service) isExist(ctx context.Context, taskname string) (exits bool) {
	var (
		resp []*pbjob.Task
		err  error
	)

	if resp, err = s.GetJob(ctx, taskname); err != nil {
		logging.Errorf("GetJob Err: %v", err)
		return false
	}
	if len(resp) == 0 {
		return false
	}
	return true

}
