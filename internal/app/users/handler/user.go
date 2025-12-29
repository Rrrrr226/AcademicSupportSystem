package handler

import (
	"HelpStudent/core/auth"
	"HelpStudent/core/logx"
	"HelpStudent/core/middleware/response"
	"HelpStudent/internal/app/permission/service/rbac"
	"HelpStudent/internal/app/users/dao"
	"HelpStudent/internal/app/users/dto"
	"HelpStudent/internal/app/users/model"

	"github.com/flamego/flamego"
)

func HandleGetPersonInfo(r flamego.Render, c flamego.Context, auth auth.Info) {
	logx.SystemLogger.Infof("HandleGetPersonInfo: auth.Uid = %s, auth.StaffId = %s", auth.Uid, auth.StaffId)

	var user model.Users
	result := dao.Users.WithContext(c.Request().Context()).Model(&model.Users{}).
		Where("id = ?", auth.Uid).Find(&user)

	userInfo := dto.UserInfoResponse{
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
	response.HTTPSuccess(r, userInfo)
}
