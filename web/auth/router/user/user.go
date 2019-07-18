package user

import (
	"context"
	"crocodile/common/bind"
	"crocodile/common/cfg"
	"crocodile/common/e"
	"crocodile/common/middle"
	"crocodile/common/registry"
	"crocodile/common/response"
	"crocodile/common/util"
	"crocodile/common/wrapper"
	pbauth "crocodile/service/auth/proto/auth"

	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/labulaka521/logging"
	"github.com/micro/go-micro/client"
	"strings"
	"time"
)

var (
	AuthClient pbauth.AuthService
)

func Init() {
	c := client.NewClient(
		client.Retries(3),
		client.Registry(registry.Etcd(cfg.EtcdConfig.Endpoints...)),
	)

	AuthClient = pbauth.NewAuthService("crocodile.srv.auth", c)
}

type User struct {
	Id       uint32   `json:"id"`
	UserName string   `json:"username" validate:"required"`
	PassWord string   `json:"password,omitempty"`
	Email    string   `json:"email" validate:"required"`
	Super    bool     `json:"super"`
	Forbid   bool     `json:"forbid"`
	Avatar   string   `json:"avatar" validate:"required"`
	Roles    []string `json:"roles"`
}

func GetUser(c *gin.Context) {
	var (
		app         response.Gin
		ctx         context.Context
		err         error
		code        int32
		reqUser     *pbauth.User
		respAuthSrv *pbauth.Response
		exists      bool
		loginuser   string
		tmpuser     *User
	)
	ctx, ok := wrapper.ContextWithSpan(c)
	if ok == false {
		logging.Error("get context err")
		ctx = context.Background()
	}
	tmpuser = &User{}
	app = response.Gin{c}
	ctx, _ = context.WithTimeout(ctx, time.Duration(cfg.MysqlConfig.MaxQueryTime)*time.Second)
	code = e.SUCCESS
	if loginuser, exists = c.Keys["user"].(string); !exists {
		code = e.ERR_TOKEN_INVALID
		app.Response(code, nil)
		return
	}

	reqUser = &pbauth.User{Username: loginuser}
	if respAuthSrv, err = AuthClient.GetUser(ctx, reqUser); err != nil {
		logging.Errorf("GetUser Err:%v", err)
		code = e.ERR_GET_USER_FAIL
		app.Response(code, nil)
		return
	}

	tmpuser.UserName = respAuthSrv.Users[0].Username
	tmpuser.Email = respAuthSrv.Users[0].Email
	tmpuser.Avatar = respAuthSrv.Users[0].Avatar
	tmpuser.Forbid = respAuthSrv.Users[0].Forbid
	tmpuser.Super = respAuthSrv.Users[0].Super
	if tmpuser.Super {
		tmpuser.Roles = []string{"admin"}
	} else {
		tmpuser.Roles = []string{}
	}
	app.Response(code, tmpuser)
}

func GetUsers(c *gin.Context) {
	var (
		app         response.Gin
		ctx         context.Context
		err         error
		code        int32
		reqUser     *pbauth.User
		respAuthSrv *pbauth.Response
		us          []*User
	)
	ctx, ok := wrapper.ContextWithSpan(c)
	if ok == false {
		logging.Error("get context err")
		ctx = context.Background()
	}
	us = []*User{}
	app = response.Gin{c}
	code = e.SUCCESS
	ctx, _ = context.WithTimeout(ctx, time.Duration(cfg.MysqlConfig.MaxQueryTime)*time.Second)

	reqUser = &pbauth.User{}
	if respAuthSrv, err = AuthClient.GetUser(ctx, reqUser); err != nil {
		logging.Errorf("GetUser Err:%v", err)
		code = e.ERR_GET_USER_FAIL
		app.Response(code, nil)
		return
	}
	for _, rep := range respAuthSrv.Users {
		tmpuser := User{}
		tmpuser.Id = rep.Id
		tmpuser.UserName = rep.Username
		tmpuser.Email = rep.Email
		tmpuser.Avatar = rep.Avatar
		tmpuser.Forbid = rep.Forbid
		tmpuser.Super = rep.Super
		if tmpuser.Super {
			tmpuser.Roles = []string{"admin"}
		} else {
			tmpuser.Roles = []string{}
		}
		us = append(us, &tmpuser)
	}
	logging.Infof("GetUsers Response  Code: %d", respAuthSrv.Code)
	app.Response(code, us)
}

func ChangeUser(c *gin.Context) {
	var (
		app     response.Gin
		u       User
		ctx     context.Context
		err     error
		code    int32
		reqUser *pbauth.User
	)
	ctx, ok := wrapper.ContextWithSpan(c)
	if ok == false {
		logging.Error("get context err")
		ctx = context.Background()
	}
	app = response.Gin{c}
	u = User{}
	ctx, _ = context.WithTimeout(ctx, time.Duration(cfg.MysqlConfig.MaxQueryTime)*time.Second)
	code = e.SUCCESS
	if err = bind.BindJson(c, &u); err != nil {
		logging.Errorf("BindJson Err:%v", err)
		code = e.ERR_BAD_REQUEST
		app.Response(code, nil)
		return
	}
	reqUser = &pbauth.User{
		Username: u.UserName,
		Password: u.PassWord,
		Email:    u.Email,
		Avatar:   u.Avatar,
		Forbid:   u.Forbid,
		Super:    u.Super,
	}
	if _, err = AuthClient.ChangeUser(ctx, reqUser); err != nil {
		logging.Errorf("ChangeUser Err:%v", err)
		code = e.ERR_CHANGE_USER_FAIL
		app.Response(code, nil)
		return
	}

	app.Response(code, nil)
}

// 创建用户
type create struct {
	UserName string `json:"username" validate:"required"`
	PassWord string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Super    bool   `json:"super"`
}

func UserCreate(c *gin.Context) {
	var (
		app     response.Gin
		u       create
		ctx     context.Context
		err     error
		code    int32
		reqUser *pbauth.User
		avatar  string
	)
	ctx, ok := wrapper.ContextWithSpan(c)
	if ok == false {
		logging.Error("get context err")
		ctx = context.Background()
	}
	app = response.Gin{c}
	u = create{}
	code = e.SUCCESS
	ctx, _ = context.WithTimeout(ctx, time.Duration(cfg.MysqlConfig.MaxQueryTime)*time.Second)

	avatar, _ = util.GenerateAvatar(u.Email, 128)

	if err = bind.BindJson(c, &u); err != nil {
		code = e.ERR_BAD_REQUEST
		app.Response(code, nil)
		return
	}

	reqUser = &pbauth.User{
		Username: u.UserName,
		Password: u.PassWord,
		Email:    u.Email,
		Avatar:   avatar,
		Forbid:   false,
		Super:    u.Super,
	}
	if _, err = AuthClient.CreateUser(ctx, reqUser); err != nil {
		logging.Errorf("CreateUser Err:%v", err)
		code = e.ERR_CREATE_USER_FAIL
		app.Response(code, nil)
		return
	}

	app.Response(code, nil)
}

func UserLogin(c *gin.Context) {
	// Authorization
	var (
		authorization string
		app           response.Gin
		data          []string
		code          int32
		body          []byte
		userpass      []string
		err           error
		ctx           context.Context
		reqUser       *pbauth.User
		resp          *pbauth.Response
	)
	ctx, ok := wrapper.ContextWithSpan(c)
	if ok == false {
		logging.Error("get context err")
		ctx = context.Background()
	}

	ctx, _ = context.WithTimeout(ctx, time.Duration(cfg.MysqlConfig.MaxQueryTime)*time.Second)

	app = response.Gin{c}
	code = e.SUCCESS
	if authorization, err = middle.GetAuthor(c); err != nil {
		code = e.ERR_USER_PASS_FAIL
		return
	}

	data = strings.Split(authorization, " ")
	if len(data) != 2 {
		code = e.ERR_USER_PASS_FAIL
		app.Response(code, nil)
		return
	}
	if data[0] != "Basic" {
		code = e.ERR_USER_PASS_FAIL
		app.Response(code, nil)
		return
	}

	if body, err = base64.StdEncoding.DecodeString(data[1]); err != nil {
		code = e.ERR_USER_PASS_FAIL
		app.Response(code, nil)
		return
	}
	userpass = strings.Split(string(body), ":")
	if len(userpass) != 2 {
		code = e.ERR_USER_PASS_FAIL
		app.Response(code, nil)
		return
	}

	reqUser = &pbauth.User{
		Username: userpass[0],
		Password: userpass[1],
	}

	if resp, err = AuthClient.LoginUser(ctx, reqUser); err != nil {
		logging.Errorf("Login User Err:%v", err)
		code = e.ERR_LOGIN_USER_FAIL
		app.Response(code, nil)
		return
	}
	app.Response(code, resp.Token)
	return
}

func Logout(c *gin.Context) {
	var (
		app     response.Gin
		ctx     context.Context
		err     error
		code    int32
		reqUser *pbauth.User

		exists    bool
		loginuser string
	)
	ctx, ok := wrapper.ContextWithSpan(c)
	if ok == false {
		logging.Error("get context err")
		ctx = context.Background()
	}
	app = response.Gin{c}
	ctx, _ = context.WithTimeout(ctx, time.Duration(cfg.MysqlConfig.MaxQueryTime)*time.Second)
	code = e.SUCCESS
	if loginuser, exists = c.Keys["user"].(string); !exists {
		code = e.ERR_TOKEN_INVALID
		app.Response(code, nil)
		return
	}
	reqUser = &pbauth.User{Username: loginuser}
	if _, err = AuthClient.LogoutUser(ctx, reqUser); err != nil {
		logging.Errorf("Logout User Err:%v", err)
		code = e.ERR_LOGOUT_USER_FAIL
		app.Response(code, nil)
		return
	}
	app.Response(code, nil)
}

func DeleteUser(c *gin.Context) {
	var (
		app     response.Gin
		ctx     context.Context
		err     error
		code    int32
		reqUser *pbauth.User

		exists    bool
		loginuser string
	)
	ctx, ok := wrapper.ContextWithSpan(c)
	if ok == false {
		logging.Error("get context err")
		ctx = context.Background()
	}
	app = response.Gin{c}
	ctx, _ = context.WithTimeout(ctx, time.Duration(cfg.MysqlConfig.MaxQueryTime)*time.Second)
	code = e.SUCCESS
	if loginuser, exists = c.Keys["user"].(string); !exists {
		code = e.ERR_TOKEN_INVALID
		app.Response(code, nil)
		return
	}

	reqUser = &pbauth.User{Username: loginuser}
	if _, err = AuthClient.DeleteUser(ctx, reqUser); err != nil {
		logging.Errorf("DeleteUser Err:%v", err)
		code = e.ERR_DELETE_USER_FAIL
		app.Response(code, nil)
		return
	}

	app.Response(code, nil)
}
