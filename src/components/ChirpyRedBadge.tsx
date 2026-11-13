import { Sparkles } from 'lucide-react'

export function ChirpyRedBadge({ className = '' }: { className?: string }) {
  return (
    <span
      title="Chirpy Red member"
      className={`inline-flex items-center gap-1 border border-rust/40 px-1.5 py-0.5 font-mono text-[10px] uppercase tracking-[0.12em] text-rust ${className}`}
    >
      <Sparkles aria-hidden="true" className="h-3 w-3" />
      Red
    </span>
  )
}
