package router

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// 自定义路由日志格式
func CustomLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求前获取当前时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 请求后获取消耗时间
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 过滤掉探活接口
		if reqUri == "/healthz" {
			return
		}

		// 日志格式，模仿默认日志格式 gin.Logger()
		fmt.Printf("[Gin] %v | %3d | %13v | %15s | %-7s  \"%s\"\n",
			startTime.Format("2006/01/02 - 15:04:05"), statusCode, latencyTime, clientIP, reqMethod, reqUri)
	}
}

// 过滤掉探活接口
func LoggerNoHealthz() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/healthz"},
	})
}
