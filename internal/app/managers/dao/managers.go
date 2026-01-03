package dao

import (
	"HelpStudent/internal/app/managers/model"

	"gorm.io/gorm"
)

type managers struct {
	*gorm.DB
}

func (u *managers) Init(db *gorm.DB) (err error) {
	u.DB = db
	return db.AutoMigrate(&model.Managers{})
}

// GetAllManagers 获取所有管理员列表
func (m *managers) GetAllManagers() ([]model.Managers, int64, error) {
	var managers []model.Managers
	var count int64

	if err := m.Model(&model.Managers{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := m.Find(&managers).Error; err != nil {
		return nil, 0, err
	}

	return managers, count, nil
}

// GetManagerById 根据ID获取管理员
func (m *managers) GetManagerById(id string) (*model.Managers, error) {
	var manager model.Managers
	if err := m.Where("id = ?", id).First(&manager).Error; err != nil {
		return nil, err
	}
	return &manager, nil
}

// GetManagerByUsername 根据用户名获取管理员
func (m *managers) GetManagerByStaffID(staffId string) (*model.Managers, error) {
	var manager model.Managers
	if err := m.Where("staff_id = ?", staffId).First(&manager).Error; err != nil {
		return nil, err
	}
	return &manager, nil
}

// DeleteManagerById 删除管理员
func (m *managers) DeleteManagerById(id string) error {
	result := m.Where("id = ?", id).Delete(&model.Managers{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (m *managers) DeleteManagetByStaffId(staffId string) error {
	result := m.Where("staff_id = ?", staffId).Delete(&model.Managers{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// IsManager 根据 StaffId 检查是否是管理员
func (m *managers) IsManager(staffId string) bool {
	var count int64
	m.Model(&model.Managers{}).Where("staff_id = ?", staffId).Count(&count)
	return count > 0
}
