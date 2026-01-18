package dto

import (
	"HelpStudent/internal/app/subject/model"
)

type AddSubjectReq struct {
	SubjectName string `json:"subject_name"`
	SubjectLink string `json:"subject_link"`
}

type AddSubjectResp struct {
	SubjectName string `json:"subject_name"`
	SubjectLink string `json:"subject_link"`
}

type UpdateSubjectReq struct {
	SubjectId   int    `json:"subject_id"`
	SubjectName string `json:"subject_name"`
	SubjectLink string `json:"subject_link"`
}

type SubjectItem struct {
	AppName      string `json:"app_name"`
	AppID        string `json:"app_id"`         // 我们系统的 ID
	FastgptAppId string `json:"fastgpt_app_id"` // FastGPT 的应用 ID
	ShareId      string `json:"share_id"`
}

type GetSubjectResp struct {
	Subjects []SubjectItem `json:"subjects"`
}

type GetSubjectListResp struct {
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
	Subjects []model.Subject `json:"subjects"`
}

// 学生科目关联相关的 DTO
type GetUserSubjectListResp struct {
	Total        int64               `json:"total"`
	Page         int                 `json:"page"`
	PageSize     int                 `json:"page_size"`
	UserSubjects []model.UserSubject `json:"user_subjects"`
}

type AddUserSubjectReq struct {
	StaffId     string `json:"staff_id"`
	SubjectName string `json:"subject_name"`
}

type UpdateUserSubjectReq struct {
	ID          string `json:"id"`
	StaffId     string `json:"staffId"`
	SubjectName string `json:"subjectName"`
}
