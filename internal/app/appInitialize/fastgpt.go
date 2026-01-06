package appInitialize

import "HelpStudent/internal/app/fastgpt"

func init() {
	apps = append(apps, &fastgpt.Fastgpt{Name: "Fastgpt module"})
}
