package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/readr/api/internal/services"
)

type CollectionHandler struct {
	svc *services.CollectionService
}

func NewCollectionHandler(svc *services.CollectionService) *CollectionHandler {
	return &CollectionHandler{svc: svc}
}

func (h *CollectionHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := mustUserID(w, r)
	if !ok {
		return
	}

	collections, err := h.svc.List(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch collections")
		return
	}

	writeJSON(w, http.StatusOK, collections)
}

func (h *CollectionHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := mustUserID(w, r)
	if !ok {
		return
	}

	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	collection, err := h.svc.Create(r.Context(), userID, body.Name, body.Description)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create collection")
		return
	}

	writeJSON(w, http.StatusCreated, collection)
}

func (h *CollectionHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := mustUserID(w, r)
	if !ok {
		return
	}
	collectionID, err := parseObjectID(chi.URLParam(r, "collectionID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid collection id")
		return
	}

	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	collection, err := h.svc.Update(r.Context(), userID, collectionID, body.Name, body.Description)
	if errors.Is(err, services.ErrCollectionNotFound) {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update collection")
		return
	}

	writeJSON(w, http.StatusOK, collection)
}

func (h *CollectionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := mustUserID(w, r)
	if !ok {
		return
	}
	collectionID, err := parseObjectID(chi.URLParam(r, "collectionID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid collection id")
		return
	}

	err = h.svc.Delete(r.Context(), userID, collectionID)
	if errors.Is(err, services.ErrCollectionNotFound) {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete collection")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CollectionHandler) AddBook(w http.ResponseWriter, r *http.Request) {
	userID, ok := mustUserID(w, r)
	if !ok {
		return
	}
	collectionID, err := parseObjectID(chi.URLParam(r, "collectionID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid collection id")
		return
	}

	var body struct {
		BookID string `json:"bookId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.BookID == "" {
		writeError(w, http.StatusBadRequest, "bookId is required")
		return
	}
	bookID, err := parseObjectID(body.BookID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid book id")
		return
	}

	err = h.svc.AddBook(r.Context(), userID, collectionID, bookID)
	if errors.Is(err, services.ErrCollectionNotFound) {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to add book")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CollectionHandler) RemoveBook(w http.ResponseWriter, r *http.Request) {
	userID, ok := mustUserID(w, r)
	if !ok {
		return
	}
	collectionID, err := parseObjectID(chi.URLParam(r, "collectionID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid collection id")
		return
	}
	bookID, err := parseObjectID(chi.URLParam(r, "bookID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid book id")
		return
	}

	err = h.svc.RemoveBook(r.Context(), userID, collectionID, bookID)
	if errors.Is(err, services.ErrCollectionNotFound) {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to remove book")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
