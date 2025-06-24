package appInitialize

import "HelpStudent/internal/app/subject"

func init() {
	apps = append(apps, &subject.Subject{Name: "Subject module"})
}
