import { describe, it, expect, vi, beforeEach } from 'vitest'
import { saveRecipe, listRecipes, getRecipe } from './recipes'
import type { Recipe, RecipeDetail } from './types'

const mockFetch = vi.fn()
vi.stubGlobal('fetch', mockFetch)

beforeEach(() => mockFetch.mockReset())

const recipe: Recipe = {
  id: 'abc-123',
  name: 'Test Recipe',
  total_weight_g: 200,
  serving_size_g: 100,
  created_at: '2026-01-01T00:00:00Z',
}

describe('saveRecipe', () => {
  it('returns id on success', async () => {
    mockFetch.mockResolvedValue({ ok: true, json: async () => ({ id: 'abc-123' }) })
    const result = await saveRecipe({ name: 'Test', ingredients: [] })
    expect(result.id).toBe('abc-123')
    expect(mockFetch).toHaveBeenCalledWith('/api/recipes', expect.objectContaining({ method: 'POST' }))
  })

  it('throws on non-ok response', async () => {
    mockFetch.mockResolvedValue({ ok: false, status: 400 })
    await expect(saveRecipe({ name: '', ingredients: [] })).rejects.toThrow('400')
  })
})

describe('listRecipes', () => {
  it('returns list of recipes', async () => {
    mockFetch.mockResolvedValue({ ok: true, json: async () => [recipe] })
    const result = await listRecipes()
    expect(result).toHaveLength(1)
    expect(result[0].name).toBe('Test Recipe')
  })

  it('returns empty array when no recipes', async () => {
    mockFetch.mockResolvedValue({ ok: true, json: async () => [] })
    const result = await listRecipes()
    expect(result).toHaveLength(0)
  })

  it('throws on non-ok response', async () => {
    mockFetch.mockResolvedValue({ ok: false, status: 500 })
    await expect(listRecipes()).rejects.toThrow('500')
  })
})

describe('getRecipe', () => {
  it('returns recipe detail', async () => {
    const detail: RecipeDetail = { ...recipe, ingredients: [] }
    mockFetch.mockResolvedValue({ ok: true, json: async () => detail })
    const result = await getRecipe('abc-123')
    expect(result.id).toBe('abc-123')
    expect(mockFetch).toHaveBeenCalledWith('/api/recipes/abc-123')
  })

  it('throws on not found', async () => {
    mockFetch.mockResolvedValue({ ok: false, status: 404 })
    await expect(getRecipe('bad-id')).rejects.toThrow('404')
  })
})
