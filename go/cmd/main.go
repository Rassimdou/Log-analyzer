package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
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
		public.DELETE("/log/delete", IngestDELETEHandler(storage))
		public.PUT("/log/put", IngestPUTHandler(storage))
		public.GET("/log/get", IngestGETHandler(storage))
		public.POST("/log/post", IngestPOSTHandler(storage))
	}

	// Protected routes (require API key)

	log.Println("🚀 Security Analyzer starting on :8080")
	r.Run(":8080")
}
