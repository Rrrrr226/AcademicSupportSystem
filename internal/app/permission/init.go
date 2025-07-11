package permission

import (
	"HelpStudent/core/kernel"
	"HelpStudent/internal/app"
	"HelpStudent/internal/app/permission/dao"
	"HelpStudent/internal/app/permission/router"
	"HelpStudent/internal/app/permission/service/rbac"
	"context"
	"github.com/pkg/errors"
	"sync"
)

type (
	Permission struct {
		Name string
		app.UnimplementedModule
	}
)

var (
	ErrEmptyDatabase = errors.New("database pointer is nil")
)

func (p *Permission) Info() string {
	return p.Name
}

func (p *Permission) PreInit(engine *kernel.Engine) error {
	if engine.MainPG == nil {
		return ErrEmptyDatabase
	}
	err := dao.InitPG(engine.MainPG.DB)
	if err != nil {
		return err
	}
	err = rbac.Init(engine.MainPG.DB)
	if err != nil {
		return err
	}
	return nil
}

func (p *Permission) Init(*kernel.Engine) error {
	return nil
}

func (p *Permission) PostInit(*kernel.Engine) error {
	return nil
}

func (p *Permission) Load(engine *kernel.Engine) error {
	// 加载flamego api
	router.AppPermissionInit(engine.Fg)
	return nil
}

func (p *Permission) Start(engine *kernel.Engine) error {
	return nil
}

func (p *Permission) Stop(wg *sync.WaitGroup, ctx context.Context) error {
	defer wg.Done()
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (p *Permission) OnConfigChange() func(*kernel.Engine) error {
	return func(engine *kernel.Engine) error {

		return nil
	}
}
