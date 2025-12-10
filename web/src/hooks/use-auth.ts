import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useRouter } from '@tanstack/react-router'
import { z } from 'zod'

// Types
export interface User {
  id: string
  email: string
  username: string
  role: string
  is_verified: boolean
  created_at: string
}

export interface AuthResponse {
  user: User
  token: string
  session: unknown
}

export const LoginSchema = z.object({
  email: z.string().email(),
  password: z.string().min(1, 'Password is required'),
})

export type LoginCredentials = z.infer<typeof LoginSchema>

export const RegisterSchema = z.object({
  username: z.string().min(3, 'Username must be at least 3 characters'),
  email: z.string().email(),
  password: z.string().min(8, 'Password must be at least 8 characters'),
})

export type RegisterCredentials = z.infer<typeof RegisterSchema>

const AUTH_KEY = 'dysv_auth_token'

// API Helpers
async function fetchWithAuth(url: string, options: RequestInit = {}) {
  const token = localStorage.getItem(AUTH_KEY)
  const headers = {
    'Content-Type': 'application/json',
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
    ...options.headers,
  }

  const res = await fetch(url, { ...options, headers })
  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: 'Unknown error' }))
    throw new Error(error.error || error.message || 'Request failed')
  }
  return res.json()
}

export function useAuth() {
  const queryClient = useQueryClient()
  const router = useRouter()

  // Load user
  const { data: user, isLoading, error } = useQuery({
    queryKey: ['auth', 'me'],
    queryFn: async () => {
      const token = localStorage.getItem(AUTH_KEY)
      if (!token) return null
      try {
        const res = await fetchWithAuth('/api/auth/me')
        return res.user as User
      } catch (_err) {
        // If 401, clear token
        localStorage.removeItem(AUTH_KEY)
        return null
      }
    },
    retry: false,
    staleTime: 1000 * 60 * 5, // 5 minutes
  })

  // Login
  const loginMutation = useMutation({
    mutationFn: async (creds: LoginCredentials) => {
      const res = await fetch('/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(creds),
      })
      if (!res.ok) {
        const err = await res.json()
        throw new Error(err.error || 'Login failed')
      }
      return res.json() as Promise<AuthResponse>
    },
    onSuccess: (data) => {
      localStorage.setItem(AUTH_KEY, data.token)
      queryClient.setQueryData(['auth', 'me'], data.user)
    },
  })

  // Register
  const registerMutation = useMutation({
    mutationFn: async (creds: RegisterCredentials) => {
      const res = await fetch('/api/auth/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(creds),
      })
      if (!res.ok) {
        const err = await res.json()
        throw new Error(err.error || 'Registration failed')
      }
      return res.json() as Promise<AuthResponse>
    },
    onSuccess: (data) => {
      localStorage.setItem(AUTH_KEY, data.token)
      queryClient.setQueryData(['auth', 'me'], data.user)
    },
  })

  // Logout
  const logout = () => {
    localStorage.removeItem(AUTH_KEY)
    queryClient.setQueryData(['auth', 'me'], null)
    queryClient.invalidateQueries({ queryKey: ['auth'] })
    router.invalidate()
  }

  return {
    user,
    isLoading,
    error,
    login: loginMutation,
    register: registerMutation,
    logout,
    isAuthenticated: !!user,
  }
}
