package logx

import (
	"fmt"
	"go.uber.org/zap/zapcore"
)

type Mock struct {
}

func (m *Mock) Open() {
	fmt.Println("Mock Open")
}

func (m *Mock) Close() {
	fmt.Println("Mock Close")
}

func (m *Mock) Write(ent zapcore.Entry, fields []zapcore.Field) {
	fmt.Println("Mock Write", ent, fields)
}

func (m *Mock) WithSource(source string) Sink {
	fmt.Println("Mock WithSource", source)
	return m
}

func (m *Mock) WithTopic(topic string) Sink {
	fmt.Println("Mock WithTopic", topic)
	return m
}

func NewMock() *Mock {
	return &Mock{}
}
