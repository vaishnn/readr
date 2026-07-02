package services

import (
	"context"
	"errors"
	"time"

	"github.com/readr/api/internal/database"
	"github.com/readr/api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrHighlightNotFound = errors.New("highlight not found")

type HighlightService struct {
	db *database.DB
}

func NewHighlightService(db *database.DB) *HighlightService {
	return &HighlightService{db: db}
}

func (s *HighlightService) List(ctx context.Context, userID, bookID primitive.ObjectID, page *int) ([]models.Highlight, error) {
	filter := bson.M{"user_id": userID, "book_id": bookID}
	if page != nil {
		filter["page"] = *page
	}

	cursor, err := s.db.Highlights().Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	highlights := make([]models.Highlight, 0)
	return highlights, cursor.All(ctx, &highlights)
}

func (s *HighlightService) Create(ctx context.Context, userID, bookID primitive.ObjectID, page int, cfiRange, text, color, note string) (*models.Highlight, error) {
	now := time.Now()
	h := &models.Highlight{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		BookID:    bookID,
		Page:      page,
		CFIRange:  cfiRange,
		Text:      text,
		Color:     color,
		Note:      note,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if _, err := s.db.Highlights().InsertOne(ctx, h); err != nil {
		return nil, err
	}

	s.writeHistory(ctx, h.ID, "create", nil)
	return h, nil
}

// Update patches a highlight's color or note. The state before the change is
// snapshotted so the edit can be undone.
func (s *HighlightService) Update(ctx context.Context, userID, highlightID primitive.ObjectID, color, note string) (*models.Highlight, error) {
	existing, err := s.findOwned(ctx, highlightID, userID)
	if err != nil {
		return nil, err
	}

	snapshot := highlightToMap(existing)

	update := bson.M{"$set": bson.M{
		"color":      color,
		"note":       note,
		"updated_at": time.Now(),
	}}
	if _, err := s.db.Highlights().UpdateOne(ctx, bson.M{"_id": highlightID}, update); err != nil {
		return nil, err
	}

	s.writeHistory(ctx, highlightID, "update", snapshot)

	existing.Color = color
	existing.Note = note
	return existing, nil
}

// Delete records the final state in history before removing the highlight,
// so the deletion can be reversed if needed.
func (s *HighlightService) Delete(ctx context.Context, userID, highlightID primitive.ObjectID) error {
	existing, err := s.findOwned(ctx, highlightID, userID)
	if err != nil {
		return err
	}

	s.writeHistory(ctx, highlightID, "delete", highlightToMap(existing))

	_, err = s.db.Highlights().DeleteOne(ctx, bson.M{"_id": highlightID})
	return err
}

func (s *HighlightService) GetHistory(ctx context.Context, userID, highlightID primitive.ObjectID) ([]models.HighlightHistory, error) {
	cursor, err := s.db.HighlightHistory().Find(ctx,
		bson.M{"highlight_id": highlightID},
		&options.FindOptions{Sort: bson.M{"timestamp": -1}},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	history := make([]models.HighlightHistory, 0)
	return history, cursor.All(ctx, &history)
}

func (s *HighlightService) findOwned(ctx context.Context, highlightID, userID primitive.ObjectID) (*models.Highlight, error) {
	var h models.Highlight
	err := s.db.Highlights().FindOne(ctx, bson.M{"_id": highlightID, "user_id": userID}).Decode(&h)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrHighlightNotFound
	}
	return &h, err
}

func (s *HighlightService) writeHistory(ctx context.Context, highlightID primitive.ObjectID, action string, snapshot primitive.M) {
	entry := models.HighlightHistory{
		ID:          primitive.NewObjectID(),
		HighlightID: highlightID,
		Action:      action,
		Snapshot:    snapshot,
		Timestamp:   time.Now(),
	}
	s.db.HighlightHistory().InsertOne(ctx, entry)
}

func highlightToMap(h *models.Highlight) primitive.M {
	return primitive.M{
		"_id":       h.ID,
		"user_id":   h.UserID,
		"book_id":   h.BookID,
		"page":      h.Page,
		"cfi_range": h.CFIRange,
		"text":      h.Text,
		"color":     h.Color,
		"note":      h.Note,
		"updated_at": h.UpdatedAt,
	}
}
