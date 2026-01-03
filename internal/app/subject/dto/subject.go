package dto

import "HelpStudent/internal/app/subject/model"

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

type GetSubjectResp struct {
	Subjects []model.Subject `json:"subjects"`
}

type GetSubjectListResp struct {
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
	Subjects []model.Subject `json:"subjects"`
}
