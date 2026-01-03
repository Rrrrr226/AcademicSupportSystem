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
		// 三方登录
		e.Group("/third", func() {
			e.Get("/jump", handler.HandleThirdPlatLogin)
			e.Post("/callback", binding.JSON(dto.ThirdPlatLoginCallbackReq{}), handler.HandleThirdPlatCallback)
		})

		// Token 刷新
		e.Post("/refresh", binding.JSON(dto.RefreshTokenRequest{}), handler.HandleRefreshToken)

		// 用户信息（需要授权）
		e.Get("/info", web.Authorization, handler.HandleGetPersonInfo)
	})

	e.Get("/api/upload/users", handler.HandleUploadUserXLSX, web.Authorization)
}

func UsersGroup(e *flamego.Flame) {

}
