package utils

import (
	"net"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var (
	errBadFormat        = errors.New("invalid format")
	errUnresolvableHost = errors.New("unresolvable smtp host")
	emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

func validateFormat(email string) error {
	if !emailRegexp.MatchString(email) {
		return errBadFormat
	}
	return nil
}

func validateHost(email string) error {
	_, host := split(email)
	_, err := net.LookupMX(host)
	if err != nil {
		return errUnresolvableHost
	}
	return nil
}
func split(email string) (account, host string) {
	i := strings.LastIndexByte(email, '@')
	account = email[:i]
	host = email[i+1:]
	return
}

// CheckEmail will check a email is valid
func CheckEmail(email string) error {
	err := validateFormat(email)
	if err != nil {
		return err
	}
	err = validateHost(email)
	if err != nil {
		return err
	}
	return nil
}
