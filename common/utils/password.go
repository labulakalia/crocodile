package utils

import (
	"golang.org/x/crypto/bcrypt"
)

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

func CheckHashPass(hashpass string, password string) (err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashpass), []byte(password)); err != nil {
		return err
	}
	return nil
}
