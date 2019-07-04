package util

import (
	"crocodile/common/cfg"
	"crocodile/common/log"

	auth "crocodile/service/auth/proto/auth"
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	var (
		err    error
		token  string
		claims *Claims
	)
	log.Init()
	cfg.Init()
	username := "test1"
	user := &auth.User{Username: username}
	if token, err = GenerateToken(user); err != nil {
		t.Error(err)
	}

	claims, err = ParseToken(token)
	if claims.ExpiresAt > time.Now().Unix() != true {
		t.Errorf("Token Should Not Expires ")
	}
	// true 为没有过期 false 为已经过期
	if claims.VerifyExpiresAt(time.Now().Unix()-10, false) != true {
		t.Errorf("Token Should Should Expires ")
	}

	if err != nil {

		t.Fatalf("Parse Token Fail:%v\n", err)
	}

	if username != claims.Username {
		t.Error(username, claims.Username)
	}
}
