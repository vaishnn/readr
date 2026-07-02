package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Note is a markdown note attached to either a whole book (Page=nil) or a specific page.
type Note struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"userId"`
	BookID    primitive.ObjectID `bson:"book_id" json:"bookId"`
	Page      *int               `bson:"page" json:"page"`
	ContentMD string             `bson:"content_md" json:"contentMd"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
}
