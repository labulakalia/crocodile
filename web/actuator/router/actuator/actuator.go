package actuator

import (
	"context"
	"crocodile/common/bind"
	"crocodile/common/cfg"
	"crocodile/common/e"
	"crocodile/common/registry"
	"crocodile/common/response"
	"crocodile/common/wrapper"
	pbactuator "crocodile/service/actuator/proto/actuator"
	"github.com/gin-gonic/gin"
	"github.com/labulaka521/logging"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-plugins/wrapper/trace/opentracing"
	"time"
)

var (
	ActuatorClient pbactuator.ActuatorService
)

func Init() {
	c := client.NewClient(
		client.Retries(3),
		client.Registry(registry.Etcd(cfg.EtcdConfig.Endpoints...)),
		client.Wrap(opentracing.NewClientWrapper()),
	)
	ActuatorClient = pbactuator.NewActuatorService("crocodile.srv.actuator", c)

}

type QueryActuat struct {
	Name string `json:"name" validate:"required"`
}
type Actuat struct {
	Name      string `json:"name" validate:"required"`
	Address   []Addr `json:"address" validate:"required"`
	Createdby string `json:"createdby" validate:"required"`
}

type Addr struct {
	Ip string `json:"ip" validate:"required"`
}

// ""
// {
//    "name": "",
// 	  "address": [
// 	  		{"ip": "ip1"}
// 	  	]
// }
func CreateActuator(c *gin.Context) {
	var (
		app         response.Gin
		ctx         context.Context
		err         error
		loginuser   string
		exists      bool
		code        int32
		reqactuator pbactuator.Actuat
	)
	ctx, ok := wrapper.ContextWithSpan(c)
	if ok == false {
		logging.Error("get context err")
		ctx = context.Background()
	}
	ctx, _ = context.WithTimeout(ctx, time.Duration(cfg.MysqlConfig.MaxQueryTime)*time.Second)
	app = response.Gin{c}
	code = e.SUCCESS
	reqactuator = pbactuator.Actuat{}
	if err = bind.BindJson(c, &reqactuator); err != nil {
		code = e.ERR_BAD_REQUEST
		app.Response(code, nil)
		return
	}
	if loginuser, exists = c.Keys["user"].(string); !exists {
		code = e.ERR_TOKEN_INVALID
		app.Response(code, nil)
		return
	}

	reqactuator.Createdby = loginuser

	_, err = ActuatorClient.CreateActuator(ctx, &reqactuator)
	if err != nil {
		logging.Errorf("CreateActuator Err: %v", err)
		code = e.ERR_CREATE_ACTUAT_FAIL
	}
	app.Response(code, nil)
}
func DeleteActuator(c *gin.Context) {
	var (
		app            response.Gin
		deleteactuator pbactuator.Actuat
		ctx            context.Context
		err            error
		code           int32
		exits          bool
	)
	ctx, ok := wrapper.ContextWithSpan(c)
	if ok == false {
		logging.Error("get context err")
		ctx = context.Background()
	}
	ctx, _ = context.WithTimeout(ctx, time.Duration(cfg.MysqlConfig.MaxQueryTime)*time.Second)
	app = response.Gin{c}
	code = e.SUCCESS
	deleteactuator, exits = c.Keys["data"].(pbactuator.Actuat)
	if !exits {
		code = e.ERR_BAD_REQUEST
		app.Response(code, nil)
		return
	}

	_, err = ActuatorClient.DeleteActuator(ctx, &deleteactuator)
	if err != nil {
		logging.Errorf("DeleteActuator Err: %v", err)
		code = e.ERR_DELETE_ACTUAT_FAIL
	}
	app.Response(code, nil)
}

func ChangeActuator(c *gin.Context) {
	var (
		app            response.Gin
		changeactuator pbactuator.Actuat
		ctx            context.Context
		err            error
		code           int32
		exits          bool
	)
	ctx, ok := wrapper.ContextWithSpan(c)
	if ok == false {
		logging.Error("get context err")
		ctx = context.Background()
	}
	ctx, _ = context.WithTimeout(ctx, time.Duration(cfg.MysqlConfig.MaxQueryTime)*time.Second)
	app = response.Gin{c}
	code = e.SUCCESS
	changeactuator = pbactuator.Actuat{}
	changeactuator, exits = c.Keys["data"].(pbactuator.Actuat)
	if !exits {
		logging.Errorf("Not Exits Actuator Data")
		code = e.ERR_BAD_REQUEST
		app.Response(code, nil)
		return
	}

	_, err = ActuatorClient.ChangeActuator(ctx, &changeactuator)
	if err != nil {
		logging.Errorf("CreateActuator Err: %v", err)
		code = e.ERR_CHANGE_ACTUAT_FAIL
	}
	app.Response(code, nil)
}

func GetActuator(c *gin.Context) {
	var (
		app          response.Gin
		queryctuator pbactuator.Actuat
		ctx          context.Context
		err          error
		code         int32
		rsp          *pbactuator.Response
	)
	ctx, ok := wrapper.ContextWithSpan(c)
	if ok == false {
		logging.Error("get context err")
		ctx = context.Background()
	}
	queryctuator = pbactuator.Actuat{}
	ctx, _ = context.WithTimeout(ctx, time.Duration(cfg.MysqlConfig.MaxQueryTime)*time.Second)
	app = response.Gin{c}
	code = e.SUCCESS
	if err = bind.BindQuery(c, &queryctuator); err != nil {
		code = e.ERR_BAD_REQUEST
		app.Response(code, nil)
		return
	}

	rsp, err = ActuatorClient.GetActuator(ctx, &queryctuator)
	if err != nil {
		logging.Errorf("CreateActuator Err: %v", err)
		code = e.ERR_CHANGE_ACTUAT_FAIL
	}
	app.Response(code, rsp.Actuators)
}

func GetALLExecuteIP(c *gin.Context) {
	var (
		app  response.Gin
		ctx  context.Context
		err  error
		code int32
		rsp  *pbactuator.Response
	)
	ctx, ok := wrapper.ContextWithSpan(c)
	if ok == false {
		logging.Error("get context err")
		ctx = context.Background()
	}
	ctx, _ = context.WithTimeout(ctx, time.Duration(cfg.MysqlConfig.MaxQueryTime)*time.Second)
	app = response.Gin{c}
	code = e.SUCCESS
	rsp, err = ActuatorClient.GetAllExecutorIP(ctx, new(pbactuator.Actuat))
	if err != nil {
		logging.Errorf("GetAllExecutorIP Err:%v", err)
		code = e.ERR_GET_EXECUTOR_IP_FAIL
	}
	app.Response(code, rsp.ExecutorIps)
}
