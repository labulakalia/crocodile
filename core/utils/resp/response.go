package resp

import (
	"github.com/gin-gonic/gin"
)

// Response api response
type Response struct {
	Code  int         `json:"code" comment:"111"`        // msg
	Msg   string      `json:"msg"`                       // code
	Data  interface{} `json:"data,omitempty" form:"111"` // data
	Count int         `json:"count,omitempty"`           // data count
}

// JSON gin resp to json
func JSON(c *gin.Context, code int, data ...interface{}) {
	resp := Response{
		Code: code,
		Msg:  GetMsg(code),
		Data: data[0],
	}
	if len(data) == 2 {
		resp.Count = data[1].(int)
	}
	c.JSON(200, resp)
	c.Set("statuscode", code)

	return
}
