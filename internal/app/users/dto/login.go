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
