import { LogOut } from 'lucide-react'
import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
export function Navbar() {
  const { isLoggedIn, logout } = useAuth(); const navigate = useNavigate()
  async function signOut() { await logout(); navigate('/') }
  return <header className="sticky top-0 z-10 border-b border-border bg-paper"><div className="mx-auto flex max-w-[760px] items-center justify-between px-5 py-5"><Link to="/" className="font-serif text-2xl tracking-tight text-charcoal">Chirpy</Link><nav className="flex items-center gap-5 font-mono text-[11px] uppercase tracking-[0.08em] text-muted"><Link className="hover:text-rust" to="/">Feed</Link>{isLoggedIn && <Link className="hover:text-rust" to="/profile">Profile</Link>}{isLoggedIn ? <button onClick={signOut} className="inline-flex items-center gap-1.5 hover:text-rust"><LogOut aria-hidden="true" className="h-3.5 w-3.5" />Logout</button> : <Link className="hover:text-rust" to="/login">Login</Link>}</nav></div></header>
}