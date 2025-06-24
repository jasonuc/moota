package dto

type GetUserPlantsReq struct {
	Coordinates
}

type ActionOnPlantReq struct {
	Coordinates
	Action int `json:"action" validate:"required,number"`
}

type ChangePlantNicknameReq struct {
	NewNickname string `json:"newNickname" validate:"required"`
}
