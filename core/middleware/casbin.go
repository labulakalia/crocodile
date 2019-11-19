package middleware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/labulaka521/crocodile/common/jwt"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/model"
	"github.com/labulaka521/crocodile/core/utils/resp"
	"github.com/pkg/errors"
	"go.uber.org/zap"
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
	if err != nil || claims.UId == "" {
		return false, err
	}
	if !claims.VerifyExpiresAt(time.Now().Unix(), false) {
		return false, err
	}

	c.Set("uid", claims.UId)
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	ok, err := model.Check(ctx, model.TB_user, model.UID, claims.UId)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}

	role, err := model.QueryUserRule(ctx, claims.UId)
	if err != nil {
		log.Error("QueryUserRule failed", zap.String("error", err.Error()))
		resp.Json(c, resp.ErrInternalServer, nil)
		return
	}
	c.Set("role", role)

	requrl := c.Request.URL.Path
	method := c.Request.Method
	enforcer := model.GetEnforcer()
	return enforcer.Enforce(claims.UId, requrl, method)
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
		if err != nil {
			code = resp.ErrUnauthorized
			goto ERR
		}
		if !pass {
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
