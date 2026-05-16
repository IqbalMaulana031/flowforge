package service

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"flowforge-api/resource"
)

type HealthUseCase interface {
	Status(ctx context.Context) resource.HealthResource
	Ping() resource.PingResource
}

type HealthService struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewHealthService(db *gorm.DB, redis *redis.Client) *HealthService {
	return &HealthService{db: db, redis: redis}
}

func (s *HealthService) Status(ctx context.Context) resource.HealthResource {
	services := map[string]string{
		"api":      "ok",
		"database": "unknown",
		"redis":    "unknown",
	}

	if s.db != nil {
		if sqlDB, err := s.db.DB(); err == nil {
			pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
			defer cancel()
			if err := sqlDB.PingContext(pingCtx); err == nil {
				services["database"] = "ok"
			} else {
				services["database"] = "error"
			}
		}
	}

	if s.redis != nil {
		pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
		if err := s.redis.Ping(pingCtx).Err(); err == nil {
			services["redis"] = "ok"
		} else {
			services["redis"] = "error"
		}
	}

	status := "ok"
	for _, serviceStatus := range services {
		if serviceStatus == "error" {
			status = "degraded"
			break
		}
	}

	return resource.HealthResource{Status: status, Services: services}
}

func (s *HealthService) Ping() resource.PingResource {
	return resource.PingResource{Message: "pong"}
}
