package v1

import (
	"HelpStudent/core/auth"
	"HelpStudent/core/logx"
	"HelpStudent/core/middleware/response"
	"HelpStudent/internal/app/fastgpt/dao"
	"HelpStudent/internal/app/fastgpt/dto"
	"HelpStudent/internal/app/fastgpt/model"
	dao2 "HelpStudent/internal/app/managers/dao"
	"errors"

	"github.com/flamego/binding"
	"github.com/flamego/flamego"
	"gorm.io/gorm"
)

// HandleCreateApp 创建应用
func HandleCreateApp(c flamego.Context, r flamego.Render, req dto.CreateAppRequest, errs binding.Errors, authInfo auth.Info) {
	if errs != nil {
		response.InValidParam(r, errs)
		return
	}

	// 检查是否是管理员
	if !dao2.Managers.IsManager(authInfo.Uid) {
		response.HTTPFail(r, 400013, "非管理员用户无法创建应用")
		return
	}

	// 检查 AppID 是否已存在
	exists, err := dao.FastgptApp.CheckAppIDExists(req.AppID, 0)
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}
	if exists {
		response.HTTPFail(r, 400014, "应用ID已存在")
		return
	}

	// 创建应用
	app := &model.FastgptApp{
		AppID:       req.AppID,
		AppName:     req.AppName,
		APIKey:      req.APIKey,
		Description: req.Description,
		Status:      1,
		CreatedBy:   authInfo.Uid,
	}

	if err := dao.FastgptApp.CreateApp(app); err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	response.HTTPSuccess(r, dto.CreateAppResponse{ID: app.ID})
}

// HandleGetAppList 获取应用列表
func HandleGetAppList(c flamego.Context, r flamego.Render, req dto.GetAppListRequest, errs binding.Errors, authInfo auth.Info) {
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	// 检查是否是管理员
	if !dao2.Managers.IsManager(authInfo.Uid) {
		response.HTTPFail(r, 400013, "非管理员用户无法创建应用")
		return
	}

	apps, total, err := dao.FastgptApp.GetAllApps(req.Offset, req.Limit)
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	// 转换为 DTO
	var appItems []dto.AppItem
	for _, app := range apps {
		appItems = append(appItems, dto.AppItem{
			ID:          app.ID,
			AppID:       app.AppID,
			AppName:     app.AppName,
			APIKey:      app.APIKey,
			Description: app.Description,
			Status:      app.Status,
			CreatedBy:   app.CreatedBy,
			CreatedAt:   app.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   app.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	response.HTTPSuccess(r, dto.AppListResponse{
		Apps:  appItems,
		Total: total,
	})
}

// HandleUpdateApp 更新应用
func HandleUpdateApp(c flamego.Context, r flamego.Render, req dto.UpdateAppRequest, errs binding.Errors, authInfo auth.Info) {
	if errs != nil {
		response.InValidParam(r, errs)
		return
	}

	// 检查是否是管理员
	if !dao2.Managers.IsManager(authInfo.Uid) {
		response.HTTPFail(r, 400013, "非管理员用户无法创建应用")
		return
	}

	// 检查应用是否存在
	_, err := dao.FastgptApp.GetAppByPrimaryID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.HTTPFail(r, 404001, "应用不存在")
			return
		}
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	// 构建更新数据
	updates := make(map[string]interface{})
	if req.AppName != "" {
		updates["app_name"] = req.AppName
	}
	if req.APIKey != "" {
		updates["api_key"] = req.APIKey
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) == 0 {
		response.HTTPFail(r, 400015, "没有需要更新的字段")
		return
	}

	// 更新应用
	if err := dao.FastgptApp.UpdateApp(req.ID, updates); err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	response.HTTPSuccess(r, nil)
}

// HandleDeleteApp 删除应用
func HandleDeleteApp(c flamego.Context, r flamego.Render, req dto.DeleteAppRequest, errs binding.Errors, authInfo auth.Info) {
	if errs != nil {
		response.InValidParam(r, errs)
		return
	}

	// 检查是否是管理员
	if !dao2.Managers.IsManager(authInfo.Uid) {
		response.HTTPFail(r, 400013, "非管理员用户无法创建应用")
		return
	}

	// 检查应用是否存在
	_, err := dao.FastgptApp.GetAppByPrimaryID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.HTTPFail(r, 404001, "应用不存在")
			return
		}
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	// 删除应用
	if err := dao.FastgptApp.DeleteApp(req.ID); err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	response.HTTPSuccess(r, nil)
}
