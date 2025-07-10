package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type LogEntry struct {
	Source    string `json:"source"`
	Message   string `json:"message"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Timestamp string `json:"timestamp"`
}

func IngestPOSTHandler(storage *LogStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var log LogEntry
		if err := c.ShouldBindJSON(&log); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}

		log.Method = "POST"
		log.Path = c.Request.URL.Path
		log.Timestamp = time.Now().Format(time.RFC3339)
		_ = storage.WriteLog(log.Method, c.ClientIP(), log.Source, log.Path, log.Message)
		c.JSON(http.StatusOK, gin.H{"status": "POST received"})
	}
}

func IngestGETHandler(storage *LogStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := LogEntry{
			Source:    c.Query("source"),
			Message:   c.Query("message"),
			Method:    "GET",
			Path:      c.Request.URL.Path,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		_ = storage.WriteLog(log.Message, c.ClientIP(), log.Source, log.Path, log.Message)
		c.JSON(http.StatusOK, gin.H{"status": "GET received"})

	}
}

func IngestPUTHandler(storage *LogStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var log LogEntry
		if err := c.ShouldBindJSON(&log); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "Invalid JSON"})

			return
		}
		log.Method = "PUT"
		log.Path = c.Request.URL.Path
		log.Timestamp = time.Now().Format(time.RFC3339)

		_ = storage.WriteLog(log.Method, c.ClientIP(), log.Source, log.Path, log.Message)
		c.JSON(http.StatusOK, gin.H{"status": "PUT received"})
	}
}

func IngestDELETEHandler(storage *LogStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := LogEntry{
			Source:    c.Query("source"),
			Message:   c.Query("message"),
			Method:    "DELETE",
			Path:      c.Request.URL.Path,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		_ = storage.WriteLog(log.Method, c.ClientIP(), log.Message, log.Path, log.Source)
		c.JSON(http.StatusAccepted, gin.H{"status": "DELETE received"})

	}

}
