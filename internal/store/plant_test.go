package store

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/jasonuc/moota/internal/models"
)

func TestPlantStore(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	pgContainer, err := createPostgresContainer(ctx)
	fmt.Println(pgContainer.connectionString)
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
		plants, err := store.GetByOwnerID(ownerID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if len(plants) != 5 {
			t.Errorf("expected 5 plants, got %d", len(plants))
		}

		fmt.Println(plants[0].ID)
	})

	t.Run("Get", func(t *testing.T) {
		var plantID = "00000000-0000-4000-a000-000000000201"
		plant, err := store.Get(plantID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if plant.ID != plantID {
			t.Errorf("expected plant ID %s, got %s", plantID, plant.ID)
		}
	})

	t.Run("GetBySoilIDAndProximity", func(t *testing.T) {
		soilID := "00000000-0000-4000-a000-000000000101"
		coords := models.Coordinates{
			Lat: 40.782865,
			Lng: -73.965355,
		}

		plants, err := store.GetBySoilIDAndProximity(soilID, coords, models.SoilRadiusMMedium)
		fmt.Println(len(plants))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if len(plants) != 3 {
			t.Errorf("expected 3 plants, got %d", len(plants))
		}

		plants, err = store.GetBySoilIDAndProximity(soilID, coords, models.PlantInteractionRadius+1)
		fmt.Println(len(plants))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if len(plants) != 1 {
			t.Errorf("expected 1 plant with tight radius, got %d", len(plants))
		}
	})

	t.Run("GetByOwnerIDAndProximity_CentralPark", func(t *testing.T) {
		userCentralParkCoords := models.Coordinates{
			Lat: 40.782865,
			Lng: -73.965355,
		}

		plants, err := store.GetByOwnerIDAndProximity(ownerID, userCentralParkCoords)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if len(plants) != 5 {
			t.Errorf("expected 5 plants, got %d", len(plants))
		}

		firstPlant := plants[0]
		if firstPlant.ID != "00000000-0000-4000-a000-000000000201" {
			t.Errorf("expected first plant to be Oak Tree in Central Park, got %s", firstPlant.Nickname)
		}
	})

	t.Run("GetByOwnerIDAndProximity_MiamiBeach", func(t *testing.T) {
		userMiamiCoords := models.Coordinates{
			Lat: 25.792236,
			Lng: -80.134358,
		}

		plants, err := store.GetByOwnerIDAndProximity(ownerID, userMiamiCoords)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if len(plants) != 5 {
			t.Errorf("expected 5 plants, got %d", len(plants))
		}

		firstPlant := plants[0]
		if firstPlant.ID != "00000000-0000-4000-a000-000000000204" {
			t.Errorf("expected first plant to be Palm Tree in Miami Beach, got %s", firstPlant.Nickname)
		}
	})
}
