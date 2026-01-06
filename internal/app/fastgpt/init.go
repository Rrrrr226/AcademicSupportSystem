package fastgpt

import (
	"HelpStudent/core/kernel"
	"HelpStudent/internal/app"
	"HelpStudent/internal/app/fastgpt/router"
	"context"
	"sync"
)

type (
	Fastgpt struct {
		Name string
		app.UnimplementedModule
	}
)

func (p *Fastgpt) Info() string {
	return p.Name
}

func (p *Fastgpt) PreInit(engine *kernel.Engine) error {
	return nil
}

func (p *Fastgpt) Init(*kernel.Engine) error {
	return nil
}

func (p *Fastgpt) PostInit(*kernel.Engine) error {
	return nil
}

func (p *Fastgpt) Load(engine *kernel.Engine) error {
	// 加载flamego api
	router.AppFastgptInit(engine.Fg)
	return nil
}

func (p *Fastgpt) Start(engine *kernel.Engine) error {
	return nil
}

func (p *Fastgpt) Stop(wg *sync.WaitGroup, ctx context.Context) error {
	defer wg.Done()
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (p *Fastgpt) OnConfigChange() func(*kernel.Engine) error {
	return func(engine *kernel.Engine) error {

		return nil
	}
}
