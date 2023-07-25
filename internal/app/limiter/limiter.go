/*
 * @Author: Monve
 * @Date: 2023-07-24 17:55:19
 * @LastEditors: Monve
 * @LastEditTime: 2023-07-25 12:26:06
 * @FilePath: /web-service-gin/utils/limiter/limiter.go
 */
package limiter

import (
	"fmt"
	"time"
	"web-service-gin/internal/pkg/redis"

	gredis "github.com/go-redis/redis"
)

// Limiter 结构体定义
type Limiter struct {
	key       string
	rateLimit int
	interval  time.Duration
}

// NewLimiter 创建一个新的限流器实例
func NewLimiter(key string, rateLimit int, interval time.Duration) *Limiter {
	return &Limiter{
		key:       key,
		rateLimit: rateLimit,
		interval:  interval,
	}
}

// Allow 检查是否允许处理请求
func (limiter *Limiter) Allow() bool {
	// 使用 Redis 原子操作获取当前令牌数量
	pipe := redis.Db.TxPipeline()
	pipe.ZRemRangeByScore(limiter.key, "0", fmt.Sprint(time.Now().Add(-limiter.interval).Unix()))
	pipe.ZCard(limiter.key)

	// 获取令牌数量
	cmds, err := pipe.Exec()
	if err != nil {
		fmt.Println("Error checking rate limit:", err)
		return false
	}

	tokens := cmds[1].(*gredis.IntCmd).Val()

	// 判断是否允许处理请求
	if tokens < int64(limiter.rateLimit) {
		// 尝试补充令牌
		pipe.ZAdd(limiter.key, gredis.Z{
			Score:  float64(time.Now().Unix()),
			Member: time.Now().UnixNano(),
		})
		pipe.ZRemRangeByScore(limiter.key, "0", fmt.Sprint(time.Now().Add(-limiter.interval).Unix()))
		_, err := pipe.Exec()
		if err != nil {
			fmt.Println("Error updating rate limit:", err)
			return false
		}
		return true
	}

	return false
}
