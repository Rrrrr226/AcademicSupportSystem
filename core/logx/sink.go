package logx

import "go.uber.org/zap/zapcore"

type Sink interface {
	Open()
	Close()
	Write(ent zapcore.Entry, fields []zapcore.Field)
	WithSource(source string) Sink
	WithTopic(topic string) Sink
}
