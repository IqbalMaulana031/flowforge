package builder

import (
	"go.uber.org/fx"

	"flowforge-api/app"
	"flowforge-api/modules/schedule/v1/handler"
	"flowforge-api/modules/schedule/v1/repository"
	"flowforge-api/modules/schedule/v1/service"
)

var ScheduleModule = fx.Options(fx.Provide(fx.Annotate(repository.NewScheduleRepository, fx.As(new(repository.ScheduleRepositoryUseCase))), fx.Annotate(service.NewScheduleService, fx.As(new(service.ScheduleUseCase))), handler.NewScheduleHandler), fx.Invoke(app.ScheduleHTTPHandler))
