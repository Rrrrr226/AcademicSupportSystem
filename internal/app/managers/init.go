package managers

import (
	"context"
	"HelpStudent/core/kernel"
	"HelpStudent/internal/app"
	"HelpStudent/internal/app/managers/router"
	"sync"
)

type (
	Managers struct {
		Name string
		app.UnimplementedModule
	}
)

func (p *Managers) Info() string {
	return p.Name
}

func (p *Managers) PreInit(engine *kernel.Engine) error {
	return nil
}

func (p *Managers) Init(*kernel.Engine) error {
	return nil
}

func (p *Managers) PostInit(*kernel.Engine) error {
	return nil
}

func (p *Managers) Load(engine *kernel.Engine) error {
	// 加载flamego api
	router.AppManagersInit(engine.Fg)
	return nil
}

func (p *Managers) Start(engine *kernel.Engine) error {
	return nil
}

func (p *Managers) Stop(wg *sync.WaitGroup, ctx context.Context) error {
	defer wg.Done()
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (p *Managers) OnConfigChange() func(*kernel.Engine) error {
	return func(engine *kernel.Engine) error {

		return nil
	}
}
