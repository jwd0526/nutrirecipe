import { useState } from 'react'
import type { ClarificationQ, QAPair } from '../api/types'

interface Props {
  questions: ClarificationQ[]
  onSubmit: (answers: QAPair[]) => void
  loading?: boolean
}

export function ClarificationDialog({ questions, onSubmit, loading = false }: Props) {
  const [answers, setAnswers] = useState<Record<string, string>>({})

  const allAnswered = questions.every(q => Boolean(answers[q.id]))

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    const pairs: QAPair[] = questions.map(q => ({
      question: q.question,
      answer: answers[q.id] ?? '',
    }))
    onSubmit(pairs)
  }

  const setAnswer = (id: string, value: string) =>
    setAnswers(prev => ({ ...prev, [id]: value }))

  return (
    <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
      <h2 style={{ margin: 0 }}>A few quick questions</h2>

      {questions.map(q => (
        <fieldset key={q.id} style={{ border: '1px solid #ccc', borderRadius: '6px', padding: '1rem' }}>
          <legend style={{ fontWeight: 600 }}>
            {q.question}{' '}
            <span style={{ color: '#888', fontWeight: 400 }}>({q.ingredient_raw})</span>
          </legend>

          <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem', marginTop: '0.5rem' }}>
            {q.options.map(opt => (
              <label key={opt} style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', cursor: 'pointer' }}>
                <input
                  type="radio"
                  name={q.id}
                  value={opt}
                  checked={answers[q.id] === opt}
                  onChange={() => setAnswer(q.id, opt)}
                />
                {opt}
              </label>
            ))}
          </div>
        </fieldset>
      ))}

      <button type="submit" disabled={!allAnswered || loading}>
        {loading ? 'Parsing…' : 'Continue'}
      </button>
    </form>
  )
}
