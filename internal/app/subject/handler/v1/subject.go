package handler

import (
	"HelpStudent/core/logx"
	"HelpStudent/internal/app/subject/dao"
	userDAO "HelpStudent/internal/app/users/dao"
	user "HelpStudent/internal/app/users/model"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/flamego/flamego"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func GetSubjectLink(r flamego.Render, c flamego.Context) ([]string, error) {
	staffId := c.Query("staff_id")
	if staffId == "" {
		return nil, fmt.Errorf("staff_id参数不能为空")
	}

	var userModel user.Users
	if err := userDAO.Users.Where("staff_id = ?", staffId).First(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	var subjectNames []string
	if userModel.NeedSubjectsDB != "" {
		if err := json.Unmarshal([]byte(userModel.NeedSubjectsDB), &subjectNames); err != nil {
			return nil, fmt.Errorf("解析NeedSubjectsDB失败: %v", err)
		}
	}

	if len(subjectNames) == 0 {
		return []string{}, nil
	}

	linkMap, err := dao.Subject.GetLinksByNames(subjectNames)
	if err != nil {
		logx.SystemLogger.Errorw("Failed to get links by names from DAO",
			zap.String("staffId", staffId),
			zap.Any("subjectNames", subjectNames),
			zap.Error(err))
		return nil, fmt.Errorf("获取科目链接失败: %v", err)
	}

	fmt.Println("NeedSubjectsDB:", userModel.NeedSubjectsDB)

	var result []string
	for _, name := range subjectNames {
		if link, exists := linkMap[name]; exists {
			result = append(result, link)
		}
	}

	return result, nil
}
