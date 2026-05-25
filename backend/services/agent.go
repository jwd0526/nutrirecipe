package services

import (
	"strconv"
	"strings"

	"github.com/jwd0526/nutrirecipe/models"
)

type AgentService struct{}

func NewAgentService() *AgentService {
	return &AgentService{}
}

// Parse implements Agent 1: clarify ambiguous syrup, else resolve with placeholders.
func (a *AgentService) Parse(req models.AgentParseRequest) models.AgentParseResponse {
	if containsSyrupWithoutQualifier(req.Input) {
		return models.AgentParseResponse{
			Status: "needs_clarification",
			Questions: []models.ClarificationQ{
				{
					ID:            "syrup_type",
					IngredientRaw: "syrup",
					Question:      "What type of syrup?",
					Options:       []string{"maple syrup", "corn syrup", "agave syrup", "simple syrup"},
				},
			},
		}
	}
	return models.AgentParseResponse{
		Status:      "resolved",
		Ingredients: parsePlaceholders(req.Input),
	}
}

func containsSyrupWithoutQualifier(input string) bool {
	lower := strings.ToLower(input)
	if !strings.Contains(lower, "syrup") {
		return false
	}
	qualifiers := []string{"maple", "corn", "agave", "simple", "golden", "rice"}
	for _, q := range qualifiers {
		if strings.Contains(lower, q+" syrup") {
			return false
		}
	}
	return true
}

var knownUnits = map[string]bool{
	"cup": true, "cups": true, "tbsp": true, "tsp": true,
	"tablespoon": true, "tablespoons": true, "teaspoon": true, "teaspoons": true,
	"g": true, "gram": true, "grams": true, "oz": true, "ounce": true, "ounces": true,
	"lb": true, "lbs": true, "pound": true, "pounds": true,
	"ml": true, "l": true, "liter": true, "liters": true,
	"piece": true, "pieces": true, "slice": true, "slices": true,
	"clove": true, "cloves": true, "whole": true,
}

func parsePlaceholders(input string) []models.ParsedIngredient {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	out := make([]models.ParsedIngredient, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		out = append(out, parseIngredientLine(line))
	}
	return out
}

func parseIngredientLine(line string) models.ParsedIngredient {
	words := strings.Fields(line)
	portion := 1.0
	unit := "serving"
	nameStart := 0

	if len(words) > 0 {
		if n, err := strconv.ParseFloat(words[0], 64); err == nil {
			portion = n
			nameStart = 1
			if len(words) > 1 && knownUnits[strings.ToLower(words[1])] {
				unit = strings.ToLower(words[1])
				nameStart = 2
			}
		}
	}

	name := strings.Join(words[nameStart:], " ")
	if name == "" {
		name = line
	}
	return models.ParsedIngredient{
		Name:      name,
		Portion:   portion,
		Unit:      unit,
		QuantityG: 100,
		SearchQueries: models.SearchQueries{
			Primary:      name,
			Alternatives: []string{},
		},
		Confidence: "high",
	}
}

// evalMatch implements Agent 2: < 2 shared words → low_confidence.
func evalMatch(query, result string) string {
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
