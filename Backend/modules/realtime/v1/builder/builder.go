package builder

import (
	"go.uber.org/fx"

	"flowforge-api/app"
	"flowforge-api/common/realtime"
	"flowforge-api/modules/realtime/v1/handler"
)

var RealtimeModule = fx.Options(fx.Provide(realtime.NewHub, handler.NewRealtimeHandler), fx.Invoke(app.RealtimeHTTPHandler))
