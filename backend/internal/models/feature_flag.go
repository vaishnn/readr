package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FeatureFlag is a single row in the feature_flags collection.
// Flags are seeded from environment variables on first startup;
// after that, they must be changed directly in MongoDB.
type FeatureFlag struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"  json:"-"`
	Key         string             `bson:"key"            json:"key"`
	Enabled     bool               `bson:"enabled"        json:"enabled"`
	Description string             `bson:"description"    json:"description"`
	UpdatedAt   time.Time          `bson:"updated_at"     json:"-"`
}
