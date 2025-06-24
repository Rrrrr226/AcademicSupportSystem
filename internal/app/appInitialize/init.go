package appInitialize

import (
	"HelpStudent/internal/app"
)

var (
	apps = make([]app.Module, 0)
)

func GetApps() []app.Module {
	return apps
}
