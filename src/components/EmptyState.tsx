import { MessageCircle } from 'lucide-react'
export function EmptyState() {
  return <div className="border-y border-border py-12 text-center"><MessageCircle aria-hidden="true" className="mx-auto mb-4 h-5 w-5 text-muted" strokeWidth={1.5} /><p className="font-serif text-2xl text-charcoal">Nothing here yet.</p><p className="mt-2 text-sm text-muted">The first small thought is waiting to be shared.</p></div>
}