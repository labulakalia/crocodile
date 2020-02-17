package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	secret = "WteHkuweywy_egf7263i,ewd"
)

// Claims Jwt token
type Claims struct {
	jwt.StandardClaims
	UID      string
	UserName string
}

// GenerateToken generate get token for uid
func GenerateToken(uid string, username string) (token string, err error) {
	var (
		now         time.Time
		expireTime  time.Time
		claims      Claims
		tokenClaims *jwt.Token
	)
	now = time.Now()
	expireTime = now.Add(7 * 24 * time.Hour)

	claims = Claims{
		UID:      uid,
		UserName: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "crocodile",
		},
	}
	tokenClaims = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tokenClaims.SignedString([]byte(secret))

	return
}

// ParseToken parse token is valid
func ParseToken(token string) (claims *Claims, err error) {
	var (
		tokenClaims *jwt.Token
		ok          bool
	)
	// 解析出token的声明字段
	tokenClaims, err = jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if tokenClaims != nil {
		if claims, ok = tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return
		}
	}
	return
}
