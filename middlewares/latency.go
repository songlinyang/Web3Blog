package middlewares

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 1.Logger中间件配置
// LatencyLogger 创建一个记录请求耗时的中间件
func LatencyLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		start := time.Now()

		// 在处理请求之前可以记录一些信息
		fmt.Printf("[%s] %s %s - Request started\n",
			start.Format("2006-01-02 15:04:05"),
			c.Request.Method,
			c.Request.URL.Path)
		// 在中间件中添加调试信息
		c.Set("requestID", uuid.NewString())
		log.Printf("[%s] %s %s", c.GetString("requestID"), c.Request.Method, c.Request.URL)
		// 处理请求
		c.Next()

		// 计算耗时
		latency := time.Since(start)

		// 获取响应状态码
		status := c.Writer.Status()

		// 记录请求完成信息，包括耗时
		fmt.Printf("[%s] %s %s - Completed in %v with status %d\n",
			time.Now().Format("2006-01-02 15:04:05"),
			c.Request.Method,
			c.Request.URL.Path,
			latency,
			status)
	}
}
