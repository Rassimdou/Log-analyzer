package main

import (
	"github.com/gin-gonic/gin"
	"github.com/Rassimdou/Log-analyzer/handlers"
	"github.com/Rassimdou/Log-analyzer/rate_limiter"
	"github.com/Rassimdou/Log-analyzer/storage"
	"github.com/Rassimdou/Log-analyzer/auth"
	"github.com/Rassimdou/Log-analyzer/middleware"

	"log"
	"time"
)

func main() {
	// Initialize storage
	storage, err := NewLogStorage("logs")
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	
	// Initialize rate limiter
	limiter := NewRateLimiter(10, time.Minute) // 10 reqs/min
	
	// Create router
	r := gin.Default()
	
	// Middleware stack
	r.Use(limiter.Middleware())
	r.Use(APIKeyAuthMiddleware()) // From auth.go
	
	// Routes
	r.GET("/status", StatusHandler)
	r.POST("/ingest", IngestHandler(storage)) // Handler with storage dependency
	
	// Start server
	log.Println("Starting security analyzer on :8080")
	r.Run(":8080")
}