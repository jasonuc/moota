package models

type Stats struct {
	PlantCount int64 `json:"plantCount"`
	UserCount  int64 `json:"userCount"`
	SeedCount  int64 `json:"seedCount"`
}
