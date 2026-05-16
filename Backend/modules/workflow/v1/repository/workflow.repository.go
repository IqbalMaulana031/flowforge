package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"flowforge-api/entity"
)

type WorkflowRepositoryUseCase interface {
	Create(ctx context.Context, wf *entity.Workflow, ver *entity.WorkflowVersion) error
	List(ctx context.Context, tenantID uuid.UUID, limit, offset int, q string) ([]entity.Workflow, int64, error)
	Find(ctx context.Context, tenantID, id uuid.UUID) (*entity.Workflow, *entity.WorkflowVersion, error)
	UpdateWithVersion(ctx context.Context, wf *entity.Workflow, ver *entity.WorkflowVersion) error
	Delete(ctx context.Context, tenantID, id uuid.UUID) error
	Versions(ctx context.Context, tenantID, wfid uuid.UUID) ([]entity.WorkflowVersion, error)
	FindVersion(ctx context.Context, tenantID, wfid uuid.UUID, version int) (*entity.WorkflowVersion, error)
	SetCurrentVersion(ctx context.Context, wf *entity.Workflow, version int) error
}
type WorkflowRepository struct{ db *gorm.DB }

func NewWorkflowRepository(db *gorm.DB) *WorkflowRepository { return &WorkflowRepository{db: db} }
func (r *WorkflowRepository) Create(ctx context.Context, wf *entity.Workflow, ver *entity.WorkflowVersion) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(wf).Error; err != nil {
			return err
		}
		ver.WorkflowID = wf.ID
		ver.TenantID = wf.TenantID
		if err := tx.Create(ver).Error; err != nil {
			return err
		}
		wf.CurrentVersion = ver.VersionNumber
		return tx.Save(wf).Error
	})
}
func (r *WorkflowRepository) List(ctx context.Context, tenantID uuid.UUID, limit, offset int, q string) ([]entity.Workflow, int64, error) {
	var items []entity.Workflow
	db := r.db.WithContext(ctx).Model(&entity.Workflow{}).Where("tenant_id=? AND is_active=true", tenantID)
	if q != "" {
		db = db.Where("name ILIKE ?", "%"+q+"%")
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&items).Error
	return items, total, err
}
func (r *WorkflowRepository) Find(ctx context.Context, tenantID, id uuid.UUID) (*entity.Workflow, *entity.WorkflowVersion, error) {
	var wf entity.Workflow
	if err := r.db.WithContext(ctx).Where("tenant_id=? AND id=? AND is_active=true", tenantID, id).First(&wf).Error; err != nil {
		return nil, nil, err
	}
	var ver entity.WorkflowVersion
	err := r.db.WithContext(ctx).Where("tenant_id=? AND workflow_id=? AND version_number=?", tenantID, id, wf.CurrentVersion).First(&ver).Error
	return &wf, &ver, err
}
func (r *WorkflowRepository) UpdateWithVersion(ctx context.Context, wf *entity.Workflow, ver *entity.WorkflowVersion) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		wf.CurrentVersion = ver.VersionNumber
		if err := tx.Save(wf).Error; err != nil {
			return err
		}
		return tx.Create(ver).Error
	})
}
func (r *WorkflowRepository) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&entity.Workflow{}).Where("tenant_id=? AND id=?", tenantID, id).Update("is_active", false).Error
}
func (r *WorkflowRepository) Versions(ctx context.Context, tenantID, wfid uuid.UUID) ([]entity.WorkflowVersion, error) {
	var versions []entity.WorkflowVersion
	err := r.db.WithContext(ctx).Where("tenant_id=? AND workflow_id=?", tenantID, wfid).Order("version_number DESC").Find(&versions).Error
	return versions, err
}
func (r *WorkflowRepository) FindVersion(ctx context.Context, tenantID, wfid uuid.UUID, version int) (*entity.WorkflowVersion, error) {
	var ver entity.WorkflowVersion
	err := r.db.WithContext(ctx).Where("tenant_id=? AND workflow_id=? AND version_number=?", tenantID, wfid, version).First(&ver).Error
	return &ver, err
}
func (r *WorkflowRepository) SetCurrentVersion(ctx context.Context, wf *entity.Workflow, version int) error {
	wf.CurrentVersion = version
	return r.db.WithContext(ctx).Save(wf).Error
}
