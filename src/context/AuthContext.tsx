import { createContext, useContext, useState, type ReactNode } from 'react'
import { api, authStorage, type User } from '../lib/api'
type AuthContextValue = { user: User | null; isLoggedIn: boolean; login: (email: string, password: string) => Promise<void>; logout: () => Promise<void>; updateUser: (email: string, password: string) => Promise<void> }
const AuthContext = createContext<AuthContextValue | null>(null)
export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(() => authStorage.getUser())
  async function login(email: string, password: string) { const result = await api.login(email, password); authStorage.save(result); setUser(result) }
  async function logout() { try { await api.revoke() } finally { authStorage.clear(); setUser(null) } }
  async function updateUser(email: string, password: string) { const result = await api.updateUser(email, password); const updated = { ...user, email: result.email } as User; localStorage.setItem('chirpy_user', JSON.stringify(updated)); setUser(updated) }
  return <AuthContext.Provider value={{ user, isLoggedIn: Boolean(user), login, logout, updateUser }}>{children}</AuthContext.Provider>
}
export function useAuth() { const context = useContext(AuthContext); if (!context) throw new Error('useAuth must be used inside AuthProvider'); return context }