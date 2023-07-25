/*
 * @Author: Monve
 * @Date: 2023-07-24 15:16:08
 * @LastEditors: Monve
 * @LastEditTime: 2023-07-25 12:31:35
 * @FilePath: /web-service-gin/internal/app/auth/auth.go
 */
package auth

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"web-service-gin/configs/env"
	"web-service-gin/internal/app/limiter"
	"web-service-gin/internal/pkg/redis"

	redis_o "github.com/go-redis/redis/v8"
)

// Jwt payload 类型
type CustomClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// 添加黑名单
func AddBlack(token string, expirationTime time.Time) error {
	return redis.Db.Set(token, "blacklisted", time.Until(expirationTime)).Err()
}

// 判断是否在黑名单
func IsBlacklisted(token string) bool {
	val, err := redis.Db.Get(token).Result()
	if err == redis_o.Nil {
		return false
	} else if err != nil {
		fmt.Println("Error checking token in Redis:", err)
		return false
	}
	return val == "blacklisted"
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(env.JwtSecretKey), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.JSON(401, gin.H{"error": "Invalid token signature"})
			} else {
				c.JSON(400, gin.H{"error": "Invalid token"})
			}
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(401, gin.H{"error": "Token is not valid"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*CustomClaims)
		if !ok {
			c.JSON(500, gin.H{"error": "Failed to parse token claims"})
			c.Abort()
			return
		}

		//限流
		var limit *limiter.Limiter

		if claims.Role == "owner" {
			//不限制,但记录
			limit = limiter.NewLimiter(c.Request.URL.Path+"_"+claims.Username, int(math.MaxInt), time.Minute)
		} else if claims.Role == "normal" {
			limit = limiter.NewLimiter(c.Request.URL.Path+"_"+claims.Username, 30, time.Minute)
		} else {
			limit = limiter.NewLimiter(c.Request.URL.Path+"_"+claims.Username, 5, time.Minute)
		}
		if !limit.Allow() {
			c.JSON(403, gin.H{"error": "Exceeded the limit, please try again later"})
			c.Abort()
			return
		}

		// 检查JWT令牌是否在黑名单中
		if IsBlacklisted(tokenString) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token blacklisted"})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("expiresAt", claims.ExpiresAt)
		c.Next()

		//记录访问次数
		recordAccess(claims.Username)
	}
}

func stringInArray(str string, arr []string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}

func isString(variable interface{}) bool {
	_, ok := variable.(string)
	return ok
}

func RoleMiddleware(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get("role")
		if !exists {
			c.JSON(403, gin.H{"error": "role not found"})
			c.Abort()
			return
		}
		if isString(val) {
			if stringInArray(val.(string), roles) {
				c.Next()
			}
		} else {
			c.JSON(403, gin.H{"error": "Role Invail"})
			c.Abort()
			return
		}
	}
}

func recordAccess(visitorID string) {
	today := time.Now().UTC().Format("2006-01-02")
	key := fmt.Sprintf("access:%s", today)
	redis.Db.HIncrBy(key, visitorID, 1)
	// 设置过期时间，保留7天的统计数据
	redis.Db.Expire(key, 7*24*time.Hour)
}
