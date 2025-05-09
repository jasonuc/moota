package store

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jasonuc/moota/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestSoilStore(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping soil store integration tests")
	}

	ctx := context.Background()
	pgContainer, err := createPostgresContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err := pgContainer.Terminate(ctx)
		if err != nil {
			t.Error(err)
		}
	})

	db, err := openDB(pgContainer.connectionString)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err := db.Close()
		if err != nil {
			t.Error(err)
		}
	})

	migrationsPath := filepath.Join("..", "..", "migrations")
	err = applyMigrations(db, migrationsPath)
	if err != nil {
		t.Fatal(err)
	}

	initScriptPath := filepath.Join("testdata", "soil_store_init.sql")
	initSQL, err := os.ReadFile(initScriptPath)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(string(initSQL))
	if err != nil {
		t.Fatal(err)
	}

	store := &soilStore{db: db}

	t.Run("Get_Existing", func(t *testing.T) {
		soilID := "00000000-0000-4000-c000-000000000101"
		soil, err := store.Get(context.Background(), soilID)
		assert.NoError(t, err, "unexpected error")
		assert.Equalf(t, soilID, soil.ID, "expected soil ID %s, got %s", soilID, soil.ID)

		// Central Park coordinates
		assert.Equal(t, 40.782865, soil.Centre().Lat, "incorrect latitude")
		assert.Equal(t, -73.965355, soil.Centre().Lon, "incorrect longitude")
		assert.EqualValues(t, "loam", soil.Type, "expected type 'loam'")
		assert.Equal(t, 22.0, soil.RadiusM(), "expected radius 22.0")
	})

	t.Run("Get_NonExisting", func(t *testing.T) {
		soilID := "00000000-0000-4000-c000-999999999999"
		_, err := store.Get(context.Background(), soilID)
		assert.ErrorIs(t, err, models.ErrSoilNotFound, "expected ErrSoilNotFound")
	})

	t.Run("GetAllInProximity_SF_5km", func(t *testing.T) {
		coords := models.Coordinates{
			Lat: 37.774929,
			Lon: -122.419416,
		}

		// 5km radius - should only include San Francisco soil
		soils, err := store.GetAllInProximity(context.Background(), coords, 5000)
		assert.NoError(t, err, "unexpected error")
		assert.Len(t, soils, 1, "expected 1 soil")

		if len(soils) > 0 {
			assert.Equal(t, "00000000-0000-4000-c000-000000000105", soils[0].ID, "expected SF soil")
		}
	})

	t.Run("GetAllInProximity_SF_15km", func(t *testing.T) {
		// San Francisco coordinates
		coords := models.Coordinates{
			Lat: 37.774929,
			Lon: -122.419416,
		}

		// 15km radius - should include San Francisco and Oakland soils
		soils, err := store.GetAllInProximity(context.Background(), coords, 15000)
		assert.NoError(t, err, "unexpected error")
		assert.Len(t, soils, 2, "expected 2 soils")
	})

	t.Run("GetAllInProximity_SF_50km", func(t *testing.T) {
		// San Francisco coordinates
		coords := models.Coordinates{
			Lat: 37.774929,
			Lon: -122.419416,
		}

		// 50km radius - should include all 3 Bay Area soils
		soils, err := store.GetAllInProximity(context.Background(), coords, 50000)
		assert.NoError(t, err, "unexpected error")
		assert.Len(t, soils, 3, "expected 3 soils")
	})

	t.Run("Insert", func(t *testing.T) {
		// New soil in Los Angeles
		newSoil := &models.Soil{
			SoilMeta: models.SoilMeta{
				Type:             "sandy",
				WaterRetention:   0.30,
				NutrientRichness: 0.40,
			},
		}

		centre := models.Coordinates{
			Lat: 34.052235,
			Lon: -118.243683,
		}
		newSoil.CircleMeta = models.NewCircleMeta(centre, 25.0)

		err := store.Insert(context.Background(), newSoil)
		assert.NoError(t, err, "unexpected error")
		assert.NotEmpty(t, newSoil.ID, "expected ID to be set after insert")
		assert.False(t, newSoil.CreatedAt.IsZero(), "expected CreatedAt to be set after insert")

		retrievedSoil, err := store.Get(context.Background(), newSoil.ID)
		assert.NoError(t, err, "unexpected error retrieving inserted soil")
		assert.EqualValues(t, "sandy", retrievedSoil.Type, "expected type 'sandy'")
		assert.Equal(t, 25.0, retrievedSoil.RadiusM(), "expected radius 25.0")

		// Los Angeles coordinates
		assert.Equal(t, 34.052235, retrievedSoil.Centre().Lat, "incorrect latitude")
		assert.Equal(t, -118.243683, retrievedSoil.Centre().Lon, "incorrect longitude")
	})

	t.Run("Delete", func(t *testing.T) {
		soilID := "00000000-0000-4000-c000-000000000104"

		_, err := store.Get(context.Background(), soilID)
		assert.NoError(t, err, "unexpected error")

		err = store.Delete(context.Background(), soilID)
		assert.NoError(t, err, "unexpected error")

		_, err = store.Get(context.Background(), soilID)
		assert.ErrorIs(t, err, models.ErrSoilNotFound, "expected ErrSoilNotFound")
	})

	t.Run("Delete_NonExisting", func(t *testing.T) {
		soilID := "00000000-0000-4000-c000-999999999999"
		err := store.Delete(context.Background(), soilID)
		assert.ErrorIs(t, err, models.ErrSoilNotFound, "expected ErrSoilNotFound")
	})
}
