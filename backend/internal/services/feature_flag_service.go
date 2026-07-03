package services

import (
	"context"
	"time"

	"github.com/readr/api/internal/config"
	"github.com/readr/api/internal/database"
	"github.com/readr/api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// seed defines every known flag: its key, human-readable description, and the
// default value read from the environment on first startup.
type flagSeed struct {
	key         string
	description string
	enabled     bool
}

type FeatureFlagService struct {
	db *database.DB
}

func NewFeatureFlagService(db *database.DB) *FeatureFlagService {
	return &FeatureFlagService{db: db}
}

// Seed inserts any flags that do not yet exist in MongoDB.
// Existing documents are never touched — changes must be made directly in the DB.
// Call this once on server startup, after the DB connection is established.
func (s *FeatureFlagService) Seed(ctx context.Context, cfg config.FeatureFlags) error {
	seeds := []flagSeed{
		{"collections", "Organize books into named collections / shelves.", cfg.Collections},
		{"public-library", "Allow unauthenticated users to browse non-private books.", cfg.PublicLibrary},
		{"reading-stats", "Reading time analytics and progress charts.", cfg.ReadingStats},
		{"highlights", "Text highlighting and inline notes inside the reader.", cfg.Highlights},
		{"registration", "Allow new users to self-register (disable for single-user instances).", cfg.Registration},
		{"social-sharing", "Generate shareable links for books and collections.", cfg.SocialSharing},
	}

	coll := s.db.FeatureFlags()
	for _, seed := range seeds {
		filter := bson.M{"key": seed.key}
		update := bson.M{
			"$setOnInsert": bson.M{
				"key":         seed.key,
				"enabled":     seed.enabled,
				"description": seed.description,
				"updated_at":  time.Now(),
			},
		}
		_, err := coll.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
		if err != nil {
			return err
		}
	}
	return nil
}

// GetAll returns all feature flags as a map[key]enabled for the frontend.
func (s *FeatureFlagService) GetAll(ctx context.Context) (map[string]bool, error) {
	cursor, err := s.db.FeatureFlags().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var flags []models.FeatureFlag
	if err := cursor.All(ctx, &flags); err != nil {
		return nil, err
	}

	result := make(map[string]bool, len(flags))
	for _, f := range flags {
		result[f.Key] = f.Enabled
	}
	return result, nil
}

// EnsureIndex creates a unique index on `key` so duplicate seeds are safe.
func (s *FeatureFlagService) EnsureIndex(ctx context.Context) error {
	_, err := s.db.FeatureFlags().Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "key", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return err
}
