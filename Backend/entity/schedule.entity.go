package entity

import (
	"time"

	"github.com/google/uuid"
)

type Schedule struct {
	Base
	WorkflowID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"workflow_id"`
	TenantID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenant_id"`
	CronExpression string     `gorm:"type:varchar(120);not null" json:"cron_expression"`
	IsActive       bool       `gorm:"not null;default:true" json:"is_active"`
	LastRunAt      *time.Time `json:"last_run_at"`
	NextRunAt      *time.Time `json:"next_run_at"`
	CreatedBy      uuid.UUID  `gorm:"type:uuid" json:"created_by"`
}

func (Schedule) TableName() string { return "scheduler.schedules" }
