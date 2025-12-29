package handler

import (
	"HelpStudent/core/auth"
	"HelpStudent/core/logx"
	"HelpStudent/core/middleware/response"
	"HelpStudent/internal/app/managers/dao"
	"HelpStudent/internal/app/managers/dto"
	"HelpStudent/internal/app/managers/model"
	"HelpStudent/internal/app/permission/service/rbac"
	"HelpStudent/pkg/utils"
	"errors"
	"github.com/flamego/binding"
	"github.com/flamego/flamego"
	"gorm.io/gorm"
	"time"
)

func HandleManagerLogin(r flamego.Render, c flamego.Context, req dto.LoginRequest, errs binding.Errors) {
	if errs != nil {
		response.InValidParam(r, errs)
		return
	}

	// 查询用户
	var user model.Managers
	result := dao.Managers.WithContext(c.Request().Context()).
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

func HandleManagerRegister(r flamego.Render, c flamego.Context, req dto.RegisterRequest, errs binding.Errors) {
	if errs != nil {
		response.InValidParam(r, errs)
		return
	}

	if dao.Managers == nil || dao.Managers.DB == nil {
		response.ServiceErr(r, errors.New("数据库连接未初始化"))
		return
	}

	// 检查用户名是否已存在
	var count int64
	result := dao.Managers.WithContext(c.Request().Context()).
		Model(&model.Managers{}).
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
	err := dao.Managers.Transaction(func(tx *gorm.DB) error {
		// 生成用户ID
		userId := utils.GenUUID()

		// 加密密码
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return err
		}

		// 创建用户
		user := &model.Managers{
			Id:       userId,
			StaffId:  req.Username, // 使用用户名作为StaffId，可根据需求调整
			Name:     req.Name,
			Username: req.Username,
			Password: hashedPassword,
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
	var user model.Managers
	dao.Managers.WithContext(c.Request().Context()).
		Where("username = ?", req.Username).
		First(&user)

	// 返回注册成功响应
	response.HTTPSuccess(r, dto.RegisterResponse{
		UserId:   user.Id,
		Username: user.Username,
		Name:     user.Name,
	})
}

func HandleManagerModify(r flamego.Render, c flamego.Context, req dto.ModifyRequest, errs binding.Errors) {
	if errs != nil {
		response.InValidParam(r, errs)
		return
	}

	if dao.Managers == nil || dao.Managers.DB == nil {
		response.ServiceErr(r, errors.New("数据库连接未初始化"))
		return
	}

	if req.UserId == "" {
		response.UnAuthorization(r)
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}

	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			logx.SystemLogger.CtxError(c.Request().Context(), err)
			response.ServiceErr(r, err)
			return
		}
		updates["password"] = hashedPassword
	}

	// 如果没有需要更新的字段
	if len(updates) == 0 {
		response.HTTPFail(r, 401001, "无更新内容")
		return
	}

	// 更新用户信息
	result := dao.Managers.WithContext(c.Request().Context()).
		Model(&model.Managers{}).
		Where("staff_id = ?", req.UserId).
		Updates(updates)

	if result.Error != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), result.Error)
		response.ServiceErr(r, result.Error)
		return
	}

	if result.RowsAffected == 0 {
		response.HTTPFail(r, 401005, "用户不存在")
		return
	}

	response.HTTPSuccess(r, nil)
}

func HandleGetManagerInfo(r flamego.Render, c flamego.Context, auth auth.Info) {
	logx.SystemLogger.Infof("HandleGetPersonInfo: auth.Uid = %s, auth.StaffId = %s", auth.Uid, auth.StaffId)

	var user model.Managers
	result := dao.Managers.WithContext(c.Request().Context()).Model(&model.Managers{}).
		Where("id = ?", auth.Uid).Find(&user)

	userInfo := dto.ManagerInfoResponse{
		Id:          user.Id,
		StaffId:     user.StaffId,
		Name:        user.Name,
		Permissions: nil,
	}

	if result.Error != nil {
		logx.SystemLogger.Errorf("HandleGetPersonInfo: 数据库链接错误: %v", result.Error)
		response.HTTPFail(r, 500, "Failed to get user info", result.Error)
		return
	} else {
		logx.SystemLogger.Info("HandleGetPersonInfo: success database connection")
	}

	if result.RowsAffected == 0 {
		logx.SystemLogger.Warnf("HandleGetPersonInfo: 找不到用户Uid %s", auth.Uid)
		response.HTTPFail(r, 404, "User not found", nil)
		return
	} else {
		logx.SystemLogger.Info("HandleGetPersonInfo: success find user info")
	}
	userInfo.Permissions = rbac.GetStaffSystemPermission(auth.StaffId)
	response.HTTPSuccess(r, user)
}
