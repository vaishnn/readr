// Integration tests for auth flows.
// These tests require a running MongoDB and Redis instance.
// Run with: go test ./tests/integration/... -tags integration
package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/readr/api/internal/config"
	"github.com/readr/api/internal/database"
	"github.com/readr/api/internal/handlers"
	"github.com/readr/api/internal/services"
	"github.com/redis/go-redis/v9"
)

var (
	testDB  *database.DB
	testRDB *redis.Client
	testCfg *config.Config
)

func TestMain(m *testing.M) {
	testCfg = &config.Config{
		MongoURI:         getEnv("TEST_MONGO_URI", "mongodb://localhost:27017/readr_test"),
		RedisAddr:        getEnv("TEST_REDIS_ADDR", "localhost:6379"),
		JWTSecret:        []byte("test-secret"),
		JWTRefreshSecret: []byte("test-refresh-secret"),
	}

	var err error
	testDB, err = database.Connect(testCfg.MongoURI)
	if err != nil {
		panic("mongodb not available: " + err.Error())
	}

	testRDB = redis.NewClient(&redis.Options{Addr: testCfg.RedisAddr})
	if err := testRDB.Ping(context.Background()).Err(); err != nil {
		panic("redis not available: " + err.Error())
	}

	code := m.Run()

	// Drop the test database after all tests complete.
	testDB.Disconnect(context.Background())
	os.Exit(code)
}

func TestRegisterAndLogin(t *testing.T) {
	svc := services.NewAuthService(testDB, testRDB, testCfg.JWTSecret, testCfg.JWTRefreshSecret)
	h := handlers.NewAuthHandler(svc)

	// Register
	body, _ := json.Marshal(map[string]string{
		"email":    "test@example.com",
		"username": "testuser",
		"password": "password123",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.Register(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rr.Code, rr.Body.String())
	}

	var result map[string]any
	json.NewDecoder(rr.Body).Decode(&result)
	if result["tokens"] == nil {
		t.Error("expected tokens in response")
	}

	// Login with same credentials
	body, _ = json.Marshal(map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	h.Login(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestDuplicateEmailIsRejected(t *testing.T) {
	svc := services.NewAuthService(testDB, testRDB, testCfg.JWTSecret, testCfg.JWTRefreshSecret)
	h := handlers.NewAuthHandler(svc)

	register := func() *httptest.ResponseRecorder {
		body, _ := json.Marshal(map[string]string{
			"email":    "duplicate@example.com",
			"username": "user_" + t.Name(),
			"password": "password123",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		h.Register(rr, req)
		return rr
	}

	register() // first registration succeeds
	rr := register() // second should conflict

	if rr.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", rr.Code)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
