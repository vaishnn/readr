package database

import (
	"context"
	"net/url"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB wraps the MongoDB database and exposes typed collection accessors.
// All collections are accessed through this struct to avoid magic string bugs.
type DB struct {
	client *mongo.Client
	db     *mongo.Database
}

// Connect establishes a connection to MongoDB and pings it to verify.
// Returns an error if the connection or ping fails.
func Connect(uri string) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	// Database name is the last path segment of the URI (e.g. mongodb://host/readr).
	dbName := dbNameFromURI(uri)

	return &DB{client: client, db: client.Database(dbName)}, nil
}

// Disconnect cleanly closes the MongoDB connection. Should be deferred after Connect.
func (d *DB) Disconnect(ctx context.Context) error {
	return d.client.Disconnect(ctx)
}

// Users stores user accounts and credentials.
func (d *DB) Users() *mongo.Collection {
	return d.db.Collection("users")
}

// Books stores book metadata (not the file itself — files live in MinIO).
func (d *DB) Books() *mongo.Collection {
	return d.db.Collection("books")
}

// Progress stores per-user reading position for each book.
func (d *DB) Progress() *mongo.Collection {
	return d.db.Collection("reading_progress")
}

// Highlights stores text highlights made by users inside books.
func (d *DB) Highlights() *mongo.Collection {
	return d.db.Collection("highlights")
}

// HighlightHistory stores an immutable log of every change made to a highlight,
// enabling full undo/redo of highlight edits.
func (d *DB) HighlightHistory() *mongo.Collection {
	return d.db.Collection("highlight_history")
}

// Notes stores user notes, which can be scoped to a book or a specific page.
func (d *DB) Notes() *mongo.Collection {
	return d.db.Collection("notes")
}

// Collections stores user-curated book collections (like playlists for books).
func (d *DB) Collections() *mongo.Collection {
	return d.db.Collection("collections")
}

// Bookmarks stores named page bookmarks within books.
func (d *DB) Bookmarks() *mongo.Collection {
	return d.db.Collection("bookmarks")
}

func dbNameFromURI(uri string) string {
	u, err := url.Parse(uri)
	if err == nil {
		if name := strings.TrimPrefix(u.Path, "/"); name != "" {
			return name
		}
	}
	return "readr"
}
