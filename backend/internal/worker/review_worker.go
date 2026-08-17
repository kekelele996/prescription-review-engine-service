package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rxcheck/rxcheck/internal/constants"
	"github.com/rxcheck/rxcheck/internal/queue"
	"github.com/rxcheck/rxcheck/internal/service"
)

// ReviewWorker Redis Stream 审核消费者：异步处理处方审核任务。
type ReviewWorker struct {
	rdb *redis.Client
	svc *service.ReviewService
	log *slog.Logger
}

func NewReviewWorker(rdb *redis.Client, svc *service.ReviewService, log *slog.Logger) *ReviewWorker {
	return &ReviewWorker{rdb: rdb, svc: svc, log: log}
}

// Start 启动消费循环（阻塞，需在 goroutine 中运行）。
func (w *ReviewWorker) Start(ctx context.Context) {
	if w.rdb == nil {
		w.log.Warn("Redis 不可用，审核 worker 未启动")
		return
	}
	consumer := "worker-" + uuid.NewString()[:8]
	// 幂等创建消费组。
	_, err := w.rdb.XGroupCreateMkStream(ctx, queue.ReviewStream, queue.ReviewGroup, "$").Result()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		w.log.Warn("创建审核消费组失败", "err", err)
	}
	for {
		select {
		case <-ctx.Done():
			w.log.Info("审核 worker 已退出")
			return
		default:
		}
		res, err := w.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    queue.ReviewGroup,
			Consumer: consumer,
			Streams:  []string{queue.ReviewStream, ">"},
			Count:    10,
			Block:    2 * time.Second,
		}).Result()
		if err != nil {
			if err != redis.Nil {
				w.log.Warn("读取审核任务失败", "err", err)
			}
			time.Sleep(500 * time.Millisecond)
			continue
		}
		for _, stream := range res {
			for _, msg := range stream.Messages {
				w.process(ctx, msg)
			}
		}
	}
}

// process 处理单条审核任务，成功后确认。
func (w *ReviewWorker) process(ctx context.Context, msg redis.XMessage) {
	var job queue.ReviewJob
	if err := json.Unmarshal([]byte(fmt.Sprintf("%v", msg.Values["payload"])), &job); err != nil {
		w.log.Warn("审核任务载荷解析失败", "err", err)
		_, _ = w.rdb.XAck(ctx, queue.ReviewStream, queue.ReviewGroup, msg.ID).Result()
		return
	}
	report, err := w.svc.ReviewPrescription(job.PrescriptionID, "async-worker", "", "")
	if err != nil {
		w.log.Error("审核任务处理失败", "prescription_id", job.PrescriptionID, "err", err)
		return // 不确认，交由 PEL 重试
	}
	w.log.Info(fmt.Sprintf(constants.LogQueueProcessed, job.PrescriptionID, report.Status))
	if _, err := w.rdb.XAck(ctx, queue.ReviewStream, queue.ReviewGroup, msg.ID).Result(); err != nil {
		w.log.Warn(fmt.Sprintf(constants.LogQueueAckFailed, job.PrescriptionID, err))
	}
}
