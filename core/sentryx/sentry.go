package sentryx

import (
	"HelpStudent/core/logx"
	"github.com/getsentry/sentry-go"
)

func NewSentry(r Config) {
	if !r.Available() {
		return
	}
	if err := sentry.Init(sentry.ClientOptions{
		Dsn: r.Dsn,
	}); err != nil {
		logx.SystemLogger.Error(err)
	}
}
