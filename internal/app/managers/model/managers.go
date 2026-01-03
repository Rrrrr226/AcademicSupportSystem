package model

import "HelpStudent/internal/model"

// 数据库模型

type Managers struct {
	model.Base
	StaffId string `gorm:"uniqueIndex;size:19"`
}
