package services

import (
	"context"
	"errors"
	"time"

	"github.com/readr/api/internal/auth"
	"github.com/readr/api/internal/database"
	"github.com/readr/api/internal/models"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrEmailTaken    = errors.New("email already in use")
	ErrUsernameTaken = errors.New("username already in use")
	ErrInvalidCreds  = errors.New("invalid email or password")
	ErrInvalidToken  = errors.New("invalid or expired token")
)

type AuthService struct {
	db            *database.DB
	redis         *redis.Client
	accessSecret  []byte
	refreshSecret []byte
}

func NewAuthService(db *database.DB, redis *redis.Client, accessSecret, refreshSecret []byte) *AuthService {
	return &AuthService{db: db, redis: redis, accessSecret: accessSecret, refreshSecret: refreshSecret}
}

func (s *AuthService) Register(ctx context.Context, email, username, password string) (*models.User, auth.TokenPair, error) {
	coll := s.db.Users()

	if err := s.assertUnique(ctx, coll, "email", email, ErrEmailTaken); err != nil {
		return nil, auth.TokenPair{}, err
	}
	if err := s.assertUnique(ctx, coll, "username", username, ErrUsernameTaken); err != nil {
		return nil, auth.TokenPair{}, err
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		return nil, auth.TokenPair{}, err
	}

	now := time.Now()
	user := &models.User{
		ID:           primitive.NewObjectID(),
		Email:        email,
		Username:     username,
		PasswordHash: hash,
		Settings:     models.UserSettings{Theme: "dark", DefaultView: "grid"},
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if _, err := coll.InsertOne(ctx, user); err != nil {
		return nil, auth.TokenPair{}, err
	}

	tokens, err := auth.GenerateTokenPair(user.ID.Hex(), s.accessSecret, s.refreshSecret)
	return user, tokens, err
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*models.User, auth.TokenPair, error) {
	var user models.User
	err := s.db.Users().FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) || !auth.CheckPassword(user.PasswordHash, password) {
		return nil, auth.TokenPair{}, ErrInvalidCreds
	}
	if err != nil {
		return nil, auth.TokenPair{}, err
	}

	tokens, err := auth.GenerateTokenPair(user.ID.Hex(), s.accessSecret, s.refreshSecret)
	return &user, tokens, err
}

// Refresh validates a refresh token and issues a new token pair.
// The old refresh token is blacklisted to prevent reuse.
func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (auth.TokenPair, error) {
	claims, err := auth.ValidateToken(refreshToken, s.refreshSecret)
	if err != nil || claims.Type != auth.TokenTypeRefresh {
		return auth.TokenPair{}, ErrInvalidToken
	}

	if s.isBlacklisted(ctx, refreshToken) {
		return auth.TokenPair{}, ErrInvalidToken
	}

	ttl := time.Until(claims.ExpiresAt.Time)
	s.redis.Set(ctx, blacklistKey(refreshToken), 1, ttl)

	return auth.GenerateTokenPair(claims.UserID, s.accessSecret, s.refreshSecret)
}

// Logout blacklists the access token for its remaining lifetime.
func (s *AuthService) Logout(ctx context.Context, accessToken string) error {
	claims, err := auth.ValidateToken(accessToken, s.accessSecret)
	if err != nil {
		return nil // already invalid, treat as success
	}
	ttl := time.Until(claims.ExpiresAt.Time)
	return s.redis.Set(ctx, blacklistKey(accessToken), 1, ttl).Err()
}

func (s *AuthService) isBlacklisted(ctx context.Context, token string) bool {
	return s.redis.Exists(ctx, blacklistKey(token)).Val() > 0
}

func blacklistKey(token string) string {
	return "blacklist:" + token
}

func (s *AuthService) assertUnique(ctx context.Context, coll *mongo.Collection, field, value string, errMsg error) error {
	count, err := coll.CountDocuments(ctx, bson.M{field: value})
	if err != nil {
		return err
	}
	if count > 0 {
		return errMsg
	}
	return nil
}
