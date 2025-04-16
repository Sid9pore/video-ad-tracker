package ads

import (
	"errors"
	"time"
)

// Service defines methods for ads-related business logic
type Service struct {
	repo Repository
}

// NewService creates a new ads service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetAllAds retrieves all ads from the repository
func (s *Service) GetAllAds() ([]Ad, error) {
	ads, err := s.repo.GetAllAds()
	if err != nil {
		return nil, err
	}
	return ads, nil
}

// LogClick processes and logs a click event for an ad
func (s *Service) LogClick(click ClickData) error {
	// Perform any necessary business logic here, e.g. validation
	if click.AdID <= 0 {
		return errors.New("invalid ad ID")
	}

	click.Timestamp = time.Now() // Set the current timestamp
	return s.repo.LogClick(click)
}
