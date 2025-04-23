package main

import (
	"log"
	"net/http"

	"github.com/Sid9pore/video-ad-tracker/internal/ads"
	"github.com/Sid9pore/video-ad-tracker/internal/analytics"
	"github.com/Sid9pore/video-ad-tracker/internal/metrics"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
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
	router.HandleFunc("/ads/click", adsHandler.LogClickHandler).Methods("POST")

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
