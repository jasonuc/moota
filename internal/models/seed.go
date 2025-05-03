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
	BotanicalName string   `json:"botanicalName"`
	OptimalSoil   SoilType `json:"optimalSoil"`
}

type Seed struct {
	ID        string     `json:"id"`
	Hp        float64    `json:"hp"` // used as plant's starting health
	Planted   bool       `json:"planted"`
	OwnerID   string     `json:"ownerID"`
	CreatedAt *time.Time `json:"createdAt"`
	SeedMeta
}

type SeedGroup struct {
	BotanicalName string  `json:"botanicalName"`
	Count         int     `json:"count"`
	Seeds         []*Seed `json:"seeds"`
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

func NewSeedWithMeta(ownerID string, seedMeta SeedMeta) *Seed {
	return &Seed{
		Hp:       50.0,
		Planted:  false,
		OwnerID:  ownerID,
		SeedMeta: seedMeta,
	}
}
