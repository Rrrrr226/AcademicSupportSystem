package ping

import (
	"HelpStudent/core/kernel"
	"HelpStudent/internal/app"
	"HelpStudent/internal/app/ping/router"
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
