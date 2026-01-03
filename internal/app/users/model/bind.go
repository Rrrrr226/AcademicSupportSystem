package model

import (
	"HelpStudent/internal/model"
	"time"

	"gorm.io/datatypes"
)

type UserBind struct {
	model.Base
	UserId            string `gorm:"type:char(26);not null;index"`
	Type              string `gorm:"type:varchar(16);not null;uniqueIndex:idx_user_bind"`
	UnionId           string `gorm:"type:varchar(128);not null;uniqueIndex:idx_user_bind"`
	Credential        string
	RefreshCredential string
	ExpiredAt         *time.Time
	Attr              datatypes.JSON
}
