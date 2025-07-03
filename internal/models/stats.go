package models

type Stats struct {
	Plant *PlantCount `json:"plant,omitempty"`
	Seed  *SeedCount  `json:"seed,omitempty"`
}
