package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"

	"flowforge-api/entity"
	"flowforge-api/modules/schedule/v1/repository"
	"flowforge-api/resource"
)

type ScheduleUseCase interface {
	Create(context.Context, string, string, resource.ScheduleRequest) (*entity.Schedule, error)
	List(context.Context, string) ([]entity.Schedule, error)
	Update(context.Context, string, string, resource.ScheduleRequest) (*entity.Schedule, error)
	Delete(context.Context, string, string) error
}
type ScheduleService struct {
	repo repository.ScheduleRepositoryUseCase
}

func NewScheduleService(repo repository.ScheduleRepositoryUseCase) *ScheduleService {
	return &ScheduleService{repo: repo}
}
func validate(expr string) error { _, err := cron.ParseStandard(expr); return err }
func (s *ScheduleService) Create(ctx context.Context, tid, userID string, req resource.ScheduleRequest) (*entity.Schedule, error) {
	if err := validate(req.CronExpression); err != nil {
		return nil, err
	}
	tenant, _ := uuid.Parse(tid)
	user, _ := uuid.Parse(userID)
	wf, _ := uuid.Parse(req.WorkflowID)
	active := true
	if req.IsActive != nil {
		active = *req.IsActive
	}
	sched := &entity.Schedule{TenantID: tenant, WorkflowID: wf, CronExpression: req.CronExpression, IsActive: active, CreatedBy: user}
	return sched, s.repo.Create(ctx, sched)
}
func (s *ScheduleService) List(ctx context.Context, tid string) ([]entity.Schedule, error) {
	tenant, _ := uuid.Parse(tid)
	return s.repo.List(ctx, tenant)
}
func (s *ScheduleService) Update(ctx context.Context, tid, id string, req resource.ScheduleRequest) (*entity.Schedule, error) {
	if err := validate(req.CronExpression); err != nil {
		return nil, err
	}
	tenant, _ := uuid.Parse(tid)
	sid, _ := uuid.Parse(id)
	sched, err := s.repo.Find(ctx, tenant, sid)
	if err != nil {
		return nil, err
	}
	wf, _ := uuid.Parse(req.WorkflowID)
	sched.WorkflowID = wf
	sched.CronExpression = req.CronExpression
	if req.IsActive != nil {
		sched.IsActive = *req.IsActive
	}
	return sched, s.repo.Save(ctx, sched)
}
func (s *ScheduleService) Delete(ctx context.Context, tid, id string) error {
	tenant, _ := uuid.Parse(tid)
	sid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, tenant, sid)
}
