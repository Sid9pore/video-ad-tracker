package analytics  

import (  
    "database/sql"  
    "time"  
)  

// Aggregator defines methods to analyze click data  
type Aggregator interface {  
    GetClickMetrics(startTime, endTime time.Time) (ClickMetrics, error)  
}  

// PostgresAggregator implements the Aggregator interface using PostgreSQL  
type PostgresAggregator struct {  
    db *sql.DB  
}  

// NewPostgresAggregator creates a new PostgresAggregator  
func NewPostgresAggregator(db *sql.DB) *PostgresAggregator {  
    return &PostgresAggregator{db: db}  
}  

// GetClickMetrics retrieves aggregated click data within a specified timeframe  
func (a *PostgresAggregator) GetClickMetrics(startTime, endTime time.Time) (ClickMetrics, error) {  
    clickCountQuery := `  
        SELECT COUNT(*)  
        FROM clicks  
        WHERE