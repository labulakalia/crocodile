package response

import (
	"crocodile/common/e"
	"github.com/gin-gonic/gin"
)

type Gin struct {
	C *gin.Context
}

// HTTP接口resp
type Response struct {
	Code int32       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (g *Gin) Response(code int32, data interface{}) {
	g.C.JSON(200, Response{
		Code: code,
		Msg:  e.GetMsg(code),
		Data: data,
	})
	return
}
