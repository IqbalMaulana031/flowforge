package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"flowforge-api/config"
	"flowforge-api/middleware"
	aiHandler "flowforge-api/modules/ai/v1/handler"
	authHandler "flowforge-api/modules/auth/v1/handler"
	runHandler "flowforge-api/modules/execution/v1/handler"
	healthHandler "flowforge-api/modules/health/v1/handler"
	realtimeHandler "flowforge-api/modules/realtime/v1/handler"
	scheduleHandler "flowforge-api/modules/schedule/v1/handler"
	workflowHandler "flowforge-api/modules/workflow/v1/handler"
	"flowforge-api/response"
)

const apiV1Prefix = "/v1"

func DefaultHTTPHandler(cfg *config.Config, router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		response.Success(c, gin.H{"name": cfg.App.Name, "status": "ok", "version": "v1"}, nil)
	})
	router.NoRoute(func(c *gin.Context) {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "route not found")
	})
}

func publicV1(router *gin.Engine) *gin.RouterGroup {
	return router.Group(apiV1Prefix)
}

func protectedV1(cfg *config.Config, redis *redis.Client, router *gin.Engine) *gin.RouterGroup {
	return router.Group(
		apiV1Prefix,
		middleware.Auth(cfg, redis),
		middleware.TenantContext(),
		middleware.RateLimit(cfg, redis),
	)
}

func AuthHTTPHandler(cfg *config.Config, redis *redis.Client, router *gin.Engine, h *authHandler.AuthHandler) {
	public := publicV1(router).Group("/auth")
	{
		public.POST("/register", h.Register)
		public.POST("/login", h.Login)
		public.POST("/refresh-token", h.Refresh)
	}

	protected := protectedV1(cfg, redis, router).Group("/auth")
	{
		protected.POST("/logout", h.Logout)
		protected.GET("/profile", h.Profile)
		protected.PUT("/profile", h.UpdateProfile)
	}
}

func WorkflowHTTPHandler(cfg *config.Config, redis *redis.Client, router *gin.Engine, h *workflowHandler.WorkflowHandler) {
	workflows := protectedV1(cfg, redis, router).Group("/workflows")
	{
		workflows.GET("", h.List)
		workflows.POST("", middleware.RequireRoles(middleware.RoleAdmin, middleware.RoleEditor), h.Create)
		workflows.GET("/:id", h.Detail)
		workflows.PUT("/:id", middleware.RequireRoles(middleware.RoleAdmin, middleware.RoleEditor), h.Update)
		workflows.DELETE("/:id", middleware.RequireRoles(middleware.RoleAdmin), h.Delete)
		workflows.GET("/:id/versions", h.Versions)
		workflows.POST("/:id/rollback", middleware.RequireRoles(middleware.RoleAdmin, middleware.RoleEditor), h.Rollback)
	}
}

func RunHTTPHandler(cfg *config.Config, redis *redis.Client, router *gin.Engine, h *runHandler.RunHandler) {
	runs := protectedV1(cfg, redis, router).Group("/runs")
	{
		runs.GET("", h.List)
		runs.POST("/trigger", middleware.RequireRoles(middleware.RoleAdmin, middleware.RoleEditor), h.Trigger)
		runs.GET("/:runId", h.Detail)
		runs.GET("/:runId/steps", h.Steps)
		runs.GET("/:runId/logs", h.Logs)
		runs.DELETE("/:runId", middleware.RequireRoles(middleware.RoleAdmin, middleware.RoleEditor), h.Cancel)
	}

	publicV1(router).POST("/runs/webhook/:workflowId", h.Webhook)
}

func ScheduleHTTPHandler(cfg *config.Config, redis *redis.Client, router *gin.Engine, h *scheduleHandler.ScheduleHandler) {
	schedules := protectedV1(cfg, redis, router).Group("/schedules")
	{
		schedules.GET("", h.List)
		schedules.POST("", middleware.RequireRoles(middleware.RoleAdmin, middleware.RoleEditor), h.Create)
		schedules.GET("/:id", h.Detail)
		schedules.PUT("/:id", middleware.RequireRoles(middleware.RoleAdmin, middleware.RoleEditor), h.Update)
		schedules.DELETE("/:id", middleware.RequireRoles(middleware.RoleAdmin, middleware.RoleEditor), h.Delete)
	}
}

func HealthHTTPHandler(router *gin.Engine, h *healthHandler.HealthHandler) {
	health := publicV1(router).Group("/health")
	{
		health.GET("", h.Status)
		health.GET("/ping", h.Ping)
	}
}

func RealtimeHTTPHandler(cfg *config.Config, redis *redis.Client, router *gin.Engine, h *realtimeHandler.RealtimeHandler) {
	realtime := router.Group(
		"/",
		middleware.Auth(cfg, redis),
		middleware.TenantContext(),
	)
	{
		realtime.GET("/ws", h.WebSocket)
		realtime.GET("/events", h.Events)
	}
}

func AIHTTPHandler(cfg *config.Config, redis *redis.Client, router *gin.Engine, h *aiHandler.AIHandler) {
	ai := protectedV1(cfg, redis, router).Group("/ai")
	{
		ai.POST("/generate-workflow", middleware.RequireRoles(middleware.RoleAdmin, middleware.RoleEditor), h.GenerateWorkflow)
		ai.POST("/analyze-failure/:runId", h.AnalyzeFailure)
	}
}
