package engine

import "encoding/json"

type DAGDefinition struct {
	Steps     []Step `json:"steps"`
	Edges     []Edge `json:"edges"`
	TimeoutMs int    `json:"timeout_ms"`
}
type Step struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Type        StepType       `json:"type"`
	Config      map[string]any `json:"config"`
	DependsOn   []string       `json:"depends_on"`
	RetryConfig *RetryConfig   `json:"retry_config,omitempty"`
	TimeoutMs   int            `json:"timeout_ms"`
}
type Edge struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Condition string `json:"condition,omitempty"`
}
type StepType string

const (
	StepTypeHTTPCall StepType = "http_call"
	StepTypeScript   StepType = "script"
	StepTypeDelay    StepType = "delay"
	StepTypeBranch   StepType = "conditional_branch"
)

type RetryConfig struct {
	MaxAttempts    int     `json:"max_attempts"`
	InitialDelayMs int     `json:"initial_delay_ms"`
	MaxDelayMs     int     `json:"max_delay_ms"`
	BackoffFactor  float64 `json:"backoff_factor"`
}

func Parse(data []byte) (*DAGDefinition, error) {
	var d DAGDefinition
	if err := json.Unmarshal(data, &d); err != nil {
		return nil, err
	}
	return &d, nil
}
