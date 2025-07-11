package dao

import "gorm.io/gorm"

type skl struct {
	*gorm.DB
}

func (s *skl) Init(db *gorm.DB) (err error) {
	s.DB = db
	return db.AutoMigrate()
}
