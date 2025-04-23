package store

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jasonuc/moota/internal/models"
)

func TestSoilStore(t *testing.T) {
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
	defer db.Close()

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
		soil, err := store.Get(soilID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if soil.ID != soilID {
			t.Errorf("expected soil ID %s, got %s", soilID, soil.ID)
		}

		// Central Park coordinates
		if soil.Centre().Lat != 40.782865 || soil.Centre().Lng != -73.965355 {
			t.Errorf("incorrect coordinates, expected (40.782865, -73.965355), got (%f, %f)",
				soil.Centre().Lat, soil.Centre().Lng)
		}

		if soil.Type != "loam" {
			t.Errorf("expected type 'loam', got '%s'", soil.Type)
		}

		if soil.RadiusM() != 22.0 {
			t.Errorf("expected radius 22.0, got %f", soil.RadiusM())
		}
	})

	t.Run("Get_NonExisting", func(t *testing.T) {
		soilID := "00000000-0000-4000-c000-999999999999"
		_, err := store.Get(soilID)
		if err != models.ErrSoilNotFound {
			t.Errorf("expected ErrSoilNotFound, got %v", err)
		}
	})

	t.Run("GetAllInProximity_SF_5km", func(t *testing.T) {
		coords := models.Coordinates{
			Lat: 37.774929,
			Lng: -122.419416,
		}

		// 5km radius - should only include San Francisco soil
		soils, err := store.GetAllInProximity(coords, 5000)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if len(soils) != 1 {
			t.Errorf("expected 1 soil, got %d", len(soils))
		}

		if len(soils) > 0 && soils[0].ID != "00000000-0000-4000-c000-000000000105" {
			t.Errorf("expected SF soil, got %s", soils[0].ID)
		}
	})

	t.Run("GetAllInProximity_SF_15km", func(t *testing.T) {
		// San Francisco coordinates
		coords := models.Coordinates{
			Lat: 37.774929,
			Lng: -122.419416,
		}

		// 15km radius - should include San Francisco and Oakland soils
		soils, err := store.GetAllInProximity(coords, 15000)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if len(soils) != 2 {
			t.Errorf("expected 2 soils, got %d", len(soils))
		}
	})

	t.Run("GetAllInProximity_SF_50km", func(t *testing.T) {
		// San Francisco coordinates
		coords := models.Coordinates{
			Lat: 37.774929,
			Lng: -122.419416,
		}

		// 50km radius - should include all 3 Bay Area soils
		soils, err := store.GetAllInProximity(coords, 50000)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if len(soils) != 3 {
			t.Errorf("expected 3 soils, got %d", len(soils))
		}
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
			Lng: -118.243683,
		}
		newSoil.CircleMeta = models.NewCircleMeta(centre, 25.0)

		err := store.Insert(newSoil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if newSoil.ID == "" {
			t.Error("expected ID to be set after insert")
		}

		if newSoil.CreatedAt.IsZero() {
			t.Error("expected CreatedAt to be set after insert")
		}

		retrievedSoil, err := store.Get(newSoil.ID)
		if err != nil {
			t.Errorf("unexpected error retrieving inserted soil: %v", err)
		}

		if retrievedSoil.Type != "sandy" {
			t.Errorf("expected type 'sandy', got '%s'", retrievedSoil.Type)
		}

		if retrievedSoil.RadiusM() != 25.0 {
			t.Errorf("expected radius 25.0, got %f", retrievedSoil.RadiusM())
		}

		// Los Angeles coordinates
		if retrievedSoil.Centre().Lat != 34.052235 || retrievedSoil.Centre().Lng != -118.243683 {
			t.Errorf("incorrect coordinates, expected (34.052235, -118.243683), got (%f, %f)",
				retrievedSoil.Centre().Lat, retrievedSoil.Centre().Lng)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		soilID := "00000000-0000-4000-c000-000000000104"

		_, err := store.Get(soilID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err = store.Delete(soilID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		_, err = store.Get(soilID)
		if err != models.ErrSoilNotFound {
			t.Errorf("expected ErrSoilNotFound, got %v", err)
		}
	})

	t.Run("Delete_NonExisting", func(t *testing.T) {
		soilID := "00000000-0000-4000-c000-999999999999"
		err := store.Delete(soilID)
		if err != models.ErrSoilNotFound {
			t.Errorf("expected ErrSoilNotFound, got %v", err)
		}
	})
}
