package define

import (
	"errors"
	"fmt"
)

// ErrBadRequest err bad request
type ErrBadRequest struct {
}

func (e ErrBadRequest) Error() string {
	return "bad request"
}

// ErrUserPass user pass err
type ErrUserPass struct {
	UserName string
}

func (e ErrUserPass) Error() string {
	return fmt.Sprintf("user %s password error ", e.UserName)
}

// ErrForbid user forbid login err
type ErrForbid struct {
	Name string
}

func (e ErrForbid) Error() string {
	return fmt.Sprintf("user %s forbid login", e.Name)
}

// ErrDelHostID delete host id err
type ErrDelHostID struct {
	ID string
}

func (e ErrDelHostID) Error() string {
	return fmt.Sprintf("can delete hostid %s, it use by other hostgroup", e.ID)
}

// ErrNotExist query not exist
type ErrNotExist struct {
	Type  string
	Value string
}

func (e ErrNotExist) Error() string {
	return fmt.Sprintf("%s value %s is not exist", e.Type, e.Value)
}

// ErrExist query  exist err
type ErrExist struct {
	Type  string
	Value string
}

func (e ErrExist) Error() string {
	return fmt.Sprintf("%s value %s is exist", e.Type, e.Value)
}

// ErrIsUsed operate data is used
type ErrIsUsed struct {
	Type  string
	Value string
}

func (e ErrIsUsed) Error() string {
	return fmt.Sprintf("%s value %s is used", e.Type, e.Value)
}

// ErrUnauthorized not has operate
type ErrUnauthorized struct {
	Type string
}

func (e ErrUnauthorized) Error() string {
	return fmt.Sprintf("do't have operate %s authority", e.Type)
}

// ErrDependByOther object depend by other, can not operate
type ErrDependByOther struct {
	Type string
}

func (e ErrDependByOther) Error() string {
	return fmt.Sprintf("depend by %s model,can not delete", e.Type)
}

// ErrServer if err not find, will return this error
type ErrServer struct {
}

func (e ErrServer) Error() string {
	return fmt.Sprintf("please try again later, the server is busy")
}

// GetError get first error
func GetError(err error) error {
	switch err = errors.Unwrap(err); err.(type) {
	case ErrUserPass:
		return err
	case ErrForbid:
		return err
	case ErrDelHostID:
		return err
	case ErrDependByOther:
		return err
	case ErrIsUsed:
		return err
	case ErrNotExist:
		return err
	case ErrExist:
		return err
	case ErrUnauthorized:
		return err
	default:
		// other error
		return ErrServer{}
	}
}
