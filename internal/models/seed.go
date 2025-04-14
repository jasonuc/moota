package models

import (
	"errors"
	"math/rand/v2"
	"time"
)

var (
	ErrSeedNotFound       = errors.New("seed not found")
	ErrSeedAlreadyPlanted = errors.New("seed already planted")
)

type SeedMeta struct {
	BotanicalName string
	OptimalSoil   SoilType
}

type Seed struct {
	ID        string
	Hp        float64 // used as plant's starting health
	Planted   bool
	OwnerID   string
	CreatedAt time.Time
	SeedMeta
}

func (s SeedMeta) IsCompatibleWithSoil(target SoilType) bool {
	return s.OptimalSoil == target ||
		s.OptimalSoil == SoilTypeLoam ||
		target == SoilTypeLoam
}

var SeedMetaCatalog = []SeedMeta{
	{
		BotanicalName: "Solanum lycopersicum", // Tomato
		OptimalSoil:   SoilTypeLoam,
	},
	{
		BotanicalName: "Zea mays", // Corn
		OptimalSoil:   SoilTypeLoam,
	},
	{
		BotanicalName: "Daucus carota", // Carrot
		OptimalSoil:   SoilTypeSandy,
	},
	{
		BotanicalName: "Oryza sativa", // Rice
		OptimalSoil:   SoilTypeClay,
	},
	{
		BotanicalName: "Cucumis sativus", // Cucumber
		OptimalSoil:   SoilTypeLoam,
	},
	{
		BotanicalName: "Pisum sativum", // Pea
		OptimalSoil:   SoilTypeSilt,
	},
	{
		BotanicalName: "Allium cepa", // Onion
		OptimalSoil:   SoilTypeSandy,
	},
	{
		BotanicalName: "Glycine max", // Soybean
		OptimalSoil:   SoilTypeClay,
	},
	{
		BotanicalName: "Spinacia oleracea", // Spinach
		OptimalSoil:   SoilTypeLoam,
	},
	{
		BotanicalName: "Helianthus annuus", // Sunflower
		OptimalSoil:   SoilTypeSandy,
	},
}

func NewSeed(ownerID string) *Seed {
	if ownerID == "" {
		ownerID = "user-id"
	}

	return &Seed{
		Hp:       50.0,
		Planted:  false,
		OwnerID:  ownerID,
		SeedMeta: SeedMetaCatalog[rand.IntN(len(SeedMetaCatalog))],
	}
}
