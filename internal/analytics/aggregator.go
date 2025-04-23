package analytics

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

// Aggregator defines methods to analyze click data.
type Aggregator interface {
	GetClickMetrics(startTime, endTime time.Time) (ClickMetrics, error)
}

// PostgresAggregator implements the Aggregator interface using PostgreSQL.
type PostgresAggregator struct {
	db *sql.DB
}

// NewPostgresAggregator creates a new PostgresAggregator.
func NewPostgresAggregator(db *sql.DB) *PostgresAggregator {
	return &PostgresAggregator{db: db}
}

// GetClickMetrics retrieves aggregated click data within a specified timeframe.
func (a *PostgresAggregator) GetClickMetrics(startTime, endTime time.Time) (ClickMetrics, error) {
	clickCountQuery := `  
        SELECT COUNT(*)  
        FROM clicks  
        WHERE timestamp >= $1 AND timestamp <= $2  
    `

	var totalClicks int64

	// Use QueryRow to get a single result.
	err := a.db.QueryRow(clickCountQuery, startTime, endTime).Scan(&totalClicks)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return total clicks as 0 and default click rate as 0.0
			return ClickMetrics{TotalClicks: 0, ClickRate: 0.0}, nil
		}
		return ClickMetrics{}, err // Return an error if there was an issue with the query.
	}

	// Calculate ClickRate (assumes totalImpressions is a placeholder)
	totalImpressions := int64(100) // Placeholder for total impressions; implement your own logic as needed.
	clickRate := 0.0
	if totalImpressions > 0 {
		clickRate = float64(totalClicks) / float64(totalImpressions) * 100
	}

	// Return the total number of clicks and the calculated click rate.
	return ClickMetrics{TotalClicks: totalClicks, ClickRate: clickRate}, nil
}

// AnalyticsService holds the aggregator for analytics.
type AnalyticsService struct {
	aggregator Aggregator
}

// NewAnalyticsService creates and returns a new AnalyticsService.
func NewAnalyticsService(aggregator Aggregator) *AnalyticsService {
	return &AnalyticsService{aggregator: aggregator}
}

// Handler handles requests for analytics.
type Handler struct {
	service *AnalyticsService
}

// NewHandler creates and returns a new Handler for analytics.
func NewHandler(service *AnalyticsService) *Handler {
	return &Handler{service: service}
}

// GetAnalyticsHandler handles the request for getting analytics data.
func (h *Handler) GetAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract start and end parameters from the query string.
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	// Parse the start and end times (this may need error handling).
	startTime, err := time.Parse(time.RFC3339, start)
	if err != nil {
		http.Error(w, "Invalid start time format", http.StatusBadRequest)
		return
	}
	endTime, err := time.Parse(time.RFC3339, end)
	if err != nil {
		http.Error(w, "Invalid end time format", http.StatusBadRequest)
		return
	}

	// Get click metrics from the service.
	clickMetrics, err := h.service.aggregator.GetClickMetrics(startTime, endTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the click metrics as JSON.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clickMetrics)
}
