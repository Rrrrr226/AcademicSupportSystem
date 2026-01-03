package model

import (
	"HelpStudent/internal/model"
)

type Users struct {
	model.Base
	StaffId string `gorm:"uniqueIndex;size:19"`
	Name    string
}
