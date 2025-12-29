package model

// 数据库模型

type Managers struct {
	Id       string `gorm:"primary_key;type:char(26)"`
	StaffId  string `gorm:"uniqueIndex;size:19"`
	Name     string
	Username string `gorm:"size:50"`
	Password string
}
