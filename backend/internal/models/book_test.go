package models_test

import (
	"testing"

	"github.com/readr/api/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestIsAccessibleBy(t *testing.T) {
	ownerID := primitive.NewObjectID()
	allowedID := primitive.NewObjectID()
	strangerID := primitive.NewObjectID()

	t.Run("public book is accessible by anyone", func(t *testing.T) {
		book := &models.Book{OwnerID: ownerID, IsPrivate: false}
		if !book.IsAccessibleBy(strangerID) {
			t.Error("public book should be accessible by all users")
		}
	})

	t.Run("private book is accessible by owner", func(t *testing.T) {
		book := &models.Book{OwnerID: ownerID, IsPrivate: true}
		if !book.IsAccessibleBy(ownerID) {
			t.Error("owner should always have access")
		}
	})

	t.Run("private book is accessible by allowed user", func(t *testing.T) {
		book := &models.Book{
			OwnerID:        ownerID,
			IsPrivate:      true,
			AllowedUserIDs: []primitive.ObjectID{allowedID},
		}
		if !book.IsAccessibleBy(allowedID) {
			t.Error("allowed user should have access to private book")
		}
	})

	t.Run("private book is not accessible by stranger", func(t *testing.T) {
		book := &models.Book{
			OwnerID:        ownerID,
			IsPrivate:      true,
			AllowedUserIDs: []primitive.ObjectID{allowedID},
		}
		if book.IsAccessibleBy(strangerID) {
			t.Error("stranger should not have access to private book")
		}
	})
}

func TestIsOwnedBy(t *testing.T) {
	ownerID := primitive.NewObjectID()
	otherID := primitive.NewObjectID()

	book := &models.Book{OwnerID: ownerID}

	if !book.IsOwnedBy(ownerID) {
		t.Error("should return true for the owner")
	}
	if book.IsOwnedBy(otherID) {
		t.Error("should return false for non-owner")
	}
}
