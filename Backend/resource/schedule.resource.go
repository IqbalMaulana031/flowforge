package resource

type ScheduleRequest struct {
	WorkflowID     string `json:"workflow_id" binding:"required"`
	CronExpression string `json:"cron_expression" binding:"required"`
	IsActive       *bool  `json:"is_active"`
}
