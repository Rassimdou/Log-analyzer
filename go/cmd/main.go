package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func main() {
	// Initialize components
	storage, err := NewLogStorage("logs")
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	
	limiter := NewRateLimiter(10, time.Minute)
	
	r := gin.Default()
	
	// Apply rate limiting globally
	r.Use(limiter.Middleware())
	
	// Public routes (no authentication)
	public := r.Group("/")
	{
		public.GET("/status", StatusHandler)
		public.POST("/ingest", IngestHandler)
	}
	
	// Protected routes (require API key)

	
	log.Println("🚀 Security Analyzer starting on :8080")
	r.Run(":8080")
}