import { useState } from 'react'
import { IngredientInput } from '../components/IngredientInput'
import { ClarificationDialog } from '../components/ClarificationDialog'
import { ReviewTable } from '../components/ReviewTable'
import { parseIngredients } from '../api/agent'
import { validateIngredients } from '../api/usda'
import { saveRecipe } from '../api/recipes'
import type { ClarificationQ, QAPair, ValidatedIngredient } from '../api/types'

type Step = 'input' | 'clarification' | 'review' | 'saved'

export function NewRecipe() {
  const [step, setStep] = useState<Step>('input')
  const [recipeName, setRecipeName] = useState('')
  const [rawInput, setRawInput] = useState('')
  const [questions, setQuestions] = useState<ClarificationQ[]>([])
  const [qaHistory, setQAHistory] = useState<QAPair[]>([])
  const [ingredients, setIngredients] = useState<ValidatedIngredient[]>([])
  const [loading, setLoading] = useState(false)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const doParse = async (input: string, history: QAPair[]) => {
    setLoading(true)
    setError(null)
    try {
      const parsed = await parseIngredients({ input, qa_history: history })
      if (parsed.status === 'needs_clarification') {
        setQuestions(parsed.questions ?? [])
        setStep('clarification')
      } else {
        const validated = await validateIngredients(parsed.ingredients ?? [])
        setIngredients(validated)
        setStep('review')
      }
    } catch {
      setError('Failed to parse ingredients. Please try again.')
    } finally {
      setLoading(false)
    }
  }

  const handleInputSubmit = (name: string, input: string) => {
    setRecipeName(name)
    setRawInput(input)
    setQAHistory([])
    doParse(input, [])
  }

  const handleClarification = (answers: QAPair[]) => {
    const history = [...qaHistory, ...answers]
    setQAHistory(history)
    doParse(rawInput, history)
  }

  const handleSave = async () => {
    setSaving(true)
    setError(null)
    try {
      await saveRecipe({ name: recipeName, ingredients })
      setStep('saved')
    } catch {
      setError('Failed to save recipe. Please try again.')
    } finally {
      setSaving(false)
    }
  }

  const reset = () => {
    setStep('input')
    setRecipeName('')
    setRawInput('')
    setQAHistory([])
    setIngredients([])
    setError(null)
  }

  return (
    <div style={{ maxWidth: '800px', margin: '0 auto', padding: '2rem' }}>
      <h1>New Recipe</h1>

      {error && (
        <div style={{ color: '#c0392b', background: '#fdecea', padding: '0.75rem', borderRadius: '6px', marginBottom: '1rem' }}>
          {error}
        </div>
      )}

      {step === 'input' && (
        <IngredientInput onSubmit={handleInputSubmit} loading={loading} />
      )}

      {step === 'clarification' && (
        <ClarificationDialog questions={questions} onSubmit={handleClarification} loading={loading} />
      )}

      {step === 'review' && (
        <ReviewTable
          ingredients={ingredients}
          onChange={setIngredients}
          onSave={handleSave}
          saving={saving}
        />
      )}

      {step === 'saved' && (
        <div style={{ textAlign: 'center', padding: '2rem' }}>
          <p style={{ fontSize: '1.25rem' }}>Recipe saved!</p>
          <button onClick={reset}>Create Another</button>
        </div>
      )}
    </div>
  )
}
