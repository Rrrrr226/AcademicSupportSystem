package router

import (
	"HelpStudent/core/middleware/response"
	"HelpStudent/core/middleware/web"
	"HelpStudent/internal/app/subject/dto"
	"HelpStudent/internal/app/subject/handler/v1"
	"errors"

	"github.com/flamego/binding"
	"github.com/flamego/flamego"
)

func AppSubjectInit(e *flamego.Flame) {
	e.Get("/subject/v1", func(r flamego.Render) {
		response.HTTPSuccess(r, map[string]any{
			"message": "subject Init Success",
		})
	})

	e.Get("/subject/v1/err", func(r flamego.Render) {
		response.HTTPFail(r, 500000, "subject Init test error", errors.New("this is err"))
	})

	e.Get("/subject/get/links/{staff_id}", web.Authorization, handler.GetSubjectLink)

	e.Group("/subject/v1", func() {
		e.Post("/add", binding.JSON(dto.AddSubjectReq{}), handler.AddSubject)
		e.Get("/list", handler.GetSubjectList)
		e.Delete("/delete/{subject_id}", handler.DeleteSubject)
		e.Post("/update", binding.JSON(dto.UpdateSubjectReq{}), handler.UpdateSubject)

		// 学生科目关联管理
		e.Get("/user-subjects", handler.GetUserSubjectList)
		e.Post("/user-subjects/add", binding.JSON(dto.AddUserSubjectReq{}), handler.AddUserSubjectHandler)
		e.Delete("/user-subjects/delete/{id}", handler.DeleteUserSubjectHandler)
		e.Post("/user-subjects/update", binding.JSON(dto.UpdateUserSubjectReq{}), handler.UpdateUserSubjectHandler)
	}, web.Authorization)
}

func SubjectGroup(e *flamego.Flame) {}
