package main 

import (
	"github.com/gin-gonic/gin"
	"net/http"
)


type LogEntry struct{
	Source 	string `json:"source"`
	Message string `json:"message"`
}



func IngestHandler(storage *LogStorage) gin.HandlerFunc{
	return func(c *gin.Context) {
		var log LogEntry 
		if err := c.ShouldBindJSON(&log); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}
		//Get client IP
		clientIP := c.ClientIP()

		//write to storage 
		if err := storage.WriteLog(log.Source, clientIP , log.Message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to write log"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "received",
	
		})

}


func StatusHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "operational"})}