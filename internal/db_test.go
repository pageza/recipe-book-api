package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupTestDB() (*gorm.DB, error) {
	// Use a separate test database
	dsn := "host=localhost user=testuser password=testpass dbname=recipe_book_test port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Run migrations
	err = db.AutoMigrate(&Recipe{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestCreateAndGetRecipe(t *testing.T) {
	db, err := SetupTestDB()
	assert.NoError(t, err, "Failed to connect to the test database")

	// Create a new recipe
	newRecipe := Recipe{
		Name:         "Test Recipe",
		Ingredients:  "Test Ingredients",
		Instructions: "Test Instructions",
		Calories:     250,
	}

	result := db.Create(&newRecipe)
	assert.NoError(t, result.Error, "Failed to create a new recipe")
	assert.NotZero(t, newRecipe.ID, "Expected recipe ID to be set")

	// Retrieve the recipe
	var retrievedRecipe Recipe
	result = db.First(&retrievedRecipe, newRecipe.ID)
	assert.NoError(t, result.Error, "Failed to retrieve the recipe")
	assert.Equal(t, newRecipe.Name, retrievedRecipe.Name, "Recipe name does not match")
	assert.Equal(t, newRecipe.Ingredients, retrievedRecipe.Ingredients, "Recipe ingredients do not match")
	assert.Equal(t, newRecipe.Instructions, retrievedRecipe.Instructions, "Recipe instructions do not match")
	assert.Equal(t, newRecipe.Calories, retrievedRecipe.Calories, "Recipe calories do not match")
}
