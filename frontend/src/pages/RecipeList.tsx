import { useEffect, useState } from 'react'
import { RecipeCard } from '../components/RecipeCard'
import { listRecipes } from '../api/recipes'
import type { Recipe } from '../api/types'

export function RecipeList() {
  const [recipes, setRecipes] = useState<Recipe[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    listRecipes()
      .then(setRecipes)
      .catch(() => setError('Failed to load recipes.'))
      .finally(() => setLoading(false))
  }, [])

  if (loading) return <p style={{ padding: '2rem' }}>Loading…</p>

  return (
    <div style={{ maxWidth: '800px', margin: '0 auto', padding: '2rem' }}>
      <h1>My Recipes</h1>

      {error && (
        <div style={{ color: '#c0392b', background: '#fdecea', padding: '0.75rem', borderRadius: '6px', marginBottom: '1rem' }}>
          {error}
        </div>
      )}

      {recipes.length === 0 && !error && (
        <p style={{ color: '#888' }}>No recipes yet. Create your first one!</p>
      )}

      <div style={{ display: 'grid', gap: '1rem', gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))' }}>
        {recipes.map(r => (
          <RecipeCard key={r.id} recipe={r} />
        ))}
      </div>
    </div>
  )
}
