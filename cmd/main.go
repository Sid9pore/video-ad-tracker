package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Sid9pore/video-ad-tracker/internal/ads"
	"github.com/Sid9pore/video-ad-tracker/internal/analytics"
	"github.com/Sid9pore/video-ad-tracker/internal/metrics"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
)

var (
	rabbitConn *amqp.Connection
	rabbitCh   *amqp.Channel
)

// Initialize RabbitMQ connection and channel
func initRabbitMQ() error {
	var err error
	rabbitConn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
		return err
	}

	rabbitCh, err = rabbitConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
		return err
	}

	// Declare a queue for click events
	_, err = rabbitCh.QueueDeclare(
		"clicks_queue", // Queue name
		true,           // Durable, survives server restarts
		false,          // Auto delete
		false,          // Exclusive
		false,          // No wait
		nil,            // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
		return err
	}
	return nil
}

// SendClickEvent sends a click event to RabbitMQ
func SendClickEvent(clickData ads.ClickData) error {
	// Convert click data to JSON
	clickJSON, err := json.Marshal(clickData)
	if err != nil {
		log.Printf("Error marshalling click data: %s", err)
		return err
	}

	// Publish the click event to the queue
	err = rabbitCh.Publish(
		"",             // Exchange
		"clicks_queue", // Routing key (queue name)
		false,          // Mandatory
		false,          // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        clickJSON,
		},
	)
	if err != nil {
		log.Printf("Failed to send message to RabbitMQ: %s", err)
		return err
	}
	log.Println("Sent click event to RabbitMQ")
	return nil
}

func main() {
	err := initRabbitMQ()
	if err != nil {
		log.Fatalf("Error initializing RabbitMQ: %s", err)
	}

	// Start the consumer in a separate goroutine
	go ConsumeClickEvents()

	log.Println("Connecting to DB")
	connStr := "host=postgres-server port=5432 user=admin password=myNewP@ssw0rd dbname=videoads sslmode=disable"
	log.Println(connStr)
	db, err := ads.InitializeDB(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Connected to DB")

	// Initialize ads repository and service
	adsRepo := ads.NewPostgresRepository(db)
	adsService := ads.NewService(adsRepo)

	// Initialize analytics repository and service
	analyticsRepo := analytics.NewPostgresAggregator(db)
	analyticsService := analytics.NewAnalyticsService(analyticsRepo)

	// Initialize metrics
	metrics.InitMetrics()

	// Create HTTP router
	router := mux.NewRouter()

	// Ads Handlers
	adsHandler := ads.NewHandler(adsService)
	router.HandleFunc("/ads", adsHandler.GetAdsHandler).Methods("GET")
	go router.HandleFunc("/ads/click", adsHandler.LogClickHandler).Methods("POST")

	// Analytics Handlers
	analyticsHandler := analytics.NewHandler(analyticsService) // Ensure this function is correctly defined
	router.HandleFunc("/ads/analytics", analyticsHandler.GetAnalyticsHandler).Methods("GET")

	// Metrics Endpoint
	router.Handle("/metrics", metrics.MetricsHandler()).Methods("GET")

	// Logging middleware
	router.Use(metrics.Middleware)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
