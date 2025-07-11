package dto

type UserInfoResponse struct {
	Id          string   `json:"id"`
	StaffId     string   `json:"staffId"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions" gorm:"-"`
}
