package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// GenerateHashPass generate hashpassword by text password
func GenerateHashPass(password string) (hashpassword string, err error) {
	var (
		generatepass []byte
	)
	if generatepass, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); err != nil {
		return hashpassword, err
	}
	hashpassword = string(generatepass)
	return hashpassword, nil
}

// CheckHashPass check hashpassword valid
func CheckHashPass(hashpass string, password string) (err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashpass), []byte(password)); err != nil {
		return err
	}
	return nil
}
