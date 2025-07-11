package rbac

import (
	"errors"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

var enforce *casbin.Enforcer

func InitTest() {
	m, _ := model.NewModelFromString(rbacRule)

	a := fileadapter.NewAdapter("./test.csv")
	e, _err := casbin.NewEnforcer(m, a)
	if _err != nil {
		panic(_err)
	}

	_err = e.LoadPolicy()
	if _err != nil {
		panic(_err)
	}

	enforce = e
}

func Init(db *gorm.DB) error {
	if db == nil {
		return errors.New("casbin init failed: db is nil")
	}

	m, _ := model.NewModelFromString(rbacRule)

	a, _ := gormadapter.NewAdapterByDB(db)
	e, _err := casbin.NewEnforcer(m, a)
	if _err != nil {
		return _err
	}

	// Load the policy from DB.
	_err = e.LoadPolicy()
	if _err != nil {
		return _err
	}
	e.EnableAutoSave(true)
	enforce = e
	AddBasePolice()
	return nil
}

func GetEnforcer() *casbin.Enforcer {
	return enforce
}
