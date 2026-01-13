package handler

import (
	"HelpStudent/core/auth"
	"HelpStudent/core/logx"
	"HelpStudent/core/middleware/response"
	fastgptDAO "HelpStudent/internal/app/fastgpt/dao"
	fastgptModel "HelpStudent/internal/app/fastgpt/model"
	"HelpStudent/internal/app/subject/dao"
	"HelpStudent/internal/app/subject/dto"
	"HelpStudent/internal/app/subject/model"
	userDAO "HelpStudent/internal/app/users/dao"
	userModel "HelpStudent/internal/app/users/model"
	"errors"
	"fmt"
	"strconv"

	"github.com/flamego/flamego"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func GetSubjectLink(r flamego.Render, c flamego.Context, authInfo auth.Info) {
	staffId := c.Param("staff_id")
	if staffId == "" {
		logx.ServiceLogger.Error("staff_id is empty")
		response.ServiceErr(r, "staff_id is empty")
		return
	}

	if authInfo.StaffId != staffId {
		logx.ServiceLogger.Error("staff_id not match")
		response.ServiceErr(r, "staff_id not match")
		return
	}

	// 检查用户是否存在
	var userModel userModel.Users
	if err := userDAO.Users.Where("staff_id = ?", staffId).First(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logx.ServiceLogger.Error("staff_id not found", zap.String("staff_id", staffId))
			response.ServiceErr(r, "staff_id not found")
			return
		}
		response.ServiceErr(r, "查询用户失败：%v", err)
		return
	}

	// 从 user_subjects 表获取用户的科目列表
	subjectNames, err := dao.Subject.GetUserSubjects(staffId)
	if err != nil {
		logx.SystemLogger.Errorw("获取用户科目失败",
			zap.String("staffId", staffId),
			zap.Error(err))
		response.ServiceErr(r, fmt.Sprintf("获取用户科目失败: %v", err))
		return
	}

	logx.SystemLogger.Infow("用户科目数据",
		zap.String("staffId", staffId),
		zap.Any("subjectNames", subjectNames))

	if len(subjectNames) == 0 {
		response.HTTPSuccess(r, dto.GetSubjectResp{})
		return
	}

	logx.SystemLogger.Infow("Attempting to get links by names",
		zap.String("staffId", staffId),
		zap.Any("subjectNames", subjectNames))

	linkMap, err := dao.Subject.GetLinksByNames(subjectNames)
	if err != nil {
		logx.SystemLogger.Errorw("Failed to get links by names from DAO",
			zap.String("staffId", staffId),
			zap.Any("subjectNames", subjectNames),
			zap.Error(err))
		response.ServiceErr(r, fmt.Sprintf("获取科目链接失败: %v", err))
		return
	}

	logx.SystemLogger.Infow("LinkMap result",
		zap.String("staffId", staffId),
		zap.Any("linkMap", linkMap))

	var result []model.Subject
	for _, name := range subjectNames {
		if link, exists := linkMap[name]; exists {
			result = append(result, model.Subject{
				SubjectName: name,
				SubjectLink: link,
			})
		}
	}

	var finalResult []dto.SubjectItem
	appMap := make(map[string]string)
	fastgptAppIdMap := make(map[string]string)
	shareIdMap := make(map[string]string)

	if fastgptDAO.FastgptApp != nil && len(result) > 0 {
		var appNames []string
		for _, s := range result {
			appNames = append(appNames, s.SubjectName)
		}
		var apps []fastgptModel.FastgptApp
		if err := fastgptDAO.FastgptApp.Where("app_name IN ?", appNames).Find(&apps).Error; err == nil {
			for _, app := range apps {
				appMap[app.AppName] = app.ID
				fastgptAppIdMap[app.AppName] = app.AppId
				shareIdMap[app.AppName] = app.ShareId
			}
		} else {
			logx.SystemLogger.Errorw("Failed to fetch fastgpt apps", zap.Error(err))
		}
	}

	for _, s := range result {
		item := dto.SubjectItem{
			Subject: s,
		}
		if appId, ok := appMap[s.SubjectName]; ok {
			item.AppID = appId
		}
		if fastgptAppId, ok := fastgptAppIdMap[s.SubjectName]; ok {
			item.FastgptAppId = fastgptAppId
		}
		if shareId, ok := shareIdMap[s.SubjectName]; ok {
			item.ShareId = shareId
		}
		finalResult = append(finalResult, item)
	}

	logx.SystemLogger.Infow("Final result before response",
		zap.String("staffId", staffId),
		zap.Any("result", finalResult))

	response.HTTPSuccess(r, dto.GetSubjectResp{
		Subjects: finalResult,
	})
}

func AddSubject(r flamego.Render, c flamego.Context, req dto.AddSubjectReq) {
	newSubject := model.Subject{
		SubjectName: req.SubjectName,
		SubjectLink: req.SubjectLink,
	}

	var count int64
	res := dao.Subject.WithContext(c.Request().Context()).
		Model(&model.Subject{}).
		Where("subject_name = ?", req.SubjectName).
		Count(&count)
	if res.Error != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), res.Error)
		response.ServiceErr(r, res.Error)
		return
	}
	if count > 0 {
		response.HTTPFail(r, 401004, "学科已存在")
		return
	}

	err := dao.Subject.Model(&model.Subject{}).Create(&newSubject).Error
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	var subject model.Subject
	dao.Subject.WithContext(c.Request().Context()).
		Model(&model.Subject{}).
		Where("id = ?", newSubject.ID).
		First(&subject)

	response.HTTPSuccess(r, dto.AddSubjectResp{
		SubjectLink: subject.SubjectLink,
		SubjectName: subject.SubjectName,
	})
}

func DeleteSubject(r flamego.Render, c flamego.Context) {
	subjectId := c.Param("subject_id")
	if subjectId == "" {
		response.ServiceErr(r, "subject_id不能为空")
		return
	}

	id, err := strconv.Atoi(subjectId)
	if err != nil {
		response.ServiceErr(r, "无效的subject_id格式")
		return
	}

	err = dao.Subject.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&model.Subject{}).Where("id = ?", uint(id)).Delete(&model.Subject{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("科目不存在，ID: %d", id)
		}
		return nil
	})
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	response.HTTPSuccess(r, "删除成功")
}

func UpdateSubject(r flamego.Render, c flamego.Context, req dto.UpdateSubjectReq) {
	if req.SubjectId == 0 {
		response.ServiceErr(r, "SubjectID不能为空")
		return
	}
	if req.SubjectName == "" && req.SubjectLink == "" {
		response.HTTPFail(r, 400001, "至少需要提供一个更新字段")
		return
	}

	err := dao.Subject.Transaction(func(tx *gorm.DB) error {
		// 检查科目是否存在
		var subject model.Subject
		if err := tx.Where("id = ?", req.SubjectId).First(&subject).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("科目不存在，ID: %d", req.SubjectId)
			}
			return err
		}

		updates := make(map[string]interface{})
		if req.SubjectName != "" {
			var count int64
			if err := tx.Model(&model.Subject{}).
				Where("subject_name = ? AND id <> ?", req.SubjectName, req.SubjectId).
				Count(&count).Error; err != nil {
				return err
			}
			if count > 0 {
				return fmt.Errorf("学科名称已存在")
			}
			updates["subject_name"] = req.SubjectName
		}
		if req.SubjectLink != "" {
			updates["subject_link"] = req.SubjectLink
		}

		// 执行更新
		result := tx.Model(&model.Subject{}).
			Where("id = ?", req.SubjectId).
			Updates(updates)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("未更新任何记录，ID: %d", req.SubjectId)
		}

		return nil
	})

	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	response.HTTPSuccess(r, "更新成功")
}

func GetSubjectList(r flamego.Render, c flamego.Context) {
	//分页
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	var subjects []model.Subject
	var total int64

	// 获取总数
	err = dao.Subject.Model(&model.Subject{}).Count(&total).Error
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	// 获取分页数据
	err = dao.Subject.Model(&model.Subject{}).
		WithContext(c.Request().Context()).
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&subjects).Error

	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	response.HTTPSuccess(r, dto.GetSubjectListResp{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Subjects: subjects,
	})
}

// GetUserSubjectList 获取学生科目关联列表（分页）
func GetUserSubjectList(r flamego.Render, c flamego.Context) {
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")
	staffId := c.Query("staff_id")         // 可选的学号筛选
	subjectName := c.Query("subject_name") // 可选的科目名筛选

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	var userSubjects []model.UserSubject
	var total int64

	query := dao.Subject.Model(&model.UserSubject{}).WithContext(c.Request().Context())

	// 添加筛选条件
	if staffId != "" {
		query = query.Where("staff_id LIKE ?", "%"+staffId+"%")
	}
	if subjectName != "" {
		query = query.Where("subject_name LIKE ?", "%"+subjectName+"%")
	}

	// 获取总数
	err = query.Count(&total).Error
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	// 获取分页数据
	err = query.Limit(pageSize).
		Offset((page - 1) * pageSize).
		Order("created_at DESC").
		Find(&userSubjects).Error

	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	response.HTTPSuccess(r, dto.GetUserSubjectListResp{
		Total:        total,
		Page:         page,
		PageSize:     pageSize,
		UserSubjects: userSubjects,
	})
}

// AddUserSubject 添加学生科目关联
func AddUserSubjectHandler(r flamego.Render, c flamego.Context, req dto.AddUserSubjectReq) {
	if req.StaffId == "" {
		response.HTTPFail(r, 400001, "学号不能为空")
		return
	}
	if req.SubjectName == "" {
		response.HTTPFail(r, 400002, "科目名称不能为空")
		return
	}

	// 查询用户是否存在
	var user userModel.Users
	if err := userDAO.Users.Where("staff_id = ?", req.StaffId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.HTTPFail(r, 404001, "用户不存在")
			return
		}
		response.ServiceErr(r, err)
		return
	}

	// 检查科目是否存在
	var subject model.Subject
	if err := dao.Subject.Where("subject_name = ?", req.SubjectName).First(&subject).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.HTTPFail(r, 404002, "科目不存在")
			return
		}
		response.ServiceErr(r, err)
		return
	}

	// 添加关联
	err := dao.Subject.AddUserSubject(user.ID, req.StaffId, req.SubjectName)
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	response.HTTPSuccess(r, "添加成功")
}

// DeleteUserSubject 删除学生科目关联
func DeleteUserSubjectHandler(r flamego.Render, c flamego.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		response.HTTPFail(r, 400001, "ID不能为空")
		return
	}

	// 删除记录
	result := dao.Subject.Where("id = ?", idStr).Delete(&model.UserSubject{})
	if result.Error != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), result.Error)
		response.ServiceErr(r, result.Error)
		return
	}

	if result.RowsAffected == 0 {
		response.HTTPFail(r, 404001, "记录不存在")
		return
	}

	response.HTTPSuccess(r, "删除成功")
}

// UpdateUserSubject 更新学生科目关联
func UpdateUserSubjectHandler(r flamego.Render, c flamego.Context, req dto.UpdateUserSubjectReq) {
	if req.ID == "" {
		response.HTTPFail(r, 400001, "ID不能为空")
		return
	}

	// 查询记录是否存在
	var userSubject model.UserSubject
	if err := dao.Subject.Where("id = ?", req.ID).First(&userSubject).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.HTTPFail(r, 404001, "记录不存在")
			return
		}
		response.ServiceErr(r, err)
		return
	}

	// 如果要修改学号，检查用户是否存在
	if req.StaffId != "" && req.StaffId != userSubject.StaffId {
		var user userModel.Users
		if err := userDAO.Users.Where("staff_id = ?", req.StaffId).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				response.HTTPFail(r, 404002, "用户不存在")
				return
			}
			response.ServiceErr(r, err)
			return
		}
		userSubject.UserId = user.ID
		userSubject.StaffId = req.StaffId
	}

	// 如果要修改科目名称，检查科目是否存在
	if req.SubjectName != "" && req.SubjectName != userSubject.SubjectName {
		var subject model.Subject
		if err := dao.Subject.Where("subject_name = ?", req.SubjectName).First(&subject).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				response.HTTPFail(r, 404003, "科目不存在")
				return
			}
			response.ServiceErr(r, err)
			return
		}
		userSubject.SubjectName = req.SubjectName
	}

	// 更新记录
	if err := dao.Subject.Save(&userSubject).Error; err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	response.HTTPSuccess(r, "更新成功")
}
