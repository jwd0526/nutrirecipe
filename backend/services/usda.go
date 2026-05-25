package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	pgkitx "github.com/jwd0526/pgkitx"

	"github.com/jwd0526/nutrirecipe/models"
)

const (
	usdaSearchURL = "https://api.nal.usda.gov/fdc/v1/foods/search"
	usdaFoodURL   = "https://api.nal.usda.gov/fdc/v1/food"
	pageSize      = "5"
)

type USDAService struct {
	db     *pgkitx.Pool
	apiKey string
}

func NewUSDAService(db *pgkitx.Pool, apiKey string) *USDAService {
	return &USDAService{db: db, apiKey: apiKey}
}

type usdaSearchResponse struct {
	Foods []usdaFood `json:"foods"`
}

type usdaFood struct {
	FdcID       int           `json:"fdcId"`
	Description string        `json:"description"`
	FoodCategory string       `json:"foodCategory"`
	FoodNutrients []usdaNutrient `json:"foodNutrients"`
}

type usdaNutrient struct {
	NutrientNumber string  `json:"nutrientNumber"`
	Value          float64 `json:"value"`
}

func (s *USDAService) Validate(ctx context.Context, ingredient models.ParsedIngredient) models.ValidatedIngredient {
	result := models.ValidatedIngredient{ParsedIngredient: ingredient}

	// 1. Check local_ingredients
	if found := s.lookupLocal(ctx, ingredient.SearchQueries.Primary, &result); found {
		return result
	}

	// 2. Check usda_cache
	if found := s.lookupCache(ctx, ingredient.SearchQueries.Primary, &result); found {
		return result
	}

	// 3. Try USDA API — primary query then alternatives (max 3 total)
	queries := append([]string{ingredient.SearchQueries.Primary}, ingredient.SearchQueries.Alternatives...)
	if len(queries) > 3 {
		queries = queries[:3]
	}

	for _, q := range queries {
		options, err := s.searchUSDA(ctx, q)
		if err != nil || len(options) == 0 {
			continue
		}
		best := options[0]
		result.FdcID = best.FdcID
		result.FdcName = best.Name
		result.CaloriesPer100g = best.Calories
		result.ProteinPer100g = best.Protein
		result.CarbsPer100g = best.Carbs
		result.FatPer100g = best.Fat
		result.Options = options
		result.MatchStatus = s.matchStatus(ingredient.SearchQueries.Primary, best.Name)
		s.cacheResult(ctx, best)
		return result
	}

	result.MatchStatus = "unresolved"
	return result
}

func (s *USDAService) lookupLocal(ctx context.Context, query string, out *models.ValidatedIngredient) bool {
	row := s.db.QueryRow(ctx,
		`SELECT name, calories_per_100g, protein_per_100g, carbs_per_100g, fat_per_100g
		 FROM local_ingredients WHERE LOWER(name) = LOWER($1) LIMIT 1`, query)

	var name string
	var cal, prot, carb, fat float64
	if err := row.Scan(&name, &cal, &prot, &carb, &fat); err != nil {
		return false
	}
	out.FdcName = name
	out.CaloriesPer100g = cal
	out.ProteinPer100g = prot
	out.CarbsPer100g = carb
	out.FatPer100g = fat
	out.MatchStatus = "matched"
	return true
}

func (s *USDAService) lookupCache(ctx context.Context, query string, out *models.ValidatedIngredient) bool {
	row := s.db.QueryRow(ctx,
		`SELECT fdc_id, name, calories, protein, carbs, fat
		 FROM usda_cache WHERE LOWER(name) = LOWER($1) LIMIT 1`, query)

	var fdcID, name string
	var cal, prot, carb, fat float64
	if err := row.Scan(&fdcID, &name, &cal, &prot, &carb, &fat); err != nil {
		return false
	}
	out.FdcID = fdcID
	out.FdcName = name
	out.CaloriesPer100g = cal
	out.ProteinPer100g = prot
	out.CarbsPer100g = carb
	out.FatPer100g = fat
	out.MatchStatus = s.matchStatus(query, name)
	return true
}

func (s *USDAService) searchUSDA(ctx context.Context, query string) ([]models.USDAOption, error) {
	params := url.Values{}
	params.Set("query", query)
	params.Set("api_key", s.apiKey)
	params.Set("pageSize", pageSize)

	resp, err := http.Get(fmt.Sprintf("%s?%s", usdaSearchURL, params.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var sr usdaSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return nil, err
	}

	options := make([]models.USDAOption, 0, len(sr.Foods))
	for _, f := range sr.Foods {
		options = append(options, extractOption(s, f))
	}
	return options, nil
}

func extractOption(_ *USDAService, f usdaFood) models.USDAOption {
	opt := models.USDAOption{
		FdcID:    fmt.Sprintf("%d", f.FdcID),
		Name:     f.Description,
		Category: f.FoodCategory,
	}
	for _, n := range f.FoodNutrients {
		switch n.NutrientNumber {
		case "208":
			opt.Calories = n.Value
		case "203":
			opt.Protein = n.Value
		case "204":
			opt.Fat = n.Value
		case "205":
			opt.Carbs = n.Value
		}
	}
	return opt
}

func (s *USDAService) cacheResult(ctx context.Context, opt models.USDAOption) {
	s.db.Exec(ctx,
		`INSERT INTO usda_cache (fdc_id, name, calories, protein, carbs, fat)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 ON CONFLICT (fdc_id) DO NOTHING`,
		opt.FdcID, opt.Name, opt.Calories, opt.Protein, opt.Carbs, opt.Fat)
}

func (s *USDAService) matchStatus(query, result string) string {
	queryWords := strings.Fields(strings.ToLower(query))
	resultWords := strings.Fields(strings.ToLower(result))

	shared := 0
	for _, qw := range queryWords {
		for _, rw := range resultWords {
			if qw == rw {
				shared++
				break
			}
		}
	}
	if shared < 2 {
		return "low_confidence"
	}
	return "matched"
}
