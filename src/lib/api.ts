export type User = { id: string; email: string; is_chirpy_red?: boolean; created_at?: string; updated_at?: string }
export type SortOrder = 'asc' | 'desc'
export type ChirpQuery = { authorId?: string; sort?: SortOrder }
export type Chirp = { id: string; body: string; user_id: string; created_at: string; updated_at: string }
export type AuthResponse = User & { token: string; refresh_token: string }

const API_URL = (import.meta.env.VITE_API_URL || 'http://localhost:8080').replace(/\/$/, '')
const keys = { token: 'chirpy_token', refresh: 'chirpy_refresh_token', user: 'chirpy_user' }
export const authStorage = {
  getToken: () => localStorage.getItem(keys.token), getRefreshToken: () => localStorage.getItem(keys.refresh),
  getUser: (): User | null => JSON.parse(localStorage.getItem(keys.user) || 'null'),
  save: (auth: AuthResponse) => { localStorage.setItem(keys.token, auth.token); localStorage.setItem(keys.refresh, auth.refresh_token); localStorage.setItem(keys.user, JSON.stringify({ id: auth.id, email: auth.email, is_chirpy_red: auth.is_chirpy_red ?? false, created_at: auth.created_at, updated_at: auth.updated_at })) },
  saveUser: (user: User) => localStorage.setItem(keys.user, JSON.stringify(user)),
  saveToken: (token: string) => localStorage.setItem(keys.token, token),
  clear: () => Object.values(keys).forEach((key) => localStorage.removeItem(key)),
}
export class ApiError extends Error { constructor(message: string, public status: number) { super(message) } }

async function parseResponse(response: Response): Promise<unknown> {
  if (response.status === 204) return null
  const data = await response.json().catch(() => null)
  if (!response.ok) throw new ApiError((data as { error?: string } | null)?.error || `Request failed (${response.status})`, response.status)
  return data
}
async function refreshToken() {
  const refresh = authStorage.getRefreshToken(); if (!refresh) throw new ApiError('Your session has expired.', 401)
  const data = await parseResponse(await fetch(`${API_URL}/api/refresh`, { method: 'POST', headers: { Authorization: `Bearer ${refresh}` } })) as { token: string }
  authStorage.saveToken(data.token); return data.token
}
export async function request<T>(path: string, options: RequestInit = {}, protectedRequest = false): Promise<T> {
  const headers = new Headers(options.headers); if (options.body && !headers.has('Content-Type')) headers.set('Content-Type', 'application/json')
  const token = authStorage.getToken(); if (protectedRequest && token) headers.set('Authorization', `Bearer ${token}`)
  let response = await fetch(`${API_URL}${path}`, { ...options, headers })
  if (protectedRequest && response.status === 401) { try { headers.set('Authorization', `Bearer ${await refreshToken()}`); response = await fetch(`${API_URL}${path}`, { ...options, headers }) } catch (error) { authStorage.clear(); throw error } }
  return await parseResponse(response) as T
}
export const api = {
  getChirps: ({ authorId, sort }: ChirpQuery = {}) => {
    const params = new URLSearchParams()
    if (authorId) params.set('author_id', authorId)
    if (sort) params.set('sort', sort)
    const query = params.toString()
    // The Go handler encodes a nil slice as `null` when there are no chirps, so normalise to [].
    return request<Chirp[] | null>(`/api/chirps${query ? `?${query}` : ''}`).then((items) => items ?? [])
  },
  register: (email: string, password: string) => request<User>('/api/users', { method: 'POST', body: JSON.stringify({ email, password }) }),
  login: (email: string, password: string) => request<AuthResponse>('/api/login', { method: 'POST', body: JSON.stringify({ email, password }) }),
  createChirp: (body: string) => request<Chirp>('/api/chirps', { method: 'POST', body: JSON.stringify({ body }) }, true),
  deleteChirp: (id: string) => request<null>(`/api/chirps/${id}`, { method: 'DELETE' }, true),
  updateUser: (email: string, password: string) => request<{ email: string }>('/api/users', { method: 'PUT', body: JSON.stringify({ email, password }) }, true),
  revoke: () => request<null>('/api/revoke', { method: 'POST', headers: { Authorization: `Bearer ${authStorage.getRefreshToken() || ''}` } }),
}