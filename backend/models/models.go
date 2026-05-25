package models

type AgentParseRequest struct {
	Input     string   `json:"input"`
	QAHistory []QAPair `json:"qa_history,omitempty"`
}

type QAPair struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type AgentParseResponse struct {
	Status      string            `json:"status"`
	Questions   []ClarificationQ  `json:"questions,omitempty"`
	Ingredients []ParsedIngredient `json:"ingredients,omitempty"`
}

type ClarificationQ struct {
	ID            string   `json:"id"`
	IngredientRaw string   `json:"ingredient_raw"`
	Question      string   `json:"question"`
	Options       []string `json:"options"`
}

type ParsedIngredient struct {
	Name          string        `json:"name"`
	Portion       float64       `json:"portion"`
	Unit          string        `json:"unit"`
	QuantityG     float64       `json:"quantity_g"`
	SearchQueries SearchQueries `json:"search_queries"`
	Confidence    string        `json:"confidence"`
	Notes         string        `json:"notes,omitempty"`
}

type SearchQueries struct {
	Primary      string   `json:"primary"`
	Alternatives []string `json:"alternatives"`
}

type ValidatedIngredient struct {
	ParsedIngredient
	FdcID           string       `json:"fdc_id,omitempty"`
	FdcName         string       `json:"fdc_name,omitempty"`
	MatchStatus     string       `json:"match_status"`
	MatchWarning    string       `json:"match_warning,omitempty"`
	Options         []USDAOption `json:"options,omitempty"`
	CaloriesPer100g float64      `json:"calories_per_100g,omitempty"`
	ProteinPer100g  float64      `json:"protein_per_100g,omitempty"`
	CarbsPer100g    float64      `json:"carbs_per_100g,omitempty"`
	FatPer100g      float64      `json:"fat_per_100g,omitempty"`
}

type USDAOption struct {
	FdcID    string  `json:"fdc_id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Calories float64 `json:"calories"`
	Protein  float64 `json:"protein"`
	Carbs    float64 `json:"carbs"`
	Fat      float64 `json:"fat"`
}

type SaveRecipeRequest struct {
	Name        string               `json:"name"`
	Ingredients []ValidatedIngredient `json:"ingredients"`
}
