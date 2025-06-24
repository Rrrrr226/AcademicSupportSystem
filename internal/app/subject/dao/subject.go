package dao

import (
	"HelpStudent/internal/app/subject/model"
	"gorm.io/gorm"
)

type subject struct {
	*gorm.DB
}

func (u *subject) Init(db *gorm.DB) (err error) {
	u.DB = db
	return db.AutoMigrate(&model.Subject{})
}

func (d *subject) GetLinksByNames(names []string) (map[string]string, error) {
	if len(names) == 0 {
		return map[string]string{}, nil
	}

	var subjects []model.Subject
	if err := d.DB.Where("subject_name IN ?", names).Find(&subjects).Error; err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, s := range subjects {
		result[s.SubjectName] = s.SubjectLink
	}

	return result, nil
}
