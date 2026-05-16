package builder

import (
	"go.uber.org/fx"

	"flowforge-api/app"
	"flowforge-api/modules/health/v1/handler"
	"flowforge-api/modules/health/v1/service"
)

var HealthModule = fx.Options(
	fx.Provide(
		fx.Annotate(
			service.NewHealthService,
			fx.As(new(service.HealthUseCase)),
		),
		handler.NewHealthHandler,
	),
	fx.Invoke(app.HealthHTTPHandler),
)
