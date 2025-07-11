package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type UserBind struct {
	gorm.Model
	UserId            string `gorm:"type:char(26);not null;index"`
	Type              string `gorm:"type:varchar(16);not null;uniqueIndex:idx_user_bind"`
	UnionId           string `gorm:"type:varchar(128);not null;uniqueIndex:idx_user_bind"`
	Credential        string
	RefreshCredential string
	ExpiredAt         *time.Time
	Attr              datatypes.JSON
}
