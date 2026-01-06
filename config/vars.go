package config

import (
	"HelpStudent/core/store/pg"
)

type GlobalConfig struct {
	MODE         string     `yaml:"Mode"`
	ProgramName  string     `yaml:"ProgramName"`
	BaseURL      string     `yaml:"BaseURL"`
	AUTHOR       string     `yaml:"Author"`
	Listen       string     `yaml:"Listen"`
	Port         string     `yaml:"Port"`
	MainPostgres pg.OrmConf `yaml:"MainPostgres"`
	Auth         struct {
		Secret string `yaml:"Secret"`
		Issuer string `yaml:"Issuer"`
	} `yaml:"Auth"`
	OAuth   []OAuth `yaml:"OAuth"`
	FastGPT FastGPT `yaml:"FastGPT"`
}

type FastGPT struct {
	BaseURL string `yaml:"BaseURL"`
}

type OAuth struct {
	CallbackURL string `yaml:"CallbackURL"`
	HDUHelp     struct {
		ClientID     string `yaml:"ClientID"`
		ClientSecret string `yaml:"ClientSecret"`
	}
}
