package store

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jasonuc/moota/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestPlantStore(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping plant store integration tests")
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

	initScriptPath := filepath.Join("testdata", "plant_store_init.sql")
	initSQL, err := os.ReadFile(initScriptPath)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(string(initSQL))
	if err != nil {
		t.Fatal(err)
	}

	store := &plantStore{db: db}
	ownerID := "00000000-0000-4000-a000-000000000001"

	t.Run("GetByOwnerID", func(t *testing.T) {
		plants, err := store.GetByOwnerID(context.Background(), ownerID, false)
		assert.NoError(t, err, "unexpected error")
		assert.Len(t, plants, 5, "expected 5 plants")
	})

	t.Run("Get", func(t *testing.T) {
		var plantID = "00000000-0000-4000-a000-000000000201"
		plant, err := store.Get(context.Background(), plantID, false)
		assert.NoError(t, err, "unexpected error")
		assert.Equalf(t, plantID, plant.ID, "expected plant ID %s, got %s", plantID, plant.ID)
	})

	t.Run("GetBySoilIDAndProximity", func(t *testing.T) {
		soilID := "00000000-0000-4000-a000-000000000101"
		coords := models.Coordinates{
			Lat: 40.782865,
			Lon: -73.965355,
		}

		plants, err := store.GetBySoilIDAndProximity(context.Background(), soilID, coords, models.SoilRadiusMMedium)
		assert.NoError(t, err, "unexpected error")
		assert.Len(t, plants, 3, "expected 3 plants")

		plants, err = store.GetBySoilIDAndProximity(context.Background(), soilID, coords, models.PlantInteractionRadius+1)
		assert.NoError(t, err, "unexpected error")
		assert.Len(t, plants, 1, "expected 1 plant with tight radius")
	})

	t.Run("GetByOwnerIDAndProximity_CentralPark", func(t *testing.T) {
		userCentralParkCoords := models.Coordinates{
			Lat: 40.782865,
			Lon: -73.965355,
		}

		plants, err := store.GetByOwnerIDAndProximity(context.Background(), ownerID, userCentralParkCoords)
		assert.NoError(t, err, "unexpected error")
		assert.Len(t, plants, 5, "expected 5 plants")

		firstPlant := plants[0]
		assert.Equal(t, "00000000-0000-4000-a000-000000000201", firstPlant.ID,
			"expected first plant to be Oak Tree in Central Park")
	})

	t.Run("GetByOwnerIDAndProximity_MiamiBeach", func(t *testing.T) {
		userMiamiCoords := models.Coordinates{
			Lat: 25.792236,
			Lon: -80.134358,
		}

		plants, err := store.GetByOwnerIDAndProximity(context.Background(), ownerID, userMiamiCoords)
		assert.NoError(t, err, "unexpected error")
		assert.Len(t, plants, 5, "expected 5 plants")

		firstPlant := plants[0]
		assert.Equal(t, "00000000-0000-4000-a000-000000000204", firstPlant.ID,
			"expected first plant to be Palm Tree in Miami Beach")
	})

	t.Run("ActivatePlant", func(t *testing.T) {
		plantID := "00000000-0000-4000-a000-000000000207"
		err := store.ActivatePlant(context.Background(), plantID)
		assert.NoError(t, err, "unexpected error")

		plant, err := store.Get(context.Background(), plantID, false)
		assert.NoError(t, err, "unexpected error")
		assert.True(t, plant.Activated, "expected plant to be activated")
	})

	t.Run("Insert", func(t *testing.T) {
		seed := models.NewSeed(ownerID)
		soil := models.NewLargeSizedSoil(models.RandomSoilMeta(), models.Coordinates{Lat: 29.153291, Lon: -89.254120})
		soil.ID = "00000000-0000-4000-a000-000000000103"
		plant, err := models.NewPlant(seed, soil, models.Coordinates{Lat: 29.153291, Lon: -89.254120})
		assert.NoError(t, err, "unexpected error")

		plant.ID = "00000000-0000-4000-a000-000000000999"

		err = store.Insert(context.Background(), plant)
		assert.NoError(t, err, "unexpected error")
	})

	t.Run("Update", func(t *testing.T) {
		t.Run("Update Plant xp", func(t *testing.T) {
			plantID := "00000000-0000-4000-a000-000000000201"
			plant, err := store.Get(context.Background(), plantID, false)
			assert.NoError(t, err, "unexpected error")

			plant.XP += 10
			err = store.Update(context.Background(), plant)
			assert.NoError(t, err, "unexpected error")

			updatedPlant, err := store.Get(context.Background(), plantID, false)
			assert.NoError(t, err, "unexpected error")
			assert.Equal(t, plant.XP, updatedPlant.XP, "expected plant XP to be updated")
		})

		t.Run("Update Plant nickname", func(t *testing.T) {
			plantID := "00000000-0000-4000-a000-000000000201"
			plant, err := store.Get(context.Background(), plantID, false)
			assert.NoError(t, err, "unexpected error")

			plant.Nickname = "Nickname"
			err = store.Update(context.Background(), plant)
			assert.NoError(t, err, "unexpected error")

			updatedPlant, err := store.Get(context.Background(), plantID, false)
			assert.NoError(t, err, "unexpected error")
			assert.Equal(t, plant.Nickname, updatedPlant.Nickname, "expected plant nickname to be updated")
		})
	})

	t.Run("Delete", func(t *testing.T) {
		plantID := "00000000-0000-4000-a000-000000000201"
		err := store.Delete(context.Background(), plantID)
		assert.NoError(t, err, "unexpected error")

		_, err = store.Get(context.Background(), plantID, false)
		assert.Error(t, err, "expected plant to be deleted, but it was found")
	})
}
