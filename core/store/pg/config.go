package pg

import "errors"

var (
	ErrEmptyHost     = errors.New("empty postgres host")
	ErrEmptyPort     = errors.New("empty postgres port")
	ErrEmptyUser     = errors.New("empty postgres user")
	ErrEmptyPass     = errors.New("empty postgres pass")
	ErrEmptyDB       = errors.New("empty postgres database")
	ErrorEmptySchema = errors.New("empty postgres schema")
)

type (
	// A OrmConf is a mysql config.
	OrmConf struct {
		Host     string `yaml:"Host"`
		Port     string `yaml:"Port"`
		User     string `yaml:"User"`
		Pass     string `yaml:"Pass"`
		Database string `yaml:"Database"`
		Schema   string `yaml:"Schema"`
		Debug    bool   `yaml:"Debug"`
		Trace    bool   `yaml:"Trace"`
	}
)

// Validate validates the MysqlConf.
func (rc OrmConf) Validate() error {
	if len(rc.Host) == 0 {
		return ErrEmptyHost
	}
	if len(rc.Port) == 0 {
		return ErrEmptyPort
	}
	if len(rc.User) == 0 {
		return ErrEmptyUser
	}
	if len(rc.Pass) == 0 {
		return ErrEmptyPass
	}
	if len(rc.Database) == 0 {
		return ErrEmptyDB
	}
	if len(rc.Schema) == 0 {
		return ErrorEmptySchema
	}

	return nil
}
