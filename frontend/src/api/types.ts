export interface QAPair {
  question: string
  answer: string
}

export interface AgentParseRequest {
  input: string
  qa_history?: QAPair[]
}

export interface ClarificationQ {
  id: string
  ingredient_raw: string
  question: string
  options: string[]
}

export interface SearchQueries {
  primary: string
  alternatives: string[]
}

export interface ParsedIngredient {
  name: string
  portion: number
  unit: string
  quantity_g: number
  search_queries: SearchQueries
  confidence: 'high' | 'low'
  notes?: string
}

export interface AgentParseResponse {
  status: 'resolved' | 'needs_clarification'
  questions?: ClarificationQ[]
  ingredients?: ParsedIngredient[]
}

export interface USDAOption {
  fdc_id: string
  name: string
  category: string
  calories: number
  protein: number
  carbs: number
  fat: number
}

export interface ValidatedIngredient extends ParsedIngredient {
  fdc_id?: string
  fdc_name?: string
  match_status: 'matched' | 'low_confidence' | 'unresolved'
  match_warning?: string
  options?: USDAOption[]
  calories_per_100g?: number
  protein_per_100g?: number
  carbs_per_100g?: number
  fat_per_100g?: number
}

export interface SaveRecipeRequest {
  name: string
  ingredients: ValidatedIngredient[]
}

export interface Recipe {
  id: string
  name: string
  total_weight_g: number
  serving_size_g: number
  created_at: string
}

export interface RecipeDetail extends Recipe {
  ingredients: IngredientDetail[]
}

export interface IngredientDetail {
  id: string
  name: string
  quantity_g: number
  weight_ratio: number
  portion?: number
  unit?: string
  fdc_id?: string
  source: string
  notes?: string
  calories_per_100g?: number
  protein_per_100g?: number
  carbs_per_100g?: number
  fat_per_100g?: number
}
