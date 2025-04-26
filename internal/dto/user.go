package dto

type UserRegisterReq struct {
	Username string `json:"username" validate:"lowercase,min=3,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type UserLoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type TokenRefreshReq struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}
