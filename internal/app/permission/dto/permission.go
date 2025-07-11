package dto

type KVMap struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type AddProjectManagerRequest struct {
	StaffId     string   `json:"staffId" validate:"required"`
	Permissions []string `json:"permissions" validate:"required"`
}

type RemoveProjectManagerRequest struct {
	StaffIds []string `json:"staffIds" validate:"required"`
}

type ProjectManager struct {
	StaffId     string  `json:"staffId"`
	Permissions []KVMap `json:"permissions"`
}
