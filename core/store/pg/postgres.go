package pg

import (
	"HelpStudent/core/logx"
	"fmt"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

type (
	Orm struct {
		Host     string
		Port     string
		User     string
		Pass     string
		Database string
		Schema   string
		Debug    bool
		Trace    bool

		Conf *gorm.Config
		*gorm.DB
	}
	Option func(r *Orm)
)

func (r *Orm) GetOrm() *gorm.DB {
	return r.DB
}

func (r *Orm) OrmConnectionUpdate(conf OrmConf) *Orm {
	orm, err := NewPostgresOrm(conf)
	if err != nil {
		return r
	}
	return orm
}

func MustNewPGOrm(conf OrmConf, opts ...Option) *Orm {
	orm, err := NewPostgresOrm(conf, opts...)
	if err != nil {
		logx.SystemLogger.Errorw("fail to load pg orm", zap.Field{Key: "err", Type: zapcore.StringType, String: err.Error()})
		os.Exit(1)
	}
	return orm
}

func NewPostgresOrm(conf OrmConf, opts ...Option) (*Orm, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}
	opts = append([]Option{WithAddr(conf.Host, conf.Port)}, opts...)
	opts = append([]Option{WithAuth(conf.User, conf.Pass)}, opts...)
	opts = append([]Option{WithTrace(conf.Trace)}, opts...)
	opts = append([]Option{WithDBName(conf.Database)}, opts...)
	opts = append([]Option{WithDBSchema(conf.Schema)}, opts...)
	opts = append([]Option{WithDebug(conf.Debug)}, opts...)
	return newOrm(opts...)
}

func WithGormConf(conf *gorm.Config) Option {
	return func(r *Orm) {
		r.Conf = conf
	}
}

func WithTrace(trace bool) Option {
	return func(r *Orm) {
		r.Trace = trace
	}
}

func WithDebug(debug bool) Option {
	return func(r *Orm) {
		r.Debug = debug
	}
}

func WithAddr(host, port string) Option {
	return func(r *Orm) {
		r.Host = host
		r.Port = port
	}
}

func WithAuth(user, pass string) Option {
	return func(r *Orm) {
		r.Pass = pass
		r.User = user
	}
}

func WithDBName(db string) Option {
	return func(r *Orm) {
		r.Database = db
	}
}

func WithDBSchema(schema string) Option {
	return func(r *Orm) {
		r.Schema = schema
	}
}

func newOrm(opts ...Option) (*Orm, error) {
	m := &Orm{}
	for _, opt := range opts {
		opt(m)
	}
	conf := m.Conf
	if conf == nil {
		conf = &gorm.Config{}
	}
	var dsn = fmt.Sprintf("host=%s user=%s password=%s "+
		"dbname=%s port=%s sslmode=disable search_path=%s TimeZone=Asia/Shanghai",
		m.Host, m.User, m.Pass, m.Database, m.Port, m.Schema)
	db, err := gorm.Open(postgres.Open(dsn), conf)
	if m.Trace {
		if _err := db.Use(otelgorm.NewPlugin()); _err != nil {
			logx.SystemLogger.Error(_err)
		}
	}
	if m.Debug {
		db = db.Debug()
	}
	m.DB = db
	return m, err
}
