package services

import (
	"context"
	"errors"
	"time"

	"github.com/readr/api/internal/auth"
	"github.com/readr/api/internal/database"
	"github.com/readr/api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrWrongPassword = errors.New("current password is incorrect")
	ErrWeakPassword  = errors.New("password must be at least 8 characters")
)

type UserService struct {
	db *database.DB
}

func NewUserService(db *database.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetMe(ctx context.Context, userID primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := s.db.Users().FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrUserNotFound
	}
	return &user, err
}

func (s *UserService) ChangePassword(ctx context.Context, userID primitive.ObjectID, currentPassword, newPassword string) error {
	if len(newPassword) < 8 {
		return ErrWeakPassword
	}

	var user models.User
	err := s.db.Users().FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return ErrUserNotFound
	}
	if err != nil {
		return err
	}

	if !auth.CheckPassword(user.PasswordHash, currentPassword) {
		return ErrWrongPassword
	}

	hash, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}

	_, err = s.db.Users().UpdateOne(ctx,
		bson.M{"_id": userID},
		bson.M{"$set": bson.M{
			"password_hash": hash,
			"updated_at":    time.Now(),
		}},
	)
	return err
}

func (s *UserService) UpdateSettings(ctx context.Context, userID primitive.ObjectID, settings models.UserSettings) (*models.User, error) {
	_, err := s.db.Users().UpdateOne(ctx,
		bson.M{"_id": userID},
		bson.M{"$set": bson.M{
			"settings":   settings,
			"updated_at": time.Now(),
		}},
	)
	if err != nil {
		return nil, err
	}
	return s.GetMe(ctx, userID)
}
