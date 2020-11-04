package resp

import (
	"github.com/gin-gonic/gin"
)

// Response api response
type Response struct {
	Code      int         `json:"code" comment:"111"`        // msg success:0 failed:1
	Msg       string      `json:"msg"`                       // code
	Data      interface{} `json:"data,omitempty" form:"111"` // data
	Count     int         `json:"count,omitempty"`           // data count
	Releation interface{} `json:"releation,omitempty"`       // releation info like taskid:taskname
}

// JSON gin resp to json
// data
// 废弃
func JSON(c *gin.Context, code int, data ...interface{}) {
	resp := Response{
		Code: code,
		Msg:  GetMsg(code),
	}
	if len(data) == 1 {
		resp.Data = data[0]
	}
	if len(data) == 2 {
		resp.Count = data[1].(int)
	}

	if len(data) == 3 {
		resp.Releation = data[2]
	}
	c.JSON(200, resp)
	c.Set("statuscode", code)

	return
}

// JSONv2 json response
func JSONv2(c *gin.Context, err error, data ...interface{}) {
	resp := Response{
		Msg: "ok",
	}

	if err != nil {
		resp.Msg = err.Error()
		resp.Code = 1
	}
	if len(data) == 1 {
		resp.Data = data[0]
	}
	if len(data) == 2 {
		resp.Count = data[1].(int)
	}

	if len(data) == 3 {
		resp.Releation = data[2]
	}
	c.JSON(200, resp)
	c.Set("statuscode", 1)
}
