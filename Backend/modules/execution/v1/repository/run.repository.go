package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"flowforge-api/entity"
)

type RunRepositoryUseCase interface {
	ActiveWorkflow(ctx context.Context, tenantID, wfid uuid.UUID) (*entity.Workflow, *entity.WorkflowVersion, error)
	CreateRun(ctx context.Context, run *entity.WorkflowRun) error
	UpdateRun(ctx context.Context, run *entity.WorkflowRun) error
	CreateStep(ctx context.Context, step *entity.RunStep) error
	UpdateStep(ctx context.Context, step *entity.RunStep) error
	CreateLog(ctx context.Context, log *entity.ExecutionLog) error
	ListRuns(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]entity.WorkflowRun, int64, error)
	FindRun(ctx context.Context, tenantID, runID uuid.UUID) (*entity.WorkflowRun, error)
	Steps(ctx context.Context, tenantID, runID uuid.UUID) ([]entity.RunStep, error)
	Logs(ctx context.Context, tenantID, runID uuid.UUID) ([]entity.ExecutionLog, error)
}
type RunRepository struct{ db *gorm.DB }

func NewRunRepository(db *gorm.DB) *RunRepository { return &RunRepository{db: db} }
func (r *RunRepository) ActiveWorkflow(ctx context.Context, tenantID, wfid uuid.UUID) (*entity.Workflow, *entity.WorkflowVersion, error) {
	var wf entity.Workflow
	if err := r.db.WithContext(ctx).Where("tenant_id=? AND id=? AND is_active=true", tenantID, wfid).First(&wf).Error; err != nil {
		return nil, nil, err
	}
	var ver entity.WorkflowVersion
	err := r.db.WithContext(ctx).Where("tenant_id=? AND workflow_id=? AND version_number=?", tenantID, wfid, wf.CurrentVersion).First(&ver).Error
	return &wf, &ver, err
}
func (r *RunRepository) CreateRun(ctx context.Context, run *entity.WorkflowRun) error {
	return r.db.WithContext(ctx).Create(run).Error
}
func (r *RunRepository) UpdateRun(ctx context.Context, run *entity.WorkflowRun) error {
	return r.db.WithContext(ctx).Save(run).Error
}
func (r *RunRepository) CreateStep(ctx context.Context, step *entity.RunStep) error {
	return r.db.WithContext(ctx).Create(step).Error
}
func (r *RunRepository) UpdateStep(ctx context.Context, step *entity.RunStep) error {
	return r.db.WithContext(ctx).Save(step).Error
}
func (r *RunRepository) CreateLog(ctx context.Context, log *entity.ExecutionLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}
func (r *RunRepository) ListRuns(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]entity.WorkflowRun, int64, error) {
	var items []entity.WorkflowRun
	db := r.db.WithContext(ctx).Model(&entity.WorkflowRun{}).Where("tenant_id=?", tenantID)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&items).Error
	return items, total, err
}
func (r *RunRepository) FindRun(ctx context.Context, tenantID, runID uuid.UUID) (*entity.WorkflowRun, error) {
	var run entity.WorkflowRun
	err := r.db.WithContext(ctx).Where("tenant_id=? AND id=?", tenantID, runID).First(&run).Error
	return &run, err
}
func (r *RunRepository) Steps(ctx context.Context, tenantID, runID uuid.UUID) ([]entity.RunStep, error) {
	var steps []entity.RunStep
	err := r.db.WithContext(ctx).Where("tenant_id=? AND run_id=?", tenantID, runID).Find(&steps).Error
	return steps, err
}
func (r *RunRepository) Logs(ctx context.Context, tenantID, runID uuid.UUID) ([]entity.ExecutionLog, error) {
	var logs []entity.ExecutionLog
	err := r.db.WithContext(ctx).Where("tenant_id=? AND run_id=?", tenantID, runID).Order("logged_at ASC").Find(&logs).Error
	return logs, err
}
