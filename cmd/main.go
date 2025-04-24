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

func main() {
	log.Println("Connecting to DB")
	connStr := "host=postgres-server port=5432 user=admin password=myNewP@ssw0rd dbname=videoads sslmode=disable"
	log.Println(connStr)
	db, err := ads.InitializeDB(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Connected to DB")

	// Setup RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("rabbit dial: %v", err)
	}
	mq, err := conn.Channel()
	if err != nil {
		log.Fatalf("rabbit channel: %v", err)
	}
	mq.QueueDeclare("clicks_queue", true, false, false, false, nil)
	// Initialize ads repository and service
	adsRepo := ads.NewPostgresRepository(db)
	adsService := ads.NewService(adsRepo)
	// Setup service & handler for insert request
	adsInsertService := ads.NewInsertService(adsRepo, mq)

	// Start consumer
	go startConsumer(mq, adsRepo)

	// Initialize analytics repository and service
	analyticsRepo := analytics.NewPostgresAggregator(db)
	analyticsService := analytics.NewAnalyticsService(analyticsRepo)

	// Initialize metrics
	metrics.InitMetrics()

	// Create HTTP router
	router := mux.NewRouter()

	// Ads Handlers
	adsHandler := ads.NewHandler(adsService)
	adsInsertHandler := ads.NewHandler(adsInsertService)
	router.HandleFunc("/ads", adsHandler.GetAdsHandler).Methods("GET")
	go router.HandleFunc("/ads/click", adsInsertHandler.LogClickHandler).Methods("POST")

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

func startConsumer(mq *amqp.Channel, repo ads.Repository) {
	msgs, err := mq.Consume("clicks_queue", "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("consume error: %v", err)
	}
	for msg := range msgs {
		var click ads.ClickData
		if err := json.Unmarshal(msg.Body, &click); err != nil {
			log.Printf("unmarshal error: %v", err)
			continue
		}
		// Persist via repo
		if err := repo.LogClick(click); err != nil {
			log.Printf("db insert error: %v", err)
			// Optionally: requeue or dead-letter
		}
	}
}
