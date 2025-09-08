package middlewares

import (
	"context"
	"fmt"
	"log"
	"myblog/tools"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuth中间件配置
func JWTAuth(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		//strings.TrimPrefix 移除字符串前缀
		tokenString := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid token"})
			return
		}

		// 先解析JWT token获取用户名
		token, err := jwt.ParseWithClaims(tokenString, &tools.MyClaims{}, func(token *jwt.Token) (interface{}, error) {
			log.Println(token)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { //作用是为了验证是基于HMAC（基于哈希的消息认证码）算法系列，而不是RSA或者AES加密。另外不用验证也可以，也不会报错。最好需要进行验证操作
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SIGNATURE_KEY")), nil
		})

		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid token"})
			return
		}

		// 检查Redis中token是否有效
		if claims, ok := token.Claims.(*tools.MyClaims); ok && token.Valid {
			// 查询Redis中存储的token
			storedToken, redisErr := rdb.Get(ctx, claims.Username).Result()
			if redisErr != nil {
				log.Println("Redis token not found or expired:", redisErr)
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "token expired or invalid"})
				return
			}

			// 比较请求中的token和Redis中存储的token是否一致
			if storedToken != tokenString {
				log.Println("Token mismatch - possible token refresh or invalid token")
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "token expired or invalid"})
				return
			}

			log.Println("打印获取的用户信息：", claims.Username, claims.Roles, claims.Exp)
			c.Set("userId", claims.UserId)
			c.Set("username", claims.Username)
			c.Set("roles", claims.Roles)
			c.Set("exp", claims.Exp)
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid token"})
		}
	}
}
