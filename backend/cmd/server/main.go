package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/readr/api/internal/config"
	"github.com/readr/api/internal/database"
	"github.com/readr/api/internal/handlers"
	"github.com/readr/api/internal/middleware"
	"github.com/readr/api/internal/services"
	"github.com/readr/api/internal/storage"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	cfg := config.Load()

	// --- MongoDB ---
	db, err := database.Connect(cfg.MongoURI)
	if err != nil {
		log.Fatalf("mongodb: %v", err)
	}
	defer db.Disconnect(context.Background())
	ensureIndexes(db)

	// --- Redis ---
	rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("redis: %v", err)
	}

	// --- MinIO ---
	store, err := storage.NewMinioClient(cfg.MinioEndpoint, cfg.MinioAccessKey, cfg.MinioSecretKey, cfg.MinioBucket, cfg.MinioUseSSL, cfg.MinioPublicURL)
	if err != nil {
		log.Fatalf("minio: %v", err)
	}

	// --- Services ---
	authSvc := services.NewAuthService(db, rdb, cfg.JWTSecret, cfg.JWTRefreshSecret)
	bookSvc := services.NewBookService(db, store)
	progressSvc := services.NewProgressService(db)
	highlightSvc := services.NewHighlightService(db)
	noteSvc := services.NewNoteService(db)
	collectionSvc := services.NewCollectionService(db)
	bookmarkSvc := services.NewBookmarkService(db)
	userSvc := services.NewUserService(db)

	// --- Handlers ---
	authH := handlers.NewAuthHandler(authSvc)
	bookH := handlers.NewBookHandler(bookSvc, progressSvc, bookmarkSvc)
	userH := handlers.NewUserHandler(userSvc)
	highlightH := handlers.NewHighlightHandler(highlightSvc)
	noteH := handlers.NewNoteHandler(noteSvc)
	collectionH := handlers.NewCollectionHandler(collectionSvc)

	// --- Router ---
	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.CORS(cfg.Env))

	r.Get("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Route("/api/v1", func(r chi.Router) {
		// Public auth routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authH.Register)
			r.Post("/login", authH.Login)
			r.Post("/refresh", authH.Refresh)
		})

		// All routes below require a valid access token
		r.Group(func(r chi.Router) {
			r.Use(middleware.Authenticate(cfg.JWTSecret))

			r.Delete("/auth/logout", authH.Logout)

			// User / settings
			r.Get("/users/me", userH.GetMe)
			r.Patch("/users/me/password", userH.ChangePassword)
			r.Patch("/users/me/settings", userH.UpdateSettings)

			// Books
			r.Route("/books", func(r chi.Router) {
				r.Get("/", bookH.List)
				r.Post("/", bookH.Upload)

				r.Route("/{bookID}", func(r chi.Router) {
					r.Get("/", bookH.Get)
					r.Delete("/", bookH.Delete)
					r.Get("/stream", bookH.Stream)
					r.Patch("/access", bookH.UpdateAccess)

					r.Get("/progress", bookH.GetProgress)
					r.Put("/progress", bookH.SaveProgress)

					r.Get("/bookmarks", bookH.ListBookmarks)
					r.Post("/bookmarks", bookH.CreateBookmark)
					r.Delete("/bookmarks/{bookmarkID}", bookH.DeleteBookmark)

					r.Get("/highlights", highlightH.List)
					r.Post("/highlights", highlightH.Create)
					r.Patch("/highlights/{highlightID}", highlightH.Update)
					r.Delete("/highlights/{highlightID}", highlightH.Delete)
					r.Get("/highlights/{highlightID}/history", highlightH.GetHistory)

					r.Get("/notes", noteH.List)
					r.Post("/notes", noteH.Create)
					r.Patch("/notes/{noteID}", noteH.Update)
					r.Delete("/notes/{noteID}", noteH.Delete)
				})
			})

			// Collections
			r.Route("/collections", func(r chi.Router) {
				r.Get("/", collectionH.List)
				r.Post("/", collectionH.Create)
				r.Patch("/{collectionID}", collectionH.Update)
				r.Delete("/{collectionID}", collectionH.Delete)
				r.Post("/{collectionID}/books", collectionH.AddBook)
				r.Delete("/{collectionID}/books/{bookID}", collectionH.RemoveBook)
			})
		})
	})

	// --- Server with graceful shutdown ---
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("server listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

// ensureIndexes creates MongoDB indexes on startup.
// Using compound indexes on (user_id, book_id) for the hot query paths,
// and a text index on books for full-text search.
func ensureIndexes(db *database.DB) {
	ctx := context.Background()

	compound := func(coll *mongo.Collection) {
		coll.Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "book_id", Value: 1}},
		})
	}

	compound(db.Progress())
	compound(db.Highlights())
	compound(db.Notes())
	compound(db.Bookmarks())

	db.HighlightHistory().Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "highlight_id", Value: 1}, {Key: "timestamp", Value: -1}},
	})

	// Text index powers the /books?search= query.
	db.Books().Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "title", Value: "text"},
			{Key: "author", Value: "text"},
			{Key: "tags", Value: "text"},
		},
	})

	db.Users().Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "email", Value: 1}},
	})
}
