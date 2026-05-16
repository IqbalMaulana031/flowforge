package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Base
	TenantID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Email        string     `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
	PasswordHash string     `gorm:"type:text;not null" json:"-"`
	Name         string     `gorm:"type:varchar(150);not null" json:"name"`
	Role         string     `gorm:"type:varchar(30);not null" json:"role"`
	IsActive     bool       `gorm:"not null;default:true" json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at"`
}

func (User) TableName() string {
	return "auth.users"
}
