package subject

import (
	"HelpStudent/core/kernel"
	"HelpStudent/core/logx"
	"HelpStudent/internal/app"
	"HelpStudent/internal/app/subject/dao"
	"HelpStudent/internal/app/subject/router"
	"context"
	"os"
	"sync"

	"go.uber.org/zap"
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

func (p *Subject) Init(engine *kernel.Engine) error {
	// 只初始化 subject DAO，users DAO 由 users 模块负责初始化
	if err := dao.InitPG(engine.MainPG.GetOrm()); err != nil {
		logx.SystemLogger.Errorw("科目DAO初始化失败", zap.Error(err))
		os.Exit(1)
	}
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
