package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"gorm.io/gorm"

	"flowforge-api/app"
	"flowforge-api/common/logger"
	"flowforge-api/config"
	"flowforge-api/middleware"
	aiBuilder "flowforge-api/modules/ai/v1/builder"
	authBuilder "flowforge-api/modules/auth/v1/builder"
	executionBuilder "flowforge-api/modules/execution/v1/builder"
	healthBuilder "flowforge-api/modules/health/v1/builder"
	realtimeBuilder "flowforge-api/modules/realtime/v1/builder"
	scheduleBuilder "flowforge-api/modules/schedule/v1/builder"
	workflowBuilder "flowforge-api/modules/workflow/v1/builder"
	"flowforge-api/utils"
)

var Module = fx.Options(fx.Provide(config.NewConfig, logger.NewLogger, NewRouter, utils.NewPostgresGormDB, utils.NewRedisClient))

func main() {
	fx.New(Module, authBuilder.AuthModule, workflowBuilder.WorkflowModule, executionBuilder.ExecutionModule, scheduleBuilder.ScheduleModule, healthBuilder.HealthModule, realtimeBuilder.RealtimeModule, aiBuilder.AIModule, fx.Invoke(app.DefaultHTTPHandler), fx.Invoke(StartServer)).Run()
}

func NewRouter(cfg *config.Config) *gin.Engine {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(gin.Recovery(), middleware.RequestID(), middleware.CORS(cfg))
	return router
}

type ServerParams struct {
	fx.In
	Lifecycle fx.Lifecycle
	Config    *config.Config
	Router    *gin.Engine
	Logger    *slog.Logger
	DB        *gorm.DB
	Redis     *redis.Client
}

func StartServer(params ServerParams) {
	server := &http.Server{Addr: fmt.Sprintf(":%s", params.Config.App.Port), Handler: params.Router, ReadHeaderTimeout: 5 * time.Second}
	params.Lifecycle.Append(fx.Hook{OnStart: func(ctx context.Context) error {
		params.Logger.Info("starting HTTP server", "addr", server.Addr)
		go func() {
			if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				params.Logger.Error("HTTP server failed", "error", err)
			}
		}()
		return nil
	}, OnStop: func(ctx context.Context) error {
		params.Logger.Info("stopping HTTP server")
		if err := server.Shutdown(ctx); err != nil {
			return err
		}
		if params.Redis != nil {
			_ = params.Redis.Close()
		}
		if params.DB != nil {
			sqlDB, err := params.DB.DB()
			if err == nil {
				_ = sqlDB.Close()
			}
		}
		return nil
	}})
}
