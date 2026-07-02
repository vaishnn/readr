package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/readr/api/internal/middleware"
	"github.com/readr/api/internal/models"
	"github.com/readr/api/internal/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookHandler struct {
	svc      *services.BookService
	progress *services.ProgressService
	bookmark *services.BookmarkService
}

func NewBookHandler(svc *services.BookService, progress *services.ProgressService, bookmark *services.BookmarkService) *BookHandler {
	return &BookHandler{svc: svc, progress: progress, bookmark: bookmark}
}

func (h *BookHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, _ := mustUserID(w, r)

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	books, total, err := h.svc.List(r.Context(), userID, services.BookListParams{
		Page:   page,
		Limit:  limit,
		Search: r.URL.Query().Get("search"),
		Tag:    r.URL.Query().Get("tag"),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch books")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"books": books, "total": total})
}

func (h *BookHandler) Upload(w http.ResponseWriter, r *http.Request) {
	userID, _ := mustUserID(w, r)

	// Limit upload size to 500MB.
	r.ParseMultipartForm(500 << 20)

	file, header, err := r.FormFile("book")
	if err != nil {
		writeError(w, http.StatusBadRequest, "book file is required")
		return
	}
	defer file.Close()

	var coverReader interface{ Read([]byte) (int, error) }
	var coverSize int64
	coverFile, coverHeader, err := r.FormFile("cover")
	if err == nil {
		defer coverFile.Close()
		coverReader = coverFile
		coverSize = coverHeader.Size
	}

	var meta models.BookMetadata
	if raw := r.FormValue("metadata"); raw != "" {
		json.Unmarshal([]byte(raw), &meta)
	}

	book, err := h.svc.Upload(r.Context(), userID, header.Filename, file, header.Size, coverReader, coverSize, meta)
	if errors.Is(err, services.ErrUnsupportedFmt) {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "upload failed")
		return
	}

	writeJSON(w, http.StatusCreated, book)
}

func (h *BookHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := mustUserID(w, r)
	if !ok {
		return
	}
	bookID, err := parseObjectID(chi.URLParam(r, "bookID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid book id")
		return
	}

	book, err := h.svc.Get(r.Context(), userID, bookID)
	if errors.Is(err, services.ErrBookNotFound) {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	if errors.Is(err, services.ErrAccessDenied) {
		writeError(w, http.StatusForbidden, err.Error())
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch book")
		return
	}

	writeJSON(w, http.StatusOK, book)
}

func (h *BookHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := mustUserID(w, r)
	if !ok {
		return
	}
	bookID, err := parseObjectID(chi.URLParam(r, "bookID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid book id")
		return
	}

	err = h.svc.Delete(r.Context(), userID, bookID)
	if errors.Is(err, services.ErrBookNotFound) {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	if errors.Is(err, services.ErrNotOwner) {
		writeError(w, http.StatusForbidden, err.Error())
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "delete failed")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Stream proxies the book file through the backend so the browser never needs
// to reach MinIO directly.
func (h *BookHandler) Stream(w http.ResponseWriter, r *http.Request) {
	userID, ok := mustUserID(w, r)
	if !ok {
		return
	}
	bookID, err := parseObjectID(chi.URLParam(r, "bookID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid book id")
		return
	}

	reader, size, format, err := h.svc.Stream(r.Context(), userID, bookID)
	if errors.Is(err, services.ErrAccessDenied) {
		writeError(w, http.StatusForbidden, err.Error())
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to stream book")
		return
	}
	defer reader.Close()

	contentTypes := map[string]string{
		"pdf":  "application/pdf",
		"epub": "application/epub+zip",
		"cbz":  "application/x-cbz",
	}
	ct := contentTypes[format]
	if ct == "" {
		ct = "application/octet-stream"
	}

	w.Header().Set("Content-Type", ct)
	w.Header().Set("Content-Length", strconv.FormatInt(size, 10))
	w.Header().Set("Content-Disposition", "inline")
	w.Header().Set("Cache-Control", "private, max-age=300")
	if _, err := io.Copy(w, reader); err != nil {
		// Headers already sent; nothing to do.
		return
	}
}

func (h *BookHandler) UpdateAccess(w http.ResponseWriter, r *http.Request) {
	userID, ok := mustUserID(w, r)
	if !ok {
		return
	}
	bookID, err := parseObjectID(chi.URLParam(r, "bookID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid book id")
		return
	}

	var body struct {
		IsPrivate      bool     `json:"isPrivate"`
		AllowedUserIDs []string `json:"allowedUserIds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	allowedIDs := make([]primitive.ObjectID, 0, len(body.AllowedUserIDs))
	for _, id := range body.AllowedUserIDs {
		oid, err := parseObjectID(id)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid user id: "+id)
			return
		}
		allowedIDs = append(allowedIDs, oid)
	}

	err = h.svc.UpdateAccess(r.Context(), userID, bookID, body.IsPrivate, allowedIDs)
	if errors.Is(err, services.ErrNotOwner) {
		writeError(w, http.StatusForbidden, err.Error())
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update access")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *BookHandler) GetProgress(w http.ResponseWriter, r *http.Request) {
	userID, ok := mustUserID(w, r)
	if !ok {
		return
	}
	bookID, err := parseObjectID(chi.URLParam(r, "bookID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid book id")
		return
	}

	progress, err := h.progress.Get(r.Context(), userID, bookID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch progress")
		return
	}

	writeJSON(w, http.StatusOK, progress)
}

func (h *BookHandler) SaveProgress(w http.ResponseWriter, r *http.Request) {
	userID, ok := mustUserID(w, r)
	if !ok {
		return
	}
	bookID, err := parseObjectID(chi.URLParam(r, "bookID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid book id")
		return
	}

	var body struct {
		Page           int     `json:"page"`
		CFI            string  `json:"cfi"`
		Percentage     float64 `json:"percentage"`
		Zoom           float64 `json:"zoom"`
		SessionSeconds int64   `json:"sessionSeconds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err = h.progress.Save(r.Context(), userID, bookID, body.Page, body.CFI, body.Percentage, body.Zoom, body.SessionSeconds)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save progress")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *BookHandler) ListBookmarks(w http.ResponseWriter, r *http.Request) {
	userID, ok := mustUserID(w, r)
	if !ok {
		return
	}
	bookID, err := parseObjectID(chi.URLParam(r, "bookID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid book id")
		return
	}

	bookmarks, err := h.bookmark.List(r.Context(), userID, bookID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch bookmarks")
		return
	}

	writeJSON(w, http.StatusOK, bookmarks)
}

func (h *BookHandler) CreateBookmark(w http.ResponseWriter, r *http.Request) {
	userID, ok := mustUserID(w, r)
	if !ok {
		return
	}
	bookID, err := parseObjectID(chi.URLParam(r, "bookID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid book id")
		return
	}

	var body struct {
		Page  int    `json:"page"`
		CFI   string `json:"cfi"`
		Label string `json:"label"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	bookmark, err := h.bookmark.Create(r.Context(), userID, bookID, body.Page, body.CFI, body.Label)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create bookmark")
		return
	}

	writeJSON(w, http.StatusCreated, bookmark)
}

func (h *BookHandler) DeleteBookmark(w http.ResponseWriter, r *http.Request) {
	userID, ok := mustUserID(w, r)
	if !ok {
		return
	}
	bookmarkID, err := parseObjectID(chi.URLParam(r, "bookmarkID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid bookmark id")
		return
	}

	err = h.bookmark.Delete(r.Context(), userID, bookmarkID)
	if errors.Is(err, services.ErrBookmarkNotFound) {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete bookmark")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// mustUserID extracts the authenticated user ID from context and writes a 401 if missing.
func mustUserID(w http.ResponseWriter, r *http.Request) (primitive.ObjectID, bool) {
	raw, ok := middleware.UserIDFromCtx(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return primitive.NilObjectID, false
	}
	id, err := parseObjectID(raw)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "invalid user id in token")
		return primitive.NilObjectID, false
	}
	return id, true
}
