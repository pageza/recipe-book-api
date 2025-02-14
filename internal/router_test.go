// router_test.go
package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

// setupRouter initializes the router with routes and middleware for testing.
func setupRouter() *gin.Engine {
	router := gin.Default()
	SetupRoutes(router)
	return router
}

// generateTestJWT generates a JWT token for testing purposes.
func generateTestJWT(userID uint) string {
	// Use the same secret as in router.go
	secret := getEnv("JWT_SECRET", "your_secret_key")
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 72).Unix(), // Token expires after 72 hours.
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Printf("Error generating test JWT: %v\n", err)
		return ""
	}
	return tokenString
}

// TestGetRecipes verifies the GET /recipes endpoint.
func TestGetRecipes(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/recipes", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var recipes []Recipe
	err := json.Unmarshal(w.Body.Bytes(), &recipes)
	assert.NoError(t, err)
	assert.NotEmpty(t, recipes, "Recipes should not be empty")
}

// TestCreateRecipe verifies the POST /recipes endpoint.
func TestCreateRecipe(t *testing.T) {
	router := setupRouter()

	// Generate JWT token
	userID := uint(1)
	token := generateTestJWT(userID)
	if token == "" {
		t.Fatalf("Failed to generate JWT token")
	}

	// Create recipe payload
	payload := Recipe{
		Title:        "Test Recipe",
		Ingredients:  "Test Ingredients",
		Instructions: "Test Instructions",
		Calories:     100,
		UserID:       userID,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/recipes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var recipe Recipe
	err := json.Unmarshal(w.Body.Bytes(), &recipe)
	assert.NoError(t, err)
	assert.Equal(t, payload.Title, recipe.Title)
	assert.Equal(t, payload.Ingredients, recipe.Ingredients)
	assert.Equal(t, payload.Instructions, recipe.Instructions)
	assert.Equal(t, payload.Calories, recipe.Calories)
	assert.Equal(t, payload.UserID, recipe.UserID)
}

// TestGetIngredient verifies the GET /ingredients/:id endpoint.
func TestGetIngredient(t *testing.T) {
	router := setupRouter()

	// First, create an ingredient to retrieve
	ingredient := Ingredient{
		Name:     "Test Ingredient",
		Quantity: "2 cups",
		RecipeID: 1, // Assuming a recipe with ID 1 exists
	}
	if err := DB.Create(&ingredient).Error; err != nil {
		t.Fatalf("Failed to create ingredient: %v", err)
	}

	req, _ := http.NewRequest("GET", "/ingredients/"+strconv.Itoa(int(ingredient.ID)), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var fetchedIngredient Ingredient
	err := json.Unmarshal(w.Body.Bytes(), &fetchedIngredient)
	assert.NoError(t, err)
	assert.Equal(t, ingredient.Name, fetchedIngredient.Name)
	assert.Equal(t, ingredient.Quantity, fetchedIngredient.Quantity)
	assert.Equal(t, ingredient.RecipeID, fetchedIngredient.RecipeID)
}

// TestUpdateIngredient verifies the PUT /ingredients/:id endpoint.
func TestUpdateIngredient(t *testing.T) {
	router := setupRouter()

	// First, create an ingredient to update
	ingredient := Ingredient{
		Name:     "Old Ingredient",
		Quantity: "1 cup",
		RecipeID: 1, // Assuming a recipe with ID 1 exists
	}
	if err := DB.Create(&ingredient).Error; err != nil {
		t.Fatalf("Failed to create ingredient: %v", err)
	}

	// Generate JWT token
	userID := uint(1)
	token := generateTestJWT(userID)
	if token == "" {
		t.Fatalf("Failed to generate JWT token")
	}

	// Define update payload
	updatePayload := Ingredient{
		Name:     "Updated Ingredient",
		Quantity: "3 cups",
	}
	body, _ := json.Marshal(updatePayload)

	req, _ := http.NewRequest("PUT", "/ingredients/"+strconv.Itoa(int(ingredient.ID)), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updatedIngredient Ingredient
	err := json.Unmarshal(w.Body.Bytes(), &updatedIngredient)
	assert.NoError(t, err)
	assert.Equal(t, updatePayload.Name, updatedIngredient.Name)
	assert.Equal(t, updatePayload.Quantity, updatedIngredient.Quantity)
}

// TestDeleteIngredient verifies the DELETE /ingredients/:id endpoint.
func TestDeleteIngredient(t *testing.T) {
	router := setupRouter()

	// First, create an ingredient to delete
	ingredient := Ingredient{
		Name:     "Ingredient to Delete",
		Quantity: "5 grams",
		RecipeID: 1, // Assuming a recipe with ID 1 exists
	}
	if err := DB.Create(&ingredient).Error; err != nil {
		t.Fatalf("Failed to create ingredient: %v", err)
	}

	// Generate JWT token
	userID := uint(1)
	token := generateTestJWT(userID)
	if token == "" {
		t.Fatalf("Failed to generate JWT token")
	}

	req, _ := http.NewRequest("DELETE", "/ingredients/"+strconv.Itoa(int(ingredient.ID)), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify deletion
	var deletedIngredient Ingredient
	err := DB.First(&deletedIngredient, ingredient.ID).Error
	assert.Error(t, err, "Ingredient should be deleted")
}

// TestGetRecipe verifies the GET /recipes/:id endpoint.
func TestGetRecipe(t *testing.T) {
	router := setupRouter()

	// First, create a recipe to retrieve
	recipe := Recipe{
		Title:        "Test Recipe for Get",
		Ingredients:  "Ingredient A, Ingredient B",
		Instructions: "Step 1, Step 2",
		Calories:     250,
		UserID:       1, // Assuming a user with ID 1 exists
	}
	if err := DB.Create(&recipe).Error; err != nil {
		t.Fatalf("Failed to create recipe: %v", err)
	}

	req, _ := http.NewRequest("GET", "/recipes/"+strconv.Itoa(int(recipe.ID)), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var fetchedRecipe Recipe
	err := json.Unmarshal(w.Body.Bytes(), &fetchedRecipe)
	assert.NoError(t, err)
	assert.Equal(t, recipe.Title, fetchedRecipe.Title)
	assert.Equal(t, recipe.Ingredients, fetchedRecipe.Ingredients)
	assert.Equal(t, recipe.Instructions, fetchedRecipe.Instructions)
	assert.Equal(t, recipe.Calories, fetchedRecipe.Calories)
	assert.Equal(t, recipe.UserID, fetchedRecipe.UserID)
}
