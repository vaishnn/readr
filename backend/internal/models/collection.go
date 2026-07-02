package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Collection is a user-curated group of books.
// BookIDs preserves insertion order, which becomes the display order.
type Collection struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID   `bson:"user_id" json:"userId"`
	Name        string               `bson:"name" json:"name"`
	Description string               `bson:"description" json:"description"`
	CoverKey    string               `bson:"cover_key" json:"-"`
	CoverURL    string               `bson:"-" json:"coverUrl,omitempty"`
	BookIDs     []primitive.ObjectID `bson:"book_ids" json:"bookIds"`
	CreatedAt   time.Time            `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time            `bson:"updated_at" json:"updatedAt"`
}
