package router

import (
	"HelpStudent/core/middleware/response"
	"HelpStudent/core/middleware/web"
	"HelpStudent/internal/app/users/dto"
	"HelpStudent/internal/app/users/handler"

	"errors"
	"github.com/flamego/binding"
	"github.com/flamego/flamego"
)

func AppUsersInit(e *flamego.Flame) {
	e.Get("/users/v1", func(r flamego.Render) {
		response.HTTPSuccess(r, map[string]any{
			"message": "users Init Success",
		})
	})

	e.Get("/users/v1/err", func(r flamego.Render) {
		response.HTTPFail(r, 500000, "users Init test error", errors.New("this is err"))
	})

	e.Group("/user/v1", func() {
		e.Group("/third", func() {
			e.Get("/jump", handler.HandleThirdPlatLogin)
			e.Post("/callback", binding.JSON(dto.ThirdPlatLoginCallbackReq{}), handler.HandleThirdPlatCallback)
		})

		e.Group("/direct", func() {
			e.Post("/register", binding.JSON(dto.RegisterRequest{}), handler.HandleRegister)
			e.Post("/login", binding.JSON(dto.LoginRequest{}), handler.HandleLogin)
		})

		e.Post("/refresh", binding.JSON(dto.RefreshTokenRequest{}), handler.HandleRefreshToken)
	})

	e.Group("/user/v1", func() {
		e.Get("", handler.HandleGetPersonInfo)
	}, web.Authorization)

	e.Get("/api/upload/users", handler.HandleUploadUserXLSX)
}

func UsersGroup(e *flamego.Flame) {

}
