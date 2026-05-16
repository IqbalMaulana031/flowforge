package builder

import (
	"go.uber.org/fx"

	"flowforge-api/app"
	"flowforge-api/modules/execution/v1/engine"
	"flowforge-api/modules/execution/v1/handler"
	"flowforge-api/modules/execution/v1/repository"
	"flowforge-api/modules/execution/v1/service"
)

var ExecutionModule = fx.Options(fx.Provide(engine.NewExecutor, fx.Annotate(repository.NewRunRepository, fx.As(new(repository.RunRepositoryUseCase))), fx.Annotate(service.NewRunService, fx.As(new(service.RunUseCase))), handler.NewRunHandler), fx.Invoke(app.RunHTTPHandler))
