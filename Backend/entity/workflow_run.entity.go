package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type WorkflowRun struct {
	ID                uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	WorkflowID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"workflow_id"`
	TenantID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenant_id"`
	WorkflowVersionID uuid.UUID      `gorm:"type:uuid;not null;index" json:"workflow_version_id"`
	TriggerType       string         `gorm:"type:varchar(30);not null" json:"trigger_type"`
	Status            string         `gorm:"type:varchar(30);not null;index" json:"status"`
	StartedAt         *time.Time     `json:"started_at"`
	FinishedAt        *time.Time     `json:"finished_at"`
	DurationMs        int64          `json:"duration_ms"`
	ErrorMessage      string         `gorm:"type:text" json:"error_message"`
	TriggerPayload    datatypes.JSON `gorm:"type:jsonb" json:"trigger_payload"`
	CreatedAt         time.Time      `json:"created_at"`
}

func (WorkflowRun) TableName() string { return "execution.workflow_runs" }
