import { X } from 'lucide-react'
import { useEffect, useState } from 'react'
import { Link, useLocation, useSearchParams } from 'react-router-dom'
import { ChirpComposer } from '../components/ChirpComposer'
import { ChirpItem } from '../components/ChirpItem'
import { EmptyState } from '../components/EmptyState'
import { SectionLabel } from '../components/SectionLabel'
import { useAuth } from '../context/AuthContext'
import { api, ApiError, type Chirp, type SortOrder } from '../lib/api'

const pill = (active: boolean) =>
  `font-mono text-[11px] uppercase tracking-[0.08em] underline-offset-4 ${active ? 'text-rust underline' : 'text-muted hover:text-rust'}`

export function FeedPage() {
  const { isLoggedIn, user } = useAuth()
  const location = useLocation()
  const [searchParams, setSearchParams] = useSearchParams()
  const authorId = searchParams.get('author_id') || ''
  const sort: SortOrder = searchParams.get('sort') === 'asc' ? 'asc' : 'desc'

  const [chirps, setChirps] = useState<Chirp[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  // Filtering and sorting are done by the API: GET /api/chirps?author_id=...&sort=asc|desc
  useEffect(() => {
    let cancelled = false
    setLoading(true)
    setError('')
    api
      .getChirps({ authorId: authorId || undefined, sort })
      .then((items) => { if (!cancelled) setChirps(items) })
      .catch((err) => { if (!cancelled) setError(err instanceof ApiError ? err.message : 'The feed could not be loaded.') })
      .finally(() => { if (!cancelled) setLoading(false) })
    return () => { cancelled = true }
  }, [authorId, sort])

  // A chirp just published from the composer arrives via router state.
  useEffect(() => {
    const created = (location.state as { created?: Chirp } | null)?.created
    if (!created) return
    const fitsFilter = !authorId || authorId === created.user_id
    if (fitsFilter) {
      setChirps((current) => {
        const rest = current.filter((item) => item.id !== created.id)
        return sort === 'desc' ? [created, ...rest] : [...rest, created]
      })
    }
    window.history.replaceState({}, '', location.pathname + location.search)
  }, [location, authorId, sort])

  function updateParams(changes: Record<string, string | null>) {
    const next = new URLSearchParams(searchParams)
    Object.entries(changes).forEach(([key, value]) => (value ? next.set(key, value) : next.delete(key)))
    setSearchParams(next)
  }

  const isMine = Boolean(user && authorId === user.id)
  const heading = !authorId ? 'Recent chirps' : isMine ? 'Your chirps' : `Chirps by ${authorId.slice(0, 8)}…`

  return (
    <main className="mx-auto max-w-[760px] px-5 py-16 sm:py-24">
      <section>
        <SectionLabel>ソーシャルフィード · Social feed</SectionLabel>
        <h1 className="mt-5 max-w-xl font-serif text-5xl leading-[1.05] tracking-tight text-charcoal sm:text-6xl">Short thoughts,<br />simply shared.</h1>
        <p className="mt-6 max-w-lg text-base leading-7 text-muted">Chirpy is a small social API built with Go, made for sharing the ideas that do not need a stage.</p>
      </section>

      {isLoggedIn ? (
        <section className="mt-16 border-y border-border py-7"><ChirpComposer /></section>
      ) : (
        <section className="mt-16 flex flex-wrap items-center justify-between gap-5 border-y border-border py-6">
          <p className="text-sm text-muted">Join the quiet conversation.</p>
          <div className="flex gap-5 font-mono text-[11px] uppercase tracking-[0.08em]">
            <Link className="text-rust underline-offset-4 hover:underline" to="/login">Login</Link>
            <Link className="text-rust underline-offset-4 hover:underline" to="/register">Create account</Link>
          </div>
        </section>
      )}

      <section className="mt-16">
        <div className="mb-5 flex items-center justify-between">
          <SectionLabel>{heading}</SectionLabel>
          <span className="font-mono text-[11px] text-muted">{chirps.length} total</span>
        </div>

        <div className="mb-6 flex flex-wrap items-center justify-between gap-4 border-b border-border pb-4">
          <div className="flex items-center gap-4" role="group" aria-label="Filter by author">
            <button className={pill(!authorId)} onClick={() => updateParams({ author_id: null })}>All</button>
            {user && <button className={pill(isMine)} onClick={() => updateParams({ author_id: user.id })}>Mine</button>}
            {authorId && !isMine && (
              <button className="inline-flex items-center gap-1 font-mono text-[11px] uppercase tracking-[0.08em] text-rust" onClick={() => updateParams({ author_id: null })}>
                Author {authorId.slice(0, 8)}… <X aria-label="Clear author filter" className="h-3 w-3" />
              </button>
            )}
          </div>
          <div className="flex items-center gap-4" role="group" aria-label="Sort order">
            <button className={pill(sort === 'desc')} onClick={() => updateParams({ sort: null })}>Newest</button>
            <button className={pill(sort === 'asc')} onClick={() => updateParams({ sort: 'asc' })}>Oldest</button>
          </div>
        </div>

        {loading && (
          <div className="space-y-5" aria-label="Loading feed">
            <div className="h-28 animate-pulse border-y border-border bg-border/20" />
            <div className="h-24 animate-pulse border-b border-border bg-border/20" />
          </div>
        )}
        {error && <p role="alert" className="border-y border-rust/40 py-5 text-sm text-rust">{error}</p>}
        {!loading && !error && chirps.length === 0 && <EmptyState />}
        {!loading && !error && chirps.map((chirp) => (
          <ChirpItem
            key={chirp.id}
            chirp={chirp}
            onAuthorClick={(id) => updateParams({ author_id: id })}
            onDeleted={(id) => setChirps((current) => current.filter((item) => item.id !== id))}
          />
        ))}
      </section>
    </main>
  )
}
