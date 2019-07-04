package util

import (
	"github.com/labulaka521/logging"
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
	logging.Info(hashpass, password)
	if err = bcrypt.CompareHashAndPassword([]byte(hashpass), []byte(password)); err != nil {
		logging.Errorf("CompareHashAndPassword Err: %v", err)
		return err
	}
	return nil
}
