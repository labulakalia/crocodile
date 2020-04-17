package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/common/utils"
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/model"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/labulaka521/crocodile/core/utils/resp"
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
	// TODO only admin

	hashpassword, err = utils.GenerateHashPass(ruser.Password)
	if err != nil {
		log.Error("GenerateHashPass failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}

	exist, err := model.Check(ctx, model.TBUser, model.Name, ruser.Name)
	if err != nil {
		log.Error("IsExist failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	if exist {
		resp.JSON(c, resp.ErrUserNameExist, nil)
		return
	}

	err = model.AddUser(ctx, ruser.Name, hashpassword, ruser.Role)
	if err != nil {
		log.Error("AddUser failed", zap.Error(err))
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
	fmt.Println(uid)
	// check uid exist
	exist, err := model.Check(ctx, model.TBUser, model.ID, uid)
	if err != nil {
		log.Error("IsExist failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	if !exist {
		resp.JSON(c, resp.ErrUserNotExist, nil)
		return
	}

	user, err := model.GetUserByID(ctx, uid)
	if err != nil {
		log.Error("GetUserByID failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	user.Password = ""
	if user.Role == 2 {
		user.Roles = []string{"admin"}
	} else {
		user.Roles = []string{}
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
	// TODO only admin

	err = c.BindQuery(&q)
	if err != nil {
		log.Error("BindQuery offset failed", zap.Error(err))
	}

	if q.Limit == 0 {
		q.Limit = define.DefaultLimit
	}
	users, count, err := model.GetUsers(ctx, nil, q.Offset, q.Limit)
	if err != nil {
		log.Error("GetUsers failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	// remove password
	for i, user := range users {
		user.Password = ""
		users[i] = user
	}

	resp.JSON(c, resp.Success, users, count)
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
	if len(newinfo.Password) > 0 && len(newinfo.Password) < 8 {
		log.Error("password is short 8")
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	uid := c.GetString("uid")
	if uid != newinfo.ID {
		log.Error("uid is error", zap.String("uid", uid), zap.String("infoid", newinfo.ID))
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	err = model.ChangeUserInfo(ctx,
		uid,
		newinfo.Email,
		newinfo.WeChat,
		newinfo.DingPhone,
		newinfo.Telegram,
		newinfo.Password,
		newinfo.Remark)
	if err != nil {
		log.Error("ChangeUserInfo failed", zap.Error(err))
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
	if len(user.Password) > 0 && len(user.Password) < 8 {
		log.Error("password is short 8")
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	// TODO only admin
	exist, err := model.Check(ctx, model.TBUser, model.ID, user.ID)
	if err != nil {
		log.Error("IsExist failed", zap.Error(err))
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
		log.Error("AdminChangeUser failed", zap.Error(err))
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
	if err != nil {
		log.Error("model.LoginUser failed", zap.Error(err))
	}
	switch err := errors.Unwrap(err); err.(type) {
	case nil:
		resp.JSON(c, resp.Success, token)
	case define.ErrUserPass:
		resp.JSON(c, resp.ErrUserPassword, nil)
	case define.ErrForbid:
		resp.JSON(c, resp.ErrUserForbid, nil)
	default:
		resp.JSON(c, resp.ErrInternalServer, nil)
	}
}

// LogoutUser logout user
// @Summary logout user
// @Tags User
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/user/logout [post]
// @Security BasicAuth
func LogoutUser(c *gin.Context) {
	resp.JSON(c, resp.Success, nil)
}

// GetSelect return name,id
// @Summary return name,id
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/user/select [post]
// @Security BasicAuth
func GetSelect(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	data, err := model.GetNameID(ctx, model.TBUser)
	if err != nil {
		log.Error("model.GetNameID failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	resp.JSON(c, resp.Success, data)
}

// GetAlarmStatus return enable alarm notify
func GetAlarmStatus(c *gin.Context) {
	type NotifyStatus struct {
		Email    bool `json:"email"`
		DingDing bool `json:"dingphone"`
		Slack    bool `json:"slack"`
		Telegram bool `json:"telegram"`
		WeChat   bool `json:"wechat"`
		WebHook  bool `json:"wehook"`
	}
	notifycfg := config.CoreConf.Notify
	notifystatus := NotifyStatus{
		Email:    notifycfg.Email.Enable,
		DingDing: notifycfg.DingDing.Enable,
		Slack:    notifycfg.Slack.Enable,
		Telegram: notifycfg.Telegram.Enable,
		WeChat:   notifycfg.WeChat.Enable,
		WebHook:  notifycfg.WebHook.Enable,
	}
	resp.JSON(c, resp.Success, notifystatus)
}

// GetOperateLog get user operate log
func GetOperateLog(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	type queryparams struct {
		define.Query
		UserName string `form:"username"`
		Method   string `form:"method"`
		Module   string `form:"module"`
	}

	q := queryparams{}
	err := c.ShouldBindQuery(&q)
	if err != nil {
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	if q.Limit == 0 {
		q.Limit = define.DefaultLimit
	}

	// uid, method, module, limit, offset
	oplogs, count, err := model.GetOperate(ctx, "", q.UserName, q.Method, q.Module, q.Limit, q.Offset)
	if err != nil {
		log.Error("model.GetOperate filed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	resp.JSON(c, resp.Success, oplogs, count)
}
