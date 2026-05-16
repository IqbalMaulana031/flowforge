package service

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	"flowforge-api/entity"
	"flowforge-api/modules/execution/v1/engine"
	"flowforge-api/modules/workflow/v1/repository"
	"flowforge-api/resource"
	"flowforge-api/response"
	"flowforge-api/utils"
)

type WorkflowUseCase interface {
	Create(ctx context.Context, tenantID, userID string, req resource.CreateWorkflowRequest) (*resource.WorkflowResource, error)
	List(ctx context.Context, tenantID string, page, limit int, q string) ([]resource.WorkflowResource, *response.Meta, error)
	Detail(ctx context.Context, tenantID, id string) (*resource.WorkflowResource, error)
	Update(ctx context.Context, tenantID, userID, id string, req resource.UpdateWorkflowRequest) (*resource.WorkflowResource, error)
	Delete(ctx context.Context, tenantID, id string) error
	Versions(ctx context.Context, tenantID, id string) ([]entity.WorkflowVersion, error)
	Rollback(ctx context.Context, tenantID, id string, version int) (*resource.WorkflowResource, error)
}
type WorkflowService struct {
	repo repository.WorkflowRepositoryUseCase
}

func NewWorkflowService(repo repository.WorkflowRepositoryUseCase) *WorkflowService {
	return &WorkflowService{repo: repo}
}
func validate(raw []byte) error {
	dag, err := engine.Parse(raw)
	if err != nil {
		return err
	}
	return engine.Validate(dag)
}
func (s *WorkflowService) Create(ctx context.Context, tenantID, userID string, req resource.CreateWorkflowRequest) (*resource.WorkflowResource, error) {
	if err := validate(req.DAGDefinition); err != nil {
		return nil, err
	}
	tid, _ := uuid.Parse(tenantID)
	uid, _ := uuid.Parse(userID)
	wf := &entity.Workflow{TenantID: tid, Name: req.Name, Description: req.Description, IsActive: true, CreatedBy: uid, CurrentVersion: 1}
	ver := &entity.WorkflowVersion{VersionNumber: 1, DAGDefinition: datatypes.JSON(req.DAGDefinition), Changelog: req.Changelog, CreatedBy: uid}
	if err := s.repo.Create(ctx, wf, ver); err != nil {
		return nil, err
	}
	return toWorkflow(wf, req.DAGDefinition), nil
}
func (s *WorkflowService) List(ctx context.Context, tenantID string, page, limit int, q string) ([]resource.WorkflowResource, *response.Meta, error) {
	tid, _ := uuid.Parse(tenantID)
	p := utils.NewPagination(page, limit)
	items, total, err := s.repo.List(ctx, tid, p.Limit, p.Offset(), q)
	if err != nil {
		return nil, nil, err
	}
	res := make([]resource.WorkflowResource, 0, len(items))
	for _, wf := range items {
		res = append(res, *toWorkflow(&wf, nil))
	}
	pages := int((total + int64(p.Limit) - 1) / int64(p.Limit))
	return res, &response.Meta{Page: p.Page, Limit: p.Limit, Total: total, TotalPages: pages}, nil
}
func (s *WorkflowService) Detail(ctx context.Context, tenantID, id string) (*resource.WorkflowResource, error) {
	tid, _ := uuid.Parse(tenantID)
	wid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	wf, ver, err := s.repo.Find(ctx, tid, wid)
	if err != nil {
		return nil, err
	}
	return toWorkflow(wf, []byte(ver.DAGDefinition)), nil
}
func (s *WorkflowService) Update(ctx context.Context, tenantID, userID, id string, req resource.UpdateWorkflowRequest) (*resource.WorkflowResource, error) {
	if err := validate(req.DAGDefinition); err != nil {
		return nil, err
	}
	tid, _ := uuid.Parse(tenantID)
	wid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	uid, _ := uuid.Parse(userID)
	wf, _, err := s.repo.Find(ctx, tid, wid)
	if err != nil {
		return nil, err
	}
	wf.Name = req.Name
	wf.Description = req.Description
	ver := &entity.WorkflowVersion{WorkflowID: wf.ID, TenantID: tid, VersionNumber: wf.CurrentVersion + 1, DAGDefinition: datatypes.JSON(req.DAGDefinition), Changelog: req.Changelog, CreatedBy: uid}
	if err := s.repo.UpdateWithVersion(ctx, wf, ver); err != nil {
		return nil, err
	}
	return toWorkflow(wf, req.DAGDefinition), nil
}
func (s *WorkflowService) Delete(ctx context.Context, tenantID, id string) error {
	tid, _ := uuid.Parse(tenantID)
	wid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, tid, wid)
}
func (s *WorkflowService) Versions(ctx context.Context, tenantID, id string) ([]entity.WorkflowVersion, error) {
	tid, _ := uuid.Parse(tenantID)
	wid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.repo.Versions(ctx, tid, wid)
}
func (s *WorkflowService) Rollback(ctx context.Context, tenantID, id string, version int) (*resource.WorkflowResource, error) {
	tid, _ := uuid.Parse(tenantID)
	wid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	wf, _, err := s.repo.Find(ctx, tid, wid)
	if err != nil {
		return nil, err
	}
	ver, err := s.repo.FindVersion(ctx, tid, wid, version)
	if err != nil {
		return nil, err
	}
	if err := s.repo.SetCurrentVersion(ctx, wf, version); err != nil {
		return nil, err
	}
	return toWorkflow(wf, []byte(ver.DAGDefinition)), nil
}
func toWorkflow(wf *entity.Workflow, dag json.RawMessage) *resource.WorkflowResource {
	return &resource.WorkflowResource{ID: wf.ID.String(), Name: wf.Name, Description: wf.Description, CurrentVersion: wf.CurrentVersion, IsActive: wf.IsActive, DAGDefinition: dag}
}
