package user

import (
	"context"
	"crocodile/common/util"
	pbauth "crocodile/service/auth/proto/auth"
	"database/sql"
	"fmt"
	"github.com/labulaka521/logging"
)

type Servicer interface {
	CreateUser(ctx context.Context, user *pbauth.User) (err error)
	ChangeUser(ctx context.Context, user *pbauth.User) (err error)
	GetUser(ctx context.Context, username string) (resp *pbauth.Response, err error)
	LoginUser(ctx context.Context, user *pbauth.User) (resp *pbauth.Response, err error)
	LogoutUser(ctx context.Context, username string) (err error)
	DeleteUser(ctx context.Context, username string) (err error)
}

type Service struct {
	DB *sql.DB
}

// 创建用户
// required
// 		username	    用户名
// 		password	    用户密码
//		email			用户邮箱
//		avatar			用户图像
//		super			管理员
func (s *Service) CreateUser(ctx context.Context, user *pbauth.User) (err error) {
	var (
		stmt           *sql.Stmt
		createuser_sql string
		hashpassword   string
		exists         bool
	)
	createuser_sql = "INSERT INTO crocodile_user (username,hashpassword,email,avatar,super,forbid) VALUE(?,?,?,?,?,?)"

	if hashpassword, err = util.GenerateHashPass(user.Password); err != nil {
		logging.Errorf("GeneratehashPass Err: %v", err)
		return
	}

	if exists, err = s.userExist(ctx, user.Username); err != nil {
		logging.Errorf("Query user Err: %v", err)
		return
	}
	if exists {
		err = fmt.Errorf("User %s Already Exist", user.Username)
		return
	}

	if stmt, err = s.DB.PrepareContext(ctx, createuser_sql); err != nil {
		logging.Errorf("Prepare SQL %s Err:%v", createuser_sql, err)
		return

	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, user.Username, hashpassword, user.Email, user.Avatar, user.Super, user.Forbid); err != nil {
		logging.Errorf("Exec Err: %v", err)
		return

	}
	return
}

// 修改用户信息
// Required
//   username
//	 email
//   avatar
//   forbid
//   super
// Option
//   password
func (s *Service) ChangeUser(ctx context.Context, user *pbauth.User) (err error) {
	var (
		stmt           *sql.Stmt
		changeuser_sql string
		hashpassword   string
		exists         bool
	)

	changeuser_sql = "UPDATE crocodile_user SET email=?,avatar=?,forbid=?,super=? WHERE username=?"

	if exists, err = s.userExist(ctx, user.Username); err != nil {
		return
	}
	if !exists {
		err = fmt.Errorf("User %s Is Not Exits", user.Username)
		return
	}

	// 存在密码时
	if user.Password != "" {
		changeuser_sql = "UPDATE crocodile_user SET hashpassword=?,email=?,avatar=?,forbid=?,super=? WHERE username=?"
		if hashpassword, err = util.GenerateHashPass(user.Password); err != nil {
			logging.Errorf("GeneratehashPass Err: %v", err)
			return
		}
	}

	if stmt, err = s.DB.PrepareContext(ctx, changeuser_sql); err != nil {
		logging.Errorf("Prepare SQL %s Err:%v", changeuser_sql, err)
		return
	}
	defer stmt.Close()
	if user.Password != "" {
		if _, err = stmt.ExecContext(ctx, hashpassword, user.Email, user.Avatar, user.Forbid, user.Super, user.Username); err != nil {
			logging.Errorf("Exec Err: %v", err)
			return
		}
	} else {
		if _, err = stmt.ExecContext(ctx, user.Email, user.Avatar, user.Forbid, user.Super, user.Username); err != nil {
			logging.Errorf("Exec Err: %v", err)
			return
		}
	}

	return
}

// 获取用户信息
// 存在用户名就返回这个用户的信息
// 否则返回所有用户的信息
// Option
//   username
func (s *Service) GetUser(ctx context.Context, username string) (resp *pbauth.Response, err error) {
	var (
		stmt        *sql.Stmt
		getuser_sql string
		rows        *sql.Rows
		users       []*pbauth.User
		args        []interface{}
	)
	resp = &pbauth.Response{}
	args = make([]interface{}, 0)
	getuser_sql = "SELECT id,username, email, avatar, forbid, super FROM crocodile_user"
	if username != "" {
		getuser_sql += " WHERE username=?"
		args = append(args, username)
	}

	if stmt, err = s.DB.PrepareContext(ctx, getuser_sql); err != nil {
		logging.Errorf("Prepare SQL %s Err:%v", getuser_sql, err)
		return
	}
	defer stmt.Close()
	if rows, err = stmt.QueryContext(ctx, args...); err != nil {
		logging.Errorf("SQL %s Query Err: %v", getuser_sql, err)
		return
	}

	for rows.Next() {
		user := pbauth.User{}
		// 这里取出
		if err = rows.Scan(&user.Id, &user.Username, &user.Email, &user.Avatar, &user.Forbid, &user.Super); err != nil {
			logging.Errorf("rows Scan To pbauth.User Err: %v", err)
			continue
		}

		users = append(users, &user)
	}
	resp.Users = users
	logging.Infof("Get User code: %d ", resp.Code)
	return
}

// 用户登录
// Required
//   username
//   password
func (s *Service) LoginUser(ctx context.Context, user *pbauth.User) (resp *pbauth.Response, err error) {
	var (
		hashpassword    string
		gethashpass_sql string
		stmt            *sql.Stmt
		respuser        *pbauth.User
	)
	respuser = &pbauth.User{}
	resp = &pbauth.Response{}
	gethashpass_sql = "SELECT hashpassword,email,forbid,super FROM crocodile_user WHERE username=?"

	if stmt, err = s.DB.PrepareContext(ctx, gethashpass_sql); err != nil {
		logging.Errorf("Prepare SQL %s Err:%v", gethashpass_sql, err)
		return
	}
	defer stmt.Close()
	if err = stmt.QueryRowContext(ctx, user.Username).Scan(&hashpassword, &respuser.Email, &respuser.Forbid, &respuser.Super); err != nil {
		logging.Errorf("Query User %s Err:%v", user.Username, err)
		return
	}

	if respuser.Forbid {
		logging.Infof("User %s Forbid Login", user.Username)
		err = fmt.Errorf("User %s is Forbid Login", user.Username)
		return
	}
	respuser.Username = user.Username

	if err = util.CheckHashPass(hashpassword, user.Password); err != nil {
		logging.Errorf("CheckHashPass Err: %v", err)
		return
	}
	logging.Infof("User %s Login Success", user.Username)
	if resp.Token, err = util.GenerateToken(respuser); err != nil {
		logging.Errorf("GenerateToken Err: %v", err)
		return
	}

	return

}

// 用户注销
// Require
//   username
func (s *Service) LogoutUser(ctx context.Context, username string) (err error) {
	var (
		exists bool
	)
	if exists, err = s.userExist(ctx, username); err != nil {
		return
	}
	if !exists {
		logging.Errorf("User Not Exists")
		err = fmt.Errorf("User %s Is Not Exits", username)
	}
	return
}

// 删除用户
// Required
//   username
func (s *Service) DeleteUser(ctx context.Context, username string) (err error) {
	var (
		deleteuser_sql string
		stmt           *sql.Stmt
		exists         bool
	)

	deleteuser_sql = "DELETE FROM crocodile_user WHERE username=?"
	if exists, err = s.userExist(ctx, username); err != nil {
		return
	}
	if !exists {
		logging.Errorf("User Not Exists")
		err = fmt.Errorf("User %s Is Not Exits", username)
	}
	if stmt, err = s.DB.PrepareContext(ctx, deleteuser_sql); err != nil {
		logging.Errorf("Prepare SQL %s Err:%v", deleteuser_sql, err)
		return
	}
	if _, err = stmt.ExecContext(ctx, username); err != nil {
		logging.Errorf("Exec SQL %s Err: %v", deleteuser_sql, err)
		return
	}
	return
}

// 查询用户是存在
// Required
//  username
func (s *Service) userExist(ctx context.Context, username string) (exists bool, err error) {
	var (
		resp *pbauth.Response
	)
	if resp, err = s.GetUser(ctx, username); err != nil {
		logging.Errorf("GetUser %s Err: %v", username, err)
		return false, err
	}
	if len(resp.Users) == 0 {
		return false, nil

	}
	return true, nil
}
