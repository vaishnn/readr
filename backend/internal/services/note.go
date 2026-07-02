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
)

var ErrNoteNotFound = errors.New("note not found")

type NoteService struct {
	db *database.DB
}

func NewNoteService(db *database.DB) *NoteService {
	return &NoteService{db: db}
}

// List returns notes for a book. Pass a non-nil page to get page-specific notes;
// pass nil to get book-level notes.
func (s *NoteService) List(ctx context.Context, userID, bookID primitive.ObjectID, page *int) ([]models.Note, error) {
	filter := bson.M{"user_id": userID, "book_id": bookID, "page": page}

	cursor, err := s.db.Notes().Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	notes := make([]models.Note, 0)
	return notes, cursor.All(ctx, &notes)
}

func (s *NoteService) Create(ctx context.Context, userID, bookID primitive.ObjectID, page *int, contentMD string) (*models.Note, error) {
	now := time.Now()
	note := &models.Note{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		BookID:    bookID,
		Page:      page,
		ContentMD: contentMD,
		CreatedAt: now,
		UpdatedAt: now,
	}
	_, err := s.db.Notes().InsertOne(ctx, note)
	return note, err
}

func (s *NoteService) Update(ctx context.Context, userID, noteID primitive.ObjectID, contentMD string) (*models.Note, error) {
	filter := bson.M{"_id": noteID, "user_id": userID}
	update := bson.M{"$set": bson.M{"content_md": contentMD, "updated_at": time.Now()}}

	res, err := s.db.Notes().UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	if res.MatchedCount == 0 {
		return nil, ErrNoteNotFound
	}

	var note models.Note
	err = s.db.Notes().FindOne(ctx, bson.M{"_id": noteID}).Decode(&note)
	return &note, err
}

func (s *NoteService) Delete(ctx context.Context, userID, noteID primitive.ObjectID) error {
	res, err := s.db.Notes().DeleteOne(ctx, bson.M{"_id": noteID, "user_id": userID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrNoteNotFound
	}
	return nil
}

func isNoDocuments(err error) bool {
	return errors.Is(err, mongo.ErrNoDocuments)
}
