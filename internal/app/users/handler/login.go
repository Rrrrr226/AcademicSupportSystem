package handler

import (
	"HelpStudent/config"
	"HelpStudent/core/middleware/response"
	"HelpStudent/core/store/rds"
	"HelpStudent/internal/app/users/dto"
	"HelpStudent/internal/app/users/model/thirdPlat"
	"HelpStudent/internal/app/users/service/oauth"
	"HelpStudent/pkg/utils"
	"github.com/flamego/flamego"
	"strings"
)

func HandleThirdPlatLogin(r flamego.Render, c flamego.Context) {
	req := dto.ThirdPlatLoginReq{
		Callback: c.Query("callback"),
		Platform: c.Query("platform"),
		From:     c.Query("from"),
	}

	if req.Callback == "" || req.Platform == "" || req.From == "" {
		response.HTTPFail(r, 401001, "参数错误")
		return
	}

	var callbackExist bool
	for _, oAuth := range config.GetConfig().OAuth {
		if oAuth.CallbackURL == req.Callback || (strings.Index(req.Callback, oAuth.CallbackURL) == 0) {
			callbackExist = true
			break
		}
	}
	if !callbackExist {
		response.HTTPFail(r, 401002, "回调地址不合法")
		return
	}

	platType := thirdPlat.FromString(req.Platform)
	if platType == thirdPlat.NotExists {
		response.HTTPFail(r, 401001, "平台暂不支持")
		return
	}

	urlParams := map[string][]string{}
	if req.From != "" {
		urlParams["from"] = []string{req.From}
	}
	callbackUrl := utils.UrlAppend(req.Callback, urlParams)
	var redirectURL, mark string
	if oauth.PlatformExists(req.Callback, platType) {
		redirectURL, mark = oauth.GetRedirectUrl(req.Callback, platType, callbackUrl)
	} else {
		response.HTTPFail(r, 401001, "平台暂不支持")
		return
	}
	err := cache.Setex(rds.Key("oauth", "mark", mark), "", 60*15)
	if err != nil {
		response.ServiceErr(r, err)
		return
	}

	response.HTTPSuccess(r, dto.ThirdPlatLoginResp{
		URL: redirectURL,
	})
}
