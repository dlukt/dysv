import { useForm } from '@tanstack/react-form'
import { useId } from 'react'
import { z } from 'zod'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

// ISO 3166-1 alpha-2 codes (subset for brevity, can enable full list later)
const countries = [
  { code: 'DE', name: 'Germany' },
  { code: 'US', name: 'United States' },
  { code: 'GB', name: 'United Kingdom' },
  { code: 'FR', name: 'France' },
  { code: 'ES', name: 'Spain' },
  { code: 'IT', name: 'Italy' },
  { code: 'HR', name: 'Croatia' },
  { code: 'AT', name: 'Austria' },
  { code: 'CH', name: 'Switzerland' },
  { code: 'NL', name: 'Netherlands' },
]

export const AddressSchema = z.object({
  label: z.string().min(1, 'Label is required'),
  line1: z.string().min(1, 'Address line 1 is required'),
  line2: z.string().optional(),
  city: z.string().min(1, 'City is required'),
  postalCode: z.string().min(1, 'Postal code is required'),
  state: z.string().optional(),
  country: z.string().length(2, 'Select a country'),
  isDefault: z.boolean().default(false),
})

export type AddressFormData = z.infer<typeof AddressSchema>

interface AddressFormProps {
  initialData?: Partial<AddressFormData>
  onSubmit: (data: AddressFormData) => Promise<void>
  onCancel: () => void
  isSubmitting?: boolean
}

// Simple helper to bridge Zod and TanStack Form
// definition: ({ value }) => string | undefined
function validateWith<T>(schema: z.ZodType<T>) {
  return ({ value }: { value: T }) => {
    const result = schema.safeParse(value)
    if (result.success) return undefined
    return result.error.issues[0].message
  }
}

export function AddressForm({ initialData, onSubmit, onCancel, isSubmitting }: AddressFormProps) {
  const id = useId()
  
  const form = useForm({
    defaultValues: {
      label: initialData?.label || '',
      line1: initialData?.line1 || '',
      line2: initialData?.line2 || '',
      city: initialData?.city || '',
      postalCode: initialData?.postalCode || '',
      state: initialData?.state || '',
      country: initialData?.country || 'DE',
      isDefault: initialData?.isDefault || false,
    },
    onSubmit: async ({ value }) => {
       await onSubmit(value)
    },
  })

  return (
    <form
      onSubmit={(e) => {
        e.preventDefault()
        e.stopPropagation()
        form.handleSubmit()
      }}
      className="space-y-4"
    >
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <form.Field
          name="label"
          validators={{
            onChange: validateWith(AddressSchema.shape.label)
          }}
        >
          {(field) => (
            <div className="space-y-2">
                <Label htmlFor={`${id}-label`}>Label (e.g. Home)</Label>
                <Input 
                    id={`${id}-label`} 
                    name={field.name}
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
            name="country"
            validators={{
                onChange: validateWith(AddressSchema.shape.country)
            }}
        >
            {(field) => (
                <div className="space-y-2">
                    <Label htmlFor={`${id}-country`}>Country</Label>
                    <select 
                        id={`${id}-country`}
                        name={field.name}
                        className="flex h-10 w-full rounded-md border border-slate-700 bg-slate-900 px-3 py-2 text-sm text-white focus:outline-none focus:ring-2 focus:ring-cyan-500"
                        value={field.state.value}
                        onBlur={field.handleBlur}
                        onChange={(e) => field.handleChange(e.target.value)}
                    >
                        {countries.map(c => (
                            <option key={c.code} value={c.code}>{c.name}</option>
                        ))}
                    </select>
                    {field.state.meta.errors ? (
                        <p className="text-xs text-red-500">{field.state.meta.errors.join(', ')}</p>
                    ) : null}
                </div>
            )}
        </form.Field>
      </div>

      <form.Field
        name="line1"
        validators={{
            onChange: validateWith(AddressSchema.shape.line1)
        }}
      >
        {(field) => (
            <div className="space-y-2">
                <Label htmlFor={`${id}-line1`}>Address Line 1</Label>
                <Input 
                    id={`${id}-line1`} 
                    name={field.name}
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
        name="line2"
        validators={{
            onChange: validateWith(AddressSchema.shape.line2)
        }}
      >
        {(field) => (
             <div className="space-y-2">
                <Label htmlFor={`${id}-line2`}>Address Line 2 (Optional)</Label>
                <Input 
                    id={`${id}-line2`} 
                    name={field.name}
                    value={field.state.value || ''} 
                    onBlur={field.handleBlur}
                    onChange={(e) => field.handleChange(e.target.value)} 
                />
                 {field.state.meta.errors ? (
                    <p className="text-xs text-red-500">{field.state.meta.errors.join(', ')}</p>
                ) : null}
            </div>
        )}
      </form.Field>


      <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
        <form.Field
            name="postalCode"
            validators={{ onChange: validateWith(AddressSchema.shape.postalCode) }}
        >
            {(field) => (
                <div className="space-y-2">
                    <Label htmlFor={`${id}-postalCode`}>Postal Code</Label>
                    <Input 
                        id={`${id}-postalCode`} 
                        name={field.name}
                        value={field.state.value} 
                        onBlur={field.handleBlur}
                        onChange={e => field.handleChange(e.target.value)} 
                        required 
                    />
                     {field.state.meta.errors ? (
                        <p className="text-xs text-red-500">{field.state.meta.errors.join(', ')}</p>
                    ) : null}
                </div>
            )}
        </form.Field>
        
        <form.Field
            name="city"
            validators={{ onChange: validateWith(AddressSchema.shape.city) }}
        >
            {(field) => (
                 <div className="space-y-2">
                    <Label htmlFor={`${id}-city`}>City</Label>
                    <Input 
                        id={`${id}-city`} 
                        name={field.name}
                        value={field.state.value} 
                        onBlur={field.handleBlur}
                        onChange={e => field.handleChange(e.target.value)} 
                        required 
                    />
                     {field.state.meta.errors ? (
                        <p className="text-xs text-red-500">{field.state.meta.errors.join(', ')}</p>
                    ) : null}
                </div>
            )}
        </form.Field>
        
        <form.Field
            name="state"
            validators={{ onChange: validateWith(AddressSchema.shape.state) }}
        >
            {(field) => (
                <div className="space-y-2 col-span-2 md:col-span-1">
                    <Label htmlFor={`${id}-state`}>State/Province</Label>
                    <Input 
                        id={`${id}-state`} 
                        name={field.name}
                        value={field.state.value || ''} 
                        onBlur={field.handleBlur}
                        onChange={e => field.handleChange(e.target.value)} 
                    />
                     {field.state.meta.errors ? (
                        <p className="text-xs text-red-500">{field.state.meta.errors.join(', ')}</p>
                    ) : null}
                </div>
            )}
        </form.Field>
      </div>
      
      <form.Field
        name="isDefault"
        validators={{ onChange: validateWith(AddressSchema.shape.isDefault) }}
      >
        {(field) => (
            <div className="flex items-center space-x-2 pt-2">
                <input 
                    type="checkbox" 
                    id={`${id}-isDefault`} 
                    name={field.name}
                    checked={field.state.value} 
                    onBlur={field.handleBlur}
                    onChange={e => field.handleChange(e.target.checked)}
                    className="h-4 w-4 rounded border-slate-700 bg-slate-900 text-cyan-500 focus:ring-cyan-500"
                />
                <Label htmlFor={`${id}-isDefault`}>Set as default address</Label>
                 {field.state.meta.errors ? (
                    <p className="text-xs text-red-500">{field.state.meta.errors.join(', ')}</p>
                ) : null}
            </div>
        )}
      </form.Field>

      <div className="flex justify-end gap-3 pt-4">
        <Button type="button" variant="ghost" onClick={onCancel} disabled={isSubmitting}>Cancel</Button>
        <form.Subscribe
            selector={(state) => [state.canSubmit, state.isSubmitting]}
        >
            {([canSubmit, isSubmittingForm]) => (
                <Button type="submit" disabled={!canSubmit || isSubmitting || isSubmittingForm}>
                    {(isSubmitting || isSubmittingForm) ? 'Saving...' : 'Save Address'}
                </Button>
            )}
        </form.Subscribe>
      </div>
    </form>
  )
}
