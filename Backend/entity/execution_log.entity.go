package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type ExecutionLog struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	RunID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"run_id"`
	StepRowID *uuid.UUID     `gorm:"type:uuid;index" json:"step_row_id"`
	TenantID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Level     string         `gorm:"type:varchar(20);not null" json:"level"`
	Message   string         `gorm:"type:text;not null" json:"message"`
	Metadata  datatypes.JSON `gorm:"type:jsonb" json:"metadata"`
	LoggedAt  time.Time      `gorm:"not null;index" json:"logged_at"`
}

func (ExecutionLog) TableName() string { return "execution.execution_logs" }
