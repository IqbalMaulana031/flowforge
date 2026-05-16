package builder

import (
	"go.uber.org/fx"

	"flowforge-api/app"
	"flowforge-api/modules/auth/v1/handler"
	"flowforge-api/modules/auth/v1/repository"
	"flowforge-api/modules/auth/v1/service"
)

var AuthModule = fx.Options(fx.Provide(fx.Annotate(repository.NewAuthRepository, fx.As(new(repository.AuthRepositoryUseCase))), fx.Annotate(service.NewAuthService, fx.As(new(service.AuthUseCase))), handler.NewAuthHandler), fx.Invoke(app.AuthHTTPHandler))
