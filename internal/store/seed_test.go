package store

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jasonuc/moota/internal/models"
)

func TestSeedStore(t *testing.T) {
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

	initScriptPath := filepath.Join("testdata", "seed_store_init.sql")
	initSQL, err := os.ReadFile(initScriptPath)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(string(initSQL))
	if err != nil {
		t.Fatal(err)
	}

	store := &seedStore{db: db}

	t.Run("GetAllByOwnerID", func(t *testing.T) {
		ownerID := "00000000-0000-4000-a000-000000000001"
		seeds, err := store.GetAllByOwnerID(ownerID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if len(seeds) != 3 {
			t.Errorf("expected 3 seeds, got %d", len(seeds))
		}

		for _, seed := range seeds {
			if seed.Planted {
				t.Errorf("expected unplanted seed, got planted seed %s", seed.ID)
			}
		}
	})

	t.Run("Get_Existing", func(t *testing.T) {
		seedID := "00000000-0000-4000-b000-000000000001"
		seed, err := store.Get(seedID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if seed.ID != seedID {
			t.Errorf("expected seed ID %s, got %s", seedID, seed.ID)
		}

		if seed.OptimalSoil != "loam" {
			t.Errorf("expected optimal soil 'loam', got '%s'", seed.OptimalSoil)
		}
	})

	t.Run("Get_NonExisting", func(t *testing.T) {
		seedID := "00000000-0000-4000-b000-999999999999"
		_, err := store.Get(seedID)
		if err != models.ErrSeedNotFound {
			t.Errorf("expected ErrSeedNotFound, got %v", err)
		}
	})

	t.Run("Insert", func(t *testing.T) {
		newSeed := &models.Seed{
			OwnerID: "00000000-0000-4000-a000-000000000001",
			Hp:      88.0,
			Planted: false,
			SeedMeta: models.SeedMeta{
				OptimalSoil:   "clay",
				BotanicalName: "Ficus benjamina",
			},
		}

		err := store.Insert(newSeed)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if newSeed.ID == "" {
			t.Error("expected ID to be set after insert")
		}

		if newSeed.CreatedAt.IsZero() {
			t.Error("expected CreatedAt to be set after insert")
		}

		retrievedSeed, err := store.Get(newSeed.ID)
		if err != nil {
			t.Errorf("unexpected error retrieving inserted seed: %v", err)
		}

		if retrievedSeed.BotanicalName != "Ficus benjamina" {
			t.Errorf("expected botanical name 'Ficus benjamina', got '%s'", retrievedSeed.BotanicalName)
		}
	})

	t.Run("MarkAsPlanted", func(t *testing.T) {
		seedID := "00000000-0000-4000-b000-000000000002"

		seed, err := store.Get(seedID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if seed.Planted {
			t.Fatal("seed is already planted before test")
		}

		err = store.MarkAsPlanted(seedID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		seed, err = store.Get(seedID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !seed.Planted {
			t.Error("seed was not marked as planted")
		}
	})

	t.Run("MarkAsPlanted_NonExisting", func(t *testing.T) {
		seedID := "00000000-0000-4000-b000-999999999999"
		err := store.MarkAsPlanted(seedID)
		if err != models.ErrSeedNotFound {
			t.Errorf("expected ErrSeedNotFound, got %v", err)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		seedID := "00000000-0000-4000-b000-000000000003"

		_, err := store.Get(seedID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err = store.Delete(seedID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		_, err = store.Get(seedID)
		if err != models.ErrSeedNotFound {
			t.Errorf("expected ErrSeedNotFound, got %v", err)
		}
	})

	t.Run("Delete_NonExisting", func(t *testing.T) {
		seedID := "00000000-0000-4000-b000-999999999999"
		err := store.Delete(seedID)
		if err != models.ErrSeedNotFound {
			t.Errorf("expected ErrSeedNotFound, got %v", err)
		}
	})
}
