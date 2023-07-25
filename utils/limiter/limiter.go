/*
 * @Author: Monve
 * @Date: 2023-07-24 17:55:19
 * @LastEditors: Monve
 * @LastEditTime: 2023-07-24 20:28:32
 * @FilePath: /web-service-gin/utils/limiter/limiter.go
 */
package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// Limiter 结构体定义
type Limiter struct {
	client    *redis.Client
	key       string
	rateLimit int
	interval  time.Duration
}

// NewLimiter 创建一个新的限流器实例
func NewLimiter(client *redis.Client, key string, rateLimit int, interval time.Duration) *Limiter {
	return &Limiter{
		client:    client,
		key:       key,
		rateLimit: rateLimit,
		interval:  interval,
	}
}

// Allow 检查是否允许处理请求
func (limiter *Limiter) Allow() bool {
	// 使用 Redis 原子操作获取当前令牌数量
	pipe := limiter.client.TxPipeline()
	pipe.ZRemRangeByScore(ctx, limiter.key, "0", fmt.Sprint(time.Now().Add(-limiter.interval).Unix()))
	pipe.ZCard(ctx, limiter.key)

	// 获取令牌数量
	cmds, err := pipe.Exec(ctx)
	if err != nil {
		fmt.Println("Error checking rate limit:", err)
		return false
	}

	tokens := cmds[1].(*redis.IntCmd).Val()

	// 判断是否允许处理请求
	if tokens < int64(limiter.rateLimit) {
		// 尝试补充令牌
		pipe.ZAdd(ctx, limiter.key, &redis.Z{
			Score:  float64(time.Now().Unix()),
			Member: time.Now().UnixNano(),
		})
		pipe.ZRemRangeByScore(ctx, limiter.key, "0", fmt.Sprint(time.Now().Add(-limiter.interval).Unix()))
		_, err := pipe.Exec(ctx)
		if err != nil {
			fmt.Println("Error updating rate limit:", err)
			return false
		}
		return true
	}

	return false
}
