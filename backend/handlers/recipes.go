package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	pgkitx "github.com/jwd0526/pgkitx"
	"github.com/jwd0526/nutrirecipe/models"
)

type RecipeHandler struct {
	db *pgkitx.Pool
}

func NewRecipeHandler(db *pgkitx.Pool) *RecipeHandler {
	return &RecipeHandler{db: db}
}

type recipeRow struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	TotalWeightG float64   `json:"total_weight_g"`
	ServingSizeG float64   `json:"serving_size_g"`
	CreatedAt    time.Time `json:"created_at"`
}

type ingredientRow struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	QuantityG        float64 `json:"quantity_g"`
	WeightRatio      float64 `json:"weight_ratio"`
	Portion          float64 `json:"portion,omitempty"`
	Unit             string  `json:"unit,omitempty"`
	FdcID            string  `json:"fdc_id,omitempty"`
	Source           string  `json:"source"`
	Notes            string  `json:"notes,omitempty"`
	CaloriesPer100g  float64 `json:"calories_per_100g,omitempty"`
	ProteinPer100g   float64 `json:"protein_per_100g,omitempty"`
	CarbsPer100g     float64 `json:"carbs_per_100g,omitempty"`
	FatPer100g       float64 `json:"fat_per_100g,omitempty"`
}

type recipeDetail struct {
	recipeRow
	Ingredients []ingredientRow `json:"ingredients"`
}

func (h *RecipeHandler) Save(c *gin.Context) {
	var req models.SaveRecipeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Name == "" || len(req.Ingredients) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and ingredients are required"})
		return
	}

	var totalWeight float64
	for _, ing := range req.Ingredients {
		totalWeight += ing.QuantityG
	}

	ctx := c.Request.Context()
	var recipeID string
	err := h.db.QueryRow(ctx,
		`INSERT INTO recipes (name, total_weight_g) VALUES ($1, $2) RETURNING id`,
		req.Name, totalWeight,
	).Scan(&recipeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save recipe"})
		return
	}

	for _, ing := range req.Ingredients {
		ratio := 0.0
		if totalWeight > 0 {
			ratio = ing.QuantityG / totalWeight
		}
		_, err := h.db.Exec(ctx,
			`INSERT INTO recipe_ingredients
			 (recipe_id, name, quantity_g, weight_ratio, portion, unit, fdc_id, source, notes,
			  calories_per_100g, protein_per_100g, carbs_per_100g, fat_per_100g)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
			recipeID, ing.Name, ing.QuantityG, ratio, ing.Portion, ing.Unit,
			ing.FdcID, "usda", ing.Notes,
			ing.CaloriesPer100g, ing.ProteinPer100g, ing.CarbsPer100g, ing.FatPer100g,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save ingredient"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"id": recipeID})
}

func (h *RecipeHandler) List(c *gin.Context) {
	rows, err := h.db.Query(c.Request.Context(),
		`SELECT id, name, total_weight_g, serving_size_g, created_at FROM recipes ORDER BY created_at DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list recipes"})
		return
	}
	defer rows.Close()

	recipes := []recipeRow{}
	for rows.Next() {
		var r recipeRow
		if err := rows.Scan(&r.ID, &r.Name, &r.TotalWeightG, &r.ServingSizeG, &r.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan recipe"})
			return
		}
		recipes = append(recipes, r)
	}
	c.JSON(http.StatusOK, recipes)
}

func (h *RecipeHandler) Get(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()

	var r recipeRow
	err := h.db.QueryRow(ctx,
		`SELECT id, name, total_weight_g, serving_size_g, created_at FROM recipes WHERE id=$1`, id,
	).Scan(&r.ID, &r.Name, &r.TotalWeightG, &r.ServingSizeG, &r.CreatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
		return
	}

	rows, err := h.db.Query(ctx,
		`SELECT id, name, quantity_g, weight_ratio, COALESCE(portion,0), COALESCE(unit,''),
		        COALESCE(fdc_id,''), source, COALESCE(notes,''),
		        COALESCE(calories_per_100g,0), COALESCE(protein_per_100g,0),
		        COALESCE(carbs_per_100g,0), COALESCE(fat_per_100g,0)
		 FROM recipe_ingredients WHERE recipe_id=$1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load ingredients"})
		return
	}
	defer rows.Close()

	ingredients := []ingredientRow{}
	for rows.Next() {
		var ing ingredientRow
		if err := rows.Scan(&ing.ID, &ing.Name, &ing.QuantityG, &ing.WeightRatio,
			&ing.Portion, &ing.Unit, &ing.FdcID, &ing.Source, &ing.Notes,
			&ing.CaloriesPer100g, &ing.ProteinPer100g, &ing.CarbsPer100g, &ing.FatPer100g,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan ingredient"})
			return
		}
		ingredients = append(ingredients, ing)
	}

	c.JSON(http.StatusOK, recipeDetail{recipeRow: r, Ingredients: ingredients})
}
