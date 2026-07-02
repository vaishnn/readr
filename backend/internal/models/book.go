package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookMetadata struct {
	Publisher   string `bson:"publisher,omitempty" json:"publisher,omitempty"`
	Year        int    `bson:"year,omitempty" json:"year,omitempty"`
	Language    string `bson:"language,omitempty" json:"language,omitempty"`
	PageCount   int    `bson:"page_count,omitempty" json:"pageCount,omitempty"`
	Description string `bson:"description,omitempty" json:"description,omitempty"`
	ISBN        string `bson:"isbn,omitempty" json:"isbn,omitempty"`
}

// Book stores metadata and access rules. The actual file lives in MinIO.
//
// Access rules:
//   - IsPrivate=false: any authenticated user can find and read the book.
//   - IsPrivate=true: only the owner and AllowedUserIDs can access it.
//   - Only the owner can change IsPrivate or AllowedUserIDs.
type Book struct {
	ID             primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	OwnerID        primitive.ObjectID   `bson:"owner_id" json:"ownerId"`
	Title          string               `bson:"title" json:"title"`
	Author         string               `bson:"author" json:"author"`
	Format         string               `bson:"format" json:"format"` // pdf, epub, cbz
	FileKey        string               `bson:"file_key" json:"-"`
	CoverKey       string               `bson:"cover_key" json:"-"`
	CoverURL       string               `bson:"-" json:"coverUrl,omitempty"`
	Metadata       BookMetadata         `bson:"metadata" json:"metadata"`
	Tags           []string             `bson:"tags" json:"tags"`
	IsPrivate      bool                 `bson:"is_private" json:"isPrivate"`
	AllowedUserIDs []primitive.ObjectID `bson:"allowed_user_ids" json:"allowedUserIds"`
	UploadedAt     time.Time            `bson:"uploaded_at" json:"uploadedAt"`
}

// IsAccessibleBy checks if the given user can read this book.
// The owner always has access; for private books, the user must be in AllowedUserIDs.
func (b *Book) IsAccessibleBy(userID primitive.ObjectID) bool {
	if !b.IsPrivate || b.OwnerID == userID {
		return true
	}
	for _, id := range b.AllowedUserIDs {
		if id == userID {
			return true
		}
	}
	return false
}

func (b *Book) IsOwnedBy(userID primitive.ObjectID) bool {
	return b.OwnerID == userID
}
