package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	"flowforge-api/common/realtime"
	"flowforge-api/entity"
	"flowforge-api/modules/execution/v1/engine"
	"flowforge-api/modules/execution/v1/repository"
	"flowforge-api/resource"
	"flowforge-api/response"
	"flowforge-api/utils"
)

type RunUseCase interface {
	Trigger(ctx context.Context, tenantID string, workflowID string, payload map[string]any, triggerType string) (*resource.RunResource, error)
	List(ctx context.Context, tenantID string, page, limit int) ([]resource.RunResource, *response.Meta, error)
	Detail(ctx context.Context, tenantID, runID string) (*resource.RunResource, error)
	Steps(ctx context.Context, tenantID, runID string) ([]entity.RunStep, error)
	Logs(ctx context.Context, tenantID, runID string) ([]entity.ExecutionLog, error)
	Cancel(ctx context.Context, tenantID, runID string) error
}
type RunService struct {
	repo     repository.RunRepositoryUseCase
	executor *engine.Executor
	hub      *realtime.Hub
}

func NewRunService(repo repository.RunRepositoryUseCase, executor *engine.Executor, hub *realtime.Hub) *RunService {
	return &RunService{repo: repo, executor: executor, hub: hub}
}
func (s *RunService) Trigger(ctx context.Context, tenantID string, workflowID string, payload map[string]any, triggerType string) (*resource.RunResource, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}
	wid, err := uuid.Parse(workflowID)
	if err != nil {
		return nil, err
	}
	_, ver, err := s.repo.ActiveWorkflow(ctx, tid, wid)
	if err != nil {
		return nil, err
	}
	rawPayload, _ := json.Marshal(payload)
	now := time.Now()
	run := &entity.WorkflowRun{WorkflowID: wid, TenantID: tid, WorkflowVersionID: ver.ID, TriggerType: triggerType, Status: "running", StartedAt: &now, TriggerPayload: datatypes.JSON(rawPayload)}
	if err := s.repo.CreateRun(ctx, run); err != nil {
		return nil, err
	}
	s.publish(tid.String(), run.ID.String(), "", "run", "running", nil)
	go s.execute(context.Background(), run, []byte(ver.DAGDefinition))
	return toRun(run), nil
}
func (s *RunService) execute(ctx context.Context, run *entity.WorkflowRun, raw []byte) {
	dag, err := engine.Parse(raw)
	if err != nil {
		s.finish(ctx, run, "failed", err)
		return
	}
	stepRows := map[string]*entity.RunStep{}
	for _, st := range dag.Steps {
		row := &entity.RunStep{RunID: run.ID, TenantID: run.TenantID, StepID: st.ID, StepName: st.Name, StepType: string(st.Type), Status: "pending"}
		_ = s.repo.CreateStep(ctx, row)
		stepRows[st.ID] = row
	}
	err = s.executor.Execute(ctx, dag, func(st engine.Step, status string, attempt int, output map[string]any, stepErr error) {
		row := stepRows[st.ID]
		now := time.Now()
		if row.StartedAt == nil {
			row.StartedAt = &now
		}
		row.Status = status
		row.AttemptNumber = attempt
		if output != nil {
			b, _ := json.Marshal(output)
			row.Output = datatypes.JSON(b)
		}
		if stepErr != nil {
			row.ErrorMessage = stepErr.Error()
		}
		if status == "completed" || status == "failed" {
			row.FinishedAt = &now
			row.DurationMs = now.Sub(*row.StartedAt).Milliseconds()
		}
		_ = s.repo.UpdateStep(ctx, row)
		level := "info"
		msg := st.Name + " " + status
		if stepErr != nil {
			level = "error"
			msg = stepErr.Error()
		}
		_ = s.repo.CreateLog(ctx, &entity.ExecutionLog{RunID: run.ID, StepRowID: &row.ID, TenantID: run.TenantID, Level: level, Message: msg, LoggedAt: time.Now()})
		s.publish(run.TenantID.String(), run.ID.String(), st.ID, "step", status, output)
	})
	if err != nil {
		s.finish(ctx, run, "failed", err)
		return
	}
	s.finish(ctx, run, "completed", nil)
}
func (s *RunService) finish(ctx context.Context, run *entity.WorkflowRun, status string, err error) {
	now := time.Now()
	run.Status = status
	run.FinishedAt = &now
	if run.StartedAt != nil {
		run.DurationMs = now.Sub(*run.StartedAt).Milliseconds()
	}
	if err != nil {
		run.ErrorMessage = err.Error()
	}
	_ = s.repo.UpdateRun(ctx, run)
	_ = s.repo.CreateLog(ctx, &entity.ExecutionLog{RunID: run.ID, TenantID: run.TenantID, Level: "info", Message: "run " + status, LoggedAt: time.Now()})
	s.publish(run.TenantID.String(), run.ID.String(), "", "run", status, nil)
}
func (s *RunService) publish(tenant, run, step, typ, status string, payload any) {
	if s.hub != nil {
		s.hub.Publish(realtime.Event{TenantID: tenant, RunID: run, StepID: step, Type: typ, Status: status, Payload: payload})
	}
}
func (s *RunService) List(ctx context.Context, tenantID string, page, limit int) ([]resource.RunResource, *response.Meta, error) {
	tid, _ := uuid.Parse(tenantID)
	p := utils.NewPagination(page, limit)
	items, total, err := s.repo.ListRuns(ctx, tid, p.Limit, p.Offset())
	if err != nil {
		return nil, nil, err
	}
	res := make([]resource.RunResource, 0, len(items))
	for _, r := range items {
		res = append(res, *toRun(&r))
	}
	pages := int((total + int64(p.Limit) - 1) / int64(p.Limit))
	return res, &response.Meta{Page: p.Page, Limit: p.Limit, Total: total, TotalPages: pages}, nil
}
func (s *RunService) Detail(ctx context.Context, tenantID, runID string) (*resource.RunResource, error) {
	tid, _ := uuid.Parse(tenantID)
	rid, err := uuid.Parse(runID)
	if err != nil {
		return nil, err
	}
	run, err := s.repo.FindRun(ctx, tid, rid)
	if err != nil {
		return nil, err
	}
	return toRun(run), nil
}
func (s *RunService) Steps(ctx context.Context, tenantID, runID string) ([]entity.RunStep, error) {
	tid, _ := uuid.Parse(tenantID)
	rid, err := uuid.Parse(runID)
	if err != nil {
		return nil, err
	}
	return s.repo.Steps(ctx, tid, rid)
}
func (s *RunService) Logs(ctx context.Context, tenantID, runID string) ([]entity.ExecutionLog, error) {
	tid, _ := uuid.Parse(tenantID)
	rid, err := uuid.Parse(runID)
	if err != nil {
		return nil, err
	}
	return s.repo.Logs(ctx, tid, rid)
}
func (s *RunService) Cancel(ctx context.Context, tenantID, runID string) error {
	tid, _ := uuid.Parse(tenantID)
	rid, err := uuid.Parse(runID)
	if err != nil {
		return err
	}
	run, err := s.repo.FindRun(ctx, tid, rid)
	if err != nil {
		return err
	}
	if run.Status != "running" && run.Status != "pending" {
		return errors.New("run is not active")
	}
	s.finish(ctx, run, "cancelled", nil)
	return nil
}
func toRun(run *entity.WorkflowRun) *resource.RunResource {
	return &resource.RunResource{ID: run.ID.String(), WorkflowID: run.WorkflowID.String(), Status: run.Status, TriggerType: run.TriggerType, DurationMs: run.DurationMs, ErrorMessage: run.ErrorMessage}
}
