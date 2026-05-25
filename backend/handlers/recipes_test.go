package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jwd0526/nutrirecipe/models"
)

func recipeRouter() *gin.Engine {
	r := gin.New()
	h := NewRecipeHandler(nil)
	r.POST("/api/recipes", h.Save)
	r.GET("/api/recipes", h.List)
	r.GET("/api/recipes/:id", h.Get)
	return r
}

func TestSaveRecipe_BadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/recipes", nil)
	req.Header.Set("Content-Type", "application/json")
	recipeRouter().ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestSaveRecipe_MissingName(t *testing.T) {
	body, _ := json.Marshal(models.SaveRecipeRequest{
		Name:        "",
		Ingredients: []models.ValidatedIngredient{{ParsedIngredient: models.ParsedIngredient{Name: "chicken", QuantityG: 100}}},
	})
	w := httptest.NewRecorder()
	recipeRouter().ServeHTTP(w, newJSONRequest(http.MethodPost, "/api/recipes", body))

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestSaveRecipe_EmptyIngredients(t *testing.T) {
	body, _ := json.Marshal(models.SaveRecipeRequest{
		Name:        "My Recipe",
		Ingredients: []models.ValidatedIngredient{},
	})
	w := httptest.NewRecorder()
	recipeRouter().ServeHTTP(w, newJSONRequest(http.MethodPost, "/api/recipes", body))

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
