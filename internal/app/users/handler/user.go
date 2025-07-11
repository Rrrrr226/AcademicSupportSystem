package handler

import (
	"HelpStudent/core/auth"
	"HelpStudent/core/middleware/response"
	"HelpStudent/internal/app/permission/service/rbac"
	"HelpStudent/internal/app/users/dao"
	"HelpStudent/internal/app/users/dto"
	"HelpStudent/internal/app/users/model"
	"github.com/flamego/flamego"
)

func HandleGetPersonInfo(r flamego.Render, c flamego.Context, auth auth.Info) {
	var user dto.UserInfoResponse
	dao.Users.WithContext(c.Request().Context()).Model(&model.Users{}).
		Where("id = ?", auth.Uid).Find(&user)
	user.Permissions = rbac.GetStaffSystemPermission(auth.StaffId)
	response.HTTPSuccess(r, user)
}
