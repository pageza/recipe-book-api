// router.go
package internal

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(getEnv("JWT_SECRET", "your_secret_key")) // Replace with a secure key in production

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

		// GET endpoint for retrieving a specific recipe.
		recipes.GET("/:id", GetRecipe)

		// PUT endpoint for updating a specific recipe.
		recipes.PUT("/:id", UpdateRecipe)

		// DELETE endpoint for deleting a specific recipe.
		recipes.DELETE("/:id", DeleteRecipe)

		// POST endpoint for creating a new recipe.
		recipes.POST("", CreateRecipe)

		// You can add more recipe-related routes here (e.g., GET /recipes/:id, PUT /recipes/:id, DELETE /recipes/:id)
	}

	// Group routes related to ingredients
	ingredients := router.Group("/ingredients")
	{
		// GET endpoint for listing ingredients.
		ingredients.GET("", GetIngredients)

		// GET endpoint for retrieving a specific ingredient.
		ingredients.GET("/:id", GetIngredient)

		// PUT endpoint for updating a specific ingredient.
		ingredients.PUT("/:id", UpdateIngredient)

		// DELETE endpoint for deleting a specific ingredient.
		ingredients.DELETE("/:id", DeleteIngredient)

		// POST endpoint for creating a new ingredient.
		ingredients.POST("", CreateIngredient)
	}

	// Group routes related to authentication
	auth := router.Group("/auth")
	{
		auth.POST("/signup", Signup)
		auth.POST("/login", Login)
		auth.GET("/profile", Profile).Use(JWTMiddleware()) // Protected route
	}
}

// GetRecipes handles the GET /recipes endpoint.
func GetRecipes(c *gin.Context) {
	title := c.Query("title")
	ingredient := c.Query("ingredient")

	var recipes []Recipe

	query := DB

	if title != "" {
		query = query.Where("title ILIKE ?", fmt.Sprintf("%%%s%%", title))
	}

	if ingredient != "" {
		query = query.Where("ingredients ILIKE ?", fmt.Sprintf("%%%s%%", ingredient))
	}

	if err := query.Find(&recipes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve recipes"})
		return
	}

	c.JSON(http.StatusOK, recipes)
}

// GetRecipe handles the GET /recipes/:id endpoint.
func GetRecipe(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe
	if err := DB.First(&recipe, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}
	c.JSON(http.StatusOK, recipe)
}

// UpdateRecipe handles the PUT /recipes/:id endpoint.
func UpdateRecipe(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe
	if err := DB.First(&recipe, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	var input struct {
		Title        string `json:"title"`
		Ingredients  string `json:"ingredients"`
		Instructions string `json:"instructions"`
		Calories     int    `json:"calories"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated := Recipe{
		Title:        input.Title,
		Ingredients:  input.Ingredients,
		Instructions: input.Instructions,
		Calories:     input.Calories,
	}

	if err := DB.Model(&recipe).Updates(updated).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update recipe"})
		return
	}

	c.JSON(http.StatusOK, recipe)
}

// DeleteRecipe handles the DELETE /recipes/:id endpoint.
func DeleteRecipe(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe
	if err := DB.First(&recipe, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	if err := DB.Delete(&recipe).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete recipe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Recipe deleted"})
}

// CreateRecipe handles the POST /recipes endpoint.
func CreateRecipe(c *gin.Context) {
	// Define a struct to bind incoming JSON data.
	var input struct {
		Title        string `json:"title" binding:"required"`
		Ingredients  string `json:"ingredients" binding:"required"`
		Instructions string `json:"instructions" binding:"required"`
		Calories     int    `json:"calories" binding:"required,min=0"`
	}

	// Bind JSON input to the input struct.
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a new recipe instance.
	recipe := Recipe{
		Title:        input.Title,
		Ingredients:  input.Ingredients,
		Instructions: input.Instructions,
		Calories:     input.Calories,
		UserID:       uint(c.GetUint("userID")), // Assuming userID is set in context
	}

	// Insert the new recipe into the database.
	if err := DB.Create(&recipe).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create recipe"})
		return
	}

	// Return the created recipe to the client.
	c.JSON(http.StatusCreated, recipe)
}

// GetIngredients handles the GET /ingredients endpoint.
func GetIngredients(c *gin.Context) {
	var ingredients []Ingredient
	if err := DB.Find(&ingredients).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve ingredients"})
		return
	}
	c.JSON(http.StatusOK, ingredients)
}

// GetIngredient handles the GET /ingredients/:id endpoint.
func GetIngredient(c *gin.Context) {
	id := c.Param("id")
	var ingredient Ingredient
	if err := DB.First(&ingredient, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ingredient not found"})
		return
	}
	c.JSON(http.StatusOK, ingredient)
}

// UpdateIngredient handles the PUT /ingredients/:id endpoint.
func UpdateIngredient(c *gin.Context) {
	id := c.Param("id")
	var ingredient Ingredient
	if err := DB.First(&ingredient, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ingredient not found"})
		return
	}

	var input struct {
		Name     string `json:"name" binding:"required"`
		Quantity string `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated := Ingredient{
		Name:     input.Name,
		Quantity: input.Quantity,
	}

	if err := DB.Model(&ingredient).Updates(updated).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ingredient"})
		return
	}

	c.JSON(http.StatusOK, ingredient)
}

// DeleteIngredient handles the DELETE /ingredients/:id endpoint.
func DeleteIngredient(c *gin.Context) {
	id := c.Param("id")
	var ingredient Ingredient
	if err := DB.First(&ingredient, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ingredient not found"})
		return
	}

	if err := DB.Delete(&ingredient).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete ingredient"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Ingredient deleted"})
}

// Signup handles the POST /auth/signup endpoint.
func Signup(c *gin.Context) {
	// Define a struct to bind incoming JSON data.
	var input struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	// Bind JSON input to the input struct.
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password using bcrypt.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create a new user instance.
	user := User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	// Insert the new user into the database.
	if err := DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate a JWT token for the newly created user.
	token, err := generateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return the JWT token to the client.
	c.JSON(http.StatusCreated, gin.H{"token": token})
}

// Login handles the POST /auth/login endpoint.
func Login(c *gin.Context) {
	// Define a struct to bind incoming JSON data.
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	// Bind JSON input to the input struct.
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve user from the database.
	var user User
	if err := DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Compare the provided password with the hashed password.
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate a JWT token for the authenticated user.
	token, err := generateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return the JWT token to the client.
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Profile handles the GET /auth/profile endpoint.
func Profile(c *gin.Context) {
	// Retrieve the user ID from the JWT token.
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Fetch the user from the database.
	var user User
	if err := DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	// Return the user profile information.
	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}

// generateJWT generates a JWT token for a given user ID.
func generateJWT(userID uint) (string, error) {
	// Define token claims.
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 72).Unix(), // Token expires after 72 hours.
	}

	// Create the token.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key.
	return token.SignedString(jwtSecret)
}

// JWTMiddleware is a middleware function for validating JWT tokens.
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header.
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// Expecting header value in the format "Bearer <token>"
		var tokenString string
		fmt.Sscanf(authHeader, "Bearer %s", &tokenString)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token missing"})
			c.Abort()
			return
		}

		// Parse the token.
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the token method conforms to "SigningMethodHMAC".
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Extract user ID from claims.
		userIDFloat, ok := claims["userID"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			c.Abort()
			return
		}
		userID := uint(userIDFloat)

		// Set user ID in context.
		c.Set("userID", userID)

		c.Next()
	}
}
