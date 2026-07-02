package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/readr/api/internal/services"
)

type HighlightHandler struct {
	svc *services.HighlightService
}

func NewHighlightHandler(svc *services.HighlightService) *HighlightHandler {
	return &HighlightHandler{svc: svc}
}

func (h *HighlightHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := mustUserID(w, r)
	if !ok {
		return
	}
	bookID, err := parseObjectID(chi.URLParam(r, "bookID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid book id")
		return
	}

	var page *int
	if p := r.URL.Query().Get("page"); p != "" {
		if n, err := strconv.Atoi(p); err == nil {
			page = &n
		}
	}

	highlights, err := h.svc.List(r.Context(), userID, bookID, page)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch highlights")
		return
	}

	writeJSON(w, http.StatusOK, highlights)
}

func (h *HighlightHandler) Create(w http.ResponseWriter, r *http.Request) {
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
		Page     int    `json:"page"`
		CFIRange string `json:"cfiRange"`
		Text     string `json:"text"`
		Color    string `json:"color"`
		Note     string `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Text == "" {
		writeError(w, http.StatusBadRequest, "text is required")
		return
	}

	highlight, err := h.svc.Create(r.Context(), userID, bookID, body.Page, body.CFIRange, body.Text, body.Color, body.Note)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create highlight")
		return
	}

	writeJSON(w, http.StatusCreated, highlight)
}

func (h *HighlightHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := mustUserID(w, r)
	if !ok {
		return
	}
	highlightID, err := parseObjectID(chi.URLParam(r, "highlightID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid highlight id")
		return
	}

	var body struct {
		Color string `json:"color"`
		Note  string `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	highlight, err := h.svc.Update(r.Context(), userID, highlightID, body.Color, body.Note)
	if errors.Is(err, services.ErrHighlightNotFound) {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update highlight")
		return
	}

	writeJSON(w, http.StatusOK, highlight)
}

func (h *HighlightHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := mustUserID(w, r)
	if !ok {
		return
	}
	highlightID, err := parseObjectID(chi.URLParam(r, "highlightID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid highlight id")
		return
	}

	err = h.svc.Delete(r.Context(), userID, highlightID)
	if errors.Is(err, services.ErrHighlightNotFound) {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete highlight")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *HighlightHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	userID, ok := mustUserID(w, r)
	if !ok {
		return
	}
	highlightID, err := parseObjectID(chi.URLParam(r, "highlightID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid highlight id")
		return
	}

	history, err := h.svc.GetHistory(r.Context(), userID, highlightID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch history")
		return
	}

	writeJSON(w, http.StatusOK, history)
}
