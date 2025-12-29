package dto

type GeneralLoginResponse struct {
	AccessToken          string `json:"token"`
	AccessTokenExpireIn  int64  `json:"expireIn"` // sec
	RefreshToken         string `json:"refreshToken"`
	RefreshTokenExpireIn int64  `json:"refreshTokenExpireIn"` // sec
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

type ModifyRequest struct {
	UserId   string `json:"userId"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type ManagerInfoResponse struct {
	Id          string   `json:"id"`
	StaffId     string   `json:"staffId"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions" gorm:"-"`
}
