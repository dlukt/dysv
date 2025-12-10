import { createFileRoute, Link, useNavigate } from '@tanstack/react-router'
import { useState, useId } from 'react'
import { z } from 'zod'
import { useAuth } from '../hooks/use-auth'
import { Button } from '../components/ui/button'
import { Input } from '../components/ui/input'
import { Label } from '../components/ui/label'

const searchSchema = z.object({
  redirect: z.string().optional(),
})

export const Route = createFileRoute('/login')({
  component: LoginComponent,
  validateSearch: (search) => searchSchema.parse(search),
})

function LoginComponent() {
  const navigate = useNavigate()
  const search = Route.useSearch()
  const { login } = useAuth()
  const id = useId()
  
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    
    try {
      await login.mutateAsync({ email, password })
      // Redirect
      navigate({ to: search.redirect || '/' })
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to login'
      setError(message)
    }
  }

  return (
    <div className="container mx-auto max-w-md py-20 px-4">
      <div className="rounded-lg border bg-card p-8 shadow-sm">
        <h1 className="text-2xl font-bold mb-6 text-center">Welcome Back</h1>
        
        {error && (
          <div className="bg-red-50 text-red-600 p-3 rounded mb-4 text-sm">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor={`${id}-email`}>Email</Label>
            <Input
              id={`${id}-email`}
              type="email"
              placeholder="you@example.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </div>
          
          <div className="space-y-2">
            <Label htmlFor={`${id}-password`}>Password</Label>
            <Input
              id={`${id}-password`}
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>

          <Button type="submit" className="w-full" disabled={login.isPending}>
            {login.isPending ? 'Logging in...' : 'Login'}
          </Button>
        </form>

        <div className="mt-6 text-center text-sm text-gray-500">
          Don't have an account?{' '}
          <Link to="/register" className="text-primary hover:underline font-medium">
            Sign up
          </Link>
        </div>
      </div>
    </div>
  )
}
