package dao

import (
	"HelpStudent/core/store/mysql"
	"HelpStudent/internal/app/ping/model"
)

var (
	Ping *mysql.Orm
)

func AutoMigrate() error {
	return Ping.AutoMigrate(&model.Ping{})
}
