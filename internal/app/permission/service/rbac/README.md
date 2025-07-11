# CASBIN权限管理

模型：含域的RBAC权限

域：分为两个域，一个是系统域，一个是Project域
系统域：系统域是指系统级别的权限，比如用户管理，系统审查，权限管理等
Project域：Project域是指项目级别的权限，比如项目管理，项目审查等

## Project域
Domain: project/<id>

Object-action:
- info:write 项目信息修改
- project:delete 项目删除
- staff:read 项目成员查看
- staff:write 项目成员修改
- manager:read 项目管理员查看
- manager:write 项目管理员修改
- rule:read 项目规则查看
- rule:write 项目规则修改
- control:write 项目控制修改(开放申请，关闭申请，开放审查，关闭审查，开放规则补充申请，关闭规则补充申请)
- form:design 项目表单设计
- review:read 审批只读
- review:write 审批修改
- review:design 项目审批流设计
- library:read 项目成绩库查看
- library:write 项目成绩库修改
- log:read 项目日志查看
- data:export 项目数据导出
