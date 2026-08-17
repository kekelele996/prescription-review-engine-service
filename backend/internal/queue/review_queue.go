package queue

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/redis/go-redis/v9"
)

// 审核任务 Redis Stream 常量。
const (
	// ReviewStream 审核任务流。
	ReviewStream = "rxcheck:review:stream"
	// ReviewGroup 消费组。
	ReviewGroup = "rxcheck:review:workers"
)

// ReviewJob 审核任务载荷。
type ReviewJob struct {
	PrescriptionID uint `json:"prescription_id"`
}

// EnqueueReview 将处方审核任务写入 Redis Stream（异步队列）。
func EnqueueReview(ctx context.Context, rdb *redis.Client, prescriptionID uint) error {
	if rdb == nil {
		return errors.New("redis unavailable")
	}
	payload, err := json.Marshal(ReviewJob{PrescriptionID: prescriptionID})
	if err != nil {
		return err
	}
	return rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: ReviewStream,
		Values: map[string]any{"payload": string(payload)},
	}).Err()
}
