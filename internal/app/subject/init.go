package subject

import (
	"context"
	"HelpStudent/core/kernel"
	"HelpStudent/internal/app"
	"HelpStudent/internal/app/subject/router"
	"sync"
)

type (
	Subject struct {
		Name string
		app.UnimplementedModule
	}
)

func (p *Subject) Info() string {
	return p.Name
}

func (p *Subject) PreInit(engine *kernel.Engine) error {
	return nil
}

func (p *Subject) Init(*kernel.Engine) error {
	return nil
}

func (p *Subject) PostInit(*kernel.Engine) error {
	return nil
}

func (p *Subject) Load(engine *kernel.Engine) error {
	// 加载flamego api
	router.AppSubjectInit(engine.Fg)
	return nil
}

func (p *Subject) Start(engine *kernel.Engine) error {
	return nil
}

func (p *Subject) Stop(wg *sync.WaitGroup, ctx context.Context) error {
	defer wg.Done()
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (p *Subject) OnConfigChange() func(*kernel.Engine) error {
	return func(engine *kernel.Engine) error {

		return nil
	}
}
