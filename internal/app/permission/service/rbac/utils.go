package rbac

import (
	"HelpStudent/config"
	"HelpStudent/core/logx"
	"strings"
)

func AddBasePolice() {
	enforce.AddPolicy("admin", "project/*", "info", "write")
	enforce.AddPolicy("admin", "project/*", "project", "delete")
	enforce.AddPolicy("admin", "project/*", "staff", "read")
	enforce.AddPolicy("admin", "project/*", "staff", "write")
	enforce.AddPolicy("admin", "project/*", "manager", "read")
	enforce.AddPolicy("admin", "project/*", "manager", "write")
	enforce.AddPolicy("admin", "project/*", "rule", "read")
	enforce.AddPolicy("admin", "project/*", "rule", "write")
	enforce.AddPolicy("admin", "project/*", "control", "write")
	enforce.AddPolicy("admin", "project/*", "form", "write")
	enforce.AddPolicy("admin", "project/*", "review", "read")
	enforce.AddPolicy("admin", "project/*", "review", "write")
	enforce.AddPolicy("admin", "project/*", "review", "design")
	enforce.AddPolicy("admin", "project/*", "library", "read")
	enforce.AddPolicy("admin", "project/*", "library", "write")
	enforce.AddPolicy("admin", "project/*", "log", "read")
	enforce.AddPolicy("admin", "project/*", "data", "export")
	enforce.AddPolicy("admin", "project/*", "result", "read")
	enforce.AddPolicy("admin", "project/*", "result", "write")

	enforce.AddPolicy("admin", "system", "project", "create")
	enforce.AddPolicy("admin", "system", "lib", "read")
	enforce.AddPolicy("admin", "system", "lib", "write")
	enforce.AddPolicy("admin", "system", "permission", "read")
	enforce.AddPolicy("admin", "system", "permission", "write")
	enforce.AddPolicy("admin", "system", "log", "read")

	// 添加系统 admin权限
	newAdmins := config.GetConfig().AdminStaffID
	oldAdmins, _ := enforce.GetAllUsersByDomain("system")

	for _, admin := range newAdmins {
		if !strings.Contains(strings.Join(oldAdmins, ","), admin) {
			_, err := enforce.AddGroupingPolicy(admin, "admin", "system")
			if err != nil {
				logx.SystemLogger.Error("AddGroupingPolicy", err)
			} else {
				logx.SystemLogger.Info("AddGroupingPolicy", admin, "admin", "system")
			}
		}
	}

	// 删除系统 admin权限
	for _, admin := range oldAdmins {
		if !strings.Contains(strings.Join(newAdmins, ","), admin) {
			_, err := enforce.RemoveGroupingPolicy(admin, "admin", "system")
			if err != nil {
				logx.SystemLogger.Error("RemoveGroupingPolicy", err)
			} else {
				logx.SystemLogger.Info("RemoveGroupingPolicy", admin, "admin", "system")
			}
		}
	}

}

func GetStaffRelatedProjects(staffId string) (res []string) {
	domains, _err := enforce.GetDomainsForUser(staffId)
	if _err != nil {
		logx.SystemLogger.Error("GetDomainsForUser", _err)
	}
	polices, err := enforce.GetFilteredPolicy(0, staffId)
	if err != nil {
		logx.SystemLogger.Error("GetFilteredPolicy", err)
	}

	for _, police := range polices {
		domains = append(domains, police[1])
	}

	for _, domain := range domains {
		if !strings.HasPrefix(domain, "project/") {
			continue
		}
		res = append(res, strings.TrimPrefix(domain, "project/"))
	}
	return
}

func LinkStaffToProjectRole(staffId string, projectId string, role string) error {
	_, err := enforce.AddGroupingPolicy(staffId, role, "project/"+projectId)
	return err
}

func UnlinkStaffToProjectRole(staffId string, projectId string, role string) error {
	_, err := enforce.RemoveGroupingPolicy(staffId, role, "project/"+projectId)
	return err
}

func LinkStaffToProjectPermission(staffId, projectId, resource, action string) error {
	_, err := enforce.AddPolicy(staffId, "project/"+projectId, resource, action)
	return err
}

func UnlinkStaffToProjectPermission(staffId, projectId, resource, action string) error {
	_, err := enforce.RemovePolicy(staffId, "project/"+projectId, resource, action)
	return err
}

func UnlinkStaffToProjectAllPermission(staffId, projectId string) error {
	_, err := enforce.RemoveFilteredPolicy(0, staffId, "project/"+projectId)
	return err
}

func GetStaffProjectPermissions(staffId, projectId string) (res []string) {
	polices, _ := enforce.GetImplicitResourcesForUser(staffId, "project/"+projectId)
	for _, police := range polices {
		res = append(res, police[2]+":"+police[3])
	}
	return
}

func GetStaffProjectPermissionsArray(staffId, projectId string) (res [][]string) {
	polices, _ := enforce.GetImplicitResourcesForUser(staffId, "project/"+projectId)
	for _, police := range polices {
		res = append(res, []string{police[2], police[3]})
	}
	return
}

// CheckAdminProjectPermission 检查用户是否有 Admin 权限
func CheckAdminProjectPermission(staffId, projectId string) bool {
	user, err := enforce.GetImplicitRolesForUser(staffId, "project/"+projectId)
	if err != nil {
		return false
	}

	for _, role := range user {
		if role == "admin" {
			return true
		}
	}

	return false
}

func CheckStaffProjectPermission(staffId, projectId, resource, action string) bool {
	ok, _ := enforce.Enforce(staffId, "project/"+projectId, resource, action)
	return ok
}

func GetProjectManager(projectId string) (res []string) {

	r, _ := enforce.GetAllUsersByDomain("project/" + projectId)
	res = append(res, r[0])
	return res
}

func CheckStaffInDomain(projectId string, staffId string) bool {
	res, _ := enforce.GetImplicitResourcesForUser(staffId, "project/"+projectId)
	return len(res) > 0
}

func GetSystemManager() (res []string) {
	r, _ := enforce.GetAllUsersByDomain("system")
	res = append(res, r[0])
	return res
}

func GetStaffSystemPermission(staffId string) (res []string) {
	polices, _ := enforce.GetImplicitResourcesForUser(staffId, "system")
	for _, police := range polices {
		res = append(res, police[2]+":"+police[3])
	}
	return
}

func GetStaffSystemPermissionArray(staffId string) (res [][]string) {
	polices, _ := enforce.GetImplicitResourcesForUser(staffId, "system")
	for _, police := range polices {
		res = append(res, []string{police[2], police[3]})
	}
	return
}

func CheckStaffSystemPermission(staffId, resource, action string) bool {
	ok, _ := enforce.Enforce(staffId, "system", resource, action)
	return ok
}

func UnlinkStaffToSystemAllPermission(staffId string) error {
	_, err := enforce.RemoveFilteredPolicy(0, staffId, "system")
	return err
}

func LinkStaffToSystemPermission(staffId, resource, action string) error {
	_, err := enforce.AddPolicy(staffId, "system", resource, action)
	return err
}
