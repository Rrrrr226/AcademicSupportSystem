package dto

type GeneralLoginResponse struct {
	AccessToken          string `json:"token"`
	AccessTokenExpireIn  int64  `json:"expireIn"` // sec
	RefreshToken         string `json:"refreshToken"`
	RefreshTokenExpireIn int64  `json:"refreshTokenExpireIn"` // sec
}

type ThirdPlatLoginReq struct {
	Callback string `json:"callback"`
	Platform string `json:"platform" validate:"required"`
	From     string `json:"from"`
}

type ThirdPlatLoginResp struct {
	URL string `json:"url"`
}

type ThirdPlatLoginCallbackReq struct {
	Callback string `json:"callback" validate:"required"` // 也就是当前回调页面的地址 用于后端解析使用哪一套secret
	Code     string `json:"code"`                         // OAuth code
	Ticket   string `json:"ticket"`                       // CAS ticket
	State    string `json:"state" validate:"required"`
}

type ThirdPlatLoginCallbackResp struct {
	AccessToken          string `json:"token"`
	AccessTokenExpireIn  int64  `json:"expireIn"` // sec
	RefreshToken         string `json:"refreshToken"`
	RefreshTokenExpireIn int64  `json:"refreshTokenExpireIn"` // sec
}

type ThirdPlatBindReq struct {
	Platform string `json:"platform" validate:"required"`
	Redirect string `json:"redirect"`
	From     string `json:"from"`
}

type ThirdPlatBindResp struct {
	URL string `json:"url"`
}

type ThirdPlatUnbindReq struct {
	BindID uint `json:"bindID" validate:"required"`
}

type ThirdPlatUnbindResp struct {
	Success bool `json:"success"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenResponse struct {
	AccessToken         string `json:"token"`
	AccessTokenExpireIn int64  `json:"expireIn"` // sec
	RefreshToken        string `json:"refreshToken"`
}

// 在文件末尾添加以下结构体

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
