package services

import (
	"context"
	"time"

	"github.com/readr/api/internal/database"
	"github.com/readr/api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProgressService struct {
	db *database.DB
}

func NewProgressService(db *database.DB) *ProgressService {
	return &ProgressService{db: db}
}

func (s *ProgressService) Get(ctx context.Context, userID, bookID primitive.ObjectID) (*models.ReadingProgress, error) {
	var p models.ReadingProgress
	err := s.db.Progress().FindOne(ctx, bson.M{"user_id": userID, "book_id": bookID}).Decode(&p)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &p, err
}

// Save upserts the reading position. sessionSeconds is the time spent in the
// current session and is added to the running total.
func (s *ProgressService) Save(ctx context.Context, userID, bookID primitive.ObjectID, page int, cfi string, percentage float64, zoom float64, sessionSeconds int64) error {
	filter := bson.M{"user_id": userID, "book_id": bookID}
	update := bson.M{
		"$set": bson.M{
			"page":         page,
			"cfi":          cfi,
			"percentage":   percentage,
			"zoom":         zoom,
			"last_read_at": time.Now(),
		},
		"$inc": bson.M{"total_seconds": sessionSeconds},
		"$setOnInsert": bson.M{
			"_id":     primitive.NewObjectID(),
			"user_id": userID,
			"book_id": bookID,
		},
	}
	_, err := s.db.Progress().UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return err
}
