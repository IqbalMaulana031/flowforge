package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"flowforge-api/entity"
)

type ScheduleRepositoryUseCase interface {
	Create(context.Context, *entity.Schedule) error
	List(context.Context, uuid.UUID) ([]entity.Schedule, error)
	Find(context.Context, uuid.UUID, uuid.UUID) (*entity.Schedule, error)
	Save(context.Context, *entity.Schedule) error
	Delete(context.Context, uuid.UUID, uuid.UUID) error
}
type ScheduleRepository struct{ db *gorm.DB }

func NewScheduleRepository(db *gorm.DB) *ScheduleRepository { return &ScheduleRepository{db: db} }
func (r *ScheduleRepository) Create(ctx context.Context, s *entity.Schedule) error {
	return r.db.WithContext(ctx).Create(s).Error
}
func (r *ScheduleRepository) List(ctx context.Context, tid uuid.UUID) ([]entity.Schedule, error) {
	var s []entity.Schedule
	err := r.db.WithContext(ctx).Where("tenant_id=?", tid).Order("created_at DESC").Find(&s).Error
	return s, err
}
func (r *ScheduleRepository) Find(ctx context.Context, tid, id uuid.UUID) (*entity.Schedule, error) {
	var s entity.Schedule
	err := r.db.WithContext(ctx).Where("tenant_id=? AND id=?", tid, id).First(&s).Error
	return &s, err
}
func (r *ScheduleRepository) Save(ctx context.Context, s *entity.Schedule) error {
	return r.db.WithContext(ctx).Save(s).Error
}
func (r *ScheduleRepository) Delete(ctx context.Context, tid, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("tenant_id=? AND id=?", tid, id).Delete(&entity.Schedule{}).Error
}
