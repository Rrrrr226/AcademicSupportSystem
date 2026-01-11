package model

import (
	"HelpStudent/internal/model"

	"gorm.io/gorm"
)

// 数据库模型

// FastgptApp FastGPT 应用配置
type FastgptApp struct {
	model.Base
	DeletedAt   gorm.DeletedAt `gorm:"uniqueIndex:idx_app_name"`
	AppName     string         `gorm:"uniqueIndex:idx_app_name;not null;type:varchar(200);comment:应用名称"`
	APIKey      string         `gorm:"not null;type:varchar(200);comment:FastGPT API密钥"`
	Description string         `gorm:"type:text;comment:应用描述"`
	CreatedBy   string         `gorm:"type:varchar(50);comment:创建者"`
}
