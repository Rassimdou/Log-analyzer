package main

import (
	"github.com/gin-gonic/gin"
)

func APIKeyAuthMiddleware() gin.HandlerFunc {

	api_key := header.Get("X-API-Key")

}
