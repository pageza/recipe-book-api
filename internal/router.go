// router.go
package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRoutes initializes all the API routes.
func SetupRoutes(router *gin.Engine) {
	// Root endpoint: returns a welcome message.
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the Recipe Book API!"})
	})

	// Group routes related to recipes
	recipes := router.Group("/recipes")
	{
		// GET endpoint for listing recipes.
		recipes.GET("", GetRecipes)

		// POST endpoint for creating a new recipe.
		recipes.POST("", CreateRecipe)

		// You can add more recipe-related routes here (e.g., GET /recipes/:id, PUT /recipes/:id, DELETE /recipes/:id)
	}

	// Future groups for authentication, user management, etc., can be added here.
}

// GetRecipes handles the GET /recipes endpoint.
func GetRecipes(c *gin.Context) {
	// This is a placeholder; later you'll query your PostgreSQL database.
	recipes := []gin.H{
		{"id": 1, "name": "Spaghetti Bolognese", "calories": 400},
		{"id": 2, "name": "Chicken Salad", "calories": 300},
	}
	c.JSON(http.StatusOK, recipes)
}

// CreateRecipe handles the POST /recipes endpoint.
func CreateRecipe(c *gin.Context) {
	// Define a struct to bind incoming JSON data.
	var newRecipe struct {
		Name         string `json:"name"`
		Ingredients  string `json:"ingredients"`
		Instructions string `json:"instructions"`
		Calories     int    `json:"calories"`
	}

	// Bind JSON input to the newRecipe struct.
	if err := c.ShouldBindJSON(&newRecipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Here, you would add logic to insert the recipe into your PostgreSQL database.

	c.JSON(http.StatusCreated, gin.H{"status": "recipe created", "recipe": newRecipe})
}
