package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Bookmark is a named location in a book, explicitly saved by the user.
// Unlike reading progress (automatic), bookmarks are intentional markers.
type Bookmark struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"userId"`
	BookID    primitive.ObjectID `bson:"book_id" json:"bookId"`
	Page      int                `bson:"page" json:"page"`
	CFI       string             `bson:"cfi" json:"cfi"` // EPUB position
	Label     string             `bson:"label" json:"label"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
}
