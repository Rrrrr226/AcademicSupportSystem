package dao

import (
	"HelpStudent/internal/app/fastgpt/model"
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
	err := u.Where("app_id = ? AND status = 1", appID).First(&app).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("应用不存在或已禁用")
		}
		return nil, err
	}
	return &app, nil
}

// GetAppByPrimaryID 根据主键ID获取应用
func (u *fastgpt) GetAppByPrimaryID(id uint) (*model.FastgptApp, error) {
	var app model.FastgptApp
	err := u.First(&app, id).Error
	return &app, err
}

// GetAllApps 获取所有应用列表
func (u *fastgpt) GetAllApps(offset, limit int) ([]model.FastgptApp, int64, error) {
	var apps []model.FastgptApp
	var total int64

	// 统计总数
	if err := u.Model(&model.FastgptApp{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询列表
	err := u.Offset(offset).Limit(limit).Order("created_at DESC").Find(&apps).Error
	return apps, total, err
}

// UpdateApp 更新应用
func (u *fastgpt) UpdateApp(id uint, updates map[string]interface{}) error {
	return u.Model(&model.FastgptApp{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteApp 删除应用（软删除）
func (u *fastgpt) DeleteApp(id uint) error {
	return u.Delete(&model.FastgptApp{}, id).Error
}

// CheckAppIDExists 检查 AppID 是否已存在
func (u *fastgpt) CheckAppIDExists(appID string, excludeID uint) (bool, error) {
	var count int64
	query := u.Model(&model.FastgptApp{}).Where("app_id = ?", appID)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}
