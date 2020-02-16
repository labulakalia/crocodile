package middleware

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/labulaka521/crocodile/common/jwt"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/model"
	"github.com/labulaka521/crocodile/core/utils/resp"
)

const (
	tokenpre = "Bearer "
)

// CheckToken check token is valid
func CheckToken(token string) (string, bool) {
	claims, err := jwt.ParseToken(token)
	if err != nil || claims.UID == "" {
		log.Error("ParseToken failed", zap.Error(err))
		return "", false
	}
	if !claims.VerifyExpiresAt(time.Now().Unix(), false) {
		log.Error("Token is Expire", zap.String("token", token))
		return "", false
	}

	return claims.UID, true
}

// 权限检查
func checkAuth(c *gin.Context) (pass bool, err error) {
	token := strings.TrimPrefix(c.GetHeader("Authorization"), tokenpre)

	if token == "" {
		err = errors.New("invalid token")
		return
	}

	// claims, err := jwt.ParseToken(token)
	// if err != nil || claims.UID == "" {
	// 	return false, err
	// }
	// if !claims.VerifyExpiresAt(time.Now().Unix(), false) {
	// 	return false, err
	// }
	// fmt.Printf("%+v", claims)
	id, pass := CheckToken(token)
	if !pass {
		return false, errors.New("CheckToken failed")
	}
	// fmt.Println(id)
	// id := claims.Id
	c.Set("uid", id)
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	ok, err := model.Check(ctx, model.TBUser, model.UID, id)
	if err != nil {
		return false, err
	}
	if !ok {
		log.Error("Check UID not exist", zap.String("uid", id))
		return false, nil
	}

	role, err := model.QueryUserRule(ctx, id)
	if err != nil {
		log.Error("QueryUserRule failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	c.Set("role", role)

	requrl := c.Request.URL.Path
	method := c.Request.Method
	enforcer := model.GetEnforcer()
	return enforcer.Enforce(id, requrl, method)
}

var excludepath = []string{"login", "swagger", "websocket","/debug/pprof"}

// PermissionControl 权限控制middle
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
			log.Error("checkAuth failed", zap.String("error", err.Error()))
			code = resp.ErrUnauthorized
			goto ERR
		}
		if !pass {
			log.Error("checkAuth not pass ")
			code = resp.ErrUnauthorized
			goto ERR
		}

		c.Next()
		return

	ERR:
		// 解析失败返回错误
		c.Writer.Header().Add("WWW-Authenticate", fmt.Sprintf("Bearer realm='%s'", resp.GetMsg(code)))
		resp.JSON(c, resp.ErrUnauthorized, nil)
		c.Abort()
		return
	}
}
