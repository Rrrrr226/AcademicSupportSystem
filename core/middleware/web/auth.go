package web

import (
	"HelpStudent/core/auth"
	"HelpStudent/core/middleware/response"
	"github.com/flamego/flamego"
	"strings"
)

func Authorization(c flamego.Context, r flamego.Render) {
	token := c.Request().Header.Get("Authorization")
	if token == "" || strings.Index(token, "Bearer") != 0 {
		response.UnAuthorization(r)
		return
	}
	token = strings.Replace(token, "Bearer ", "", 1)
	entity, err := auth.ParseToken(token)
	if err != nil {
		response.UnAuthorization(r)
		return
	}
	c.Map(entity.Info)
}
