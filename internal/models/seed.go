package models

type SeedMeta struct {
	BotanicalName string
	OptimalSoil   SoilType
}

type Seed struct {
	Health  float64 // used as plant's starting health
	Planted bool
	SeedMeta
}

func (s SeedMeta) IsCompatibleWithSoil(target SoilType) bool {
	return s.OptimalSoil == target ||
		s.OptimalSoil == SoilTypeLoam ||
		target == SoilTypeLoam
}

// TODO: Add Seeds
