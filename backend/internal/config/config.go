package config

import "os"

// Config holds all runtime configuration loaded from environment variables.
// In development, sensible defaults are used so the app runs without any setup.
type Config struct {
	Port string
	// Env controls behaviour like verbose logging and relaxed CORS. Values: "development", "production"
	Env string

	// MongoURI is the full MongoDB connection string, including database name.
	MongoURI string

	RedisAddr string

	// MinioEndpoint is host:port of the MinIO server (no scheme).
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	// MinioBucket is the single bucket used for all file storage (books, covers).
	MinioBucket string
	// MinioUseSSL should be true in production when MinIO is behind TLS.
	MinioUseSSL bool
	// MinioPublicURL is the browser-reachable base URL for MinIO (used in presigned URLs).
	// In dev this is http://localhost:9000; in prod it's your public MinIO domain.
	MinioPublicURL string

	// JWTSecret signs short-lived access tokens (15 min expiry).
	JWTSecret []byte
	// JWTRefreshSecret signs long-lived refresh tokens (7 day expiry).
	JWTRefreshSecret []byte
}

// Load reads configuration from environment variables, falling back to
// development defaults when a variable is not set.
func Load() *Config {
	return &Config{
		Port:             getEnv("PORT", "8080"),
		Env:              getEnv("ENV", "development"),
		MongoURI:         getEnv("MONGO_URI", "mongodb://localhost:27017/readr"),
		RedisAddr:        getEnv("REDIS_ADDR", "localhost:6379"),
		MinioEndpoint:    getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinioAccessKey:   getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinioSecretKey:   getEnv("MINIO_SECRET_KEY", "minioadmin"),
		MinioBucket:      getEnv("MINIO_BUCKET", "readr"),
		MinioUseSSL:      getEnv("MINIO_USE_SSL", "false") == "true",
		MinioPublicURL:   getEnv("MINIO_PUBLIC_URL", "http://localhost:9000"),
		JWTSecret:        []byte(getEnv("JWT_SECRET", "dev-jwt-secret")),
		JWTRefreshSecret: []byte(getEnv("JWT_REFRESH_SECRET", "dev-refresh-secret")),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
