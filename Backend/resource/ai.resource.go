package resource

type GenerateWorkflowRequest struct {
	Requirements string `json:"requirements" binding:"required"`
}
type AnalyzeFailureRequest struct {
	Question string `json:"question"`
}
type AIResource struct {
	Provider      string `json:"provider"`
	Result        string `json:"result"`
	DAGDefinition any    `json:"dag_definition,omitempty"`
}
