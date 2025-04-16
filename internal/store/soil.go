package store

import (
	"database/sql"
	"errors"

	"github.com/jasonuc/moota/internal/models"
)

type soilStore struct {
	db Querier
}

func (s *soilStore) Get(id string) (*models.Soil, error) {
	q := `SELECT id, ST_AsText(centre) as centre, radius_m, soil_type, water_retention, nutrient_richness, created_at FROM soils
            WHERE id = $1;`

	var centreText string
	var radiusM float64
	soil := new(models.Soil)

	err := s.db.QueryRow(q, id).Scan(
		&soil.ID, &centreText, &radiusM, &soil.Type, &soil.WaterRetention, &soil.NutrientRichness, &soil.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrSoilNotFound
		}
		return nil, err
	}

	centre, err := models.CoordinatesFromPostGIS(centreText)
	if err != nil {
		return nil, err
	}

	soil.CircleMeta = models.NewCircleMeta(centre, radiusM)

	return soil, nil
}

func (s *soilStore) GetAllInProximity(point models.Coordinates, distanceM float64) ([]*models.Soil, error) {
	q := `SELECT id, ST_AsText(centre) as centre, radius_m, soil_type, water_retention, nutrient_richness, created_at FROM soils
			WHERE ST_DWithin(centre, ST_Point($1, $2)::GEOGRAPHY, $3);`

	rows, err := s.db.Query(q, point.Lng, point.Lat, distanceM)
	if err != nil {
		return nil, err
	}

	//nolint:errcheck
	defer rows.Close()

	soils := make([]*models.Soil, 0)
	for rows.Next() {
		var centreText string
		var radiusM float64

		soil := new(models.Soil)
		err := rows.Scan(&soil.ID, &centreText, &radiusM, &soil.WaterRetention, &soil.NutrientRichness, &soil.CreatedAt)
		if err != nil {
			return nil, err
		}

		centre, err := models.CoordinatesFromPostGIS(centreText)
		if err != nil {
			return nil, err
		}

		soil.CircleMeta = models.NewCircleMeta(centre, radiusM)
		soils = append(soils, soil)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return soils, nil
}

func (s *soilStore) Insert(soil *models.Soil) error {
	q := `INSERT INTO soils (centre, radius_m, soil_type, nutrient_richness, water_retention)
			VALUES (ST_SetSRID(ST_MakePoint($1, $2), 4326), $3, $4, $5, $6)
			RETURNING id, created_at;`

	err := s.db.QueryRow(
		q, soil.Centre().Lng, soil.Centre().Lat, soil.RadiusM(), soil.Type, soil.NutrientRichness, soil.WaterRetention,
	).Scan(
		soil.ID, soil.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *soilStore) Delete(id string) error {
	q := `DELETE from soils
			WHERE ID = $1;`

	res, err := s.db.Exec(q, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrSoilNotFound
	}

	return nil
}
