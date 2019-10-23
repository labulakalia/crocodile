package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/labulaka521/crocodile/common/db"
	"github.com/labulaka521/crocodile/common/jwt"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/common/utils"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

type reqtype uint

const (
	Email reqtype = iota
	UserName
	Uid
)

// 检查存在的项
func IsExist(ctx context.Context, reqType reqtype, value interface{}) (bool, error) {
	check := "select COUNT(id) FROM crocodile_user WHERE "
	switch reqType {
	case Email:
		check += "email=?"
	case UserName:
		check += "name=?"
	case Uid:
		check += "id=?"
	default:
		return false, errors.New("reqType Only Support email username")
	}
	conn, err := db.GetConn(ctx)
	if err != nil {
		return false, errors.Wrap(err, "sqlDb.GetConn")
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, check)
	if err != nil {
		return false, errors.Wrap(err, "conn.PrepareContext")
	}
	defer stmt.Close()

	res := 0
	err = stmt.QueryRowContext(ctx, value).Scan(&res)
	if err != nil && err != sql.ErrNoRows {
		return false, errors.Wrap(err, "stmt.QueryRowContext")
	}
	if err == sql.ErrNoRows || res == 0 {
		return false, nil
	}
	return true, nil
}

// 登录用户
func LoginUser(ctx context.Context, name string, password string) (string, error) {
	var (
		hashpassword string
		uid          int64
		forbid       int
	)
	loguser := `SELECT id, hashpassword, forbid FROM crocodile_user WHERE name=?`

	conn, err := db.GetConn(ctx)
	if err != nil {
		return "", errors.Wrap(err, "db.Db.GetConn")
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, loguser)
	if err != nil {
		return "", errors.Wrap(err, "conn.PrepareContext")
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, name).Scan(&uid, &hashpassword, &forbid)
	if err != nil && err != sql.ErrNoRows {
		return "", errors.Wrap(err, "stmt.QueryRowContext Scan")
	}
	if forbid != 0 {
		return "", errors.Wrap(define.ErrForbid{name}, "")
	}

	err = utils.CheckHashPass(hashpassword, password)
	if err != nil {
		return "", errors.Wrap(define.ErrUserPass{err}, "utils.CheckHashPass")
	}
	token, err := jwt.GenerateToken(uid)
	if err != nil {
		return "", errors.Wrap(err, "jwt.GenerateToken")
	}

	return token, nil
}

// 添加用户
func AddUser(ctx context.Context, u *define.User) error {
	adduser := `INSERT INTO crocodile_user (id,name,email,hashpassword,role,forbid,remark,createTime,updateTime)VALUES(?,?,?,?,?,?,?,?,?)`
	conn, err := db.GetConn(ctx)
	if err != nil {
		return errors.Wrap(err, "db.Db.GetConn")
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, adduser)
	if err != nil {
		return errors.Wrap(err, "conn.PrepareContext")
	}
	defer stmt.Close()

	now := time.Now().Unix()

	_, err = stmt.ExecContext(ctx, u.Id, u.Name, u.Email, u.Password,
		u.Role, u.Forbid, u.Remark, now, now)
	if err != nil {
		return errors.Wrap(err, "stmt.ExecContext")
	}

	pass, err := enforcer.AddRoleForUser(fmt.Sprintf("%d", u.Id),
		define.GetUserByRole(u.Role))
	if err != nil {
		return err
	}
	if !pass {
		return errors.New("AddRoleForUser failed")
	}

	return nil
}

// 获取用户
func GetUser(ctx context.Context, uid int64) ([]*define.User, error) {
	getuser := `select id,name,role,forbid,email,createTime,updateTime,remark FROM crocodile_user `
	args := []interface{}{}
	users := []*define.User{}
	if uid > 0 {
		getuser += "WHERE id=?"
		args = append(args, uid)
	}
	conn, err := db.GetConn(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "db.Db.GetConn")
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, getuser)
	if err != nil {
		return nil, errors.Wrap(err, "conn.PrepareContext")
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, errors.Wrap(err, "stmt.QueryContext")
	}
	for rows.Next() {
		var (
			createTime int64
			updateTime int64
		)
		user := define.User{}
		err := rows.Scan(&user.Id, &user.Name, &user.Role,
			&user.Forbid, &user.Email, &createTime,
			&updateTime, &user.Remark,
		)
		if err != nil {
			log.Error("Scan Err", zap.Error(err))
			continue
		}
		user.CreateTime = utils.UnixToStr(createTime)
		user.UpdateTime = utils.UnixToStr(updateTime)
		users = append(users, &user)
	}
	return users, nil
}

// 修改用户
func ChangeUser(ctx context.Context, u *define.User) error {
	changeuser := `UPDATE crocodile_user SET role=?,forbid=?,email=?,updateTime=?,remark=? WHERE id=? AND name=?`
	conn, err := db.GetConn(ctx)
	if err != nil {
		return errors.Wrap(err, "db.Db.GetConn")
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, changeuser)
	if err != nil {
		return errors.Wrap(err, "conn.PrepareContext")
	}
	defer stmt.Close()
	updateTime := time.Now().Unix()
	_, err = stmt.ExecContext(ctx, u.Role, u.Forbid, u.Email, updateTime, u.Remark, u.Id, u.Name)
	if err != nil {
		return errors.Wrap(err, "stmt.ExecContext")
	}
	return nil
}
