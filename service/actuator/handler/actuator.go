package handler

import (
	"context"
	"crocodile/service/actuator/model/actuator"
	pbactuat "crocodile/service/actuator/proto/actuator"
	"github.com/labulaka521/logging"
)

type Actua struct {
	Service actuator.Servicer
}

func (a *Actua) CreateActuator(ctx context.Context, actuat *pbactuat.Actuat, resp *pbactuat.Response) (err error) {
	logging.Debugf("CreateActuator %s", actuat.Name)
	if err = a.Service.CreateActuator(ctx, actuat); err != nil {
		logging.Errorf("CreateActuator Err:%v", err)
	}
	return
}
func (a *Actua) DeleteActuator(ctx context.Context, actuat *pbactuat.Actuat, resp *pbactuat.Response) (err error) {
	logging.Debugf("DeleteActuator %s", actuat.Name)
	if err = a.Service.DeleteActuator(ctx, actuat.Name); err != nil {
		logging.Errorf("DeleteActuator Err:%v", err)
	}
	return
}

func (a *Actua) ChangeActuator(ctx context.Context, actuat *pbactuat.Actuat, resp *pbactuat.Response) (err error) {
	logging.Debugf("ChangeActuator %s", actuat.Name)
	if err = a.Service.ChangeActuator(ctx, actuat); err != nil {
		logging.Errorf("ChangeActuator Err:%v", err)
	}
	return
}
func (a *Actua) GetActuator(ctx context.Context, actuat *pbactuat.Actuat, resp *pbactuat.Response) (err error) {
	logging.Debugf("GetActuator %s", actuat.Name)
	resp.Actuators = []*pbactuat.Actuat{}
	if resp.Actuators, err = a.Service.GetActuator(ctx, actuat.Name); err != nil {
		logging.Errorf("GetActuator Err:%v", err)
	}
	return
}

func (a *Actua) GetAllExecutorIP(ctx context.Context, actuat *pbactuat.Actuat, resp *pbactuat.Response) (err error) {
	logging.Debugf("Get All ExecutorIP")

	if resp.ExecutorIps, err = a.Service.GetAllExecutorIP(ctx); err != nil {
		logging.Errorf("GetAllExecutorIP Err:%v", err)
	}
	return
}
