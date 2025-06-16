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

type ChangeUsernameReq struct {
	NewUsername string `json:"newUsername" validate:"required"`
}

type ChangeEmailReq struct {
	NewEmail string `json:"newEmail" validate:"required,email"`
}

type ChangePasswordReq struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=8,max=72"`
}
