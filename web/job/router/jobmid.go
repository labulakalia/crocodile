package router

import (
	"context"
	"crocodile/common/bind"
	"crocodile/common/e"
	"crocodile/common/response"
	pbjob "crocodile/service/job/proto/job"
	"crocodile/web/job/router/job"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/labulaka521/logging"
	"strings"
)

// 普通用户只可以修改删除自已创建的任务
// 也可以使用GET操作获取全部的任务
var (
	joballowurl = []string{"/job/create"}
)

func JobControl() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err     error
			app     response.Gin
			rooturl string
			url     string
			super   bool
			exits   bool
			code    int32
			user    string
			resp    *pbjob.Response
			qtask   pbjob.Task
		)

		app = response.Gin{c}
		// GET操作允许
		if c.Request.Method == "GET" {
			c.Next()
			return
		}
		// 创建任务
		rooturl = strings.Split(c.Request.RequestURI, "?")[0]
		for _, url = range joballowurl {
			if rooturl == url {
				c.Next()
				return
			}
		}

		qtask = pbjob.Task{}
		if err = bind.BindJson(c, &qtask); err != nil {
			code = e.ERR_BAD_REQUEST
			goto ERR
		}
		c.Set("data", qtask)
		// 管理员
		if super, exits = c.Keys["super"].(bool); !exits {
			code = e.ERR_TOKEN_INVALID
			goto ERR
		}
		if super {
			c.Next()
			return
		}

		// 检查当前任务的创建人是否与当前执行的用户是否一致
		if user, exits = c.Keys["user"].(string); !exits {
			code = e.ERR_TOKEN_INVALID
			goto ERR
		}

		if resp, err = job.JobClient.GetJob(context.Background(), &qtask); err != nil {
			code = e.ERR_GET_JOB_FAIL
			goto ERR
		}

		if len(resp.Tasks) != 1 {
			code = e.ERR_JOB_NOT_EXITS
			goto ERR
		}

		if resp.Tasks[0].Createdby != user {
			code = e.ERR_NOT_PERMISSION
			goto ERR
		}
		// 设置task名称

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
