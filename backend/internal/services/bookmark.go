package services

import (
	"context"
	"errors"
	"time"

	"github.com/readr/api/internal/database"
	"github.com/readr/api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrBookmarkNotFound = errors.New("bookmark not found")

type BookmarkService struct {
	db *database.DB
}

func NewBookmarkService(db *database.DB) *BookmarkService {
	return &BookmarkService{db: db}
}

func (s *BookmarkService) List(ctx context.Context, userID, bookID primitive.ObjectID) ([]models.Bookmark, error) {
	cursor, err := s.db.Bookmarks().Find(ctx, bson.M{"user_id": userID, "book_id": bookID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	bookmarks := make([]models.Bookmark, 0)
	return bookmarks, cursor.All(ctx, &bookmarks)
}

func (s *BookmarkService) Create(ctx context.Context, userID, bookID primitive.ObjectID, page int, cfi, label string) (*models.Bookmark, error) {
	b := &models.Bookmark{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		BookID:    bookID,
		Page:      page,
		CFI:       cfi,
		Label:     label,
		CreatedAt: time.Now(),
	}
	_, err := s.db.Bookmarks().InsertOne(ctx, b)
	return b, err
}

func (s *BookmarkService) Delete(ctx context.Context, userID, bookmarkID primitive.ObjectID) error {
	res, err := s.db.Bookmarks().DeleteOne(ctx, bson.M{"_id": bookmarkID, "user_id": userID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrBookmarkNotFound
	}
	return nil
}
