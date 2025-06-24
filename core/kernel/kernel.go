package kernel

import (
	"HelpStudent/config"
	"HelpStudent/core/logx/sls"
	"HelpStudent/core/store/mysql"
	"HelpStudent/core/store/pg"
	"HelpStudent/core/store/rds"
	"context"
	"github.com/flamego/flamego"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"net/http"
)

type (
	Engine struct {
		MainPG     *pg.Orm
		MainV3PG   *pg.Orm
		SKLMySQL   *mysql.Orm
		MainCache  *rds.Redis
		Fg         *flamego.Flame
		Grpc       *grpc.Server
		Conn       *grpc.ClientConn
		Mux        *runtime.ServeMux
		HttpServer *http.Server
		SlsClient  *sls.Client

		CurrentIpList []string

		Ctx    context.Context
		Cancel context.CancelFunc

		ConfigListener []func(*config.GlobalConfig)
	}
)

var Kernel *Engine
