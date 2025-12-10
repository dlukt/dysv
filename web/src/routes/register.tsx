import { createFileRoute, Link, useNavigate } from '@tanstack/react-router'
import { useId } from 'react'
import type { z } from 'zod'
import { useForm } from '@tanstack/react-form'
import { useAuth, RegisterSchema } from '../hooks/use-auth'
import { Button } from '../components/ui/button'
import { Input } from '../components/ui/input'
import { Label } from '../components/ui/label'

export const Route = createFileRoute('/register')({
  component: RegisterComponent,
})

function validateWith<T>(schema: z.ZodType<T>) {
  return ({ value }: { value: T }) => {
    const result = schema.safeParse(value)
    if (result.success) return undefined
    return result.error.issues[0].message
  }
}

function RegisterComponent() {
  const navigate = useNavigate()
  const { register } = useAuth()
  const id = useId()
  
  const form = useForm({
    defaultValues: {
      username: '',
      email: '',
      password: '',
    },
    onSubmit: async ({ value }) => {
        await register.mutateAsync(value)
        navigate({ to: '/' })
    },
  })

  return (
    <div className="container mx-auto max-w-md py-20 px-4">
      <div className="rounded-lg border bg-card p-8 shadow-sm">
        <h1 className="text-2xl font-bold mb-6 text-center">Create Account</h1>
        
        {register.error && (
          <div className="bg-red-50 text-red-600 p-3 rounded mb-4 text-sm">
            {register.error.message || 'Failed to register'}
          </div>
        )}

        <form 
            onSubmit={(e) => {
                e.preventDefault()
                e.stopPropagation()
                form.handleSubmit()
            }}
            className="space-y-4"
        >
          <form.Field
            name="username"
            validators={{ onChange: validateWith(RegisterSchema.shape.username) }}
          >
            {(field) => (
                <div className="space-y-2">
                    <Label htmlFor={`${id}-username`}>Username</Label>
                    <Input
                        id={`${id}-username`}
                        name={field.name}
                        type="text"
                        placeholder="johndoe"
                        value={field.state.value}
                        onBlur={field.handleBlur}
                        onChange={(e) => field.handleChange(e.target.value)}
                        required
                    />
                     {field.state.meta.errors ? (
                        <p className="text-xs text-red-500">{field.state.meta.errors.join(', ')}</p>
                    ) : null}
                </div>
            )}
          </form.Field>

          <form.Field
            name="email"
            validators={{ onChange: validateWith(RegisterSchema.shape.email) }}
          >
            {(field) => (
                <div className="space-y-2">
                    <Label htmlFor={`${id}-email`}>Email</Label>
                    <Input
                        id={`${id}-email`}
                        name={field.name}
                        type="email"
                        placeholder="you@example.com"
                        value={field.state.value}
                        onBlur={field.handleBlur}
                        onChange={(e) => field.handleChange(e.target.value)}
                        required
                    />
                     {field.state.meta.errors ? (
                        <p className="text-xs text-red-500">{field.state.meta.errors.join(', ')}</p>
                    ) : null}
                </div>
            )}
          </form.Field>
          
          <form.Field
            name="password"
            validators={{ onChange: validateWith(RegisterSchema.shape.password) }}
          >
            {(field) => (
                <div className="space-y-2">
                    <Label htmlFor={`${id}-password`}>Password</Label>
                    <Input
                        id={`${id}-password`}
                        name={field.name}
                        type="password"
                        value={field.state.value}
                        onBlur={field.handleBlur}
                        onChange={(e) => field.handleChange(e.target.value)}
                        required
                    />
                     {field.state.meta.errors ? (
                        <p className="text-xs text-red-500">{field.state.meta.errors.join(', ')}</p>
                    ) : null}
                     <p className="text-xs text-gray-500">Must be at least 8 characters long</p>
                </div>
            )}
          </form.Field>

          <form.Subscribe
            selector={(state) => [state.canSubmit, state.isSubmitting]}
          >
            {([canSubmit, isSubmitting]) => (
                <Button type="submit" className="w-full" disabled={!canSubmit || isSubmitting || register.isPending}>
                    {isSubmitting || register.isPending ? 'Creating Account...' : 'Sign Up'}
                </Button>
            )}
          </form.Subscribe>
        </form>

        <div className="mt-6 text-center text-sm text-gray-500">
          Already have an account?{' '}
          <Link to="/login" className="text-primary hover:underline font-medium">
            Login
          </Link>
        </div>
      </div>
    </div>
  )
}
