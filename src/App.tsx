import { GitBranch } from 'lucide-react'
import { Navigate, Route, Routes, useLocation } from 'react-router-dom'
import { Navbar } from './components/Navbar'
import { ChirpComposer } from './components/ChirpComposer'
import { SectionLabel } from './components/SectionLabel'
import { useAuth } from './context/AuthContext'
import { FeedPage } from './pages/FeedPage'
import { LoginPage } from './pages/LoginPage'
import { ProfilePage } from './pages/ProfilePage'
import { RegisterPage } from './pages/RegisterPage'
function Protected({ children }: { children: React.ReactNode }) { const { isLoggedIn } = useAuth(); const location = useLocation(); return isLoggedIn ? children : <Navigate to="/login" state={{ from: location.pathname }} replace /> }
function Footer() { return <footer className="mx-auto flex max-w-[760px] flex-col gap-4 border-t border-border px-5 py-8 font-mono text-[11px] uppercase tracking-[0.08em] text-muted sm:flex-row sm:items-center sm:justify-between"><a href="https://github.com" target="_blank" rel="noreferrer" className="inline-flex items-center gap-2 hover:text-rust"><GitBranch aria-hidden="true" className="h-3.5 w-3.5" />GitHub</a><span>Built with Go · net/http · PostgreSQL</span><span>{new Date().getFullYear()}</span></footer> }
export function App() { return <><Navbar /><Routes><Route path="/" element={<FeedPage />} /><Route path="/login" element={<LoginPage />} /><Route path="/register" element={<RegisterPage />} /><Route path="/compose" element={<Protected><main className="mx-auto max-w-[760px] px-5 py-16"><SectionLabel>New chirp</SectionLabel><h1 className="mt-5 font-serif text-5xl text-charcoal">Say something.</h1><div className="mt-12 border-y border-border py-7"><ChirpComposer /></div></main></Protected>} /><Route path="/profile" element={<Protected><ProfilePage /></Protected>} /><Route path="*" element={<Navigate to="/" replace />} /></Routes><Footer /></> }