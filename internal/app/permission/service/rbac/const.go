package rbac

const (
	rbacRule = `
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub, r.dom) && keyMatch(r.dom, p.dom) && r.obj == p.obj && r.act == p.act 
`
)

var projectNameMap = map[string]map[string]string{
	"info": {
		"write": "项目信息修改",
	},
	"staff": {
		"read":  "员工浏览",
		"write": "员工浏览",
	},
	"project": {
		"delete": "项目删除",
	},
	"manager": {
		"read":  "项目管理员查看",
		"write": "项目管理员修改",
	},
	"form": {
		"write": "项目表单修改",
	},
	"control": {
		"write": "项目控制修改",
	},
	"review": {
		"read":   "审批只读",
		"write":  "审批修改",
		"design": "项目审批流设计",
	},
	"library": {
		"read":  "项目成绩库查看",
		"write": "项目成绩库修改",
	},
	"log": {
		"read": "项目日志查看",
	},
	"data": {
		"export": "项目数据导出",
	},
	"rule": {
		"read":  "项目规则查看",
		"write": "项目规则修改",
	},
	"result": {
		"read":  "项目结果查看",
		"write": "项目结果修改",
	},
}

func GetProjectActionName(obj, act string) string {
	if projectNameMap[obj] == nil {
		return ""
	}
	return projectNameMap[obj][act]
}

func ProjectActionExist(obj, act string) bool {
	if projectNameMap[obj] == nil {
		return false
	}
	if projectNameMap[obj][act] == "" {
		return false
	}
	return true
}

var systemNameMap = map[string]map[string]string{
	"project": {
		"create": "项目创建",
	},
	"lib": {
		"read":  "库查看",
		"write": "库修改",
	},
	"permission": {
		"read":  "权限查看",
		"write": "权限修改",
	},
	"log": {
		"read": "日志查看",
	},
}

func GetSystemActionName(obj, act string) string {
	if systemNameMap[obj] == nil {
		return ""
	}
	return systemNameMap[obj][act]
}

func SystemActionExist(obj, act string) bool {
	if systemNameMap[obj] == nil {
		return false
	}
	if systemNameMap[obj][act] == "" {
		return false
	}
	return true
}
