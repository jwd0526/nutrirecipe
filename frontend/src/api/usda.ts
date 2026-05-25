import type { ParsedIngredient, ValidatedIngredient } from './types'

export async function validateIngredients(ingredients: ParsedIngredient[]): Promise<ValidatedIngredient[]> {
  const res = await fetch('/api/usda/validate', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ ingredients }),
  })
  if (!res.ok) throw new Error(`usda/validate failed: ${res.status}`)
  return res.json()
}
