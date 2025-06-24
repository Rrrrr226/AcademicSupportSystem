package config

import (
	"HelpStudent/core/fileServer"
	"HelpStudent/core/logx"
	"HelpStudent/core/sentryx"
	"HelpStudent/core/store/mysql"
	"HelpStudent/core/store/pg"
	"HelpStudent/core/store/rds"
	"HelpStudent/core/tracex"
)

type GlobalConfig struct {
	MODE           string         `yaml:"Mode"`
	ProgramName    string         `yaml:"ProgramName"`
	BaseURL        string         `yaml:"BaseURL"`
	AUTHOR         string         `yaml:"Author"`
	Listen         string         `yaml:"Listen"`
	Port           string         `yaml:"Port"`
	AdminStaffID   []string       `yaml:"AdminStaffID"`
	MainPostgres   pg.OrmConf     `yaml:"MainPostgres"`
	MainV3Postgres pg.OrmConf     `yaml:"MainV3Postgres"`
	SKLMysql       mysql.OrmConf  `yaml:"SKLMysql"`
	MainCache      rds.RedisConf  `yaml:"MainCache"`
	Sentry         sentryx.Config `yaml:"Sentry"`
	Log            logx.LogConf   `yaml:"Log"`
	Trace          tracex.Config  `yaml:"Trace"`
	Auth           struct {
		Secret string `yaml:"Secret"`
		Issuer string `yaml:"Issuer"`
	} `yaml:"Auth"`
	OAuth      []OAuth             `yaml:"OAuth"`
	FileServer []fileServer.Config `yaml:"FileServer"`
}

type OAuth struct {
	CallbackURL string `yaml:"CallbackURL"`
	HDUHelp     struct {
		ClientID     string `yaml:"ClientID"`
		ClientSecret string `yaml:"ClientSecret"`
	}
}
