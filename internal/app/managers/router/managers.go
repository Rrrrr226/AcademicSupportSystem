package router

import (
	"HelpStudent/core/middleware/response"
	"HelpStudent/core/middleware/web"
	"HelpStudent/internal/app/managers/dto"
	"HelpStudent/internal/app/managers/handler/v1"
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
		e.Post("/register", binding.JSON(dto.RegisterRequest{}), handler.HandleManagerRegister)
		e.Post("/login", binding.JSON(dto.LoginRequest{}), handler.HandleManagerLogin)
		e.Post("/modify", binding.JSON(dto.ModifyRequest{}), handler.HandleManagerModify, web.Authorization)

		e.Get("/info", handler.HandleGetManagerInfo, web.Authorization)
	})

}

func ManagersGroup(e *flamego.Flame) {}
