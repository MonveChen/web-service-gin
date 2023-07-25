/*
 * @Author: Monve
 * @Date: 2023-07-25 09:49:14
 * @LastEditors: Monve
 * @LastEditTime: 2023-07-25 10:25:56
 * @FilePath: /web-service-gin/utils/env/env.go
 */
package env

import (
	"os"
	"time"
)

func If(condition bool, trueVal, falseVal string) string {
	if condition {
		return trueVal
	}
	return falseVal
}

var JwtSecretKey = If(os.Getenv("JWT_SECRET_KEY") != "", os.Getenv("JWT_SECRET_KEY"), "your_secret_key")
var JwtTokenExpiresIn = time.Hour * 2
var PostgresUrl = If(os.Getenv("POSTGRES_URL") != "", os.Getenv("POSTGRES_URL"), "user=root password=123456 dbname=postgres host=server1.dibiaozuitu.com port=5432 sslmode=disable")
var RedisAddr = If(os.Getenv("REDIS_ADDR") != "", os.Getenv("REDIS_ADDR"), "server1.dibiaozuitu.com:6379")
