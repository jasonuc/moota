package dto

import "time"

type GetAllUserPlantsReq struct {
	Longitude float64 `json:"longitude" validate:"required,longitude"`
	Latitude  float64 `json:"latitude" validate:"required,latitude"`
}

type ActionOnPlantReq struct {
	Action    int       `json:"action" validate:"required,number"`
	Longitude float64   `json:"longitude" validate:"required,longitude"`
	Latitude  float64   `json:"latitude" validate:"required,latitude"`
	Time      time.Time `json:"time" validate:"required"`
	PlantID   string    `json:"plantID" validate:"required"`
}
