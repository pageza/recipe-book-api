package main

import (
	"github.com/gin-gonic/gin"
	internal "github.com/pageza/recipe-book-api/internal"
)

func main() {
	// Initialize the database connection
	internal.InitDB()

	// Create a Gin router with default middleware (logger and recovery).
	router := gin.Default()

	// Setup API routes
	internal.SetupRoutes(router)

	// Start the server on port 8080.
	router.Run(":8080")
}
