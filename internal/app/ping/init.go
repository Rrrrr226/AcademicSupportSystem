package ping

import (
	"HelpStudent/core/kernel"
	pingV1 "HelpStudent/gen/proto/ping/v1"
	"HelpStudent/internal/app"
	"HelpStudent/internal/app/ping/router"
	"HelpStudent/internal/app/ping/service/gw"
	"context"
	"sync"
)

type (
	Ping struct {
		Name string
		app.UnimplementedModule
	}
)

func (p *Ping) Info() string {
	return p.Name
}

func (p *Ping) PreInit(engine *kernel.Engine) error {
	return nil
}

func (p *Ping) Init(*kernel.Engine) error {
	return nil
}

func (p *Ping) PostInit(*kernel.Engine) error {
	return nil
}

func (p *Ping) Load(engine *kernel.Engine) error {
	// 加载flamego api
	router.AppPingInit(engine.Fg)
	// 加载grpc gw
	pingV1.RegisterPingServiceServer(engine.Grpc, &gw.S{})
	_err := pingV1.RegisterPingServiceHandler(engine.Ctx, engine.Mux, engine.Conn)
	if _err != nil {
		return _err
	}
	return nil
}

func (p *Ping) Start(engine *kernel.Engine) error {
	return nil
}

func (p *Ping) Stop(wg *sync.WaitGroup, ctx context.Context) error {
	defer wg.Done()
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (p *Ping) OnConfigChange() func(*kernel.Engine) error {
	return func(engine *kernel.Engine) error {

		return nil
	}
}
