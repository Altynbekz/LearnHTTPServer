import { CheckCircle2, Sparkles } from 'lucide-react'
import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { ChirpyRedBadge } from '../components/ChirpyRedBadge'
import { SectionLabel } from '../components/SectionLabel'
import { useAuth } from '../context/AuthContext'
import { ApiError } from '../lib/api'
import { Error, Field } from './LoginPage'

export function ProfilePage() {
  const { user, updateUser, logout } = useAuth()
  const navigate = useNavigate()
  const [email, setEmail] = useState(user?.email || '')
  const [password, setPassword] = useState('')
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const isRed = Boolean(user?.is_chirpy_red)

  async function submit(event: React.FormEvent) {
    event.preventDefault()
    setLoading(true)
    setError('')
    setMessage('')
    try {
      await updateUser(email, password)
      setPassword('')
      setMessage('Profile updated.')
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Your profile could not be updated.')
    } finally {
      setLoading(false)
    }
  }

  async function signOut() { await logout(); navigate('/') }

  return (
    <main className="mx-auto max-w-[760px] px-5 py-16 sm:py-24">
      <SectionLabel>設定 · Settings</SectionLabel>
      <h1 className="mt-5 font-serif text-5xl text-charcoal">Your profile.</h1>

      <div className="mt-10 flex items-center gap-3 font-mono text-[11px] uppercase tracking-[0.08em] text-sage">
        <CheckCircle2 aria-hidden="true" className="h-4 w-4" />Authenticated
        {isRed && <ChirpyRedBadge />}
      </div>
      <p className="mt-3 text-sm text-muted">{user?.email}</p>
      {user && (
        <Link to={`/?author_id=${user.id}`} className="mt-4 inline-block font-mono text-[11px] uppercase tracking-[0.08em] text-rust underline-offset-4 hover:underline">
          View my chirps →
        </Link>
      )}

      <section className="mt-14 border-y border-border py-7">
        <SectionLabel>Membership</SectionLabel>
        <div className="mt-4 flex items-start gap-3">
          <Sparkles aria-hidden="true" className={`mt-1 h-5 w-5 ${isRed ? 'text-rust' : 'text-muted'}`} />
          <div>
            <p className="font-serif text-2xl text-charcoal">{isRed ? 'Chirpy Red member' : 'Free plan'}</p>
            <p className="mt-2 max-w-md text-sm leading-6 text-muted">
              {isRed
                ? 'Thanks for supporting Chirpy. Your upgrade was confirmed by our payment provider.'
                : 'Upgrades are processed by our payment provider, Polka. Once your payment is confirmed, log in again to see your Chirpy Red status.'}
            </p>
          </div>
        </div>
      </section>

      <form onSubmit={submit} className="mt-14 max-w-[500px] space-y-7">
        <Field label="Email" type="email" value={email} onChange={setEmail} />
        <Field label="New password" type="password" value={password} onChange={setPassword} />
        {error && <Error text={error} />}
        {message && <p className="text-sm text-sage">{message}</p>}
        <button disabled={loading || !password} className="font-mono text-[11px] uppercase tracking-[0.08em] text-rust underline-offset-4 hover:underline disabled:opacity-40">
          {loading ? 'Saving...' : 'Save changes →'}
        </button>
      </form>

      <div className="mt-16 border-t border-border pt-6">
        <button onClick={signOut} className="font-mono text-[11px] uppercase tracking-[0.08em] text-muted hover:text-rust">Logout</button>
      </div>
    </main>
  )
}
