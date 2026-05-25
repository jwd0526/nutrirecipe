import { useState } from 'react'
import type { ValidatedIngredient, USDAOption } from '../api/types'

interface Props {
  ingredients: ValidatedIngredient[]
  onChange: (updated: ValidatedIngredient[]) => void
  onSave: () => void
  saving?: boolean
}

const BADGE: Record<string, { label: string; color: string }> = {
  matched:        { label: 'Matched',        color: '#2d8a4e' },
  low_confidence: { label: 'Low Confidence', color: '#b07d00' },
  unresolved:     { label: 'Unresolved',     color: '#c0392b' },
}

export function ReviewTable({ ingredients, onChange, onSave, saving = false }: Props) {
  const [expanded, setExpanded] = useState<Set<number>>(() => {
    const s = new Set<number>()
    ingredients.forEach((ing, i) => { if (ing.match_status !== 'matched') s.add(i) })
    return s
  })

  const toggle = (i: number) =>
    setExpanded(prev => { const s = new Set(prev); s.has(i) ? s.delete(i) : s.add(i); return s })

  const patch = (i: number, p: Partial<ValidatedIngredient>) =>
    onChange(ingredients.map((ing, j) => j === i ? { ...ing, ...p } : ing))

  return (
    <div>
      <table style={{ width: '100%', borderCollapse: 'collapse' }}>
        <thead>
          <tr style={{ borderBottom: '2px solid #ddd' }}>
            {['Ingredient', 'Status', 'g', 'kcal', 'prot (g)', 'carbs (g)', 'fat (g)'].map(h => (
              <th key={h} style={{ textAlign: 'left', padding: '0.5rem' }}>{h}</th>
            ))}
          </tr>
        </thead>
        <tbody>
          {ingredients.map((ing, i) => (
            <IngredientRow
              key={i}
              ing={ing}
              expanded={expanded.has(i)}
              onToggle={() => toggle(i)}
              onPatch={p => patch(i, p)}
            />
          ))}
        </tbody>
      </table>

      <div style={{ marginTop: '1rem', textAlign: 'right' }}>
        <button onClick={onSave} disabled={saving}>
          {saving ? 'Saving…' : 'Save Recipe'}
        </button>
      </div>
    </div>
  )
}

interface RowProps {
  ing: ValidatedIngredient
  expanded: boolean
  onToggle: () => void
  onPatch: (p: Partial<ValidatedIngredient>) => void
}

const total = (per100g: number | null | undefined, g: number) => {
  if (per100g == null || !g) return '—'
  return (per100g * g / 100).toFixed(1)
}

function IngredientRow({ ing, expanded, onToggle, onPatch }: RowProps) {
  const badge = BADGE[ing.match_status] ?? { label: ing.match_status, color: '#888' }

  return (
    <>
      <tr style={{ borderBottom: '1px solid #eee' }}>
        <td style={{ padding: '0.5rem' }}>
          <button type="button" onClick={onToggle} style={{ background: 'none', border: 'none', cursor: 'pointer', textAlign: 'left' }}>
            {ing.name}
          </button>
        </td>
        <td style={{ padding: '0.5rem' }}>
          <span style={{ color: badge.color, fontWeight: 600, fontSize: '0.8rem' }}>{badge.label}</span>
        </td>
        <td style={{ padding: '0.5rem' }}>{ing.quantity_g}</td>
        <td style={{ padding: '0.5rem' }}>{total(ing.calories_per_100g, ing.quantity_g)}</td>
        <td style={{ padding: '0.5rem' }}>{total(ing.protein_per_100g, ing.quantity_g)}</td>
        <td style={{ padding: '0.5rem' }}>{total(ing.carbs_per_100g, ing.quantity_g)}</td>
        <td style={{ padding: '0.5rem' }}>{total(ing.fat_per_100g, ing.quantity_g)}</td>
      </tr>

      {expanded && (
        <tr>
          <td colSpan={7} style={{ padding: '0.75rem 0.5rem', background: '#fafafa' }}>
            {ing.match_status === 'low_confidence' && (
              <OptionsDropdown ing={ing} onPatch={onPatch} />
            )}
            {ing.match_status === 'unresolved' && (
              <MacroEntry ing={ing} onPatch={onPatch} />
            )}
          </td>
        </tr>
      )}
    </>
  )
}

function OptionsDropdown({ ing, onPatch }: { ing: ValidatedIngredient; onPatch: (p: Partial<ValidatedIngredient>) => void }) {
  const options = ing.options ?? []
  const handleChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const opt: USDAOption | undefined = options.find(o => o.fdc_id === e.target.value)
    if (!opt) return
    onPatch({ fdc_id: opt.fdc_id, fdc_name: opt.name, match_status: 'matched',
      calories_per_100g: opt.calories, protein_per_100g: opt.protein,
      carbs_per_100g: opt.carbs, fat_per_100g: opt.fat })
  }
  return (
    <label style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
      Select best match:
      <select value={ing.fdc_id ?? ''} onChange={handleChange} style={{ flex: 1 }}>
        <option value="">— choose —</option>
        {options.map(o => (
          <option key={o.fdc_id} value={o.fdc_id}>{o.name} ({o.calories} kcal)</option>
        ))}
      </select>
    </label>
  )
}

function MacroEntry({ ing, onPatch }: { ing: ValidatedIngredient; onPatch: (p: Partial<ValidatedIngredient>) => void }) {
  const field = (label: string, key: keyof ValidatedIngredient) => (
    <label style={{ display: 'flex', flexDirection: 'column', fontSize: '0.85rem' }}>
      {label}
      <input
        type="number"
        min={0}
        step={0.1}
        value={(ing[key] as number | undefined) ?? ''}
        onChange={e => onPatch({ [key]: parseFloat(e.target.value) || 0 })}
        style={{ width: '80px' }}
      />
    </label>
  )
  return (
    <div style={{ display: 'flex', gap: '1rem', alignItems: 'flex-end' }}>
      <span style={{ fontStyle: 'italic', color: '#888' }}>Enter macros manually:</span>
      {field('kcal/100g', 'calories_per_100g')}
      {field('protein', 'protein_per_100g')}
      {field('carbs', 'carbs_per_100g')}
      {field('fat', 'fat_per_100g')}
    </div>
  )
}
