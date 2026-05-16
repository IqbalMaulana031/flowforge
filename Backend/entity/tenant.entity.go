package entity

type Tenant struct {
	Base
	Name   string `gorm:"type:varchar(150);not null" json:"name"`
	Slug   string `gorm:"type:varchar(120);uniqueIndex;not null" json:"slug"`
	Status string `gorm:"type:varchar(30);not null;default:active" json:"status"`
}

func (Tenant) TableName() string {
	return "master.tenants"
}
