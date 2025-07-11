package rbac

import (
	"testing"
)

func Test(t *testing.T) {
	InitTest()
	AddBasePolice()
	// 非admin和read权限
	permission := CheckStaffProjectPermission("staffId", "*", "log", "read")
	if permission != false {
		t.Fatal("expect false, but got true")
	}
	// admin权限
	permission = CheckStaffProjectPermission("admin", "*", "log", "read")
	if permission != true {
		t.Fatal("expect true, but got false")
	}
	// read权限
	permission = CheckStaffProjectPermission("read", "*", "log", "read")
	if permission != true {
		t.Fatal("expect true, but got false")
	}
	// 赋予staffId权限
	err := LinkStaffToProjectRole("staffId", "projectId", "admin")
	if err != nil {
		t.Fatal(err)
	}
	res := GetStaffRelatedProjects("staffId")
	t.Log(res)
	permission = CheckStaffProjectPermission("staffId", res[0], "log", "read")
	if permission != true {
		t.Fatal("expect true, but got false")
	}
	// 获取admin关联的项目
	res = GetStaffRelatedProjects("admin")
	t.Log(res)
}
