package sls

import (
	"HelpStudent/core/logx"
	"fmt"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"go.uber.org/zap/zapcore"
	"sync"
)

var (
	topicCache  sync.Map
	sourceCache sync.Map
)

type Sink struct {
	producer *producer.Producer
	mux      sync.Mutex
	logLevel zapcore.Level

	project  string
	logstore string
	source   string
	topic    string
}

func New(conf logx.SlsSinkConf, opts ...Option) (*Sink, error) {
	u, err := ParseURL(conf.Url)
	if err != nil {
		return nil, fmt.Errorf("invalid aliyunsls-url: %w", err)
	}
	cfg := producer.GetDefaultProducerConfig()
	cfg.Endpoint = u.Endpoint
	cfg.CredentialsProvider = sls.NewStaticCredentialsProvider(u.AccessKeyID, u.AccessKeySecret, "")
	_producer := producer.InitProducer(cfg)
	sink := &Sink{
		producer: _producer,
		project:  u.Project,
		logstore: u.LogStore,
		topic:    "operation_log",
		source:   "default",
	}
	for _, opt := range opts {
		opt(sink)
	}
	return sink, nil
}

func WithSource(source string) Option {
	return func(s *Sink) {
		s.source = source
	}
}

func WithTopic(topic string) Option {
	return func(s *Sink) {
		s.topic = topic
	}
}

func (s *Sink) Open() {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.producer == nil {
		panic(fmt.Errorf("nil producer"))
	}
	s.producer.Start()
}

func (s *Sink) Close() {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.producer == nil {
		return
	}
	s.producer.SafeClose()
	s.producer = nil
}

func (s *Sink) Write(ent zapcore.Entry, fields []zapcore.Field) {
	//if s.logLevel >= ent.Level {
	//	return
	//}
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.producer == nil {
		return
	}
	l := fields2logs(ent, fields)
	_err := s.producer.SendLog(s.project, s.logstore, s.topic, s.source, l)
	if _err != nil {
		fmt.Printf("send log failed, err: %v\n", _err)
	}
}

func (s *Sink) WithSource(source string) logx.Sink {
	value, ok := sourceCache.Load(source)
	if ok {
		return value.(logx.Sink)
	}
	_sink := &Sink{
		producer: s.producer,
		project:  s.project,
		logstore: s.logstore,
		topic:    s.topic,
		logLevel: s.logLevel,
		source:   source,
	}
	sourceCache.Store(source, _sink)
	return _sink
}

func (s *Sink) WithTopic(topic string) logx.Sink {
	value, ok := topicCache.Load(topic)
	if ok {
		return value.(logx.Sink)
	}
	_sink := &Sink{
		producer: s.producer,
		project:  s.project,
		logstore: s.logstore,
		topic:    topic,
		source:   s.source,
		logLevel: s.logLevel,
	}
	topicCache.Store(topic, _sink)
	return _sink
}

func (s *Sink) GetProducer() *producer.Producer {
	return s.producer
}

func (s *Sink) GetLogStore() string {
	return s.logstore
}
