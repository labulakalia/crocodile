package router

// 用户权限控制
// 除登录注销外普通的用户的请求方法只可以为GET
// 超级用户具有全部操作

import (
	"crocodile/common/e"
	"crocodile/common/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/labulaka521/logging"
	"strings"
)

var (
	userallowurl = []string{"/auth/login", "/auth/logout"}
)

func UserControl() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err     error
			app     response.Gin
			rooturl string
			url     string
			super   bool
			exits   bool
			code    int32
		)
		app = response.Gin{c}
		if c.Request.Method == "GET" {
			c.Next()
			return
		}
		rooturl = strings.Split(c.Request.RequestURI, "?")[0]
		for _, url = range userallowurl {
			if rooturl == url {
				c.Next()
				return
			}
		}

		if super, exits = c.Keys["super"].(bool); !exits {
			code = e.ERR_TOKEN_INVALID
			goto ERR
		}
		if !super {
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
