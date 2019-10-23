package resp

import (
	"github.com/gin-gonic/gin"
)

// HTTP接口resp
type response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func Json(c *gin.Context, code int, data interface{}) {
	c.JSON(200, response{
		Code: code,
		Msg:  GetMsg(code),
		Data: data,
	})
	c.Set("statuscode", code)

	return
}
