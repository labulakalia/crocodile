package user

import (
	"context"
	"crocodile/common/cfg"
	"crocodile/common/db/mysql"
	pbauth "crocodile/service/auth/proto/auth"
	"github.com/labulaka521/logging"
	"testing"
)

func TestService_User(t *testing.T) {
	var (
		err  error
		resp *pbauth.Response
	)
	logging.SetLogLevel("FATAL")
	logging.Setup()

	cfg.Init()
	var _ Servicer = &Service{}

	s := Service{
		DB: mysql.New(cfg.MysqlConfig.DSN, cfg.MysqlConfig.MaxIdleConnection, cfg.MysqlConfig.MaxIdleConnection),
	}
	testu := &pbauth.User{
		Username: "testusername",
		Password: "testpassword",
		Email:    "testemail@email.com",
		Avatar:   "http://avatar.com",
		Super:    true,
	}

	defer func() {
		if resp, err = s.DeleteUser(context.Background(), testu); err != nil {
			t.Fatalf("Delete User %s Err:%v", testu.Username, err)
		}
		if resp.Code != 200 {
			t.Fatalf("Delete User Err Code: %d", resp.Code)
		}
	}()

	if resp, err = s.CreateUser(context.Background(), testu); err != nil {
		t.Fatalf("Create User Err: %v", err)
	}
	if resp.Code != 200 {
		t.Fatalf("Create User Err Code: %d ", resp.Code)
	}
	if resp, err = s.GetUser(context.Background(), testu); err != nil {
		t.Fatalf("Get User Err: %v", err)
	}
	if resp.Code != 200 {
		t.Fatalf("Get User Err Code: %d", resp.Code)
	}
	if len(resp.Users) != 1 {
		t.Fatalf("No Get User %s", testu.Username)
	}

	if resp, err = s.LoginUser(context.Background(), testu); err != nil {
		t.Fatalf("Login  User %s Err: %v", testu.Username, err)
	}
	if resp.Code != 200 {
		t.Fatalf("Login User Err Code: %d Msg", resp.Code)
	}
	t.Logf("Token: %s", resp.Token)

	testu.Password = "changepass"

	if resp, err = s.ChangeUser(context.Background(), testu); err != nil {
		t.Fatalf("Change  User Err: %v", err)
	}
	if resp.Code != 200 {
		t.Fatalf("ChangeUser User Err Code: %d", resp.Code)
	}

	if resp, err = s.LoginUser(context.Background(), testu); err != nil {
		t.Fatalf("Login  User %s Err: %v", testu.Username, err)
	}
	if resp.Code != 200 {
		t.Fatalf("Login User Err Code: %d", resp.Code)
	}

	if resp, err = s.LogoutUser(context.Background(), testu); err != nil {
		t.Fatalf("Login  User %s Err: %v", testu.Username, err)
	}
	if resp.Code != 200 {
		t.Fatalf("Logout User Err Code: %d", resp.Code)
	}
}
