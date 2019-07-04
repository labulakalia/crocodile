package bind

import (
	"github.com/gin-gonic/gin"
	"github.com/labulaka521/logging"
	"gopkg.in/go-playground/validator.v9"
)

// 绑定请求的JSON数据
func BindJson(c *gin.Context, data interface{}) (err error) {
	if err = c.ShouldBindJSON(data); err != nil {
		logging.Errorf("Bind Json Err:%v", err)
		return
	}
	//if err = Check(data); err != nil {
	//	return
	//}
	return
}

// 绑定请求的form数据
func BindQuery(c *gin.Context, data interface{}) (err error) {
	if err = c.ShouldBindQuery(data); err != nil {
		return
	}

	return
}

func Check(data interface{}) (err error) {
	var (
		validate *validator.Validate
	)
	validate = validator.New()

	if err = validate.Struct(data); err != nil {
		logging.Errorf("Check Err:%v", err)
		return
	}
	return
}
