import type { SaveRecipeRequest, Recipe, RecipeDetail } from './types'

export async function saveRecipe(req: SaveRecipeRequest): Promise<{ id: string }> {
  const res = await fetch('/api/recipes', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(req),
  })
  if (!res.ok) throw new Error(`recipes save failed: ${res.status}`)
  return res.json()
}

export async function listRecipes(): Promise<Recipe[]> {
  const res = await fetch('/api/recipes')
  if (!res.ok) throw new Error(`recipes list failed: ${res.status}`)
  return res.json()
}

export async function getRecipe(id: string): Promise<RecipeDetail> {
  const res = await fetch(`/api/recipes/${id}`)
  if (!res.ok) throw new Error(`recipes get failed: ${res.status}`)
  return res.json()
}
