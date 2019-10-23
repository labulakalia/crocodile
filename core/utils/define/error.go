package define

import "fmt"

// 用户密码错误
type ErrUserPass struct {
	Err error
}

func (u ErrUserPass) Error() string {
	return "username or password error: " + u.Err.Error()
}

// 用户禁止登陆
type ErrForbid struct {
	Name string
}

func (u ErrForbid) Error() string {
	return fmt.Sprintf("user %s forbid login", u.Name)
}
