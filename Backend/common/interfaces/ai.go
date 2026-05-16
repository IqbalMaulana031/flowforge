package interfaces

import "context"

type AIProvider interface {
	GenerateWorkflow(context.Context, string) (string, error)
	AnalyzeFailure(context.Context, string) (string, error)
}
