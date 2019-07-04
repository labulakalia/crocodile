package router

// 用户权限控制
// 除登录注销外普通的用户的请求方法只可以为GET
// 超级用户具有全部操作

import (
	"context"
	"crocodile/common/bind"
	"crocodile/common/e"
	"crocodile/common/response"
	pbactuator "crocodile/service/actuator/proto/actuator"
	"crocodile/web/actuator/router/actuator"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/labulaka521/logging"
	"strings"
)

var (
	actuatorallowurl = []string{"/actuator/create"}
)

func ActuatorControl() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err         error
			app         response.Gin
			rooturl     string
			url         string
			super       bool
			exits       bool
			code        int32
			user        string
			resp        *pbactuator.Response
			reqactuator pbactuator.Actuat
		)

		app = response.Gin{c}
		// GET操作允许
		if c.Request.Method == "GET" {
			c.Next()
			return
		}
		// 创建任务
		rooturl = strings.Split(c.Request.RequestURI, "?")[0]
		for _, url = range actuatorallowurl {
			if rooturl == url {
				c.Next()
				return
			}
		}

		reqactuator = pbactuator.Actuat{}
		if err = bind.BindJson(c, &reqactuator); err != nil {
			code = e.ERR_BAD_REQUEST
			goto ERR
		}

		c.Set("data", reqactuator)
		// 管理员
		if super, exits = c.Keys["super"].(bool); !exits {
			code = e.ERR_TOKEN_INVALID
			goto ERR
		}
		if super {
			c.Next()
			return
		}
		// 检查当前执行器的创建人是否与当前执行的用户是否一致
		if user, exits = c.Keys["user"].(string); !exits {
			code = e.ERR_TOKEN_INVALID
			goto ERR
		}

		resp, err = actuator.ActuatorClient.GetActuator(context.Background(), &reqactuator)

		if len(resp.Actuators) != 1 {
			code = e.ERR_ACTUATOR_NOT_EXITS
			goto ERR
		}

		if resp.Actuators[0].Createdby != user {
			code = e.ERR_NOT_PERMISSION
			goto ERR
		}

		c.Next()
		return

	ERR:
		// 解析失败返回错误
		c.Writer.Header().Add("WWW-Authenticate", fmt.Sprintf("Bearer realm='%s'", e.GetMsg(code)))
		app.Response(code, nil)
		logging.Errorf("Token Check Fail: %v %s", err, e.GetMsg(code))
		c.Abort()
		return
	}
}
