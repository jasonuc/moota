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
	RefreshToken string `json:"refreshToken" validate:"required,base64rawurl"`
}

type ChangeUsernameReq struct {
	NewUsername string `json:"username,omitempty"`
}

type ChangeEmailReq struct {
	NewEmail string `json:"email" validate:"required,email"`
}

type ChangePasswordReq struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=8,max=72"`
}
