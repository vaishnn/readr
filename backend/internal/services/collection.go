package services

import (
	"context"
	"errors"
	"time"

	"github.com/readr/api/internal/database"
	"github.com/readr/api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrCollectionNotFound = errors.New("collection not found")

type CollectionService struct {
	db *database.DB
}

func NewCollectionService(db *database.DB) *CollectionService {
	return &CollectionService{db: db}
}

func (s *CollectionService) List(ctx context.Context, userID primitive.ObjectID) ([]models.Collection, error) {
	cursor, err := s.db.Collections().Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	collections := make([]models.Collection, 0)
	return collections, cursor.All(ctx, &collections)
}

func (s *CollectionService) Create(ctx context.Context, userID primitive.ObjectID, name, description string) (*models.Collection, error) {
	now := time.Now()
	c := &models.Collection{
		ID:          primitive.NewObjectID(),
		UserID:      userID,
		Name:        name,
		Description: description,
		BookIDs:     []primitive.ObjectID{},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	_, err := s.db.Collections().InsertOne(ctx, c)
	return c, err
}

func (s *CollectionService) Update(ctx context.Context, userID, collectionID primitive.ObjectID, name, description string) (*models.Collection, error) {
	filter := bson.M{"_id": collectionID, "user_id": userID}
	update := bson.M{"$set": bson.M{"name": name, "description": description, "updated_at": time.Now()}}

	res, err := s.db.Collections().UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	if res.MatchedCount == 0 {
		return nil, ErrCollectionNotFound
	}

	var c models.Collection
	err = s.db.Collections().FindOne(ctx, bson.M{"_id": collectionID}).Decode(&c)
	return &c, err
}

func (s *CollectionService) Delete(ctx context.Context, userID, collectionID primitive.ObjectID) error {
	res, err := s.db.Collections().DeleteOne(ctx, bson.M{"_id": collectionID, "user_id": userID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrCollectionNotFound
	}
	return nil
}

func (s *CollectionService) AddBook(ctx context.Context, userID, collectionID, bookID primitive.ObjectID) error {
	filter := bson.M{"_id": collectionID, "user_id": userID}
	// $addToSet prevents duplicate entries.
	update := bson.M{
		"$addToSet": bson.M{"book_ids": bookID},
		"$set":      bson.M{"updated_at": time.Now()},
	}
	res, err := s.db.Collections().UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrCollectionNotFound
	}
	return nil
}

func (s *CollectionService) RemoveBook(ctx context.Context, userID, collectionID, bookID primitive.ObjectID) error {
	filter := bson.M{"_id": collectionID, "user_id": userID}
	update := bson.M{
		"$pull": bson.M{"book_ids": bookID},
		"$set":  bson.M{"updated_at": time.Now()},
	}
	res, err := s.db.Collections().UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrCollectionNotFound
	}
	return nil
}

func (s *CollectionService) findOwned(ctx context.Context, collectionID, userID primitive.ObjectID) (*models.Collection, error) {
	var c models.Collection
	err := s.db.Collections().FindOne(ctx, bson.M{"_id": collectionID, "user_id": userID}).Decode(&c)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrCollectionNotFound
	}
	return &c, err
}
