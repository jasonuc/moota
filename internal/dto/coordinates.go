package dto

type Coordinates struct {
	Longitude *float64 `json:"longitude" validate:"required,longitude"`
	Latitude  *float64 `json:"latitude" validate:"required,latitude"`
}
