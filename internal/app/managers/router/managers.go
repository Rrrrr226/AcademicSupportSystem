package router

import (
	"HelpStudent/core/middleware/response"
	"HelpStudent/core/middleware/web"
	"HelpStudent/internal/app/managers/dto"
	handler "HelpStudent/internal/app/managers/handler/v1"
	"errors"

	"github.com/flamego/binding"
	"github.com/flamego/flamego"
)

func AppManagersInit(e *flamego.Flame) {
	e.Get("/managers/v1", func(r flamego.Render) {
		response.HTTPSuccess(r, map[string]any{
			"message": "managers Init Success",
		})
	})

	e.Get("/managers/v1/err", func(r flamego.Render) {
		response.HTTPFail(r, 500000, "managers Init test error", errors.New("this is err"))
	})

	e.Group("/managers", func() {
		// 管理员管理相关接口（需要登录）
		e.Get("/list", handler.HandleGetManagerList)
		e.Post("/add", binding.JSON(dto.AddManagerRequest{}), handler.HandleAddManager)
		e.Post("/delete", binding.JSON(dto.DeleteManagerRequest{}), handler.HandleDeleteManager)

		// 学生科目导入接口（Excel上传）
		e.Post("/import/students", handler.HandleImportStudentSubjectsExcel)

		// 下载导入模板
		e.Get("/import/template", handler.HandleDownloadTemplate)
	}, web.Authorization)

}
