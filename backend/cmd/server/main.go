package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rxcheck/rxcheck/internal/config"
	"github.com/rxcheck/rxcheck/internal/constants"
	"github.com/rxcheck/rxcheck/internal/database"
	"github.com/rxcheck/rxcheck/internal/repository"
	"github.com/rxcheck/rxcheck/internal/router"
	"github.com/rxcheck/rxcheck/internal/service"
	"github.com/rxcheck/rxcheck/internal/storage"
	"github.com/rxcheck/rxcheck/internal/util"
	"github.com/rxcheck/rxcheck/internal/worker"
)

func main() {
	cfg := config.Load()
	level := slog.LevelInfo
	if cfg.RunMode == "debug" {
		level = slog.LevelDebug
	}
	util.InitLogger(level)

	// 初始化数据库。
	db, err := database.New(cfg)
	if err != nil {
		util.Log.Error(fmt.Sprintf(constants.LogDBInitFailed, err))
		os.Exit(1)
	}
	// 种子数据（管理员/药品/相互作用规则，幂等）。
	if err := database.Seed(db, "admin", "admin123"); err != nil {
		util.Log.Error("种子数据初始化失败", "err", err)
		os.Exit(1)
	}

	// 初始化 Redis（失败仅告警，限流降级为内存、审核队列降级为同步）。
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
		DB:       cfg.RedisDB,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		util.Log.Warn(fmt.Sprintf(constants.LogRedisUnavailable, err))
		rdb = nil
	}

	// 初始化 MinIO（失败仅告警，报告导出接口会返回存储错误）。
	minioClient, err := storage.NewMinioClient(cfg.MinIOEndpoint, cfg.MinIOAccessKey, cfg.MinIOSecretKey, cfg.MinIOBucket, cfg.MinIOUseSSL, util.Log)
	if err != nil {
		util.Log.Warn("MinIO 初始化失败，报告导出功能不可用", "err", err)
		minioClient = nil
	}

	// 装配服务。
	auditSvc := service.NewAuditService(repository.NewAuditRepository(db), util.Log)
	reviewSvc := service.NewReviewService(db,
		repository.NewPrescriptionRepository(db),
		repository.NewDrugRepository(db),
		repository.NewInteractionRepository(db),
		repository.NewReportRepository(db),
		auditSvc, util.Log)

	engine := router.New(router.Deps{DB: db, Cfg: cfg, Log: util.Log, RDB: rdb, Minio: minioClient})

	// 启动异步审核 worker（Redis Stream）。
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go worker.NewReviewWorker(rdb, reviewSvc, util.Log).Start(ctx)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: engine,
	}
	go func() {
		util.Log.Info(fmt.Sprintf(constants.LogServerStarted, cfg.ServerPort, cfg.RunMode))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			util.Log.Error("服务监听失败", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	cancel()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		util.Log.Error("服务关闭失败", "err", err)
	}
	util.Log.Info(constants.LogServerStopped)
}
