import { Trash2 } from 'lucide-react'
import { useState } from 'react'
import { useAuth } from '../context/AuthContext'
import { api, ApiError, type Chirp } from '../lib/api'
import { ChirpyRedBadge } from './ChirpyRedBadge'

type Props = { chirp: Chirp; onDeleted: (id: string) => void; onAuthorClick?: (userId: string) => void }

export function ChirpItem({ chirp, onDeleted, onAuthorClick }: Props) {
  const { user } = useAuth()
  const [deleting, setDeleting] = useState(false)
  const [error, setError] = useState('')
  const isOwner = user?.id === chirp.user_id

  async function deleteChirp() {
    if (!window.confirm('Delete this chirp?')) return
    setDeleting(true)
    setError('')
    try {
      await api.deleteChirp(chirp.id)
      onDeleted(chirp.id)
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'The chirp could not be deleted.')
      setDeleting(false)
    }
  }

  const date = new Date(chirp.created_at).toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' })
  const authorLabel = isOwner ? 'You' : `${chirp.user_id.slice(0, 8)}…`

  return (
    <article className="group border-b border-border py-7 first:border-t">
      <p className="whitespace-pre-wrap break-words text-[1.12rem] leading-8 text-charcoal">{chirp.body}</p>
      <div className="mt-4 flex items-center justify-between gap-4 font-mono text-[11px] uppercase tracking-[0.08em] text-muted">
        <span className="inline-flex items-center gap-2">
          {date} ·
          {onAuthorClick ? (
            <button title="Show only this author's chirps" onClick={() => onAuthorClick(chirp.user_id)} className="uppercase underline-offset-4 hover:text-rust hover:underline">
              {authorLabel}
            </button>
          ) : (
            <span>{authorLabel}</span>
          )}
          {isOwner && user?.is_chirpy_red && <ChirpyRedBadge />}
        </span>
        {isOwner && (
          <button aria-label="Delete chirp" disabled={deleting} onClick={deleteChirp} className="inline-flex items-center gap-1.5 text-muted opacity-0 transition hover:text-rust focus:opacity-100 group-hover:opacity-100 disabled:opacity-50">
            <Trash2 aria-hidden="true" className="h-3.5 w-3.5" />
            {deleting ? 'Deleting' : 'Delete'}
          </button>
        )}
      </div>
      {error && <p role="alert" className="mt-3 text-sm text-rust">{error}</p>}
    </article>
  )
}
