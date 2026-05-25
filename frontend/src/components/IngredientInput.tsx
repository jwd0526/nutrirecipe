import { useState } from 'react'

interface Props {
  onSubmit: (name: string, input: string) => void
  loading?: boolean
}

export function IngredientInput({ onSubmit, loading = false }: Props) {
  const [name, setName] = useState('')
  const [input, setInput] = useState('')

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (name.trim() && input.trim()) {
      onSubmit(name.trim(), input.trim())
    }
  }

  const canSubmit = !loading && name.trim().length > 0 && input.trim().length > 0

  return (
    <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
      <div>
        <label htmlFor="recipe-name" style={{ display: 'block', marginBottom: '0.25rem' }}>
          Recipe Name
        </label>
        <input
          id="recipe-name"
          value={name}
          onChange={e => setName(e.target.value)}
          placeholder="e.g. Chicken Stir Fry"
          style={{ width: '100%' }}
          required
        />
      </div>

      <div>
        <label htmlFor="ingredients" style={{ display: 'block', marginBottom: '0.25rem' }}>
          Ingredients <span style={{ color: '#888' }}>(one per line)</span>
        </label>
        <textarea
          id="ingredients"
          value={input}
          onChange={e => setInput(e.target.value)}
          placeholder={'2 cups chicken breast\n1 tbsp olive oil\n100g brown rice'}
          rows={8}
          style={{ width: '100%', fontFamily: 'monospace' }}
          required
        />
      </div>

      <button type="submit" disabled={!canSubmit}>
        {loading ? 'Parsing…' : 'Parse Ingredients'}
      </button>
    </form>
  )
}
