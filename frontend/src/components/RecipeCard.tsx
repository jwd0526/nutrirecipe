import type { Recipe } from '../api/types'

interface Props {
  recipe: Recipe
  onClick?: () => void
}

export function RecipeCard({ recipe, onClick }: Props) {
  return (
    <div
      onClick={onClick}
      style={{
        border: '1px solid #ddd',
        borderRadius: '8px',
        padding: '1rem',
        cursor: onClick ? 'pointer' : 'default',
        transition: 'box-shadow 0.15s',
      }}
      onMouseEnter={e => { if (onClick) (e.currentTarget as HTMLDivElement).style.boxShadow = '0 2px 8px rgba(0,0,0,0.12)' }}
      onMouseLeave={e => { (e.currentTarget as HTMLDivElement).style.boxShadow = 'none' }}
    >
      <h3 style={{ margin: '0 0 0.5rem' }}>{recipe.name}</h3>
      <p style={{ margin: '0 0 0.25rem', color: '#555' }}>
        {recipe.total_weight_g}g total · {recipe.serving_size_g}g per serving
      </p>
      <p style={{ margin: 0, color: '#999', fontSize: '0.85rem' }}>
        {new Date(recipe.created_at).toLocaleDateString()}
      </p>
    </div>
  )
}
