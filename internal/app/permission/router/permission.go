package router

import (
	"HelpStudent/core/middleware/web"
	"HelpStudent/internal/app/permission/dto"
	"HelpStudent/internal/app/permission/handler"
	"github.com/flamego/binding"
	"github.com/flamego/flamego"
)

func AppPermissionInit(e *flamego.Flame) {
	e.Group("/permission/v1", func() {
		e.Get("", handler.HandleManagerList)                                                       // 项目管理员列表
		e.Post("", binding.JSON(dto.AddProjectManagerRequest{}), handler.HandleAddManager)         // 添加项目管理员
		e.Delete("", binding.JSON(dto.RemoveProjectManagerRequest{}), handler.HandleRemoveManager) // 删除项目管理员
	}, web.Authorization)
}
