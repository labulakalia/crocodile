package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/common/utils"
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/model"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/labulaka521/crocodile/core/utils/resp"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

// @params
// name
// password
// email
// role	 option
// remark option
func RegistryUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.CoreConf.Db.MaxQueryTime)*time.Second)
	defer cancel()
	user := define.User{}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		log.Error("ShouldBindJSON failed", zap.Error(err))
		resp.Json(c, resp.ErrBadRequest, nil)
		return
	}
	user.Password, err = utils.GenerateHashPass(user.Password)
	if err != nil {
		log.Error("GenerateHashPass failed", zap.String("error", err.Error()))
		resp.Json(c, resp.ErrBadRequest, user)
		return
	}
	user.Id = utils.GetId()

	exist, err := model.IsExist(ctx, model.UserName, user.Name)
	if err != nil {
		log.Error("IsExist failed", zap.String("error", err.Error()))
		resp.Json(c, resp.ErrInternalServer, nil)
		return
	}
	if exist {
		resp.Json(c, resp.ErrUserNameExist, nil)
		return
	}
	exist, err = model.IsExist(ctx, model.Email, user.Email)
	if err != nil {
		log.Error("IsExist failed", zap.String("error", err.Error()))
		resp.Json(c, resp.ErrInternalServer, nil)
		return
	}
	if exist {
		resp.Json(c, resp.ErrEmailExist, nil)
		return
	}

	err = model.AddUser(ctx, &user)
	if err != nil {
		log.Error("AddUser failed", zap.String("error", err.Error()))
		resp.Json(c, resp.ErrInternalServer, nil)
		return
	}
	resp.Json(c, resp.Success, nil)
}

// @params
// 通过解析token的ID获取请求者的信息
func GetUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.CoreConf.Db.MaxQueryTime)*time.Second)
	defer cancel()

	uid := c.GetInt64("uid")

	exist, err := model.IsExist(ctx, model.Uid, uid)
	if err != nil {
		log.Error("IsExist failed", zap.String("error", err.Error()))
		resp.Json(c, resp.ErrInternalServer, nil)
	}
	if !exist {
		resp.Json(c, resp.ErrUserNotExist, nil)
		return
	}

	users, err := model.GetUser(ctx, uid)
	if err != nil {
		log.Error("GetUser failed", zap.String("error", err.Error()))
		resp.Json(c, resp.ErrInternalServer, nil)
		return
	}
	if len(users) != 1 {
		log.Error("Get many users", zap.Int64("uid", uid))
		resp.Json(c, resp.ErrBadRequest, nil)
		return
	}
	resp.Json(c, resp.Success, users[0])
}

// @params
// 通过解析token的ID获取请求者的信息
func GetUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.CoreConf.Db.MaxQueryTime)*time.Second)
	defer cancel()

	users, err := model.GetUser(ctx, 0)
	if err != nil {
		log.Error("GetUsers failed", zap.String("error", err.Error()))
		resp.Json(c, resp.ErrInternalServer, nil)
		return
	}
	resp.Json(c, resp.Success, users)
}

// @params
// id	required
// name required
// remark
// super
// email
// role
// forbid
func ChangeUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.CoreConf.Db.MaxQueryTime)*time.Second)
	defer cancel()

	user := define.User{}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		log.Error("ShouldBindJSON failed", zap.Error(err))
		resp.Json(c, resp.ErrBadRequest, nil)
		return
	}

	exist, err := model.IsExist(ctx, model.Uid, user.Id)
	if err != nil {
		log.Error("IsExist failed", zap.String("error", err.Error()))
		resp.Json(c, resp.ErrInternalServer, nil)
	}
	if !exist {
		resp.Json(c, resp.ErrUserNotExist, nil)
		return
	}

	err = model.ChangeUser(ctx, &user)
	if err != nil {
		log.Error("ChangeUser failed", zap.String("error", err.Error()))
		resp.Json(c, resp.ErrInternalServer, nil)
		return
	}
	resp.Json(c, resp.Success, nil)
}

// @params
func LoginUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.CoreConf.Db.MaxQueryTime)*time.Second)
	defer cancel()
	username, password, ok := c.Request.BasicAuth()
	if !ok {
		resp.Json(c, resp.ErrBadRequest, nil)
	}
	token, err := model.LoginUser(ctx, username, password)

	switch err := errors.Cause(err).(type) {
	case nil:
		resp.Json(c, resp.Success, token)
	case define.ErrUserPass:
		resp.Json(c, resp.ErrUserPassword, nil)
	case define.ErrForbid:
		resp.Json(c, resp.ErrUserForbid, nil)
	default:
		resp.Json(c, resp.ErrInternalServer, nil)
		log.Info("LoginUser", zap.String("error", err.Error()))
	}
}
