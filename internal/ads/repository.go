package ads

import (
	"database/sql"
)

// Repository interface defines methods for accessing ads data
type Repository interface {
	GetAllAds() ([]Ad, error)
	LogClick(click ClickData) error
}

// PostgresRepository implements the Repository interface using PostgreSQL
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgresRepository
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// GetAllAds retrieves all ads from the database
func (r *PostgresRepository) GetAllAds() ([]Ad, error) {
	rows, err := r.db.Query("SELECT id, image_url, target_url FROM ads")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ads []Ad
	for rows.Next() {
		var ad Ad
		if err := rows.Scan(&ad.ID, &ad.ImageURL, &ad.TargetURL); err != nil {
			return nil, err
		}
		ads = append(ads, ad)
	}
	return ads, nil
}

// LogClick records a click event in the database
func (r *PostgresRepository) LogClick(click ClickData) error {
	query := `  
        INSERT INTO clicks (ad_id, timestamp, ip, video_playback_time)  
        VALUES ($1, $2, $3, $4)  
    `
	_, err := r.db.Exec(query, click.AdID, click.Timestamp, click.IP, click.VideoPlaybackTime)
	return err
}

// InitializeDB initializes the database connection
func InitializeDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Verify the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
