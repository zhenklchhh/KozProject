package part_test

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/zhenklchhh/KozProject/inventory/internal/model"
	invPartRepo "github.com/zhenklchhh/KozProject/inventory/internal/repository/part"
)

func setUpDatabase(t *testing.T) (*mongo.Database, func()) {
	t.Helper()
	ctx := t.Context()

	container, err := mongodb.Run(ctx, "mongo:7")
	require.NoError(t, err)

	uri, err := container.ConnectionString(ctx)
	require.NoError(t, err)

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	require.NoError(t, err)

	dbname := strings.ReplaceAll(t.Name(), "/", "_")
	db := client.Database("test_" + dbname)
	cleanup := func() {
		_ = db.Drop(ctx)
		_ = client.Disconnect(ctx)
		_ = container.Terminate(ctx)
	}
	return db, cleanup
}

func TestGetPartSuccess(t *testing.T) {
	db, cleanup := setUpDatabase(t)
	defer cleanup()

	ctx := t.Context()
	partID := uuid.New()
	_, err := db.Collection("parts").InsertOne(ctx, map[string]any{
		"_id":   partID,
		"name":  "rocket engine",
		"price": 12.2,
	})
	require.NoError(t, err)
	repo, err := invPartRepo.NewMongoRepository(db)
	require.NoError(t, err)

	part, err := repo.GetPart(ctx, partID)
	require.NoError(t, err)
	assert.Equal(t, partID, part.GetUUID())
	assert.Equal(t, "rocket engine", part.GetName())
	assert.Equal(t, 12.2, part.Price)
}

func TestGetPartNotFound(t *testing.T) {
	db, cleanup := setUpDatabase(t)
	defer cleanup()

	ctx := t.Context()
	partID := uuid.New()
	repo, err := invPartRepo.NewMongoRepository(db)
	require.NoError(t, err)

	_, err = repo.GetPart(ctx, partID)
	require.Error(t, err)
	assert.Equal(t, err, model.ErrNotFound)
}

func TestListParts(t *testing.T) {
	type seedPart struct {
		UUID         uuid.UUID
		Name         string
		Price        float64
		Tags         []string
		Category     model.Category
		Manufacturer map[string]any
	}

	seeds := []seedPart{
		{
			UUID: uuid.New(), Name: "rocket engine", Price: 12.2,
			Tags: []string{"space", "propulsion"}, Category: model.CategoryEngine,
			Manufacturer: map[string]any{"country": "USA"},
		},
		{
			UUID: uuid.New(), Name: "fuel pump", Price: 5.5,
			Tags: []string{"fuel", "propulsion"}, Category: model.CategoryFuel,
			Manufacturer: map[string]any{"country": "Germany"},
		},
		{
			UUID: uuid.New(), Name: "heat shield", Price: 99.9,
			Tags: []string{"thermal"}, Category: model.CategoryEngine,
			Manufacturer: map[string]any{"country": "USA"},
		},
	}

	insertSeeds := func(t *testing.T, db *mongo.Database, parts []seedPart) {
		t.Helper()
		docs := make([]any, len(parts))
		for i, p := range parts {
			docs[i] = map[string]any{
				"_id":          p.UUID,
				"name":         p.Name,
				"price":        p.Price,
				"tags":         p.Tags,
				"category":     p.Category,
				"manufacturer": p.Manufacturer,
			}
		}
		_, err := db.Collection("parts").InsertMany(t.Context(), docs)
		require.NoError(t, err)
	}

	t.Run("nil filter returns all parts", func(t *testing.T) {
		db, cleanup := setUpDatabase(t)
		defer cleanup()
		insertSeeds(t, db, seeds)

		repo, err := invPartRepo.NewMongoRepository(db)
		require.NoError(t, err)

		parts, err := repo.ListParts(t.Context(), nil)
		require.NoError(t, err)
		assert.Len(t, parts, 3)
	})

	t.Run("filter by uuid", func(t *testing.T) {
		db, cleanup := setUpDatabase(t)
		defer cleanup()
		insertSeeds(t, db, seeds)

		repo, err := invPartRepo.NewMongoRepository(db)
		require.NoError(t, err)

		filter := &model.PartFilter{
			Uuids: []uuid.UUID{seeds[0].UUID, seeds[1].UUID},
		}
		parts, err := repo.ListParts(t.Context(), filter)
		require.NoError(t, err)
		assert.Len(t, parts, 2)
	})

	t.Run("filter by name", func(t *testing.T) {
		db, cleanup := setUpDatabase(t)
		defer cleanup()
		insertSeeds(t, db, seeds)

		repo, err := invPartRepo.NewMongoRepository(db)
		require.NoError(t, err)

		filter := &model.PartFilter{
			Names: []string{"fuel pump"},
		}
		parts, err := repo.ListParts(t.Context(), filter)
		require.NoError(t, err)
		require.Len(t, parts, 1)
		assert.Equal(t, "fuel pump", parts[0].GetName())
	})

	t.Run("filter by tag", func(t *testing.T) {
		db, cleanup := setUpDatabase(t)
		defer cleanup()
		insertSeeds(t, db, seeds)

		repo, err := invPartRepo.NewMongoRepository(db)
		require.NoError(t, err)

		filter := &model.PartFilter{
			Tags: []string{"propulsion"},
		}
		parts, err := repo.ListParts(t.Context(), filter)
		require.NoError(t, err)
		assert.Len(t, parts, 2)
	})

	t.Run("filter by category", func(t *testing.T) {
		db, cleanup := setUpDatabase(t)
		defer cleanup()
		insertSeeds(t, db, seeds)

		repo, err := invPartRepo.NewMongoRepository(db)
		require.NoError(t, err)

		filter := &model.PartFilter{
			Categories: []model.Category{model.CategoryEngine},
		}
		parts, err := repo.ListParts(t.Context(), filter)
		require.NoError(t, err)
		require.Len(t, parts, 2)
		assert.Equal(t, "rocket engine", parts[0].GetName())
		assert.Equal(t, "heat shield", parts[1].GetName())
	})

	t.Run("filter by manufacturer country", func(t *testing.T) {
		db, cleanup := setUpDatabase(t)
		defer cleanup()
		insertSeeds(t, db, seeds)

		repo, err := invPartRepo.NewMongoRepository(db)
		require.NoError(t, err)

		filter := &model.PartFilter{
			ManufacturerCountries: []string{"USA"},
		}
		parts, err := repo.ListParts(t.Context(), filter)
		require.NoError(t, err)
		assert.Len(t, parts, 2)
	})

	t.Run("empty filter returns all parts", func(t *testing.T) {
		db, cleanup := setUpDatabase(t)
		defer cleanup()
		insertSeeds(t, db, seeds)

		repo, err := invPartRepo.NewMongoRepository(db)
		require.NoError(t, err)

		parts, err := repo.ListParts(t.Context(), &model.PartFilter{})
		require.NoError(t, err)
		assert.Len(t, parts, 3)
	})

	t.Run("filter with no matches returns empty slice", func(t *testing.T) {
		db, cleanup := setUpDatabase(t)
		defer cleanup()
		insertSeeds(t, db, seeds)

		repo, err := invPartRepo.NewMongoRepository(db)
		require.NoError(t, err)

		filter := &model.PartFilter{
			Names: []string{"warp drive"},
		}
		parts, err := repo.ListParts(t.Context(), filter)
		require.NoError(t, err)
		assert.Empty(t, parts)
	})

	t.Run("combined filters narrow results", func(t *testing.T) {
		db, cleanup := setUpDatabase(t)
		defer cleanup()
		insertSeeds(t, db, seeds)

		repo, err := invPartRepo.NewMongoRepository(db)
		require.NoError(t, err)

		filter := &model.PartFilter{
			Tags:                  []string{"propulsion"},
			ManufacturerCountries: []string{"USA"},
		}
		parts, err := repo.ListParts(t.Context(), filter)
		require.NoError(t, err)
		require.Len(t, parts, 1)
		assert.Equal(t, "rocket engine", parts[0].GetName())
	})
}
