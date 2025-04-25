package schema

type AuthRegisterRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required,min=6,max=16"`
}

type AuthLoginRequest struct {
	Account  string  `json:"account" binding:"required"`
	Password string  `json:"password" binding:"required,min=6,max=16"`
	Otp      *string `json:"otp"`
}

type AuthLoginResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
	Otp     string `json:"otp"` // for local test only
}

type AuthVerifyOTPRequest struct {
	Account string  `json:"account"`
	Otp     *string `json:"otp" binding:"required"`
}

type AuthVerifyOTPResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
}
