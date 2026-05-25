import type { AgentParseRequest, AgentParseResponse } from './types'

export async function parseIngredients(req: AgentParseRequest): Promise<AgentParseResponse> {
  const res = await fetch('/api/agent/parse', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(req),
  })
  if (!res.ok) throw new Error(`agent/parse failed: ${res.status}`)
  return res.json()
}
