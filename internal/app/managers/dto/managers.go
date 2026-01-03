package dto

// AddManagerRequest 添加管理员请求
type AddManagerRequest struct {
	StaffId string `json:"staffId" validate:"required"`
}

// AddManagerResponse 添加管理员响应
type AddManagerResponse struct {
	Id      string `json:"id"`
	StaffId string `json:"staffId"`
}

// DeleteManagerRequest 删除管理员请求
type DeleteManagerRequest struct {
	StaffId string `json:"staffId" validate:"required"`
}

// ManagerListResponse 管理员列表响应
type ManagerListResponse struct {
	Managers []ManagerItem `json:"managers"`
	Total    int64         `json:"total"`
}

type ManagerItem struct {
	StaffId string `json:"staffId"`
}

// ImportStudentSubjectsRequest 导入学生科目请求（用于JSON格式）
type ImportStudentSubjectsRequest struct {
	Data []StudentSubjectItem `json:"data" validate:"required"`
}

type StudentSubjectItem struct {
	StaffId     string `json:"staffId" validate:"required"`     // 学号
	SubjectName string `json:"subjectName" validate:"required"` // 科目名称
}

// ImportStudentSubjectsResponse 导入学生科目响应
type ImportStudentSubjectsResponse struct {
	Total        int      `json:"total"`        // 总记录数
	SuccessCount int      `json:"successCount"` // 成功数
	FailCount    int      `json:"failCount"`    // 失败数
	Errors       []string `json:"errors"`       // 错误详情
}
