package rbac

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"testing"
)

func TestCasbin(t *testing.T) {
	InitTest()
	roles, err := GetEnforcer().GetPermissionsForUser("alice")
	t.Log(roles, err)
	roles = GetEnforcer().GetFilteredPolicy(0, "alice", "project/xxx")
	t.Log(roles)
	t.Log(GetEnforcer().Enforce("bob", "project/xxx", "data", "read"))
	t.Log(GetEnforcer().Enforce("alice", "project/xxx", "review", "read"))
	t.Log(GetEnforcer().Enforce("alice", "project/666", "review", "read"))
	t.Log(GetEnforcer().GetFilteredGroupingPolicy(0, "alice"))
	t.Log(GetEnforcer().GetImplicitResourcesForUser("alice", "project/xxx"))
	t.Log(GetEnforcer().GetDomainsForUser("tom"))
	t.Log(GetEnforcer().GetAllUsersByDomain("project/xxx"))
}

func check(e *casbin.Enforcer, sub, obj, act string) {
	ok, _err := e.Enforce(sub, obj, act)
	if ok {
		fmt.Printf("%s CAN %s %s\n", sub, act, obj)
	} else {
		fmt.Printf("%s CANNOT %s %s: %v\n", sub, act, obj, _err)
	}
}

func TestArr(t *testing.T) {
	var a []int
	b := make([]int, 0)
	fmt.Println(len(a), len(b), a, b, a == nil, b == nil)
}
