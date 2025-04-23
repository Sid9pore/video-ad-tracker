Video-ad-tracker
A lightweight Go service for tracking video ad impressions and clicks, backed by PostgreSQL and instrumented with Prometheus metrics.

📦 Features
Ad Retrieval: Fetch available ads via REST.

Click Tracking: Record ad click events (/ads/click) with resilience and queuing on DB downtime.

Analytics: Aggregate click metrics over time windows (/ads/analytics).

Prometheus Metrics: Expose /metrics endpoint (http_requests_total, Go runtime stats, custom gauges).

Dockerized: Run application, PostgreSQL, and Prometheus in separate containers without Docker Compose.

🚀 Quickstart
Prerequisites
Docker & Docker Desktop

Go (for local development)

(Optional) PgAdmin / DBeaver for DB inspection

1. Clone & Build

git clone https://github.com/Sid9pore/video-ad-tracker.git
cd video-ad-tracker
docker build -t video-ad-tracker .

2. Run Services

# Create network
docker network create video-ad-net

# PostgreSQL (exposed on host port 8084)
docker run -d --name postgres-server --network video-ad-net \
  -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=myNewP@ssw0rd \
  -e POSTGRES_DB=videoads \
  -v "$(pwd)/migrations:/docker-entrypoint-initdb.d" \
  -p 8084:5432 \
  postgres

# Application (connects to postgres-server:5432 internally)
docker run -d --name video-ad-tracker --network video-ad-net \
  -e DB_CONN="postgres://admin:myNewP@ssw0rd@postgres-server:5432/videoads?sslmode=disable" \
  -p 8080:8080 \
  video-ad-tracker

# Prometheus (scrapes video-ad-tracker)
docker run -d --name prometheus --network video-ad-net \
  -v "$(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml" \
  -p 9090:9090 \
  prom/prometheus
📑 Configuration
Environment Variables

DB_CONN: PostgreSQL connection string

PORT (optional): HTTP server port (default 8080)

prometheus.yml

global:
  scrape_interval: 10s

scrape_configs:
  - job_name: 'video-ad-tracker'
    static_configs:
      - targets: ['video-ad-tracker:8080']
🛠️ Development
Build & Run Locally

go run cmd/main.go

Migrations

Place SQL files in migrations/.

They auto-run on first container start.

Testing Endpoints

Click endpoint:

curl -X POST http://localhost:8080/ads/click \
  -H 'Content-Type: application/json' \
  -d '{"ad_id":1,"ip":"127.0.0.1","video_playback_time":30}'

Analytics:

curl -G http://localhost:8080/ads/analytics \
  --data-urlencode "start=2025-04-20T00:00:00Z" \
  --data-urlencode "end=2025-04-22T23:59:59Z"

Metrics:

curl http://localhost:8080/metrics

📊 Monitoring

Prometheus UI: http://localhost:9090

Grafana (if installed): add Prometheus at http://prometheus:9090

🔒 Security & Resilience
Click Queue: In-memory buffering with worker retries ensures no data loss during DB downtime.

Graceful Shutdown: Waits for queued events to flush on exit.

Prepared Statements: Prevents SQL injection and improves performance.