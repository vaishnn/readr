package handlers

import (
	"net/http"

	"github.com/readr/api/internal/services"
)

// FeaturesHandler exposes the current feature flag state to the frontend.
// This endpoint is intentionally read-only and public — flags can only be
// changed directly in MongoDB (kubectl exec + mongosh, or env vars on redeploy).
type FeaturesHandler struct {
	svc *services.FeatureFlagService
}

func NewFeaturesHandler(svc *services.FeatureFlagService) *FeaturesHandler {
	return &FeaturesHandler{svc: svc}
}

func (h *FeaturesHandler) Get(w http.ResponseWriter, r *http.Request) {
	flags, err := h.svc.GetAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load feature flags")
		return
	}
	writeJSON(w, http.StatusOK, flags)
}
