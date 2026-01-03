package managers

import (
	"HelpStudent/core/kernel"
	"HelpStudent/core/logx"
	"HelpStudent/internal/app"
	"HelpStudent/internal/app/managers/dao"
	"HelpStudent/internal/app/managers/router"
	"context"
	"os"
	"sync"

	"go.uber.org/zap"
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

func (p *Managers) Init(engine *kernel.Engine) error {
	if err := dao.InitPG(engine.MainPG.GetOrm()); err != nil {
		logx.SystemLogger.Errorw("管理员DAO初始化失败", zap.Error(err))
		os.Exit(1)
	}
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
