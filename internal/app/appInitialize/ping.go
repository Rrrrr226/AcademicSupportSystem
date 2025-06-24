package appInitialize

import "HelpStudent/internal/app/ping"

func init() {
	apps = append(apps, &ping.Ping{Name: "ping module"})
}
