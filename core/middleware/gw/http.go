package gw

import (
	"HelpStudent/core/kernel"
	"HelpStudent/core/logx"
	"github.com/flamego/flamego"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"time"
)

func RequestLog() flamego.Handler {
	return func(c flamego.Context, r *http.Request) {

		// 开始时间
		startTime := time.Now()

		otel.GetTextMapPropagator().Inject(r.Context(), propagation.HeaderCarrier(c.ResponseWriter().Header()))
		// 处理请求
		c.Next()

		// source来源是当前服务器的ip，topic是access_log
		logx.SystemLogger.WithTopic("access_log").WithSource(kernel.Kernel.CurrentIpList[0]).Infow("request log",
			zap.Field{Key: "path", Type: zapcore.StringType, String: c.Request().RequestURI},
			zap.Field{Key: "method", Type: zapcore.StringType, String: c.Request().Method},
			zap.Field{Key: "ip", Type: zapcore.StringType, String: c.RemoteAddr()},
			zap.Field{Key: "status", Type: zapcore.Int64Type, Integer: int64(c.ResponseWriter().Status())},
			zap.Field{Key: "duration", Type: zapcore.StringType, String: time.Now().Sub(startTime).String()},
			zap.Field{Key: "user-agent", Type: zapcore.StringType, String: c.Request().UserAgent()},
		)
	}
}
