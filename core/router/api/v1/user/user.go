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
)

// RegistryUser registry new user
// POST /api/v1/user
// @params
// name
// password
// email
// role	 option
// remark option
func RegistryUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	user := define.User{}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		log.Error("ShouldBindJSON failed", zap.Error(err))
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	if user.Name == "" {
		log.Error("User.Name is empty")
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	user.Password, err = utils.GenerateHashPass(user.Password)
	if err != nil {
		log.Error("GenerateHashPass failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrBadRequest, user)
		return
	}
	user.ID = utils.GetID()

	exist, err := model.Check(ctx, model.TBUser, model.Name, user.Name)
	if err != nil {
		log.Error("IsExist failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	if exist {
		resp.JSON(c, resp.ErrUserNameExist, nil)
		return
	}
	exist, err = model.Check(ctx, model.TBUser, model.Email, user.Email)
	if err != nil {
		log.Error("IsExist failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	if exist {
		resp.JSON(c, resp.ErrEmailExist, nil)
		return
	}

	err = model.AddUser(ctx, &user)
	if err != nil {
		log.Error("AddUser failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}

	resp.JSON(c, resp.Success, nil)
}

// GetUser get user
// GET /api/v1/user
// @params
// 通过解析token的ID获取请求者的信息
func GetUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	uid := c.GetString("uid")

	log.Info("Check Uid " + uid)
	exist, err := model.Check(ctx, model.TBUser, model.ID, uid)
	if err != nil {
		log.Error("IsExist failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
	}
	if !exist {
		resp.JSON(c, resp.ErrUserNotExist, nil)
		return
	}

	users, err := model.GetUser(ctx, uid)
	if err != nil {
		log.Error("GetUser failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	if len(users) != 1 {
		log.Error("Get user failed", zap.String("uid", uid))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	resp.JSON(c, resp.Success, users[0])
}

// GetUsers get all users
// @params
// 通过解析token的ID获取请求者的信息
func GetUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	users, err := model.GetUser(ctx, "")
	if err != nil {
		log.Error("GetUsers failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	resp.JSON(c, resp.Success, users)
}

// ChangeUser change user message
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
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	user := define.User{}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		log.Error("ShouldBindJSON failed", zap.Error(err))
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}

	exist, err := model.Check(ctx, model.TBUser, model.ID, user.ID)
	if err != nil {
		log.Error("IsExist failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
	}
	if !exist {
		resp.JSON(c, resp.ErrUserNotExist, nil)
		return
	}
	var role define.Role
	if v, ok := c.Get("role"); ok {
		role = v.(define.Role)
	}

	err = model.ChangeUser(ctx, &user, role)
	if err != nil {
		log.Error("ChangeUser failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}

	resp.JSON(c, resp.Success, nil)
}

// LoginUser login user
func LoginUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	username, password, ok := c.Request.BasicAuth()
	if !ok {
		resp.JSON(c, resp.ErrBadRequest, nil)
	}
	token, err := model.LoginUser(ctx, username, password)

	switch err := errors.Cause(err).(type) {
	case nil:
		resp.JSON(c, resp.Success, token)
	case define.ErrUserPass:
		resp.JSON(c, resp.ErrUserPassword, nil)
	case define.ErrForbid:
		resp.JSON(c, resp.ErrUserForbid, nil)
	default:
		resp.JSON(c, resp.ErrInternalServer, nil)
		log.Info("LoginUser", zap.String("error", err.Error()))
	}
}
