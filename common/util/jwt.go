package util

import (
	"crocodile/common/cfg"
	auth "crocodile/service/auth/proto/auth"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Claims struct {
	Username string
	Email    string
	Forbid   bool
	Super    bool
	jwt.StandardClaims
}

func GenerateToken(user *auth.User) (token string, err error) {
	var (
		now         time.Time
		expireTime  time.Time
		claims      Claims
		tokenClaims *jwt.Token
	)
	now = time.Now()
	expireTime = now.Add(7 * 24 * time.Hour)

	claims = Claims{
		Username: user.Username,
		Email:    user.Email,
		Forbid:   user.Forbid,
		Super:    user.Super,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "cron",
		},
	}
	tokenClaims = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tokenClaims.SignedString([]byte(cfg.JwtConfig.SecretKey))

	return
}

// 解析token
func ParseToken(token string) (claims *Claims, err error) {
	var (
		tokenClaims *jwt.Token
		ok          bool
	)
	// 解析出token的声明字段
	tokenClaims, err = jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JwtConfig.SecretKey), nil
	})

	if tokenClaims != nil {
		if claims, ok = tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return
		}
	}
	return
}
