package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/readr/api/internal/database"
	"github.com/readr/api/internal/models"
	"github.com/readr/api/internal/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrBookNotFound   = errors.New("book not found")
	ErrAccessDenied   = errors.New("access denied")
	ErrNotOwner       = errors.New("only the owner can perform this action")
	ErrUnsupportedFmt = errors.New("unsupported file format")
)

var supportedFormats = map[string]bool{"pdf": true, "epub": true, "cbz": true}

type BookListParams struct {
	Page   int
	Limit  int
	Search string
	Tag    string
}

type BookService struct {
	db      *database.DB
	storage *storage.MinioClient
}

func NewBookService(db *database.DB, storage *storage.MinioClient) *BookService {
	return &BookService{db: db, storage: storage}
}

func (s *BookService) List(ctx context.Context, userID primitive.ObjectID, p BookListParams) ([]models.Book, int64, error) {
	if p.Limit == 0 {
		p.Limit = 20
	}
	if p.Page == 0 {
		p.Page = 1
	}

	// A user can see: books they own, books shared with them, and all public books.
	filter := bson.M{
		"$or": bson.A{
			bson.M{"owner_id": userID},
			bson.M{"allowed_user_ids": userID},
			bson.M{"is_private": false},
		},
	}

	if p.Search != "" {
		filter["$text"] = bson.M{"$search": p.Search}
	}
	if p.Tag != "" {
		filter["tags"] = p.Tag
	}

	total, err := s.db.Books().CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	skip := int64((p.Page - 1) * p.Limit)
	cursor, err := s.db.Books().Find(ctx, filter, &options.FindOptions{
		Skip:  &skip,
		Limit: int64Ptr(int64(p.Limit)),
		Sort:  bson.M{"uploaded_at": -1},
	})
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	books := make([]models.Book, 0)
	if err := cursor.All(ctx, &books); err != nil {
		return nil, 0, err
	}

	s.populateCoverURLs(ctx, books)
	return books, total, nil
}

func (s *BookService) Upload(ctx context.Context, userID primitive.ObjectID, filename string, file io.Reader, fileSize int64, coverReader io.Reader, coverSize int64, meta models.BookMetadata) (*models.Book, error) {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(filename), "."))
	if !supportedFormats[ext] {
		return nil, ErrUnsupportedFmt
	}

	bookID := primitive.NewObjectID()
	fileKey := fmt.Sprintf("books/%s/%s/book.%s", userID.Hex(), bookID.Hex(), ext)
	coverKey := fmt.Sprintf("books/%s/%s/cover.jpg", userID.Hex(), bookID.Hex())

	contentTypes := map[string]string{
		"pdf":  "application/pdf",
		"epub": "application/epub+zip",
		"cbz":  "application/x-cbz",
	}

	if err := s.storage.Upload(ctx, fileKey, file, fileSize, contentTypes[ext]); err != nil {
		return nil, err
	}

	if coverReader != nil {
		if err := s.storage.Upload(ctx, coverKey, coverReader, coverSize, "image/jpeg"); err != nil {
			// Non-fatal: book is still usable without a cover.
			coverKey = ""
		}
	}

	title := strings.TrimSuffix(filename, filepath.Ext(filename))

	book := &models.Book{
		ID:             bookID,
		OwnerID:        userID,
		Title:          title,
		Author:         meta.Publisher,
		Format:         ext,
		FileKey:        fileKey,
		CoverKey:       coverKey,
		Metadata:       meta,
		Tags:           []string{},
		IsPrivate:      false,
		AllowedUserIDs: []primitive.ObjectID{},
		UploadedAt:     time.Now(),
	}

	if _, err := s.db.Books().InsertOne(ctx, book); err != nil {
		return nil, err
	}

	return book, nil
}

func (s *BookService) Get(ctx context.Context, userID, bookID primitive.ObjectID) (*models.Book, error) {
	book, err := s.findBook(ctx, bookID)
	if err != nil {
		return nil, err
	}
	if !book.IsAccessibleBy(userID) {
		return nil, ErrAccessDenied
	}
	s.populateCoverURLs(ctx, []models.Book{*book})
	return book, nil
}

func (s *BookService) Delete(ctx context.Context, userID, bookID primitive.ObjectID) error {
	book, err := s.findBook(ctx, bookID)
	if err != nil {
		return err
	}
	if !book.IsOwnedBy(userID) {
		return ErrNotOwner
	}

	s.storage.Delete(ctx, book.FileKey)
	if book.CoverKey != "" {
		s.storage.Delete(ctx, book.CoverKey)
	}

	_, err = s.db.Books().DeleteOne(ctx, bson.M{"_id": bookID})
	return err
}

// Stream returns a raw reader for the book file streamed through the backend.
// This avoids exposing internal MinIO hostnames to the browser.
func (s *BookService) Stream(ctx context.Context, userID, bookID primitive.ObjectID) (io.ReadCloser, int64, string, error) {
	book, err := s.findBook(ctx, bookID)
	if err != nil {
		return nil, 0, "", err
	}
	if !book.IsAccessibleBy(userID) {
		return nil, 0, "", ErrAccessDenied
	}
	reader, size, err := s.storage.Stream(ctx, book.FileKey)
	if err != nil {
		return nil, 0, "", err
	}
	return reader, size, book.Format, nil
}

// UpdateAccess changes the privacy setting and allowed user list.
// Only the book owner can call this.
func (s *BookService) UpdateAccess(ctx context.Context, userID, bookID primitive.ObjectID, isPrivate bool, allowedUserIDs []primitive.ObjectID) error {
	book, err := s.findBook(ctx, bookID)
	if err != nil {
		return err
	}
	if !book.IsOwnedBy(userID) {
		return ErrNotOwner
	}

	_, err = s.db.Books().UpdateOne(ctx,
		bson.M{"_id": bookID},
		bson.M{"$set": bson.M{
			"is_private":       isPrivate,
			"allowed_user_ids": allowedUserIDs,
		}},
	)
	return err
}

func (s *BookService) findBook(ctx context.Context, bookID primitive.ObjectID) (*models.Book, error) {
	var book models.Book
	err := s.db.Books().FindOne(ctx, bson.M{"_id": bookID}).Decode(&book)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrBookNotFound
	}
	return &book, err
}

func (s *BookService) populateCoverURLs(ctx context.Context, books []models.Book) {
	for i := range books {
		if books[i].CoverKey != "" {
			url, err := s.storage.PresignedURL(ctx, books[i].CoverKey, 1*time.Hour)
			if err == nil {
				books[i].CoverURL = url
			}
		}
	}
}

func int64Ptr(v int64) *int64 { return &v }
