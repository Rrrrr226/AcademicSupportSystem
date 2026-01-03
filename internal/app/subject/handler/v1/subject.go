package handler

import (
	"HelpStudent/core/logx"
	"HelpStudent/core/middleware/response"
	"HelpStudent/internal/app/subject/dao"
	"HelpStudent/internal/app/subject/dto"
	"HelpStudent/internal/app/subject/model"
	userDAO "HelpStudent/internal/app/users/dao"
	user "HelpStudent/internal/app/users/model"
	"errors"
	"fmt"
	"strconv"

	"github.com/flamego/flamego"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func GetSubjectLink(r flamego.Render, c flamego.Context) {
	staffId := c.Param("staff_id")
	if staffId == "" {
		logx.ServiceLogger.Error("staff_id is empty")
		response.ServiceErr(r, "staff_id is empty")
		return
	}

	// 检查用户是否存在
	var userModel user.Users
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

	logx.SystemLogger.Infow("Final result before response",
		zap.String("staffId", staffId),
		zap.Any("result", result))

	response.HTTPSuccess(r, dto.GetSubjectResp{
		Subjects: result,
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
