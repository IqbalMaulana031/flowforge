package resource

import "encoding/json"

type TriggerRunRequest struct {
	WorkflowID string         `json:"workflow_id" binding:"required"`
	Payload    map[string]any `json:"payload"`
}
type RunResource struct {
	ID           string `json:"id"`
	WorkflowID   string `json:"workflow_id"`
	Status       string `json:"status"`
	TriggerType  string `json:"trigger_type"`
	DurationMs   int64  `json:"duration_ms"`
	ErrorMessage string `json:"error_message"`
}
type RunStepResource struct {
	ID            string          `json:"id"`
	StepID        string          `json:"step_id"`
	StepName      string          `json:"step_name"`
	StepType      string          `json:"step_type"`
	Status        string          `json:"status"`
	AttemptNumber int             `json:"attempt_number"`
	DurationMs    int64           `json:"duration_ms"`
	ErrorMessage  string          `json:"error_message"`
	Output        json.RawMessage `json:"output,omitempty"`
}
