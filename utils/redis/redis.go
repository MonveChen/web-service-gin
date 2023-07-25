/*
 * @Author: Monve
 * @Date: 2023-07-24 14:52:31
 * @LastEditors: Monve
 * @LastEditTime: 2023-07-25 10:25:16
 * @FilePath: /web-service-gin/utils/redis/redis.go
 */
package redis

import (
	"bytes"
	"context"
	"fmt"
	"time"
	"web-service-gin/utils/env"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var Db *redis.Client

func Init() {
	// 初始化Redis客户端
	Db = redis.NewClient(&redis.Options{
		Addr: env.RedisAddr, // Redis服务器地址
		// Password: "",                             // Redis服务器密码
		DB: 0, // 选择使用的数据库，默认为0
	})

	// 检查是否可以连接到Redis服务器
	_, err := Db.Ping(context.Background()).Result()
	if err != nil {
		panic("无法连接到Redis服务器")
	} else {
		fmt.Println("【redis】连接成功")
	}
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func CacheMiddleware(c *gin.Context) {
	key := c.Request.URL.String()

	// 尝试从缓存中获取数据
	val, err := Db.Get(context.Background(), key).Result()
	if err == nil {
		c.Status(200)
		c.Writer.Header().Add("Content-Type", "application/json; charset=utf-8")
		c.Writer.Write([]byte(val))
		c.Abort()
		return
	}

	w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
	c.Writer = w

	// 缓存未命中，执行处理程序
	c.Next()

	if c.Writer.Status() == 200 {
		// 获取处理程序返回的数据
		response := w.body.String()

		// 将数据保存到缓存中，设置适当的过期时间
		Db.Set(context.Background(), key, response, time.Duration(60)*time.Second)
	}
}
