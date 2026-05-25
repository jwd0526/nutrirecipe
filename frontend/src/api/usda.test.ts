import { describe, it, expect, vi, beforeEach } from 'vitest'
import { validateIngredients } from './usda'
import type { ParsedIngredient } from './types'

const mockFetch = vi.fn()
vi.stubGlobal('fetch', mockFetch)

beforeEach(() => mockFetch.mockReset())

const ingredient: ParsedIngredient = {
  name: 'chicken breast',
  portion: 1,
  unit: 'serving',
  quantity_g: 100,
  search_queries: { primary: 'chicken breast', alternatives: [] },
  confidence: 'high',
}

describe('validateIngredients', () => {
  it('returns validated ingredients on success', async () => {
    const validated = [{ ...ingredient, match_status: 'matched', fdc_id: '123' }]
    mockFetch.mockResolvedValue({ ok: true, json: async () => validated })
    const result = await validateIngredients([ingredient])
    expect(result).toHaveLength(1)
    expect(result[0].match_status).toBe('matched')
  })

  it('returns empty array for empty input', async () => {
    mockFetch.mockResolvedValue({ ok: true, json: async () => [] })
    const result = await validateIngredients([])
    expect(result).toHaveLength(0)
  })

  it('sends ingredients wrapped in object', async () => {
    mockFetch.mockResolvedValue({ ok: true, json: async () => [] })
    await validateIngredients([ingredient])
    const body = JSON.parse(mockFetch.mock.calls[0][1].body)
    expect(body).toHaveProperty('ingredients')
    expect(body.ingredients).toHaveLength(1)
  })

  it('throws on non-ok response', async () => {
    mockFetch.mockResolvedValue({ ok: false, status: 400 })
    await expect(validateIngredients([ingredient])).rejects.toThrow('400')
  })
})
