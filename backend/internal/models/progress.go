package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ReadingProgress tracks where a user left off in a book.
// Restored automatically when the user reopens the book.
type ReadingProgress struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID       primitive.ObjectID `bson:"user_id" json:"userId"`
	BookID       primitive.ObjectID `bson:"book_id" json:"bookId"`
	Page         int                `bson:"page" json:"page"`
	CFI          string             `bson:"cfi" json:"cfi"` // EPUB position identifier
	Percentage   float64            `bson:"percentage" json:"percentage"`
	Zoom         float64            `bson:"zoom" json:"zoom"` // PDF zoom scale, 0 means default
	LastReadAt   time.Time          `bson:"last_read_at" json:"lastReadAt"`
	TotalSeconds int64              `bson:"total_seconds" json:"totalSeconds"`
}
