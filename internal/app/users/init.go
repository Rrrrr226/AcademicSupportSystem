package users

import (
	"HelpStudent/core/kernel"
	"HelpStudent/core/logx"
	"HelpStudent/internal/app"
	users "HelpStudent/internal/app/users/dao"
	"HelpStudent/internal/app/users/router"
	"HelpStudent/internal/app/users/service/oauth"
	"context"
	"os"
	"sync"

	"go.uber.org/zap"
)

type (
	Users struct {
		Name string
		app.UnimplementedModule
	}
)

func (p *Users) Info() string {
	return p.Name
}

func (p *Users) PreInit(engine *kernel.Engine) error {
	oauth.Init()
	return nil
}

func (p *Users) Init(engine *kernel.Engine) error {
	if err := users.InitPG(engine.MainPG.GetOrm()); err != nil {
		logx.SystemLogger.Errorw("用户DAO初始化失败", zap.Error(err))
		os.Exit(1)
	}
	return nil
}

func (p *Users) PostInit(*kernel.Engine) error {
	return nil
}

func (p *Users) Load(engine *kernel.Engine) error {
	// 加载flamego api
	router.AppUsersInit(engine.Fg)
	return nil
}

func (p *Users) Start(engine *kernel.Engine) error {
	return nil
}

func (p *Users) Stop(wg *sync.WaitGroup, ctx context.Context) error {
	defer wg.Done()
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (p *Users) OnConfigChange() func(*kernel.Engine) error {
	return func(engine *kernel.Engine) error {

		return nil
	}
}
