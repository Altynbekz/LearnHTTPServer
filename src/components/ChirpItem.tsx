import { Trash2 } from 'lucide-react'
import { useState } from 'react'
import { useAuth } from '../context/AuthContext'
import { api, type Chirp } from '../lib/api'
export function ChirpItem({ chirp, onDeleted }: { chirp: Chirp; onDeleted: (id: string) => void }) {
  const { user } = useAuth(); const [deleting, setDeleting] = useState(false); const isOwner = user?.id === chirp.user_id
  async function deleteChirp() { if (!window.confirm('Delete this chirp?')) return; setDeleting(true); try { await api.deleteChirp(chirp.id); onDeleted(chirp.id) } catch { setDeleting(false) } }
  return <article className="group border-b border-border py-7 first:border-t"><p className="whitespace-pre-wrap break-words text-[1.12rem] leading-8 text-charcoal">{chirp.body}</p><div className="mt-4 flex items-center justify-between gap-4 font-mono text-[11px] uppercase tracking-[0.08em] text-muted"><span>{new Date(chirp.created_at).toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' })} · {chirp.user_id.slice(0, 8)}...</span>{isOwner && <button aria-label="Delete chirp" disabled={deleting} onClick={deleteChirp} className="inline-flex items-center gap-1.5 text-muted opacity-0 transition hover:text-rust group-hover:opacity-100 disabled:opacity-50"><Trash2 aria-hidden="true" className="h-3.5 w-3.5" />{deleting ? 'Deleting' : 'Delete'}</button>}</div></article>
}