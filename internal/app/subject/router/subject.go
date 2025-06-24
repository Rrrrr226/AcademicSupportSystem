package router

import (
	"HelpStudent/core/middleware/response"
	"HelpStudent/internal/app/subject/handler/v1"
	"errors"
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

	e.Get("/subject/get/links", handler.GetSubjectLink)
}

func SubjectGroup(e *flamego.Flame) {}
