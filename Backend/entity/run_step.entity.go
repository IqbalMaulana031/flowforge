package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type RunStep struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	RunID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"run_id"`
	TenantID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenant_id"`
	StepID        string         `gorm:"type:varchar(120);not null" json:"step_id"`
	StepName      string         `gorm:"type:varchar(200);not null" json:"step_name"`
	StepType      string         `gorm:"type:varchar(50);not null" json:"step_type"`
	Status        string         `gorm:"type:varchar(30);not null;index" json:"status"`
	AttemptNumber int            `gorm:"not null;default:0" json:"attempt_number"`
	StartedAt     *time.Time     `json:"started_at"`
	FinishedAt    *time.Time     `json:"finished_at"`
	DurationMs    int64          `json:"duration_ms"`
	ErrorMessage  string         `gorm:"type:text" json:"error_message"`
	Output        datatypes.JSON `gorm:"type:jsonb" json:"output"`
}

func (RunStep) TableName() string { return "execution.run_steps" }
