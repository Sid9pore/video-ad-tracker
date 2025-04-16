package ads

import (
	"encoding/json"
	"net/http"
)

// Handler manages the ads API requests
type Handler struct {
	service *Service
}

// NewHandler creates a new Handler for ads
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GetAdsHandler handles GET requests for fetching ads
func (h *Handler) GetAdsHandler(w http.ResponseWriter, r *http.Request) {
	ads, err := h.service.GetAllAds()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ads)
}

// LogClickHandler handles POST requests to log a click
func (h *Handler) LogClickHandler(w http.ResponseWriter, r *http.Request) {
	var click ClickData
	if err := json.NewDecoder(r.Body).Decode(&click); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if err := h.service.LogClick(click); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
