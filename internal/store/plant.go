package store

import (
	"database/sql"

	"github.com/jasonuc/moota/internal/models"
)

type PlantStore interface {
	Get(string) (*models.Plant, error)
	GetAllByOwnerID(string) ([]*models.Plant, error)
	GetAllInSoilAndInProximity(string, models.Coordinates, float64) ([]*models.Plant, error)
	GetByOwnerIDAndProximity(string, models.Coordinates, float64) ([]*models.Plant, error)
	GetByOwnerIDAndOrderByProximity(string, models.Coordinates) ([]*models.Plant, error)
	ActivatePlant(string) error
	Insert(*models.Plant) error
	Update(*models.Plant) error
	Delete(string) error
}

type plantStore struct {
	db Querier
}

func (s *plantStore) GetByOwnerIDAndOrderByProximity(ownerID string, point models.Coordinates) ([]*models.Plant, error) {
	q := `SELECT id, nickname, hp, dead, owner_id, planted_at, last_watered_at, last_action_time, ST_AsText(centre) as centre, radius_m, soil_id, optimal_soil, botanical_name, level, xp, woe, frolic, dread, malice FROM plants
			WHERE owner_id = $1 AND activated = true AND dead = false
			ORDER BY ST_Distance(centre, ST_SetSRID(ST_MakePoint($2, $3), 4326)::GEOGRAPHY) DESC;`

	rows, err := s.db.Query(q, ownerID, point.Lng, point.Lat)
	if err != nil {
		return nil, err
	}
	//nolint:errcheck
	defer rows.Close()

	plants := make([]*models.Plant, 0)
	for rows.Next() {
		var centreText string
		var radiusM float64

		plant := new(models.Plant)
		err := rows.Scan(
			&plant.ID, &plant.Nickname, &plant.Hp, &plant.Dead, &plant.OwnerID,
			&plant.TimePlanted, &plant.LastWateredAt, &plant.LastActionTime, &centreText,
			&radiusM, &plant.Soil.ID, &plant.OptimalSoil, &plant.BotanicalName, &plant.Level, &plant.Xp,
			&plant.Tempers.Woe, &plant.Tempers.Frolic, &plant.Tempers.Dread, &plant.Tempers.Malice, &plant.Activated,
		)
		if err != nil {
			return nil, err
		}

		centre, err := models.CoordinatesFromPostGIS(centreText)
		if err != nil {
			return nil, err
		}

		plant.CircleMeta = models.NewCircleMeta(centre, radiusM)

		plants = append(plants, plant)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return plants, nil
}

func (s *plantStore) GetPendingPlant(plantID string, soilID string) (*models.Plant, error) {
	q := `SELECT id, nickname, hp, dead, owner_id, planted_at, last_watered_at, last_action_time, ST_AsText(centre) as centre, radius_m, soil_id, optimal_soil, botanical_name, level, xp, woe, frolic, dread, malice, activated FROM plants
			WHERE id = $1 AND soil_id = $2 AND dead = false AND activated = false;`

	var centreText string
	var radiusM float64
	plant := new(models.Plant)

	err := s.db.QueryRow(q, plantID, soilID).Scan(
		&plant.ID, &plant.Nickname, &plant.Hp, &plant.Dead, &plant.OwnerID,
		&plant.TimePlanted, &plant.LastWateredAt, &plant.LastActionTime, &centreText,
		&radiusM, &plant.Soil.ID, &plant.OptimalSoil, &plant.BotanicalName, &plant.Level, &plant.Xp,
		&plant.Tempers.Woe, &plant.Tempers.Frolic, &plant.Tempers.Dread, &plant.Tempers.Malice,
		&plant.Activated,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrPlantNotFound
		}
		return nil, err
	}

	centre, err := models.CoordinatesFromPostGIS(centreText)
	if err != nil {
		return nil, err
	}
	plant.CircleMeta = models.NewCircleMeta(centre, radiusM)
	plant.Soil.CircleMeta = models.NewCircleMeta(centre, radiusM)

	return plant, nil
}

func (s *plantStore) GetAllInSoilAndInProximity(soilID string, point models.Coordinates, distanceM float64) ([]*models.Plant, error) {
	q := `SELECT id, nickname, hp, dead, owner_id, planted_at, last_watered_at, last_action_time, ST_AsText(centre) as centre, radius_m, soil_id, optimal_soil, botanical_name, level, xp, woe, frolic, dread, malice, activated FROM plants
			WHERE soil_id = $1 AND dead = false AND activated = true
			AND ST_DWithin(centre, ST_SetSRID(ST_MakePoint($2, $3), 4326)::GEOGRAPHY, $4);`

	rows, err := s.db.Query(q, soilID, point.Lng, point.Lat, distanceM)
	if err != nil {
		return nil, err
	}
	//nolint:errcheck
	defer rows.Close()

	plants := make([]*models.Plant, 0)
	for rows.Next() {
		var centreText string
		var radiusM float64

		plant := new(models.Plant)
		err := rows.Scan(
			&plant.ID, &plant.Nickname, &plant.Hp, &plant.Dead, &plant.OwnerID,
			&plant.TimePlanted, &plant.LastWateredAt, &plant.LastActionTime, &centreText,
			&radiusM, &plant.Soil.ID, &plant.OptimalSoil, &plant.BotanicalName, &plant.Level, &plant.Xp,
			&plant.Tempers.Woe, &plant.Tempers.Frolic, &plant.Tempers.Dread, &plant.Tempers.Malice, &plant.Activated,
		)
		if err != nil {
			return nil, err
		}

		centre, err := models.CoordinatesFromPostGIS(centreText)
		if err != nil {
			return nil, err
		}

		plant.CircleMeta = models.NewCircleMeta(centre, radiusM)

		plants = append(plants, plant)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return plants, nil
}

func (s *plantStore) GetAllByOwnerID(ownerID string) ([]*models.Plant, error) {
	q := `SELECT id, nickname, hp, dead, owner_id, planted_at, last_watered_at, last_action_time, ST_AsText(centre) as centre, radius_m, soil_id, optimal_soil, botanical_name, level, xp, woe, frolic, dread, malice, activated FROM plants
			WHERE owner_id = $1 AND dead = false AND activated = true;`

	rows, err := s.db.Query(q, ownerID)
	if err != nil {
		return nil, err
	}
	//nolint:errcheck
	defer rows.Close()

	plants := make([]*models.Plant, 0)
	for rows.Next() {
		var centreText string
		var radiusM float64

		plant := new(models.Plant)
		err := rows.Scan(
			&plant.ID, &plant.Nickname, &plant.Hp, &plant.Dead, &plant.OwnerID,
			&plant.TimePlanted, &plant.LastWateredAt, &plant.LastActionTime, &centreText,
			&radiusM, &plant.Soil.ID, &plant.OptimalSoil, &plant.BotanicalName, &plant.Level, &plant.Xp,
			&plant.Tempers.Woe, &plant.Tempers.Frolic, &plant.Tempers.Dread, &plant.Tempers.Malice, &plant.Activated,
		)
		if err != nil {
			return nil, err
		}

		centre, err := models.CoordinatesFromPostGIS(centreText)
		if err != nil {
			return nil, err
		}

		plant.CircleMeta = models.NewCircleMeta(centre, radiusM)

		plants = append(plants, plant)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return plants, nil
}

func (s *plantStore) GetByOwnerIDAndProximity(ownerID string, point models.Coordinates, distanceM float64) ([]*models.Plant, error) {
	q := `SELECT id, nickname, hp, dead, owner_id, planted_at, last_watered_at, last_action_time, ST_AsText(centre) as centre, radius_m, soil_id, optimal_soil, botanical_name, level, xp, woe, frolic, dread, malice, activated FROM plants
			WHERE ST_DWithin(centre, ST_SetSRID(ST_MakePoint($1, $2), 4326)::GEOGRAPHY, $3) AND owner_id = $4 AND activated = true;`

	rows, err := s.db.Query(q, point.Lng, point.Lat, distanceM, ownerID)
	if err != nil {
		return nil, err
	}
	//nolint:errcheck
	defer rows.Close()

	plants := make([]*models.Plant, 0)
	for rows.Next() {
		var centreText string
		var radiusM float64
		plant := new(models.Plant)
		err := rows.Scan(
			&plant.ID, &plant.Nickname, &plant.Hp, &plant.Dead, &plant.OwnerID,
			&plant.TimePlanted, &plant.LastWateredAt, &plant.LastActionTime, &centreText,
			&radiusM, &plant.Soil.ID, &plant.OptimalSoil, &plant.BotanicalName, &plant.Level, &plant.Xp,
			&plant.Tempers.Woe, &plant.Tempers.Frolic, &plant.Tempers.Dread, &plant.Tempers.Malice, &plant.Activated,
		)
		if err != nil {
			return nil, err
		}
		centre, err := models.CoordinatesFromPostGIS(centreText)
		if err != nil {
			return nil, err
		}
		plant.CircleMeta = models.NewCircleMeta(centre, radiusM)
		plants = append(plants, plant)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return plants, nil
}

func (s *plantStore) Get(id string) (*models.Plant, error) {
	q := `SELECT 
			p.id, p.nickname, p.hp, p.dead, p.owner_id, p.planted_at, p.last_watered_at, p.last_action_time, ST_AsText(p.centre), 
			p.radius_m, p.soil_id, p.optimal_soil, p.botanical_name, p.level, p.xp, p.woe, p.frolic, p.dread, p.malice, p.activated, 
			ST_AsText(s.centre), s.radius_m, s.soil_type, s.water_retention, s.nutrient_richness, s.created_at
			FROM plants p JOIN soils s ON p.soil_id = s.id
			WHERE p.id = $1 AND p.activated = true AND p.dead = false;`

	var plantCentreText string
	var plantRadiusM float64

	var soilCentreText string
	var soilRadiusM float64
	plant := new(models.Plant)

	err := s.db.QueryRow(q, id).Scan(
		&plant.ID, &plant.Nickname, &plant.Hp, &plant.Dead, &plant.OwnerID,
		&plant.TimePlanted, &plant.LastWateredAt, &plant.LastActionTime, &plantCentreText,
		&plantRadiusM, &plant.Soil.ID, &plant.OptimalSoil, &plant.BotanicalName, &plant.Level, &plant.Xp,
		&plant.Tempers.Woe, &plant.Tempers.Frolic, &plant.Tempers.Dread, &plant.Tempers.Malice, &plant.Activated,
		&soilCentreText, &soilRadiusM, &plant.Soil.Type, &plant.Soil.WaterRetention, &plant.Soil.NutrientRichness, &plant.Soil.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrPlantNotFound
		}
		return nil, err
	}

	plantCentre, err := models.CoordinatesFromPostGIS(plantCentreText)
	if err != nil {
		return nil, err
	}
	plant.CircleMeta = models.NewCircleMeta(plantCentre, plantRadiusM)

	soilCentre, err := models.CoordinatesFromPostGIS(soilCentreText)
	if err != nil {
		return nil, err
	}
	plant.Soil.CircleMeta = models.NewCircleMeta(soilCentre, soilRadiusM)

	return plant, nil
}

func (s *plantStore) Insert(plant *models.Plant) error {
	q := `INSERT INTO plants (nickname, hp, owner_id, centre, radius_m, soil_id, optimal_soil, botanical_name, level, xp, woe, frolic, dread, malice)
			VALUES ($1, $2, $3, ST_SetSRID(ST_MakePoint($4, $5), 4326), $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
			RETURNING id, dead, time_planted, last_watered_at, last_action_time, activated;`

	err := s.db.QueryRow(
		q,
		plant.Nickname, plant.Hp, plant.OwnerID,
		plant.CircleMeta.Centre().Lng, plant.CircleMeta.Centre().Lat, plant.CircleMeta.RadiusM,
		plant.Soil.ID, plant.OptimalSoil, plant.BotanicalName,
		plant.Level, plant.Xp,
		plant.Tempers.Woe, plant.Tempers.Frolic, plant.Tempers.Dread, plant.Tempers.Malice,
	).Scan(
		&plant.ID, &plant.Dead, &plant.TimePlanted, &plant.LastWateredAt, &plant.LastActionTime, &plant.Activated,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *plantStore) Update(plant *models.Plant) error {
	q := `UPDATE plants 
          SET nickname = $1, hp = $2, dead = $3, 
              level = $4, xp = $5,
              last_action_time = $6, last_watered_at = $7
          WHERE id = $8;`

	res, err := s.db.Exec(q,
		plant.Nickname, plant.Hp, plant.Dead,
		plant.Level, plant.Xp,
		plant.LastActionTime, plant.LastWateredAt,
		plant.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrPlantNotFound
	}

	return nil
}

func (s *plantStore) ActivatePlant(plantID string) error {
	q := `UPDATE plants 
          SET activated = true
          WHERE id = $1 AND activated = false;`
	res, err := s.db.Exec(q, plantID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return models.ErrPlantNotFound
	}
	return nil
}

func (s *plantStore) Delete(id string) error {
	q := `DELETE from plants
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
		return models.ErrPlantNotFound
	}

	return nil
}
