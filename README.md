# 📺 Video Ad Tracker

A GoLang backend service to manage and track video advertisements. The system logs ad clicks, provides real-time analytics, handles spikes in traffic, and remains resilient under partial failures.

## 🚀 Features

- RESTful API to serve ads and record clicks  
- Real-time analytics with Redis caching  
- Asynchronous click processing via message queue  
- Scalable architecture with Docker  
- Rate limiting and structured logging  
- Resilient to DB/service failures  

## 📁 Project Structure

├── cmd/ # App entry point
├── internal/
│ ├── ads/ # Ad logic and handlers
│ ├── click/ # Click tracking logic
│ ├── analytics/ # Metrics aggregation
│ ├── middleware/ # Rate limiting, logging
│ └── db/ # DB access, migrations
├── configs/ # YAML or ENV configs
├── migrations/ # SQL migration scripts
├── Dockerfile # Container setup
├── docker-compose.yml # Service orchestration
└── README.md


## 📦 Installation

### 🐳 Docker (Recommended)


docker-compose up --build
🔧 Local Setup
Install Go 1.20+

Set up PostgreSQL and Redis

Run migrations:


migrate -path ./migrations -database postgres://user:pass@localhost:5432/ad_tracker up
Start the server:


go run ./cmd/server

🔌 API Endpoints

GET /ads
Returns a list of ad metadata.

json
[
  {
    "id": 1,
    "image_url": "https://cdn.com/ad1.png",
    "target_url": "https://example.com/product"
  }
]
POST /ads/click
Tracks ad click.

json
{
  "ad_id": 1,
  "timestamp": "2025-04-21T15:00:00Z",
  "ip": "203.0.113.1",
  "video_playback_time": 12
}
GET /ads/analytics
Returns real-time performance metrics per ad.

json
Copy
Edit
{
  "ad_id": 1,
  "clicks": 1200,
  "ctr": 0.045
}
🧱 Database Schema
See /migrations for full scripts.

ads: Stores ad metadata.

clicks: Tracks each user click with playback time and IP.

🛡 Resilience & Scalability
Message queue buffers clicks (Kafka or RabbitMQ)

Redis caches hot analytics for quick access

Goroutine worker pool for concurrent DB writes

Rate limiter middleware to prevent abuse

Dockerized for horizontal scaling

📈 Monitoring & Logging
Structured JSON logs via logrus or zap

Prometheus metrics endpoint (coming soon)

🔄 CI/CD Pipeline (Optional)
To add:

GitHub Actions or GitLab CI for:

Linting

Tests

Docker build

🧪 Testing

go test ./...
🌐 Contributing
Fork the repo

Create a new feature branch

Submit a PR with detailed description

📄 License
MIT

👤 Author
Sid9pore — GitHub


Let me know if you want a separate CONTRIBUTING.md or API Swagger spec!