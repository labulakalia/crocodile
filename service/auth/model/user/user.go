package user

import (
	"context"
	"crocodile/common/e"
	"crocodile/common/util"
	pbauth "crocodile/service/auth/proto/auth"
	"database/sql"
	"errors"
	"github.com/labulaka521/logging"
)

type Servicer interface {
	CreateUser(ctx context.Context, user *pbauth.User) (resp *pbauth.Response, err error)
	ChangeUser(ctx context.Context, user *pbauth.User) (resp *pbauth.Response, err error)
	GetUser(ctx context.Context, user *pbauth.User) (resp *pbauth.Response, err error)
	LoginUser(ctx context.Context, user *pbauth.User) (resp *pbauth.Response, err error)
	LogoutUser(ctx context.Context, user *pbauth.User) (resp *pbauth.Response, err error)
	DeleteUser(ctx context.Context, user *pbauth.User) (resp *pbauth.Response, err error)
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
func (s *Service) CreateUser(ctx context.Context, user *pbauth.User) (resp *pbauth.Response, err error) {
	var (
		stmt           *sql.Stmt
		createuser_sql string
		code           int32
		hashpassword   string
		exists         bool
	)
	resp = &pbauth.Response{}
	createuser_sql = "INSERT INTO crocodile_user (username,hashpassword,email,avatar,super,forbid) VALUE(?,?,?,?,?,?)"

	if hashpassword, err = util.GenerateHashPass(user.Password); err != nil {
		logging.Errorf("GeneratehashPass Err: %v", err)
		code = e.ERR_GENERATE_HASHPASS_FAIL
		goto EXIT
	}

	if exists, err = s.userExist(ctx, user); exists {
		logging.Errorf("User Already Exists")
		code = e.ERR_USER_EXIST
		goto EXIT
	}

	if stmt, err = s.DB.PrepareContext(ctx, createuser_sql); err != nil {
		logging.Errorf("Prepare SQL %s Err:%v", createuser_sql, err)
		code = e.ERR_SQL_FAIl
		goto EXIT

	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, user.Username, hashpassword, user.Email, user.Avatar, user.Super, user.Forbid); err != nil {
		logging.Errorf("Exec Err: %v", err)
		code = e.ERR_CREATE_USER_FAIL
		goto EXIT

	}
	code = e.SUCCESS

	goto EXIT

EXIT:
	resp.Code = code
	logging.Infof("Create User code: %d ", resp.Code)
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
func (s *Service) ChangeUser(ctx context.Context, user *pbauth.User) (resp *pbauth.Response, err error) {
	var (
		stmt           *sql.Stmt
		changeuser_sql string
		code           int32
		hashpassword   string
		exists         bool
	)
	resp = &pbauth.Response{}

	changeuser_sql = "UPDATE crocodile_user SET email=?,avatar=?,forbid=?,super=? WHERE username=?"

	if exists, err = s.userExist(ctx, user); !exists {
		code = e.ERR_USER_NOT_EXIST
		logging.Errorf("User Not Exists")
		goto EXIT
	}

	// 存在密码时
	if user.Password != "" {
		changeuser_sql = "UPDATE crocodile_user SET hashpassword=?,email=?,avatar=?,forbid=?,super=? WHERE username=?"
		if hashpassword, err = util.GenerateHashPass(user.Password); err != nil {
			logging.Errorf("GeneratehashPass Err: %v", err)
			code = e.ERR_GENERATE_HASHPASS_FAIL
			goto EXIT
		}
	}

	if stmt, err = s.DB.PrepareContext(ctx, changeuser_sql); err != nil {
		logging.Errorf("Prepare SQL %s Err:%v", changeuser_sql, err)
		code = e.ERR_SQL_FAIl
		goto EXIT
	}
	defer stmt.Close()
	if user.Password != "" {
		if _, err = stmt.ExecContext(ctx, hashpassword, user.Email, user.Avatar, user.Forbid, user.Super, user.Username); err != nil {
			logging.Errorf("Exec Err: %v", err)
			code = e.ERR_CHANGE_USER_FAIL
			goto EXIT
		}
	} else {
		if _, err = stmt.ExecContext(ctx, user.Email, user.Avatar, user.Forbid, user.Super, user.Username); err != nil {
			logging.Errorf("Exec Err: %v", err)
			code = e.ERR_CHANGE_USER_FAIL
			goto EXIT
		}
	}

	code = e.SUCCESS

	goto EXIT

EXIT:
	resp.Code = code
	logging.Infof("Change User code: %d ", resp.Code)
	return
}

// 获取用户信息
// 存在用户名就返回这个用户的信息
// 否则返回所有用户的信息
// Option
//   username
func (s *Service) GetUser(ctx context.Context, user *pbauth.User) (resp *pbauth.Response, err error) {
	var (
		stmt        *sql.Stmt
		getuser_sql string
		query_name  string
		rows        *sql.Rows
		code        int32
		users       []*pbauth.User
	)
	resp = &pbauth.Response{}

	if user.Username == "" {
		query_name = "%"
	} else {
		query_name = user.Username
	}
	getuser_sql = "SELECT id,username, email, avatar, forbid, super FROM crocodile_user WHERE username like ?"
	if stmt, err = s.DB.PrepareContext(ctx, getuser_sql); err != nil {
		logging.Errorf("Prepare SQL %s Err:%v", getuser_sql, err)
		code = e.ERR_SQL_FAIl
		goto EXIT
	}
	defer stmt.Close()
	if rows, err = stmt.QueryContext(ctx, query_name); err != nil {
		logging.Errorf("SQL %s Query Err: %v", getuser_sql, err)
		code = e.ERR_GET_USER_FAIL
		goto EXIT
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

	code = e.SUCCESS

EXIT:
	resp.Code = code
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
		code            int32
		gethashpass_sql string
		stmt            *sql.Stmt
		exists          bool
		respuser        *pbauth.User
	)
	respuser = &pbauth.User{}
	resp = &pbauth.Response{}
	gethashpass_sql = "SELECT hashpassword,email,forbid,super FROM crocodile_user WHERE username=?"

	if exists, err = s.userExist(ctx, user); !exists {

		code = e.ERR_USER_NOT_EXIST
		err = errors.New(e.GetMsg(code))
		logging.Errorf("User Not Exists")
		goto EXIT
	}

	if stmt, err = s.DB.PrepareContext(ctx, gethashpass_sql); err != nil {
		logging.Errorf("Prepare SQL %s Err:%v", gethashpass_sql, err)
		code = e.ERR_SQL_FAIl
		goto EXIT
	}
	defer stmt.Close()
	if err = stmt.QueryRowContext(ctx, user.Username).Scan(&hashpassword, &respuser.Email, &respuser.Forbid, &respuser.Super); err != nil && err != sql.ErrNoRows {
		logging.Errorf("Query User %s Err:%v", user.Username, err)
		code = e.ERR_GET_USER_FAIL
		goto EXIT
	}
	if err == sql.ErrNoRows {
		logging.Errorf("User %s Not Exist Err:%v", user.Username, err)
		code = e.ERR_USER_PASS_FAIL
		goto EXIT
	}
	if respuser.Forbid {
		logging.Infof("User %s Forbid Login", user.Username)
		code = e.ERR_NOT_ALLOW_LOGIN
		goto EXIT
	}
	respuser.Username = user.Username

	if err = util.CheckHashPass(hashpassword, user.Password); err != nil {
		logging.Errorf("CheckHashPass Err: %v", err)
		code = e.ERR_USER_PASS_FAIL
		goto EXIT
	}
	logging.Infof("User %s Login Success", user.Username)
	if resp.Token, err = util.GenerateToken(respuser); err != nil {
		logging.Errorf("GenerateToken Err: %v", err)
		code = e.ERR_GENERATE_TOKEN_FAIL
		goto EXIT
	}

	code = e.SUCCESS
	goto EXIT

EXIT:
	resp.Code = code
	logging.Infof("Login User code: %d ", resp.Code)
	return
}

// 用户注销
// Require
//   username
func (s *Service) LogoutUser(ctx context.Context, user *pbauth.User) (resp *pbauth.Response, err error) {
	var (
		code   int32
		exists bool
	)
	resp = &pbauth.Response{}
	if exists, err = s.userExist(ctx, user); !exists {
		code = e.ERR_USER_NOT_EXIST
		logging.Errorf("User Not Exists")
		goto EXIT
	}

	if resp, err = s.GetUser(ctx, user); err != nil {
		logging.Errorf("GetUser %s Err: %v", user.Username, err)
		code = e.ERR_GET_USER_FAIL
		goto EXIT
	}
	if len(resp.Users) != 1 {
		logging.Errorf("User %s Not Exists", user.Username)
		code = e.ERR_USER_NOT_EXIST
		goto EXIT
	}
	code = e.SUCCESS
	goto EXIT

EXIT:
	resp.Code = code
	logging.Infof("Logout User code: %d ", resp.Code)
	return
}

// 删除用户
// Required
//   username
func (s *Service) DeleteUser(ctx context.Context, user *pbauth.User) (resp *pbauth.Response, err error) {
	var (
		deleteuser_sql string
		stmt           *sql.Stmt
		exists         bool
		code           int32
	)
	resp = &pbauth.Response{}
	logging.Info(user.Username)
	deleteuser_sql = "DELETE FROM crocodile_user WHERE username=?"
	if exists, err = s.userExist(ctx, user); !exists {
		code = e.ERR_USER_NOT_EXIST
		logging.Errorf("User Not Exists")
		goto EXIT
	}
	if stmt, err = s.DB.PrepareContext(ctx, deleteuser_sql); err != nil {
		logging.Errorf("Prepare SQL %s Err:%v", deleteuser_sql, err)
		code = e.ERR_SQL_FAIl
		goto EXIT
	}
	if _, err = stmt.ExecContext(ctx, user.Username); err != nil {
		logging.Errorf("Exec SQL %s Err: %v", deleteuser_sql, err)
		code = e.ERR_DELETE_USER_FAIL
		goto EXIT
	}
	code = e.SUCCESS
	goto EXIT

EXIT:
	resp.Code = code
	logging.Infof("Delete User code: %d ", resp.Code)
	return
}

// 查询用户是存在
// Required
//  username
func (s *Service) userExist(ctx context.Context, user *pbauth.User) (exists bool, err error) {
	var (
		resp *pbauth.Response
	)
	if resp, err = s.GetUser(ctx, user); err != nil {
		logging.Errorf("GetUser %s Err: %v", user.Username, err)
		return false, err
	}
	if len(resp.Users) == 0 {
		return false, nil

	}
	return true, nil
}
