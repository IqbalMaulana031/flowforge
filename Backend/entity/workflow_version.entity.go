package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type WorkflowVersion struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	WorkflowID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"workflow_id"`
	TenantID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenant_id"`
	VersionNumber int            `gorm:"not null" json:"version_number"`
	DAGDefinition datatypes.JSON `gorm:"type:jsonb;not null" json:"dag_definition"`
	Changelog     string         `gorm:"type:text" json:"changelog"`
	CreatedAt     time.Time      `json:"created_at"`
	CreatedBy     uuid.UUID      `gorm:"type:uuid" json:"created_by"`
}

func (WorkflowVersion) TableName() string { return "workflow.workflow_versions" }
