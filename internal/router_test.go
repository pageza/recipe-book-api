package internal

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetupTestRouter() *gin.Engine {
	router := gin.Default()
	SetupRoutes(router)
	return router
}

func TestGetRecipes(t *testing.T) {
	router := SetupTestRouter()

	req, err := http.NewRequest("GET", "/recipes", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Perform the request
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Check status code
	assert.Equal(t, http.StatusOK, resp.Code)

	// Check response body
	var recipes []map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &recipes)
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, recipes, 2)
	assert.Equal(t, "Spaghetti Bolognese", recipes[0]["name"])
	assert.Equal(t, float64(400), recipes[0]["calories"])
}

func TestCreateRecipe(t *testing.T) {
	router := SetupTestRouter()

	newRecipe := map[string]interface{}{
		"name":         "Test Recipe",
		"ingredients":  "Test Ingredients",
		"instructions": "Test Instructions",
		"calories":     250,
	}

	jsonValue, _ := json.Marshal(newRecipe)
	req, err := http.NewRequest("POST", "/recipes", bytes.NewBuffer(jsonValue))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Check status code
	assert.Equal(t, http.StatusCreated, resp.Code)

	// Check response body
	var response map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "recipe created", response["status"])
	recipe, exists := response["recipe"].(map[string]interface{})
	assert.True(t, exists)
	assert.Equal(t, "Test Recipe", recipe["name"])
	assert.Equal(t, "Test Ingredients", recipe["ingredients"])
	assert.Equal(t, "Test Instructions", recipe["instructions"])
	assert.Equal(t, float64(250), recipe["calories"])
}
