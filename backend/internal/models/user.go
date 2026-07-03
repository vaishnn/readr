package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserSettings stores per-user UI preferences, persisted alongside the user document.
type UserSettings struct {
	// Theme is either "dark" or "light".
	Theme string `bson:"theme" json:"theme"`
	// DefaultView controls library layout: "grid" or "list".
	DefaultView string `bson:"default_view" json:"defaultView"`
	// LibrarySidebarOpen controls whether the right sidebar is expanded in the library view.
	LibrarySidebarOpen bool `bson:"library_sidebar_open" json:"librarySidebarOpen"`
}

// User represents an account in the system.
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email        string             `bson:"email" json:"email"`
	Username     string             `bson:"username" json:"username"`
	// PasswordHash is a bcrypt hash — never returned in API responses (json:"-").
	PasswordHash string             `bson:"password_hash" json:"-"`
	Settings     UserSettings       `bson:"settings" json:"settings"`
	CreatedAt    time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updatedAt"`
}
