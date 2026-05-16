package resource

import "encoding/json"

type CreateWorkflowRequest struct {
	Name          string          `json:"name" binding:"required"`
	Description   string          `json:"description"`
	DAGDefinition json.RawMessage `json:"dag_definition" binding:"required"`
	Changelog     string          `json:"changelog"`
}
type UpdateWorkflowRequest struct {
	Name          string          `json:"name" binding:"required"`
	Description   string          `json:"description"`
	DAGDefinition json.RawMessage `json:"dag_definition" binding:"required"`
	Changelog     string          `json:"changelog"`
}
type RollbackWorkflowRequest struct {
	VersionNumber int `json:"version_number" binding:"required"`
}
type WorkflowResource struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	CurrentVersion int             `json:"current_version"`
	IsActive       bool            `json:"is_active"`
	DAGDefinition  json.RawMessage `json:"dag_definition,omitempty"`
}
