package middleware

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rxcheck/rxcheck/internal/constants"
	"github.com/rxcheck/rxcheck/internal/util"
	"github.com/redis/go-redis/v9"
)

// RateLimit 接口限流中间件：优先 Redis 滑动窗口，Redis 不可用时降级内存计数。
func RateLimit(rdb *redis.Client, limitPerSecond int) gin.HandlerFunc {
	mem := &memLimiter{buckets: make(map[string]*bucket)}
	return func(c *gin.Context) {
		key := "rl:" + c.ClientIP()
		allowed := true
		if rdb != nil {
			allowed = redisAllow(c.Request.Context(), rdb, key, limitPerSecond)
		} else {
			allowed = mem.allow(key, limitPerSecond)
		}
		if !allowed {
			util.Log.Warn(fmt.Sprintf(constants.LogRateLimited, c.ClientIP(), c.Request.URL.Path))
			util.Fail(c, http.StatusTooManyRequests, constants.CodeRateLimited, constants.MsgRateLimited)
			c.Abort()
			return
		}
		c.Next()
	}
}

// redisAllow Redis ZSET 滑动窗口限流。
func redisAllow(ctx context.Context, rdb *redis.Client, key string, limit int) bool {
	now := time.Now().Unix()
	pipe := rdb.Pipeline()
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", now-1))
	pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: fmt.Sprintf("%d-%d", now, time.Now().Nanosecond())})
	pipe.Expire(ctx, key, 2*time.Second)
	cmd := pipe.ZCard(ctx, key)
	if _, err := pipe.Exec(ctx); err != nil {
		return true // Redis 异常时放行，避免阻断业务。
	}
	return cmd.Val() <= int64(limit)
}

type bucket struct {
	count int
	reset time.Time
}

type memLimiter struct {
	mu      sync.Mutex
	buckets map[string]*bucket
}

func (m *memLimiter) allow(key string, limit int) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	b, ok := m.buckets[key]
	if !ok || now.After(b.reset) {
		m.buckets[key] = &bucket{count: 1, reset: now.Add(time.Second)}
		return true
	}
	b.count++
	return b.count <= limit
}
