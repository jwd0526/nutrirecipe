package services

import (
	"strings"
	"testing"
)

func TestUSDASearchURL_Format(t *testing.T) {
	if !strings.HasPrefix(usdaSearchURL, "https://") {
		t.Error("usdaSearchURL must use HTTPS")
	}
	if !strings.Contains(usdaSearchURL, "usda.gov") {
		t.Error("usdaSearchURL must point to usda.gov")
	}
}

func TestNutrientNumbers(t *testing.T) {
	food := usdaFood{
		FdcID:       123,
		Description: "test food",
		FoodNutrients: []usdaNutrient{
			{NutrientNumber: "208", Value: 100},
			{NutrientNumber: "203", Value: 10},
			{NutrientNumber: "204", Value: 5},
			{NutrientNumber: "205", Value: 20},
			{NutrientNumber: "999", Value: 999},
		},
	}

	opt := extractOption(nil, food)

	if opt.Calories != 100 {
		t.Errorf("expected calories 100, got %f", opt.Calories)
	}
	if opt.Protein != 10 {
		t.Errorf("expected protein 10, got %f", opt.Protein)
	}
	if opt.Fat != 5 {
		t.Errorf("expected fat 5, got %f", opt.Fat)
	}
	if opt.Carbs != 20 {
		t.Errorf("expected carbs 20, got %f", opt.Carbs)
	}
}

func TestNutrientNumbers_NoNutrients(t *testing.T) {
	food := usdaFood{FdcID: 1, Description: "empty food"}
	opt := extractOption(nil, food)
	if opt.Calories != 0 || opt.Protein != 0 || opt.Fat != 0 || opt.Carbs != 0 {
		t.Error("expected all zeros for food with no nutrients")
	}
}
