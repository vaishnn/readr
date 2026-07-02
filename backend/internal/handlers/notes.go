package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/readr/api/internal/services"
)

type NoteHandler struct {
	svc *services.NoteService
}

func NewNoteHandler(svc *services.NoteService) *NoteHandler {
	return &NoteHandler{svc: svc}
}

func (h *NoteHandler) List(w http.ResponseWriter, r *http.Request) {
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

	notes, err := h.svc.List(r.Context(), userID, bookID, page)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch notes")
		return
	}

	writeJSON(w, http.StatusOK, notes)
}

func (h *NoteHandler) Create(w http.ResponseWriter, r *http.Request) {
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
		Page      *int   `json:"page"`
		ContentMD string `json:"contentMd"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	note, err := h.svc.Create(r.Context(), userID, bookID, body.Page, body.ContentMD)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create note")
		return
	}

	writeJSON(w, http.StatusCreated, note)
}

func (h *NoteHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := mustUserID(w, r)
	if !ok {
		return
	}
	noteID, err := parseObjectID(chi.URLParam(r, "noteID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid note id")
		return
	}

	var body struct {
		ContentMD string `json:"contentMd"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	note, err := h.svc.Update(r.Context(), userID, noteID, body.ContentMD)
	if errors.Is(err, services.ErrNoteNotFound) {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update note")
		return
	}

	writeJSON(w, http.StatusOK, note)
}

func (h *NoteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := mustUserID(w, r)
	if !ok {
		return
	}
	noteID, err := parseObjectID(chi.URLParam(r, "noteID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid note id")
		return
	}

	err = h.svc.Delete(r.Context(), userID, noteID)
	if errors.Is(err, services.ErrNoteNotFound) {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete note")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
