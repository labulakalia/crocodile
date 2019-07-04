package handler

import (
	"context"
	"crocodile/service/auth/model/user"
	pbauth "crocodile/service/auth/proto/auth"
	"github.com/labulaka521/logging"
)

type Auth struct {
	Service user.Servicer
}

func (auth *Auth) CreateUser(ctx context.Context, req *pbauth.User, resp *pbauth.Response) (err error) {
	var (
		rsp *pbauth.Response
	)
	logging.Debugf("CreateUser %s", req.Username)
	rsp, err = auth.Service.CreateUser(ctx, req)
	resp.Code = rsp.Code
	return
}

func (auth *Auth) ChangeUser(ctx context.Context, req *pbauth.User, resp *pbauth.Response) (err error) {
	var (
		rsp *pbauth.Response
	)
	logging.Debugf("ChangeUser %s", req.Username)
	rsp, err = auth.Service.ChangeUser(ctx, req)
	resp.Code = rsp.Code
	return
}
func (auth *Auth) DeleteUser(ctx context.Context, req *pbauth.User, resp *pbauth.Response) (err error) {
	var (
		rsp *pbauth.Response
	)
	logging.Debugf("ChangeUser %s", req.Username)
	rsp, err = auth.Service.DeleteUser(ctx, req)
	resp.Code = rsp.Code
	return
}
func (auth *Auth) GetUser(ctx context.Context, req *pbauth.User, resp *pbauth.Response) (err error) {
	var (
		rsp *pbauth.Response
	)
	logging.Debug("GetUser", req.Username)
	rsp, err = auth.Service.GetUser(ctx, req)
	resp.Code = rsp.Code
	resp.Users = rsp.Users

	return
}
func (auth *Auth) LoginUser(ctx context.Context, req *pbauth.User, resp *pbauth.Response) (err error) {
	var (
		rsp *pbauth.Response
	)
	logging.Debugf("LoginUser %s", req.Username)
	rsp, err = auth.Service.LoginUser(ctx, req)
	resp.Code = rsp.Code
	resp.Token = rsp.Token
	return
}
func (auth *Auth) LogoutUser(ctx context.Context, req *pbauth.User, resp *pbauth.Response) (err error) {
	var (
		rsp *pbauth.Response
	)
	logging.Debugf("ChangeUser %s", req.Username)
	rsp, err = auth.Service.LogoutUser(ctx, req)
	resp.Code = rsp.Code
	logging.Infof("Logout User %s Code: %d", req.Username, rsp.Code)
	return
}
