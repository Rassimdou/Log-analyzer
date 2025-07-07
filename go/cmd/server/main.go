package main

import (
	"github.com/gin-gonic/gin"
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