package model

// 数据库模型

// 确保 Users 结构体中有以下字段
// 如果已存在，则不需要添加

type Users struct {
	Id       string `gorm:"primary_key;type:char(26)"`
	StaffId  string `gorm:"uniqueIndex;size:19"`
	Name     string
	Username string `gorm:"uniqueIndex;size:50"`
	Password string
	Email    string
	Phone    string

	NeedSubjects   []string `gorm:"-"`                              // 不直接映射到数据库字段
	NeedSubjectsDB string   `gorm:"column:need_subjects;type:text"` // 实际数据库存储字段
}
