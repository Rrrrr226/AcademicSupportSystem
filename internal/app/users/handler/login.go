package handler

import (
	"HelpStudent/config"
	"HelpStudent/core/auth"
	"HelpStudent/core/logx"
	"HelpStudent/core/middleware/response"
	"HelpStudent/core/store/rds"
	"HelpStudent/internal/app/users/dao"
	"HelpStudent/internal/app/users/dto"
	"HelpStudent/internal/app/users/model"
	"HelpStudent/internal/app/users/model/thirdPlat"
	"HelpStudent/internal/app/users/service/oauth"
	"HelpStudent/pkg/utils"
	"errors"
	"github.com/flamego/binding"
	"github.com/flamego/flamego"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"strings"
	"time"
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

func HandleThirdPlatCallback(r flamego.Render, c flamego.Context, req dto.ThirdPlatLoginCallbackReq, errs binding.Errors) {
	if errs != nil {
		response.InValidParam(r, errs)
		return
	}

	if req.Code == "" && req.State == "" {
		response.HTTPFail(r, 401001, "code和state不能同时为空")
		return
	}

	states := strings.Split(req.State, "_")
	// 校验state格式 platform_mark
	if len(states) != 2 {
		response.HTTPFail(r, 401001, "state格式错误")
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

	platType := thirdPlat.FromString(states[0])
	if platType == thirdPlat.NotExists {
		response.HTTPFail(r, 401001, "平台暂不支持")
		return
	}

	mark := states[1]
	if mark == "" {
		response.HTTPFail(r, 401001, "mark不能为空")
		return
	}

	if exist, err := cache.ExistsCtx(c.Request().Context(), rds.Key("oauth", "mark", mark)); !exist || err != nil {
		response.HTTPFail(r, 401001, "mark已失效")
		return
	}
	_, err := cache.DelCtx(c.Request().Context(), rds.Key("oauth", "mark", mark))
	if err != nil {
		logx.ServiceLogger.CtxError(c.Request().Context(), err)
	}

	var (
		uid  string
		attr datatypes.JSON
	)
	if oauth.PlatformExists(req.Callback, platType) {
		uid, attr, err = oauth.Validate(req.Callback, platType, req.Code, req.State)
	} else {
		response.HTTPFail(r, 401001, "平台暂不支持")
		return
	}
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}
	b := &model.UserBind{Type: platType.String(), UnionId: uid}
	if result := dao.Users.WithContext(c.Request().Context()).Where(b).
		Find(b); result.Error != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), result.Error)
		response.ServiceErr(r, result.Error)
		return
	} else if result.RowsAffected == 0 {
		// 新用户
		b.Attr = attr
		user := &model.Users{
			StaffId: oauth.GetStaffId(*b),
			Name:    oauth.GetUserName(*b),
		}
		b.Attr = attr
		err = dao.Users.CreateWithBind(c.Request().Context(), user, b)
		if err != nil {
			logx.SystemLogger.CtxError(c.Request().Context(), err)
			response.ServiceErr(r, err)
			return
		}
	} else {
		// 老用户更新用户信息
		b.Attr = attr
		dao.Users.Model(b).Update("attr", attr)
	}

	token, err := auth.GenToken(auth.Info{Uid: b.UserId, StaffId: oauth.GetStaffId(*b), Name: oauth.GetUserName(*b)})
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}
	refreshToken, err := auth.GenToken(auth.Info{Uid: b.UserId, StaffId: oauth.GetStaffId(*b), Name: oauth.GetUserName(*b), IsRefreshToken: true}, auth.RefreshTokenExpireIn)
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}
	response.HTTPSuccess(r, dto.GeneralLoginResponse{
		AccessToken:          token,
		AccessTokenExpireIn:  int64(auth.AccessTokenExpireIn / time.Second),
		RefreshToken:         refreshToken,
		RefreshTokenExpireIn: int64(auth.RefreshTokenExpireIn / time.Second),
	})
}

func HandleRefreshToken(r flamego.Render, req dto.RefreshTokenRequest) {
	entity, err := auth.ParseToken(req.RefreshToken)
	if err != nil || !entity.Info.IsRefreshToken || entity.Info.Name == "" || entity.Info.StaffId == "" {
		response.UnAuthorization(r)
		return
	}

	token, err := auth.GenToken(auth.Info{Uid: entity.Info.Uid, StaffId: entity.Info.StaffId, Name: entity.Info.Name})
	if err != nil {
		response.ServiceErr(r, err)
		return
	}
	response.HTTPSuccess(r, dto.RefreshTokenResponse{
		AccessToken:         token,
		AccessTokenExpireIn: int64(auth.AccessTokenExpireIn / time.Second),
		RefreshToken:        req.RefreshToken,
	})
}

// 修改 HandleLogin 函数
func HandleLogin(r flamego.Render, c flamego.Context, req dto.LoginRequest, errs binding.Errors) {
	if errs != nil {
		response.InValidParam(r, errs)
		return
	}

	// 查询用户
	var user model.Users
	result := dao.Users.WithContext(c.Request().Context()).
		Where("username = ?", req.Username).
		First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			response.HTTPFail(r, 401003, "用户名或密码错误")
			return
		}
		logx.SystemLogger.CtxError(c.Request().Context(), result.Error)
		response.ServiceErr(r, result.Error)
		return
	}

	// 验证密码
	if !utils.VerifyPassword(user.Password, req.Password) {
		response.HTTPFail(r, 401003, "用户名或密码错误")
		return
	}

	// 生成令牌
	token, err := auth.GenToken(auth.Info{Uid: user.Id, StaffId: user.StaffId, Name: user.Name})
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	refreshToken, err := auth.GenToken(auth.Info{
		Uid:            user.Id,
		StaffId:        user.StaffId,
		Name:           user.Name,
		IsRefreshToken: true,
	}, auth.RefreshTokenExpireIn)

	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	// 返回登录成功响应
	response.HTTPSuccess(r, dto.GeneralLoginResponse{
		AccessToken:          token,
		AccessTokenExpireIn:  int64(auth.AccessTokenExpireIn / time.Second),
		RefreshToken:         refreshToken,
		RefreshTokenExpireIn: int64(auth.RefreshTokenExpireIn / time.Second),
	})
}

func HandleRegister(r flamego.Render, c flamego.Context, req dto.RegisterRequest, errs binding.Errors) {
	if errs != nil {
		response.InValidParam(r, errs)
		return
	}

	if dao.Users == nil || dao.Users.DB == nil {
		response.ServiceErr(r, errors.New("数据库连接未初始化"))
		return
	}

	// 检查用户名是否已存在
	var count int64
	result := dao.Users.WithContext(c.Request().Context()).
		Model(&model.Users{}).
		Where("username = ?", req.Username).
		Count(&count)

	if result.Error != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), result.Error)
		response.ServiceErr(r, result.Error)
		return
	}

	if count > 0 {
		response.HTTPFail(r, 401004, "用户名已存在")
		return
	}

	// 使用事务创建用户
	err := dao.Users.Transaction(func(tx *gorm.DB) error {
		// 生成用户ID
		userId := utils.GenUUID()

		// 加密密码
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return err
		}

		// 创建用户
		user := &model.Users{
			Id:       userId,
			StaffId:  req.Username, // 使用用户名作为StaffId，可根据需求调整
			Name:     req.Name,
			Username: req.Username,
			Password: hashedPassword,
			Email:    req.Email,
			Phone:    req.Phone,
		}

		if err := tx.Create(user).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	// 查询创建的用户
	var user model.Users
	dao.Users.WithContext(c.Request().Context()).
		Where("username = ?", req.Username).
		First(&user)

	// 返回注册成功响应
	response.HTTPSuccess(r, dto.RegisterResponse{
		UserId:   user.Id,
		Username: user.Username,
		Name:     user.Name,
	})
}
