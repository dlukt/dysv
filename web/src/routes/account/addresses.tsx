import { createFileRoute } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import { Plus, Edit2, Trash2, MapPin, Check } from 'lucide-react'
import { useTranslation } from 'react-i18next'

import { Button } from '@/components/ui/button'
import { type AddressFormData, AddressForm } from '@/components/AddressForm'

export const Route = createFileRoute('/account/addresses')({
  component: AddressManagementPage,
})

// Types mirroring backend response
interface Address {
  id: string
  label: string
  line1: string
  line2?: string
  city: string
  postalCode: string
  state?: string
  country: string
  isDefault: boolean
}

function AddressManagementPage() {
  const { t } = useTranslation()
  const queryClient = useQueryClient()
  
  const [isEditing, setIsEditing] = useState(false)
  const [editingAddress, setEditingAddress] = useState<Address | null>(null)

  // Fetch Addresses
  const { data: addresses, isLoading } = useQuery({
    queryKey: ['addresses'],
    queryFn: async () => {
      const token = localStorage.getItem('dysv_auth_token')
      const res = await fetch('/api/user/addresses', {
        headers: { Authorization: `Bearer ${token}` }
      })
      if (!res.ok) throw new Error('Failed to fetch addresses')
      const data = await res.json()
      return (data.addresses || []) as Address[]
    },
  })

  // Mutations
  const createMutation = useMutation({
    mutationFn: async (data: AddressFormData) => {
      const token = localStorage.getItem('dysv_auth_token')
      const res = await fetch('/api/user/addresses', {
        method: 'POST',
        headers: { 
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`
        },
        body: JSON.stringify(data),
      })
      if (!res.ok) throw new Error('Failed to create address')
      return res.json()
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['addresses'] })
      setIsEditing(false)
      // toast success
    },
  })

  const updateMutation = useMutation({
    mutationFn: async ({ id, data }: { id: string; data: AddressFormData }) => {
      const token = localStorage.getItem('dysv_auth_token')
      const res = await fetch(`/api/user/addresses/${id}`, {
        method: 'PUT',
        headers: { 
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`
        },
        body: JSON.stringify(data),
      })
      if (!res.ok) throw new Error('Failed to update address')
      return res.json()
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['addresses'] })
      setIsEditing(false)
      setEditingAddress(null)
    },
  })

  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      const token = localStorage.getItem('dysv_auth_token')
      const res = await fetch(`/api/user/addresses/${id}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${token}` }
      })
      if (!res.ok) throw new Error('Failed to delete address')
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['addresses'] })
    },
  })

  const handleSubmit = async (data: AddressFormData) => {
    if (editingAddress) {
      await updateMutation.mutateAsync({ id: editingAddress.id, data })
    } else {
      await createMutation.mutateAsync(data)
    }
  }

  const handleEdit = (address: Address) => {
    setEditingAddress(address)
    setIsEditing(true)
  }

  const handleDelete = async (id: string) => {
    if (window.confirm(t('address.delete_confirm'))) {
        await deleteMutation.mutateAsync(id)
    }
  }

  const handleAddNew = () => {
    setEditingAddress(null)
    setIsEditing(true)
  }

  const handleCancel = () => {
    setIsEditing(false)
    setEditingAddress(null)
  }

  if (isEditing) {
    return (
      <div className="min-h-screen bg-gradient-to-b from-slate-900 via-slate-800 to-slate-900 p-6 flex items-center">
        <div className="max-w-2xl mx-auto w-full">
          <h1 className="text-2xl font-bold text-white mb-6">
              {editingAddress ? t('address.form.title_edit') : t('address.form.title_new')}
          </h1>
          <div className="bg-slate-800 rounded-lg p-6 border border-slate-700">
             <AddressForm 
               initialData={editingAddress || undefined}
               onSubmit={handleSubmit}
               onCancel={handleCancel}
               isSubmitting={createMutation.isPending || updateMutation.isPending}
             />
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gradient-to-b from-slate-900 via-slate-800 to-slate-900 p-6">
      <div className="max-w-5xl mx-auto">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold text-white">{t('address.title')}</h1>
        <Button onClick={handleAddNew} className="bg-cyan-500 hover:bg-cyan-600 text-white">
            <Plus className="w-4 h-4 mr-2" />
            {t('address.add_new')}
        </Button>
      </div>

      {isLoading ? (
        <div className="text-slate-400 text-center py-12">Loading...</div>
      ) : addresses?.length === 0 ? (
        <div className="text-center py-12 bg-slate-800/50 rounded-xl border border-slate-700/50">
            <MapPin className="w-12 h-12 text-slate-600 mx-auto mb-4" />
            <p className="text-slate-400">{t('address.empty')}</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {addresses?.map((address) => (
            <div key={address.id} className="bg-slate-800 border border-slate-700 rounded-xl p-6 relative group">
                <div className="flex justify-between items-start mb-4">
                    <div>
                        <span className="font-semibold text-white block text-lg">{address.label}</span>
                        {address.isDefault && (
                            <span className="inline-flex items-center text-xs text-cyan-400 bg-cyan-400/10 px-2 py-0.5 rounded-full mt-1">
                                <Check className="w-3 h-3 mr-1" />
                                {t('address.default_badge')}
                            </span>
                        )}
                    </div>
                    <div className="flex gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                         <Button variant="ghost" size="icon" className="h-8 w-8 text-slate-400 hover:text-white" onClick={() => handleEdit(address)}>
                             <Edit2 className="w-4 h-4" />
                         </Button>
                         <Button variant="ghost" size="icon" className="h-8 w-8 text-red-400 hover:text-red-300 hover:bg-red-400/10" onClick={() => handleDelete(address.id)}>
                             <Trash2 className="w-4 h-4" />
                         </Button>
                    </div>
                </div>
                
                <div className="text-slate-400 text-sm space-y-1">
                    <p>{address.line1}</p>
                    {address.line2 && <p>{address.line2}</p>}
                    <p>{address.postalCode} {address.city}</p>
                    <p>{address.state}</p>
                    <p>{address.country}</p>
                </div>
            </div>
          ))}
        </div>
      )}
      </div>
    </div>
  )
}
