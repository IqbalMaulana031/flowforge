package builder

import (
	"go.uber.org/fx"

	"flowforge-api/app"
	"flowforge-api/modules/workflow/v1/handler"
	"flowforge-api/modules/workflow/v1/repository"
	"flowforge-api/modules/workflow/v1/service"
)

var WorkflowModule = fx.Options(fx.Provide(fx.Annotate(repository.NewWorkflowRepository, fx.As(new(repository.WorkflowRepositoryUseCase))), fx.Annotate(service.NewWorkflowService, fx.As(new(service.WorkflowUseCase))), handler.NewWorkflowHandler), fx.Invoke(app.WorkflowHTTPHandler))
