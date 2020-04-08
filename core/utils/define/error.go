package define

import "fmt"

// ErrUserPass user pass err
type ErrUserPass struct {
	Err error
}

func (u ErrUserPass) Error() string {
	return "username or password error: " + u.Err.Error()
}

// ErrForbid user forbid login err
type ErrForbid struct {
	Name string
}

func (u ErrForbid) Error() string {
	return fmt.Sprintf("user %s forbid login", u.Name)
}

// ErrDelHostID delete host id err
type ErrDelHostID struct {
	ID string
}

func (u ErrDelHostID) Error() string {
	return fmt.Sprintf("can delete hostid %s, it use by other hostgroup", u.ID)
}

// ErrNotExist query not exist
type ErrNotExist struct {
	Value string
}

func (u ErrNotExist) Error() string {
	return fmt.Sprintf("value %s is not exist", u.Value)
}
