package middle

// 基础的验证中间件
// 验证访问是否带有token 及验证token的有效性
import (
	"crocodile/common/e"
	"crocodile/common/response"
	"crocodile/common/util"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/labulaka521/logging"
	"strings"
	"time"
)

// 解析请求头中的Token
func GetClaims(c *gin.Context) (claims *util.Claims, err error) {
	var (
		token         string
		authorization string
		data          []string
	)

	if authorization, err = GetAuthor(c); err != nil {
		return
	}

	data = strings.Split(authorization, " ")
	if len(data) != 2 {
		return
	}
	if data[0] != "Bearer" {
		return
	}

	token = data[1]
	// 解析token
	claims, err = util.ParseToken(token)

	return
}

// 获取请求头的authorization
func GetAuthor(c *gin.Context) (authorization string, err error) {
	authorization = c.GetHeader("Authorization")

	if authorization == "" {
		err = errors.New("Invalid Token")
	}
	return
}

var excludeurl = []string{"/auth/login"}

// 验证token
func MiddleJwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			app     response.Gin
			rooturl string
			url     string
			claims  *util.Claims
			err     error
			code    int32
		)
		app = response.Gin{c}

		rooturl = strings.Split(c.Request.RequestURI, "?")[0]

		for _, url = range excludeurl {
			if rooturl == url {
				c.Next()
				return
			}
		}
		if claims, err = GetClaims(c); err != nil {
			code = e.ERR_TOKEN_INVALID
			goto ERR
		}

		if !claims.VerifyExpiresAt(time.Now().Unix(), false) {
			code = e.ERR_TOKEN_INVALID
			goto ERR
		}

		c.Set("user", claims.Username)
		c.Set("super", claims.Super)
		c.Set("email", claims.Email)
		c.Set("forbid", claims.Forbid)
		c.Next()
		return

	ERR:
		// 解析失败返回错误
		c.Writer.Header().Add("WWW-Authenticate", fmt.Sprintf("Bearer realm='%s'", e.GetMsg(code)))
		app.Response(code, nil)
		logging.Errorf("Token Verify Fail: %v %s", err, e.GetMsg(code))
		c.Abort()
		return
	}
}
