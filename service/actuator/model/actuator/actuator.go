package actuator

import (
	"context"
	"crocodile/common/cfg"
	"crocodile/common/registry"
	pbactuat "crocodile/service/actuator/proto/actuator"
	"database/sql"
	"fmt"
	"github.com/labulaka521/logging"
	"strings"
)

// 执行器的增删改查

type Servicer interface {
	CreateActuator(ctx context.Context, actuat *pbactuat.Actuat) (err error)
	DeleteActuator(ctx context.Context, tname string) (err error)
	ChangeActuator(ctx context.Context, actuat *pbactuat.Actuat) (err error)
	// 获取执行器 获取所有的数据 执行的模块通过自已的IP获取自已所属的执行器
	GetActuator(ctx context.Context, name string) (actuats []*pbactuat.Actuat, err error)
	// 获取Executor全部的注册IP
	GetAllExecutorIP(ctx context.Context) (resp []string, err error)
}

var _ Servicer = &Service{}

type Service struct {
	DB *sql.DB
}

func (s *Service) CreateActuator(ctx context.Context, actuat *pbactuat.Actuat) (err error) {
	var (
		createActuator_sql string
		stmt               *sql.Stmt
		alladdress         []string
		addr               *pbactuat.Addr
	)
	createActuator_sql = "INSERT INTO crocodile_actuator (name,address,createdby) VALUE (?,?,?)"

	for _, addr = range actuat.Address {
		alladdress = append(alladdress, addr.Ip)
	}

	if s.isExist(ctx, actuat.Name) {
		err = fmt.Errorf("Actuat %s alreay Exists", actuat.Name)
		return
	}

	if stmt, err = s.DB.PrepareContext(ctx, createActuator_sql); err != nil {
		logging.Errorf("Prepare SQL %s Err: %v", createActuator_sql, err)
		return
	}
	defer stmt.Close()
	if _, err = stmt.ExecContext(ctx, actuat.Name, strings.Join(alladdress, ","), actuat.Createdby); err != nil {
		logging.Errorf("Exec Err:%v", err)
		return
	}

	return
}
func (s *Service) DeleteActuator(ctx context.Context, name string) (err error) {
	var (
		deleteActuator_sql string
		stmt               *sql.Stmt
	)
	deleteActuator_sql = "DELETE FROM crocodile_actuator WHERE name=?"
	if stmt, err = s.DB.PrepareContext(ctx, deleteActuator_sql); err != nil {
		logging.Errorf("Prepare SQL %s Err: %v", deleteActuator_sql, err)
		return
	}
	if _, err = stmt.ExecContext(ctx, name); err != nil {
		logging.Errorf("Exec Err:%v", err)
		return
	}
	return
}

func (s *Service) ChangeActuator(ctx context.Context, actuat *pbactuat.Actuat) (err error) {
	var (
		changeActuator_sql string
		alladdress         []string
		addr               *pbactuat.Addr
		stmt               *sql.Stmt
	)
	changeActuator_sql = "UPDATE crocodile_actuator SET address=? WHERE name=?"
	if !s.isExist(ctx, actuat.Name) {
		err = fmt.Errorf("Actuat %s Not Exists", actuat.Name)
		return
	}
	for _, addr = range actuat.Address {
		alladdress = append(alladdress, addr.Ip)
	}

	if stmt, err = s.DB.PrepareContext(ctx, changeActuator_sql); err != nil {
		logging.Errorf("Prepare SQL %s Err: %v", changeActuator_sql, err)
	}
	defer stmt.Close()
	if _, err = stmt.ExecContext(ctx, strings.Join(alladdress, ","), actuat.Name); err != nil {
		logging.Errorf("Exec Err:%v", err)
		return
	}
	return
}

// 获取所有的执行器
func (s *Service) GetActuator(ctx context.Context, name string) (actuats []*pbactuat.Actuat, err error) {
	var (
		stmt            *sql.Stmt
		getActuator_sql string
		rows            *sql.Rows
		address         string
		allExecutorIP   []string
	)
	actuats = []*pbactuat.Actuat{}
	if name == "" {
		name = "%"
	}

	getActuator_sql = "SELECT * FROM crocodile_actuator WHERE name LIKE ?"
	// 所有在线的执行器
	if allExecutorIP, err = s.GetAllExecutorIP(ctx); err != nil {
		logging.Errorf("Get AllExecutorIP Err: %v", err)
		return
	}

	if stmt, err = s.DB.PrepareContext(ctx, getActuator_sql); err != nil {
		logging.Errorf("Prepare SQL %s Err: %v", getActuator_sql, err)
		return
	}
	if rows, err = stmt.QueryContext(ctx, name); err != nil {
		logging.Errorf("Exec Err:%v", err)
		return
	}

	for rows.Next() {
		actuat := pbactuat.Actuat{}

		if err = rows.Scan(&actuat.Id, &actuat.Name, &address, &actuat.Createdby); err != nil {
			continue
		}

		for _, ip := range strings.Split(address, ",") {
			var (
				exits bool
			)
			if ip == "" {
				continue
			}
			// 检查执行器分配的IP是否在线
			for _, eip := range allExecutorIP {
				if eip == ip {
					exits = true
				}
			}

			addr := pbactuat.Addr{
				Ip:     ip,
				Online: exits,
			}
			actuat.Address = append(actuat.Address, &addr)
		}

		actuats = append(actuats, &actuat)
	}

	return
}

// 获取所有注册的executor的IP
func (s *Service) GetAllExecutorIP(ctx context.Context) (resp []string, err error) {
	const srvname = "topic:crocodile.srv.executor"
	resp, err = registry.GetEtcdListServicesIP(srvname, cfg.EtcdConfig.Endpoints...)
	return
}

func (s *Service) isExist(ctx context.Context, name string) bool {
	var (
		actuats []*pbactuat.Actuat
		err     error
	)

	if actuats, err = s.GetActuator(ctx, name); err != nil {
		logging.Errorf("GetActuator Err: %v", err)
		return false
	}
	if len(actuats) == 0 {
		return false
	}
	return true
}
