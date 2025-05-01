package dto

import "time"

type GetAllUserPlantsReq struct {
	Coordinates
}

type ActionOnPlantReq struct {
	Coordinates
	Action int       `json:"action" validate:"required,number"`
	Time   time.Time `json:"time" validate:"required"`
}
