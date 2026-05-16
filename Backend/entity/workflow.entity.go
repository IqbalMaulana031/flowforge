package entity

import "github.com/google/uuid"

type Workflow struct {
	Base
	TenantID       uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Name           string    `gorm:"type:varchar(200);not null" json:"name"`
	Description    string    `gorm:"type:text" json:"description"`
	CurrentVersion int       `gorm:"not null;default:1" json:"current_version"`
	IsActive       bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedBy      uuid.UUID `gorm:"type:uuid" json:"created_by"`
}

func (Workflow) TableName() string { return "workflow.workflows" }
