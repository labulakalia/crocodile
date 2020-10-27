package model

import (
	"context"
	"testing"

	"github.com/labulaka521/crocodile/common/utils"
	"github.com/labulaka521/crocodile/core/utils/define"
)

func TestAddUserv2(t *testing.T) {
	gormdb = InitGormSqlite()

	password := "password"
	hashpassword, err := utils.GenerateHashPass(password)
	if err != nil {
		t.Fatalf("gener hash pass failed %v", err)
	}
	type args struct {
		ctx          context.Context
		name         string
		hashpassword string
		role         define.Role
		remark       string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "create user",
			args: args{
				ctx:          context.Background(),
				name:         "user1",
				hashpassword: hashpassword,
				role:         define.AdminUser,
				remark:       "remark",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddUserv2(tt.args.ctx, tt.args.name, tt.args.hashpassword, tt.args.role, tt.args.remark); (err != nil) != tt.wantErr {
				t.Errorf("AddUserv2() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoginUserv2(t *testing.T) {
	t.Run("create user", TestAddUserv2)

	type args struct {
		ctx      context.Context
		name     string
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "login user",
			args: args{
				ctx:      context.Background(),
				name:     "user1",
				password: "password",
			},
			wantErr: false,
		},
		{
			name: "failed login user",
			args: args{
				ctx:      context.Background(),
				name:     "user1",
				password: "error_passwd",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := LoginUserv2(tt.args.ctx, tt.args.name, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoginUserv2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestGetUserByIDv2(t *testing.T) {
	t.Run("create user", TestAddUserv2)
	users, count, err := GetUsersv2(context.Background(), []string{}, 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("users count want 1, but get count %d", count)
	}
	user := users[0]
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "get user",
			args: args{
				ctx: context.Background(),
				id:  user.ID,
			},
			wantErr: false,
		},
		{
			name: "get exist user",
			args: args{
				ctx: context.Background(),
				id:  "not_exist_user",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserByIDv2(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByIDv2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if got.Name != "user1" {
					t.Fatalf("get user failed want username user1, but get username %s", got.Name)
				}
			}
		})
	}
}

func TestGetUserByNamev2(t *testing.T) {
	t.Run("create user", TestAddUserv2)
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "get user",
			args: args{
				ctx:  context.Background(),
				name: "user1",
			},
			wantErr: false,
		},
		{
			name: "get exist user",
			args: args{
				ctx:  context.Background(),
				name: "not_exist_user",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserByNamev2(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByNamev2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if got.Name != "user1" {
					t.Fatalf("get user failed want username user1, but get username %s", got.Name)
				}
			}
		})
	}
}

func TestAdminChangeUserv2(t *testing.T) {
	t.Run("create user", TestAddUserv2)
	users, count, err := GetUsersv2(context.Background(), []string{}, 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("users count want 1, but get count %d", count)
	}
	user := users[0]

	type args struct {
		ctx          context.Context
		id           string
		role         define.Role
		forbid       bool
		hashpassword string
		remark       string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "admin change user",
			args: args{
				ctx:    context.Background(),
				id:     user.ID,
				role:   define.AdminUser,
				remark: "change remark",
			},
			wantErr: false,
		},
		{
			name: "admin change user",
			args: args{
				ctx:    context.Background(),
				id:     user.ID,
				role:   define.AdminUser,
				forbid: true,
				remark: "change remark",
			},
			wantErr: false,
		},
		{
			name: "admin change user",
			args: args{
				ctx:    context.Background(),
				id:     "not_exist_user",
				role:   define.AdminUser,
				forbid: true,
				remark: "change remark",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AdminChangeUserv2(tt.args.ctx, tt.args.id, tt.args.role, tt.args.forbid, tt.args.hashpassword, tt.args.remark); (err != nil) != tt.wantErr {
				t.Errorf("AdminChangeUserv2() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	_, err = LoginUserv2(context.Background(), user.Name, "password")
	switch err.(type) {
	case define.ErrForbid:
		t.Log("forbid")
	default:
		t.Fatal("should get error type ErrForbid")
	}
}

func TestChangeUserInfov2(t *testing.T) {
	t.Run("create user", TestAddUserv2)
	users, count, err := GetUsersv2(context.Background(), []string{}, 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("users count want 1, but get count %d", count)
	}
	user := users[0]
	type args struct {
		ctx          context.Context
		id           string
		email        string
		wechat       string
		wechatbot    string
		dingding     string
		telegram     string
		hashpassword string
		alarmTmpl    string
		remark       string
		env          Env
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "change user",
			args: args{
				ctx:   context.Background(),
				id:    user.ID,
				email: "test@email.com",
			},
			wantErr: false,
		},
		{
			name: "change not exist user",
			args: args{
				ctx:   context.Background(),
				id:    "not exist",
				email: "test@email.com",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ChangeUserInfov2(tt.args.ctx, tt.args.id, tt.args.email, tt.args.wechat, tt.args.wechatbot, tt.args.dingding, tt.args.telegram, tt.args.hashpassword, tt.args.alarmTmpl, tt.args.remark, tt.args.env); (err != nil) != tt.wantErr {
				t.Errorf("ChangeUserInfov2() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
