package builder

import (
	"go.uber.org/fx"

	"flowforge-api/app"
	"flowforge-api/modules/ai/v1/handler"
	"flowforge-api/modules/ai/v1/service"
)

var AIModule = fx.Options(fx.Provide(fx.Annotate(service.NewAIService, fx.As(new(service.AIUseCase))), handler.NewAIHandler), fx.Invoke(app.AIHTTPHandler))
