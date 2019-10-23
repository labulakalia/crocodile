package middleware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/labulaka521/crocodile/common/jwt"
	"github.com/labulaka521/crocodile/core/model"
	"github.com/labulaka521/crocodile/core/utils/resp"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

const (
	tokenpre = "Bearer "
)

// 权限检查
func checkAuth(c *gin.Context) (pass bool, err error) {

	token := strings.TrimPrefix(c.GetHeader("Authorization"), tokenpre)

	if token == "" {
		err = errors.New("invalid token")
		return
	}

	claims, err := jwt.ParseToken(token)
	if err != nil || claims.UId == 0 {
		return
	}
	if !claims.VerifyExpiresAt(time.Now().Unix(), false) {
		return
	}

	c.Set("uid", claims.UId)
	exist, err := model.IsExist(context.Background(), model.Uid, claims.UId)
	if err != nil {
		return false, err
	}
	if !exist {
		return false, nil
	}
	requrl := c.Request.URL.Path
	method := c.Request.Method
	struid := strconv.FormatInt(claims.UId, 10)
	enforcer := model.GetEnforcer()
	return enforcer.Enforce(struid, requrl, method)
}

var excludepath = []string{"login"}

// 权限控制middle
func PermissionControl() func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			code = resp.Success
			err  error
		)
		for _, url := range excludepath {
			if strings.Contains(c.Request.RequestURI, url) {
				c.Next()
				return
			}
		}
		defer func() {
			c.Set("statuscode", code)
		}()

		pass, err := checkAuth(c)
		if err != nil || !pass {
			code = resp.ErrUnauthorized
			goto ERR
		}
		c.Next()
		return

	ERR:
		// 解析失败返回错误
		c.Writer.Header().Add("WWW-Authenticate", fmt.Sprintf("Bearer realm='%s'", resp.GetMsg(code)))
		resp.Json(c, resp.ErrUnauthorized, nil)
		c.Abort()
		return
	}
}
