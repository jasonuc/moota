package store

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jasonuc/moota/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestSeedStore(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping seed store integration tests")
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

	t.Run("GetByOwnerID", func(t *testing.T) {
		ownerID := "00000000-0000-4000-a000-000000000001"
		seeds, err := store.GetByOwnerID(context.Background(), ownerID)
		assert.NoError(t, err)
		assert.Len(t, seeds, 3, "expected 3 seeds")

		for _, seed := range seeds {
			assert.Falsef(t, seed.Planted, "expected unplanted seed, got planted seed %s", seed.ID)
		}
	})

	t.Run("Get_Existing", func(t *testing.T) {
		seedID := "00000000-0000-4000-b000-000000000001"
		seed, err := store.Get(context.Background(), seedID)
		assert.NoError(t, err)
		assert.Equalf(t, seedID, seed.ID, "expected seed ID %s, got %s", seedID, seed.ID)
		assert.EqualValuesf(t, "loam", seed.OptimalSoil, "expected optimal soil 'loam', got '%s'", seed.OptimalSoil)
	})

	t.Run("Get_NonExisting", func(t *testing.T) {
		seedID := "00000000-0000-4000-b000-999999999999"
		_, err := store.Get(context.Background(), seedID)
		assert.ErrorIs(t, err, models.ErrSeedNotFound, "expected ErrSeedNotFound")
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

		err := store.Insert(context.Background(), newSeed)
		assert.NoError(t, err)
		assert.NotEmpty(t, newSeed.ID, "expected ID to be set after insert")
		assert.False(t, newSeed.CreatedAt.IsZero(), "expected CreatedAt to be set after insert")

		retrievedSeed, err := store.Get(context.Background(), newSeed.ID)
		assert.NoError(t, err, "unexpected error retrieving inserted seed")
		assert.Equalf(t, "Ficus benjamina", retrievedSeed.BotanicalName,
			"expected botanical name 'Ficus benjamina', got '%s'", retrievedSeed.BotanicalName)
	})

	t.Run("MarkAsPlanted", func(t *testing.T) {
		seedID := "00000000-0000-4000-b000-000000000002"

		seed, err := store.Get(context.Background(), seedID)
		assert.NoError(t, err)
		assert.False(t, seed.Planted, "seed is already planted before test")

		err = store.MarkAsPlanted(context.Background(), seedID)
		assert.NoError(t, err)

		seed, err = store.Get(context.Background(), seedID)
		assert.NoError(t, err)
		assert.True(t, seed.Planted, "seed was not marked as planted")
	})

	t.Run("MarkAsPlanted_NonExisting", func(t *testing.T) {
		seedID := "00000000-0000-4000-b000-999999999999"
		err := store.MarkAsPlanted(context.Background(), seedID)
		assert.ErrorIs(t, err, models.ErrSeedNotFound, "expected ErrSeedNotFound")
	})

	t.Run("Delete", func(t *testing.T) {
		seedID := "00000000-0000-4000-b000-000000000003"

		_, err := store.Get(context.Background(), seedID)
		assert.NoError(t, err)

		err = store.Delete(context.Background(), seedID)
		assert.NoError(t, err)

		_, err = store.Get(context.Background(), seedID)
		assert.ErrorIs(t, err, models.ErrSeedNotFound, "expected ErrSeedNotFound")
	})

	t.Run("Delete_NonExisting", func(t *testing.T) {
		seedID := "00000000-0000-4000-b000-999999999999"
		err := store.Delete(context.Background(), seedID)
		assert.ErrorIs(t, err, models.ErrSeedNotFound, "expected ErrSeedNotFound")
	})
}
