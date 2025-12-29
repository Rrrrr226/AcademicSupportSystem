package web

import (
	"HelpStudent/core/auth"
	"HelpStudent/core/logx"
	"HelpStudent/core/middleware/response"
	"fmt"
	"github.com/flamego/flamego"
	"strings"
)

func Authorization(c flamego.Context, r flamego.Render) {
	token := c.Request().Header.Get("Authorization")
	if token == "" || strings.Index(token, "Bearer") != 0 {
		response.UnAuthorization(r)
		return
	}
	fmt.Println("1,token:", token)
	token = strings.Replace(token, "Bearer ", "", 1)
	entity, err := auth.ParseToken(token)
	if err != nil {
		fmt.Println("2,token:", token)
		response.UnAuthorization(r)
		return
	}
	logx.SystemLogger.Infof("Authorization: Parsed auth.Info: Uid=%s, StaffId=%s", entity.Info.Uid, entity.Info.StaffId)
	c.Map(entity.Info)
}
