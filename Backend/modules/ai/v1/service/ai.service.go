package service

import (
	"context"
	"errors"

	"flowforge-api/config"
	"flowforge-api/resource"
)

type AIUseCase interface {
	GenerateWorkflow(context.Context, resource.GenerateWorkflowRequest) (resource.AIResource, error)
	AnalyzeFailure(context.Context, string, resource.AnalyzeFailureRequest) (resource.AIResource, error)
}
type AIService struct{ cfg *config.Config }

func NewAIService(cfg *config.Config) *AIService { return &AIService{cfg: cfg} }
func (s *AIService) provider() (string, error) {
	if s.cfg.AI.Provider == "anthropic" && s.cfg.AI.AnthropicAPIKey != "" {
		return "anthropic", nil
	}
	if s.cfg.AI.Provider == "openai" && s.cfg.AI.OpenAIAPIKey != "" {
		return "openai", nil
	}
	return "", errors.New("AI provider is not configured")
}
func (s *AIService) GenerateWorkflow(ctx context.Context, req resource.GenerateWorkflowRequest) (resource.AIResource, error) {
	p, err := s.provider()
	if err != nil {
		return resource.AIResource{}, err
	}
	return resource.AIResource{Provider: p, Result: "AI workflow generation is wired; configure SDK call for production. Requirements: " + req.Requirements}, nil
}
func (s *AIService) AnalyzeFailure(ctx context.Context, runID string, req resource.AnalyzeFailureRequest) (resource.AIResource, error) {
	p, err := s.provider()
	if err != nil {
		return resource.AIResource{}, err
	}
	return resource.AIResource{Provider: p, Result: "AI failure analysis is wired for run " + runID}, nil
}
