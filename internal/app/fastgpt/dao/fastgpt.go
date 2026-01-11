package dao

import (
	"HelpStudent/internal/app/fastgpt/model"
	"context"
	"errors"

	"gorm.io/gorm"
)

type fastgpt struct {
	*gorm.DB
}

var FastgptApp *fastgpt

func (u *fastgpt) Init(db *gorm.DB) (err error) {
	u.DB = db
	FastgptApp = u
	return db.AutoMigrate(&model.FastgptApp{})
}

// CreateApp 创建应用
func (u *fastgpt) CreateApp(app *model.FastgptApp) error {
	return u.Create(app).Error
}

// GetAppByID 根据 AppID 获取应用
func (u *fastgpt) GetAppByID(appID string) (*model.FastgptApp, error) {
	var app model.FastgptApp
	err := u.Where("app_id = ?", appID).First(&app).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("应用不存在或已禁用")
		}
		return nil, err
	}
	return &app, nil
}

// GetAppByPrimaryID 根据主键ID获取应用
func (u *fastgpt) GetAppByPrimaryID(ctx context.Context, id string) (*model.FastgptApp, error) {
	var app model.FastgptApp
	err := u.Model(&model.FastgptApp{}).WithContext(ctx).Where("id = ?", id).First(&app).Error
	return &app, err
}

// GetAllApps 获取所有应用列表
func (u *fastgpt) GetAllApps(ctx context.Context, offset, limit int) ([]model.FastgptApp, int64, error) {
	var apps []model.FastgptApp
	var total int64

	// 统计总数
	if err := u.Model(&model.FastgptApp{}).WithContext(ctx).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询列表
	err := u.WithContext(ctx).Offset(offset).Limit(limit).Order("created_at ASC").Find(&apps).Error
	return apps, total, err
}

// UpdateApp 更新应用
func (u *fastgpt) UpdateApp(id string, updates map[string]interface{}) error {
	return u.Model(&model.FastgptApp{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteApp 删除应用（软删除）
func (u *fastgpt) DeleteApp(ctx context.Context, id string) error {
	return u.Model(&model.FastgptApp{}).WithContext(ctx).Where("id = ?", id).Delete(&model.FastgptApp{}).Error
}

// CheckAppNameExists 检查 AppName 是否已存在
func (u *fastgpt) CheckAppNameExists(appName string) (bool, error) {
	var count int64
	query := u.Model(&model.FastgptApp{}).Where("app_name = ?", appName)
	err := query.Count(&count).Error
	return count > 0, err
}

func (u *fastgpt) SubjectsExist(subjectNames []string) ([]string, error) {
	if len(subjectNames) == 0 {
		return nil, nil
	}

	var existingSubjects []model.FastgptApp
	if err := u.Where("app_name IN ?", subjectNames).Find(&existingSubjects).Error; err != nil {
		return nil, err
	}

	existingMap := make(map[string]bool)
	for _, s := range existingSubjects {
		existingMap[s.AppName] = true
	}

	var notExist []string
	for _, name := range subjectNames {
		if !existingMap[name] {
			notExist = append(notExist, name)
		}
	}

	return notExist, nil
}
