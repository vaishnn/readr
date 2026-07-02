// Integration tests for book upload, access control, and deletion.
package integration_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/readr/api/internal/models"
	"github.com/readr/api/internal/services"
	"github.com/readr/api/internal/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func newBookService(t *testing.T) *services.BookService {
	t.Helper()
	store, err := storage.NewMinioClient(
		getEnv("TEST_MINIO_ENDPOINT", "localhost:9000"),
		getEnv("TEST_MINIO_ACCESS_KEY", "minioadmin"),
		getEnv("TEST_MINIO_SECRET_KEY", "minioadmin"),
		"readr-test",
		false,
	)
	if err != nil {
		t.Skipf("minio not available: %v", err)
	}
	return services.NewBookService(testDB, store)
}

func TestUploadAndDeleteBook(t *testing.T) {
	svc := newBookService(t)
	userID := primitive.NewObjectID()

	content := bytes.NewReader([]byte("%PDF-1.4 test content"))
	book, err := svc.Upload(context.Background(), userID, "test.pdf", content, int64(content.Len()), nil, 0, models.BookMetadata{})
	if err != nil {
		t.Fatalf("Upload: %v", err)
	}
	if book.ID.IsZero() {
		t.Error("expected non-zero book ID")
	}

	err = svc.Delete(context.Background(), userID, book.ID)
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}

	// Confirm it's gone.
	_, err = svc.Get(context.Background(), userID, book.ID)
	if err != services.ErrBookNotFound {
		t.Errorf("expected ErrBookNotFound after delete, got %v", err)
	}
}

func TestPrivateBookAccessControl(t *testing.T) {
	svc := newBookService(t)
	ownerID := primitive.NewObjectID()
	strangerID := primitive.NewObjectID()
	allowedID := primitive.NewObjectID()

	content := bytes.NewReader([]byte("%PDF-1.4 private"))
	book, err := svc.Upload(context.Background(), ownerID, "private.pdf", content, int64(content.Len()), nil, 0, models.BookMetadata{})
	if err != nil {
		t.Fatalf("Upload: %v", err)
	}

	// Make the book private with one allowed user.
	err = svc.UpdateAccess(context.Background(), ownerID, book.ID, true, []primitive.ObjectID{allowedID})
	if err != nil {
		t.Fatalf("UpdateAccess: %v", err)
	}

	// Stranger should be denied.
	_, err = svc.Get(context.Background(), strangerID, book.ID)
	if err != services.ErrAccessDenied {
		t.Errorf("expected ErrAccessDenied for stranger, got %v", err)
	}

	// Allowed user should have access.
	_, err = svc.Get(context.Background(), allowedID, book.ID)
	if err != nil {
		t.Errorf("expected no error for allowed user, got %v", err)
	}

	// Cleanup
	svc.Delete(context.Background(), ownerID, book.ID)
}

func TestOnlyOwnerCanDelete(t *testing.T) {
	svc := newBookService(t)
	ownerID := primitive.NewObjectID()
	otherID := primitive.NewObjectID()

	content := bytes.NewReader([]byte("%PDF-1.4 owner test"))
	book, _ := svc.Upload(context.Background(), ownerID, "owned.pdf", content, int64(content.Len()), nil, 0, models.BookMetadata{})

	err := svc.Delete(context.Background(), otherID, book.ID)
	if err != services.ErrNotOwner {
		t.Errorf("expected ErrNotOwner, got %v", err)
	}

	// Cleanup
	svc.Delete(context.Background(), ownerID, book.ID)
}
