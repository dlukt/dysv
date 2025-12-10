import { useState, useId } from 'react'
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

export function AddressForm({ initialData, onSubmit, onCancel, isSubmitting }: AddressFormProps) {
  const [formData, setFormData] = useState<AddressFormData>({
    label: initialData?.label || '',
    line1: initialData?.line1 || '',
    line2: initialData?.line2 || '',
    city: initialData?.city || '',
    postalCode: initialData?.postalCode || '',
    state: initialData?.state || '',
    country: initialData?.country || 'DE',
    isDefault: initialData?.isDefault || false,
  })
  
  const id = useId()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
        AddressSchema.parse(formData)
        await onSubmit(formData)
    } catch (err) {
        console.error(err)
        alert('Please check your input')
    }
  }

  const handleChange = (field: keyof AddressFormData, value: string | boolean) => {
    setFormData(prev => ({ ...prev, [field]: value }))
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="space-y-2">
            <Label htmlFor={`${id}-label`}>Label (e.g. Home)</Label>
            <Input id={`${id}-label`} value={formData.label} onChange={e => handleChange('label', e.target.value)} required />
        </div>
        <div className="space-y-2">
             <Label htmlFor={`${id}-country`}>Country</Label>
             <select 
                id={`${id}-country`} 
                className="flex h-10 w-full rounded-md border border-slate-700 bg-slate-900 px-3 py-2 text-sm text-white focus:outline-none focus:ring-2 focus:ring-cyan-500"
                value={formData.country}
                onChange={e => handleChange('country', e.target.value)}
             >
                {countries.map(c => (
                    <option key={c.code} value={c.code}>{c.name}</option>
                ))}
             </select>
        </div>
      </div>

      <div className="space-y-2">
        <Label htmlFor={`${id}-line1`}>Address Line 1</Label>
        <Input id={`${id}-line1`} value={formData.line1} onChange={e => handleChange('line1', e.target.value)} required />
      </div>

      <div className="space-y-2">
        <Label htmlFor={`${id}-line2`}>Address Line 2 (Optional)</Label>
        <Input id={`${id}-line2`} value={formData.line2} onChange={e => handleChange('line2', e.target.value)} />
      </div>

      <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
        <div className="space-y-2">
             <Label htmlFor={`${id}-postalCode`}>Postal Code</Label>
             <Input id={`${id}-postalCode`} value={formData.postalCode} onChange={e => handleChange('postalCode', e.target.value)} required />
        </div>
         <div className="space-y-2">
             <Label htmlFor={`${id}-city`}>City</Label>
             <Input id={`${id}-city`} value={formData.city} onChange={e => handleChange('city', e.target.value)} required />
        </div>
         <div className="space-y-2 col-span-2 md:col-span-1">
             <Label htmlFor={`${id}-state`}>State/Province</Label>
             <Input id={`${id}-state`} value={formData.state || ''} onChange={e => handleChange('state', e.target.value)} />
        </div>
      </div>
      
       <div className="flex items-center space-x-2 pt-2">
          <input 
            type="checkbox" 
            id={`${id}-isDefault`} 
            checked={formData.isDefault} 
            onChange={e => handleChange('isDefault', e.target.checked)}
            className="h-4 w-4 rounded border-slate-700 bg-slate-900 text-cyan-500 focus:ring-cyan-500"
          />
          <Label htmlFor={`${id}-isDefault`}>Set as default address</Label>
      </div>

      <div className="flex justify-end gap-3 pt-4">
        <Button type="button" variant="ghost" onClick={onCancel} disabled={isSubmitting}>Cancel</Button>
        <Button type="submit" disabled={isSubmitting}>
            {isSubmitting ? 'Saving...' : 'Save Address'}
        </Button>
      </div>
    </form>
  )
}
