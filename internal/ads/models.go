package ads

import "time"

// Ad represents the advertisement data structure
type Ad struct {
	ID        int64  `json:"id"`
	ImageURL  string `json:"image_url"`
	TargetURL string `json:"target_url"`
}

// ClickData represents a click event
type ClickData struct {
	AdID              int64     `json:"ad_id"`
	Timestamp         time.Time `json:"timestamp"`
	IP                string    `json:"ip"`
	VideoPlaybackTime int       `json:"video_playback_time"`
}
