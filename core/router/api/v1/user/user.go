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

// RegistryUser new user
// @Summary registry new user
// @Tags User
// @Produce json
// @Param Registry body define.RegistryUser true "registry user"
// @Success 200 {object} resp.Response
// @Router /api/v1/user/registry [post]
// @Security ApiKeyAuth
func RegistryUser(c *gin.Context) {
	var (
		hashpassword string
	)
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	ruser := define.RegistryUser{}
	err := c.ShouldBindJSON(&ruser)
	if err != nil {
		log.Error("ShouldBindJSON failed", zap.Error(err))
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}

	hashpassword, err = utils.GenerateHashPass(ruser.Password)
	if err != nil {
		log.Error("GenerateHashPass failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}

	exist, err := model.Check(ctx, model.TBUser, model.Name, ruser.Name)
	if err != nil {
		log.Error("IsExist failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	if exist {
		resp.JSON(c, resp.ErrUserNameExist, nil)
		return
	}

	err = model.AddUser(ctx, ruser.Name, hashpassword, ruser.Role)
	if err != nil {
		log.Error("AddUser failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}

	resp.JSON(c, resp.Success, nil)
}

// GetUser Get User Info By Token
// @Summary get user info by token
// @Tags User
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/user/info [get]
// @Security ApiKeyAuth
func GetUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	uid := c.GetString("uid")

	// check uid exist
	exist, err := model.Check(ctx, model.TBUser, model.ID, uid)
	if err != nil {
		log.Error("IsExist failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
	}
	if !exist {
		resp.JSON(c, resp.ErrUserNotExist, nil)
		return
	}

	user, err := model.GetUser(ctx, uid)
	if err != nil {
		log.Error("GetUser failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	user.Password = ""
	if user.Role == 2 {
		user.Roles = []string{"admin"}
	}
	resp.JSON(c, resp.Success, user)
}

// GetUsers get user info by token
// @Summary  get all users info
// @Tags User
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/user/all [get]
// @Security ApiKeyAuth
func GetUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	var (
		q   define.Query
		err error
	)

	err = c.BindQuery(&q)
	if err != nil {
		log.Error("BindQuery offset failed", zap.Error(err))
	}

	if q.Limit == 0 {
		q.Limit = define.DefaultLimit
	}
	users, err := model.GetUsers(ctx, nil, q.Offset, q.Limit)
	if err != nil {
		log.Error("GetUsers failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	// remove password
	for i, user := range users {
		user.Password = ""
		users[i] = user
	}

	resp.JSON(c, resp.Success, users)
}

// ChangeUserInfo change user self config
// @Summary user change self's config info
// @Tags User
// @Description change self config,like email,wechat,dingphone,slack,telegram,password,remark
// @Produce json
// @Param User body define.ChangeUserSelf true "Change Self User Info"
// @Success 200 {object} resp.Response
// @Router /api/v1/user/info [put]
// @Security ApiKeyAuth
func ChangeUserInfo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	newinfo := define.ChangeUserSelf{}
	err := c.ShouldBindJSON(&newinfo)
	if err != nil {
		log.Error("ShouldBindJSON failed", zap.Error(err))
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	uid := c.GetString("uid")
	err = model.ChangeUserInfo(ctx,
		uid,
		newinfo.Email,
		newinfo.WeChat,
		newinfo.DingPhone,
		newinfo.Slack,
		newinfo.Telegram,
		newinfo.Password,
		newinfo.Remark)
	if err != nil {
		log.Error("ChangeUserInfo failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}

	resp.JSON(c, resp.Success, nil)
}

// AdminChangeUser will change role,forbid,password,Remark
// @Summary admin change user info
// @Tags User
// @Description admin change user's role,forbid,password,remark
// @Produce json
// @Param User body define.AdminChangeUser true "Admin Change User"
// @Success 200 {object} resp.Response
// @Router /api/v1/user/admin [put]
// @Security ApiKeyAuth
func AdminChangeUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	user := define.AdminChangeUser{}
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
	if role != define.AdminUser {
		resp.JSON(c, resp.ErrUnauthorized, nil)
		return
	}

	err = model.AdminChangeUser(ctx, user.ID, user.Role, user.Forbid, user.Password, user.Remark)
	if err != nil {
		log.Error("AdminChangeUser failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}

	resp.JSON(c, resp.Success, nil)
}

// LoginUser login user
// @Summary login user
// @Tags User
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/user/login [post]
// @Security BasicAuth
func LoginUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	username, password, ok := c.Request.BasicAuth()
	if !ok {
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
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

// GetSelect return name,id
func GetSelect(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	data, err := model.GetNameID(ctx, model.TBUser)
	if err != nil {
		log.Error("model.GetNameID", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	resp.JSON(c, resp.Success, data)
}
