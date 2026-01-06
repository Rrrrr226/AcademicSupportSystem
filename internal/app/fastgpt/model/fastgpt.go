package model

import "gorm.io/gorm"

// 数据库模型

// FastgptApp FastGPT 应用配置
type FastgptApp struct {
	gorm.Model
	AppID       string `gorm:"uniqueIndex;not null;type:varchar(100);comment:应用ID"`
	AppName     string `gorm:"not null;type:varchar(200);comment:应用名称"`
	APIKey      string `gorm:"not null;type:varchar(200);comment:FastGPT API密钥"`
	Description string `gorm:"type:text;comment:应用描述"`
	Status      int    `gorm:"default:1;comment:状态(1:启用,0:禁用)"`
	CreatedBy   string `gorm:"type:varchar(50);comment:创建者"`
}

func (FastgptApp) TableName() string {
	return "fastgpt_apps"
}
