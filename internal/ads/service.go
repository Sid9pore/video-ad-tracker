package ads

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/streadway/amqp"
)

// Service defines methods for ads-related business logic
type Service struct {
	repo Repository
	mq   *amqp.Channel
}

// NewService creates a new ads service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func NewInsertService(repo Repository, mq *amqp.Channel) *Service {
	return &Service{repo: repo, mq: mq}
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
	if click.AdID <= 0 {
		return errors.New("invalid ad ID")
	}
	click.Timestamp = time.Now()

	// Marshal to JSON
	body, err := json.Marshal(click)
	if err != nil {
		return err
	}
	// Publish to RabbitMQ (durable queue)
	return s.mq.Publish(
		"",             // default exchange
		"clicks_queue", // routing key = queue name
		false, false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	)
}
