package analytics

// ClickMetrics holds aggregated click data
type ClickMetrics struct {
	TotalClicks int64   `json:"total_clicks"`
	ClickRate   float64 `json:"click_rate"` // Click-Through Rate (CTR)
}
