package model

import (
	"context"
	"database/sql"
	"github.com/labulaka521/crocodile/common/db"
	"github.com/labulaka521/crocodile/common/jwt"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/common/utils"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

// LoginUser login user
func LoginUser(ctx context.Context, name string, password string) (string, error) {
	var (
		hashpassword string
		uid          string
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
	if forbid == 1 {
		return "", errors.Wrap(define.ErrForbid{Name:name}, "")
	}

	err = utils.CheckHashPass(hashpassword, password)
	if err != nil {
		return "", errors.Wrap(define.ErrUserPass{Err:err}, "utils.CheckHashPass")
	}
	token, err := jwt.GenerateToken(uid)
	if err != nil {
		return "", errors.Wrap(err, "jwt.GenerateToken")
	}

	return token, nil
}

// AddUser add new user
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

	_, err = stmt.ExecContext(ctx, u.ID, u.Name, u.Email, u.Password,
		u.Role, u.Forbid, u.Remark, now, now)
	if err != nil {
		return errors.Wrap(err, "stmt.ExecContext")
	}

	ok, err := enforcer.AddRoleForUser(u.ID, u.Role.String())
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("AddRoleForUser failed")
	}

	return nil
}

// GetUser get user by id
func GetUser(ctx context.Context, uid string) ([]define.User, error) {
	getuser := `select id,name,role,forbid,email,createTime,updateTime,remark FROM crocodile_user `
	args := []interface{}{}
	users := []define.User{}
	if uid != "" {
		getuser += "WHERE id=?"
		args = append(args, uid)
	}
	conn, err := db.GetConn(ctx)
	if err != nil {
		return users, errors.Wrap(err, "db.Db.GetConn")
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, getuser)
	if err != nil {
		return users, errors.Wrap(err, "conn.PrepareContext")
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return users, errors.Wrap(err, "stmt.QueryContext")
	}
	for rows.Next() {
		var (
			createTime int64
			updateTime int64
		)
		user := define.User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Role,
			&user.Forbid, &user.Email, &createTime,
			&updateTime, &user.Remark,
		)
		if err != nil {
			log.Error("Scan Err", zap.Error(err))
			continue
		}
		user.CreateTime = utils.UnixToStr(createTime)
		user.UpdateTime = utils.UnixToStr(updateTime)
		users = append(users, user)
	}
	return users, nil
}

// ChangeUser change user message
func ChangeUser(ctx context.Context, u *define.User, role define.Role) error {
	var (
		changeuser string
		args       []interface{}
	)
	conn, err := db.GetConn(ctx)
	if err != nil {
		return errors.Wrap(err, "db.Db.GetConn")
	}
	defer conn.Close()
	updateTime := time.Now().Unix()
	switch role {
	case define.AdminUser:
		changeuser = `UPDATE crocodile_user SET role=?,forbid=?,email=?,updateTime=?,remark=? WHERE id=?`
		args = append(args, u.Role, u.Forbid, u.Email, updateTime, u.Remark, u.ID)
		// 修改权限表
		ok, err := enforcer.DeleteUser(u.ID)
		if err != nil {
			return err
		}
		if !ok {
			return errors.New("Delete user failed")
		}
		ok, err = enforcer.AddRoleForUser(u.ID, u.Role.String())
		if err != nil {
			return err
		}
		if !ok {
			return errors.New("AddRoleForUser failed")
		}
		err = enforcer.LoadPolicy()
		if err != nil {
			return errors.Wrap(err, "enforcer.LoadPolicy")
		}
	case define.NormalUser:
		changeuser = `UPDATE crocodile_user SET email=?,updateTime=?,remark=? WHERE id=?`
		args = append(args, u.Email, updateTime, u.Remark, u.ID)
	}
	stmt, err := conn.PrepareContext(ctx, changeuser)
	if err != nil {
		return errors.Wrap(err, "conn.PrepareContext")
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		return errors.Wrap(err, "stmt.ExecContext")
	}

	return nil
}
