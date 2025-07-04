package store

import (
	"context"
	"database/sql"
	"strings"

	"github.com/jasonuc/moota/internal/models"
)

type PlantStore interface {
	Get(context.Context, string, *GetPlantsOpts) (*models.Plant, error)
	GetByOwnerID(context.Context, string, *GetPlantsOpts) ([]*models.Plant, error)
	GetCountByUsername(context.Context, string) (*models.PlantCount, error)
	GetBySoilIDAndProximity(context.Context, string, models.Coordinates, float64) ([]*models.Plant, error)
	GetByOwnerIDAndProximity(context.Context, string, models.Coordinates) ([]*models.Plant, error)
	Insert(context.Context, *models.Plant) error
	Update(context.Context, *models.Plant) error
	Delete(context.Context, string) error
	GetTotalCount(context.Context) (*models.PlantCount, error)
}

var _ PlantStore = (*plantStore)(nil)

type plantStore struct {
	db Querier
}

// GetTotalCount implements PlantStore.
func (s *plantStore) GetTotalCount(ctx context.Context) (*models.PlantCount, error) {
	q := `SELECT p.dead, count(*) AS plant_count FROM plants p
			GROUP BY p.dead ORDER BY p.dead ASC;`

	plantCount := new(models.PlantCount)

	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	//nolint:errcheck
	defer rows.Close()

	for rows.Next() {
		var dead bool
		var count int64
		if err := rows.Scan(&dead, &count); err != nil {
			return nil, err
		}

		switch dead {
		case true:
			plantCount.Deceased = count
		case false:
			plantCount.Alive = count
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return plantCount, nil
}

type GetPlantsOpts struct {
	IncludeDeceased bool
}

func (s *plantStore) GetCountByUsername(ctx context.Context, userID string) (*models.PlantCount, error) {
	q := `SELECT p.dead, count(*) AS plant_count FROM plants p
			JOIN users u ON p.owner_id = u.id
			WHERE u.username = $1 
			GROUP BY p.dead ORDER BY p.dead ASC;`

	plantCount := new(models.PlantCount)

	rows, err := s.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	//nolint:errcheck
	defer rows.Close()

	for rows.Next() {
		var dead bool
		var count int64
		if err := rows.Scan(&dead, &count); err != nil {
			return nil, err
		}

		switch dead {
		case true:
			plantCount.Deceased = count
		case false:
			plantCount.Alive = count
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return plantCount, nil
}

func (s *plantStore) GetByOwnerIDAndProximity(ctx context.Context, ownerID string, point models.Coordinates) ([]*models.Plant, error) {
	q := `SELECT id, nickname, hp, dead, owner_id, time_planted, last_watered_at, last_action_at, 
         last_refreshed_at, grace_period_ends_at, ST_AsText(centre) as centre, radius_m, soil_id, 
         optimal_soil, botanical_name, level, xp, woe, frolic, dread, malice, time_of_death 
		FROM plants
		WHERE owner_id = $1 AND dead = false
		ORDER BY ST_Distance(centre, ST_SetSRID(ST_MakePoint($2, $3), 4326)::GEOGRAPHY) ASC;`

	rows, err := s.db.QueryContext(ctx, q, ownerID, point.Lon, point.Lat)
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
		plant.Soil = new(models.Soil)
		plant.Tempers = new(models.Tempers)
		err := rows.Scan(
			&plant.ID, &plant.Nickname, &plant.Hp, &plant.Dead, &plant.OwnerID,
			&plant.TimePlanted, &plant.LastWateredAt, &plant.LastActionAt,
			&plant.LastRefreshedAt, &plant.GracePeriodEndsAt, &centreText,
			&radiusM, &plant.Soil.ID, &plant.OptimalSoil, &plant.BotanicalName, &plant.Level, &plant.XP,
			&plant.Tempers.Woe, &plant.Tempers.Frolic, &plant.Tempers.Dread, &plant.Tempers.Malice, &plant.TimeOfDeath,
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

func (s *plantStore) GetBySoilIDAndProximity(ctx context.Context, soilID string, point models.Coordinates, distanceM float64) ([]*models.Plant, error) {
	q := `SELECT id, nickname, hp, dead, owner_id, time_planted, last_watered_at, last_action_at, 
         last_refreshed_at, grace_period_ends_at, ST_AsText(centre) as centre, radius_m, soil_id, 
         optimal_soil, botanical_name, level, xp, woe, frolic, dread, malice, time_of_death 
		FROM plants
		WHERE soil_id = $1 AND dead = false 
		AND ST_DWithin(centre, ST_SetSRID(ST_MakePoint($2, $3), 4326)::GEOGRAPHY, $4);`

	rows, err := s.db.QueryContext(ctx, q, soilID, point.Lon, point.Lat, distanceM)
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
		plant.Soil = new(models.Soil)
		plant.Tempers = new(models.Tempers)

		err := rows.Scan(
			&plant.ID, &plant.Nickname, &plant.Hp, &plant.Dead, &plant.OwnerID,
			&plant.TimePlanted, &plant.LastWateredAt, &plant.LastActionAt,
			&plant.LastRefreshedAt, &plant.GracePeriodEndsAt, &centreText,
			&radiusM, &plant.Soil.ID, &plant.OptimalSoil, &plant.BotanicalName, &plant.Level, &plant.XP,
			&plant.Tempers.Woe, &plant.Tempers.Frolic, &plant.Tempers.Dread, &plant.Tempers.Malice, &plant.TimeOfDeath,
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

func (s *plantStore) GetByOwnerID(ctx context.Context, ownerID string, opts *GetPlantsOpts) ([]*models.Plant, error) {
	q := `SELECT id, nickname, hp, dead, owner_id, time_planted, last_watered_at, 
         last_action_at, last_refreshed_at, grace_period_ends_at, ST_AsText(centre) as centre, 
         radius_m, soil_id, optimal_soil, botanical_name, level, xp, woe, frolic, dread, malice, time_of_death 
		FROM plants
		WHERE owner_id = $1`

	if !opts.IncludeDeceased {
		q += ` AND dead = false`
	}

	rows, err := s.db.QueryContext(ctx, q, ownerID)
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
		plant.Soil = new(models.Soil)
		plant.Tempers = new(models.Tempers)

		err := rows.Scan(
			&plant.ID, &plant.Nickname, &plant.Hp, &plant.Dead, &plant.OwnerID,
			&plant.TimePlanted, &plant.LastWateredAt, &plant.LastActionAt,
			&plant.LastRefreshedAt, &plant.GracePeriodEndsAt, &centreText,
			&radiusM, &plant.Soil.ID, &plant.OptimalSoil, &plant.BotanicalName, &plant.Level, &plant.XP,
			&plant.Tempers.Woe, &plant.Tempers.Frolic, &plant.Tempers.Dread, &plant.Tempers.Malice, &plant.TimeOfDeath,
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

func (s *plantStore) Get(ctx context.Context, id string, opts *GetPlantsOpts) (*models.Plant, error) {
	q := `SELECT 
			p.id, p.nickname, p.hp, p.dead, p.owner_id, p.time_planted, p.last_watered_at, p.last_action_at, 
			p.last_refreshed_at, p.grace_period_ends_at, ST_AsText(p.centre), 
			p.radius_m, p.soil_id, p.optimal_soil, p.botanical_name, p.level, p.xp, p.woe, p.frolic, p.dread, p.malice, 
			ST_AsText(s.centre), s.radius_m, s.soil_type, s.water_retention, s.nutrient_richness, s.created_at, p.time_of_death
			FROM plants p JOIN soils s ON p.soil_id = s.id
			WHERE p.id = $1 AND (dead = false OR dead = $2);`

	var plantCentreText string
	var plantRadiusM float64

	var soilCentreText string
	var soilRadiusM float64
	plant := new(models.Plant)
	plant.Soil = new(models.Soil)
	plant.Tempers = new(models.Tempers)

	err := s.db.QueryRowContext(ctx, q, id, opts.IncludeDeceased).Scan(
		&plant.ID, &plant.Nickname, &plant.Hp, &plant.Dead, &plant.OwnerID,
		&plant.TimePlanted, &plant.LastWateredAt, &plant.LastActionAt,
		&plant.LastRefreshedAt, &plant.GracePeriodEndsAt, &plantCentreText,
		&plantRadiusM, &plant.Soil.ID, &plant.OptimalSoil, &plant.BotanicalName, &plant.Level, &plant.XP,
		&plant.Tempers.Woe, &plant.Tempers.Frolic, &plant.Tempers.Dread, &plant.Tempers.Malice,
		&soilCentreText, &soilRadiusM, &plant.Soil.Type, &plant.Soil.WaterRetention, &plant.Soil.NutrientRichness, &plant.Soil.CreatedAt, &plant.TimeOfDeath,
	)
	if err != nil {
		if err == sql.ErrNoRows || strings.Contains(err.Error(), ErrInvalidUUIDSyntax) {
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

func (s *plantStore) Insert(ctx context.Context, plant *models.Plant) error {
	q := `INSERT INTO plants (nickname, hp, owner_id, centre, radius_m, soil_id, optimal_soil, botanical_name, level, xp, woe, frolic, dread, malice, last_refreshed_at, grace_period_ends_at)
			VALUES ($1, $2, $3, ST_SetSRID(ST_MakePoint($4, $5), 4326), $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
			RETURNING id, dead, time_planted, last_watered_at, last_action_at;`

	err := s.db.QueryRowContext(ctx,
		q,
		plant.Nickname, plant.Hp, plant.OwnerID,
		plant.CircleMeta.Centre().Lon, plant.CircleMeta.Centre().Lat, plant.CircleMeta.RadiusM(),
		plant.Soil.ID, plant.OptimalSoil, plant.BotanicalName,
		plant.Level, plant.XP,
		plant.Tempers.Woe, plant.Tempers.Frolic, plant.Tempers.Dread, plant.Tempers.Malice,
		plant.LastRefreshedAt, plant.GracePeriodEndsAt,
	).Scan(
		&plant.ID, &plant.Dead, &plant.TimePlanted, &plant.LastWateredAt, &plant.LastActionAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *plantStore) Update(ctx context.Context, plant *models.Plant) error {
	q := `UPDATE plants 
          SET nickname = $1, hp = $2, dead = $3, 
              level = $4, xp = $5,
              last_action_at = $6, last_watered_at = $7, time_of_death = $8,
              last_refreshed_at = $9, grace_period_ends_at = $10
          WHERE id = $11;`

	res, err := s.db.ExecContext(ctx, q,
		plant.Nickname, plant.Hp, plant.Dead,
		plant.Level, plant.XP,
		plant.LastActionAt, plant.LastWateredAt, plant.TimeOfDeath,
		plant.LastRefreshedAt, plant.GracePeriodEndsAt,
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

func (s *plantStore) Delete(ctx context.Context, id string) error {
	q := `DELETE from plants
			WHERE ID = $1;`

	res, err := s.db.ExecContext(ctx, q, id)
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
