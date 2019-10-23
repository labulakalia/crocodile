package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/utils/resp"
	"go.uber.org/zap"
	"strings"
	"time"
)

// 日志
func ZapLogger() func(c *gin.Context) {
	return func(c *gin.Context) {
		start := time.Now()
		url := c.Request.URL.RequestURI()

		c.Next()
		latency := time.Now().Sub(start)
		statuscode := c.GetInt("statuscode")
		reqip := c.ClientIP()
		method := c.Request.Method
		bodySize := c.Writer.Size()

		fields := []zap.Field{
			zap.Int64("uid", c.GetInt64("uid")),
			zap.String("method", strings.ToLower(method)),
			zap.Int("statuscode", statuscode),
			zap.String("reqip", reqip),
			zap.Duration("latency", latency),
			zap.String("url", url),
		}

		if bodySize > 0 {
			fields = append(fields, zap.Int("respsize", bodySize))
		}

		switch {
		case statuscode < resp.ErrBadRequest:
			log.Info("Gin", fields...)
		case statuscode < resp.ErrInternalServer:
			log.Warn("Gin", fields...)
		case statuscode >= resp.ErrInternalServer:
			log.Error("Gin", fields...)
		}
	}
}
