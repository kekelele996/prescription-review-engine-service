package router

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rxcheck/rxcheck/internal/config"
	"github.com/rxcheck/rxcheck/internal/handler"
	"github.com/rxcheck/rxcheck/internal/middleware"
	"github.com/rxcheck/rxcheck/internal/repository"
	"github.com/rxcheck/rxcheck/internal/service"
	"github.com/rxcheck/rxcheck/internal/storage"
	"github.com/rxcheck/rxcheck/internal/util"
	"gorm.io/gorm"
)

// Deps 路由装配依赖。
type Deps struct {
	DB    *gorm.DB
	Cfg   *config.Config
	Log   *slog.Logger
	RDB   *redis.Client
	Minio *storage.MinioClient
}

// New 构建 Gin 引擎并注册全部路由。
func New(deps Deps) *gin.Engine {
	if deps.Cfg.RunMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.CORS())
	r.Use(middleware.ErrorHandler(deps.Log))
	r.Use(middleware.RateLimit(deps.RDB, deps.Cfg.RateLimit))

	// 健康检查。
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, util.Resp{Code: 0, Message: "ok", Data: gin.H{"status": "healthy", "service": "rxcheck"}})
	})
	registerSwaggerRoutes(r)

	// 仓储层。
	userRepo := repository.NewUserRepository(deps.DB)
	drugRepo := repository.NewDrugRepository(deps.DB)
	interRepo := repository.NewInteractionRepository(deps.DB)
	prescRepo := repository.NewPrescriptionRepository(deps.DB)
	reportRepo := repository.NewReportRepository(deps.DB)
	auditRepo := repository.NewAuditRepository(deps.DB)

	// 服务层。
	auditSvc := service.NewAuditService(auditRepo, deps.Log)
	authSvc := service.NewAuthService(userRepo, auditSvc, deps.Cfg.JWTSecret, int(deps.Cfg.JWTExpire.Hours()), deps.Log)
	userSvc := service.NewUserService(userRepo, auditSvc, deps.Log)
	drugSvc := service.NewDrugService(drugRepo, auditSvc, deps.Log)
	interSvc := service.NewInteractionService(interRepo, drugRepo, auditSvc, deps.Log)
	reviewSvc := service.NewReviewService(deps.DB, prescRepo, drugRepo, interRepo, reportRepo, auditSvc, deps.Log)
	prescSvc := service.NewPrescriptionService(deps.DB, prescRepo, drugRepo, reviewSvc, auditSvc, deps.RDB, deps.Log)
	reportSvc := service.NewReportService(reportRepo, prescRepo, deps.Minio, auditSvc, deps.Log)
	statsSvc := service.NewStatsService(prescRepo, drugRepo, interRepo, reportRepo, deps.Log)

	// 处理器层。
	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)
	drugHandler := handler.NewDrugHandler(drugSvc)
	interHandler := handler.NewInteractionHandler(interSvc)
	prescHandler := handler.NewPrescriptionHandler(prescSvc)
	reportHandler := handler.NewReportHandler(reportSvc)
	statsHandler := handler.NewStatsHandler(statsSvc)
	auditHandler := handler.NewAuditHandler(auditSvc)

	api := r.Group("/api/v1")
	{
		registerAuthRoutes(api, authHandler)
		protected := api.Group("")
		protected.Use(middleware.RequireAnyAuth(deps.Cfg)) // API Key + JWT 双认证
		protected.Use(middleware.Audit(auditRepo))         // 操作审计
		{
			registerUserRoutes(protected, userHandler)
			registerDrugRoutes(protected, drugHandler)
			registerInteractionRoutes(protected, interHandler)
			registerPrescriptionRoutes(protected, prescHandler, reportHandler)
			registerReportRoutes(protected, reportHandler)
			registerStatsRoutes(protected, statsHandler)
			registerAuditRoutes(protected, auditHandler)
		}
	}
	return r
}
