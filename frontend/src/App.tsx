import { useState } from 'react'
import { NewRecipe } from './pages/NewRecipe'
import { RecipeList } from './pages/RecipeList'

type Page = 'new' | 'list'

export default function App() {
  const [page, setPage] = useState<Page>('list')

  return (
    <div>
      <nav style={{ borderBottom: '1px solid #ddd', padding: '0.75rem 2rem', display: 'flex', gap: '1rem', alignItems: 'center' }}>
        <strong style={{ marginRight: '1rem' }}>NutriRecipe</strong>
        <button
          onClick={() => setPage('list')}
          style={{ fontWeight: page === 'list' ? 700 : 400, background: 'none', border: 'none', cursor: 'pointer' }}
        >
          My Recipes
        </button>
        <button
          onClick={() => setPage('new')}
          style={{ fontWeight: page === 'new' ? 700 : 400, background: 'none', border: 'none', cursor: 'pointer' }}
        >
          + New Recipe
        </button>
      </nav>

      {page === 'list' && <RecipeList />}
      {page === 'new' && <NewRecipe />}
    </div>
  )
}
