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

type UpdateUserReq struct {
	UserID   string `json:"userID" validate:"required"`
	Username string `json:"username,omitempty"`
}

type ChangeEmailReq struct {
	UserID string `json:"userID" validate:"required"`
	Email  string `json:"email" validate:"required,email"`
}

type ChangePasswordReq struct {
	UserID   string `json:"userID" validate:"required"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}
