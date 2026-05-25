import { describe, it, expect, vi, beforeEach } from 'vitest'
import { parseIngredients } from './agent'

const mockFetch = vi.fn()
vi.stubGlobal('fetch', mockFetch)

beforeEach(() => mockFetch.mockReset())

describe('parseIngredients', () => {
  it('returns resolved response on success', async () => {
    mockFetch.mockResolvedValue({
      ok: true,
      json: async () => ({ status: 'resolved', ingredients: [] }),
    })
    const result = await parseIngredients({ input: 'chicken breast' })
    expect(result.status).toBe('resolved')
    expect(mockFetch).toHaveBeenCalledWith('/api/agent/parse', expect.objectContaining({ method: 'POST' }))
  })

  it('returns needs_clarification for ambiguous input', async () => {
    mockFetch.mockResolvedValue({
      ok: true,
      json: async () => ({ status: 'needs_clarification', questions: [{ id: 'syrup_type' }] }),
    })
    const result = await parseIngredients({ input: '1 cup syrup' })
    expect(result.status).toBe('needs_clarification')
    expect(result.questions).toHaveLength(1)
  })

  it('sends qa_history when provided', async () => {
    mockFetch.mockResolvedValue({ ok: true, json: async () => ({ status: 'resolved', ingredients: [] }) })
    await parseIngredients({ input: 'syrup', qa_history: [{ question: 'type?', answer: 'maple syrup' }] })
    const body = JSON.parse(mockFetch.mock.calls[0][1].body)
    expect(body.qa_history).toHaveLength(1)
  })

  it('throws on non-ok response', async () => {
    mockFetch.mockResolvedValue({ ok: false, status: 500 })
    await expect(parseIngredients({ input: 'test' })).rejects.toThrow('500')
  })
})
