package appInitialize

import "HelpStudent/internal/app/managers"

func init() {
	apps = append(apps, &managers.Managers{Name: "Managers module"})
}