package kernel

import (
	"HelpStudent/config"
	"HelpStudent/core/logx/sls"
	"HelpStudent/core/store/pg"
	"context"
	"net/http"

	"github.com/flamego/flamego"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type (
	Engine struct {
		MainPG     *pg.Orm
		Fg         *flamego.Flame
		Mux        *runtime.ServeMux
		HttpServer *http.Server
		SlsClient  *sls.Client

		Ctx            context.Context
		Cancel         context.CancelFunc
		ConfigListener []func(*config.GlobalConfig)
	}
)

var Kernel *Engine
