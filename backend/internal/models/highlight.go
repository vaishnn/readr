package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Highlight struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"userId"`
	BookID    primitive.ObjectID `bson:"book_id" json:"bookId"`
	Page      int                `bson:"page" json:"page"`
	CFIRange  string             `bson:"cfi_range" json:"cfiRange"` // EPUB range string
	Text      string             `bson:"text" json:"text"`
	Color     string             `bson:"color" json:"color"` // yellow, green, blue, pink
	Note      string             `bson:"note" json:"note"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
}

// HighlightHistory is an immutable log of every change made to a highlight.
// Snapshot holds the full highlight document before the action was applied,
// allowing any change to be replayed or reversed.
type HighlightHistory struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	HighlightID primitive.ObjectID `bson:"highlight_id" json:"highlightId"`
	Action      string             `bson:"action" json:"action"` // create, update, delete
	Snapshot    primitive.M        `bson:"snapshot" json:"snapshot"`
	Timestamp   time.Time          `bson:"timestamp" json:"timestamp"`
}
