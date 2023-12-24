package main

import (
	"net/http"

	"github.com/anthonyhungnguyen/image-maestro/middleware"
	"github.com/anthonyhungnguyen/image-maestro/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	private := r.Group("/api")
	private.Use(middleware.AuthRequired())
	routes.ImageRoutes(private)

	r.Run(":8082")
}
