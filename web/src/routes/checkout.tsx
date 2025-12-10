import { useState } from 'react'
import { createFileRoute, useNavigate, Link } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Loader2, Check } from 'lucide-react'
import { useAuth } from '@/hooks/use-auth'
import { Button } from '@/components/ui/button'
import { AddressForm, type AddressFormData } from '@/components/AddressForm'

export const Route = createFileRoute('/checkout')({
  component: CheckoutPage,
})

interface LineItem {
  itemId: string
  itemType: 'plan' | 'addon'
  name: string
  price: number
  quantity: number
}

interface Cart {
  items: LineItem[]
  billingCycle: 'monthly' | 'yearly'
  total: {
    monthly: number
    yearly: number
  }
}

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

function CheckoutPage() {
  const navigate = useNavigate()
  const { user, isLoading: isAuthLoading } = useAuth()
  const [selectedAddressId, setSelectedAddressId] = useState<string | null>(null)
  const [isAddingAddress, setIsAddingAddress] = useState(false)
  const queryClient = useQueryClient()

  // Redirect if not logged in
  if (!isAuthLoading && !user) {
    navigate({ to: '/login', search: { redirect: '/checkout' } })
  }

  // Fetch Cart
  const { data: cart, isLoading: isCartLoading } = useQuery({
    queryKey: ['cart'],
    queryFn: async () => {
      const sessionId = localStorage.getItem('dysv_session_id')
      if (!sessionId) {
          return { items: [], billingCycle: 'monthly', total: { monthly: 0, yearly: 0 } } as Cart
      }
      const res = await fetch('/api/cart', {
          headers: {
              'X-Session-ID': sessionId
          }
      })
      if (!res.ok) throw new Error('Failed to load cart')
      const data = await res.json()
      
      return {
          items: data.cart.items,
          billingCycle: data.cart.billingCycle,
          total: {
              monthly: data.monthlyTotal,
              yearly: data.yearlyTotal
          }
      } as Cart
    },
    staleTime: 1000 * 60 * 5, // 5 minutes
    retry: false
  })

  // Fetch Addresses
  const { data: addressData, isLoading: isAddressLoading } = useQuery({
    queryKey: ['addresses'],
    queryFn: async () => {
      const token = localStorage.getItem('dysv_auth_token')
      const res = await fetch('/api/user/addresses', {
          headers: { Authorization: `Bearer ${token}` }
      })
      if (!res.ok) throw new Error('Failed to load addresses')
      const data = await res.json()
      // Auto-select default address
      const addresses = data.addresses as Address[]
      const defaultAddr = addresses.find(a => a.isDefault)
      if (defaultAddr && !selectedAddressId) {
          setSelectedAddressId(defaultAddr.id)
      } else if (addresses.length > 0 && !selectedAddressId) {
          setSelectedAddressId(addresses[0].id)
      }
      return data
    },
    enabled: !!user,
    staleTime: 1000 * 60 * 5, // 5 minutes
    retry: 1
  })

  // Create Address Mutation
  const createAddressMutation = useMutation({
    mutationFn: async (data: AddressFormData) => {
       const token = localStorage.getItem('dysv_auth_token')
       const res = await fetch('/api/user/addresses', {
         method: 'POST',
         headers: { 
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}` 
         },
         body: JSON.stringify(data)
       })
       if (!res.ok) throw new Error('Failed to create address')
       return res.json()
    },
    onSuccess: (newAddr) => {
        queryClient.invalidateQueries({ queryKey: ['addresses'] })
        setIsAddingAddress(false)
        setSelectedAddressId(newAddr.id)
    }
  })

  // Checkout Mutation
  const checkoutMutation = useMutation({
    mutationFn: async () => {
      const token = localStorage.getItem('dysv_auth_token')
      const res = await fetch('/api/checkout', {
        method: 'POST',
        headers: { 
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`
        },
        body: JSON.stringify({ addressId: selectedAddressId }),
      })
      if (!res.ok) {
        const err = await res.json()
        throw new Error(err.error || 'Checkout failed')
      }
      return res.json() as Promise<{ url: string }>
    },
    onSuccess: (data) => {
      window.location.href = data.url
    },
  })

  const addresses = (addressData?.addresses as Address[]) || []
  const cartItems = cart?.items || []

  // Calculate total based on billing cycle
  const currentTotal = cart?.billingCycle === 'yearly' 
    ? cart.items.reduce((acc, item) => {
        const m = item.itemType === 'plan' ? 10 : 12
        return acc + (item.price * m * item.quantity)
      }, 0)
    : cart?.items.reduce((acc, item) => acc + (item.price * item.quantity), 0) ?? 0

  if (isAuthLoading || isCartLoading || isAddressLoading) {
    return <div className="min-h-screen flex items-center justify-center bg-slate-900 text-white"><Loader2 className="w-8 h-8 animate-spin" /></div>
  }

  return (
    <div className="min-h-screen bg-slate-900 text-slate-100 py-12 px-4">
      <div className="max-w-4xl mx-auto grid grid-cols-1 md:grid-cols-3 gap-8">
        
        {/* Left Column: Address Selection */}
        <div className="md:col-span-2 space-y-6">
          <h1 className="text-3xl font-bold bg-gradient-to-r from-cyan-400 to-blue-500 bg-clip-text text-transparent">
            Checkout
          </h1>
          
          <div className="bg-slate-800/50 border border-slate-700 rounded-xl p-6">
            <div className="flex justify-between items-center mb-4">
                <h2 className="text-xl font-semibold text-white">Billing Address</h2>
                <Button variant="ghost" size="sm" onClick={() => setIsAddingAddress(!isAddingAddress)} className="text-cyan-400 hover:text-cyan-300">
                    {isAddingAddress ? 'Cancel' : 'Add New'}
                </Button>
            </div>

            {isAddingAddress ? (
                <AddressForm 
                    onSubmit={async (data) => await createAddressMutation.mutateAsync(data)} 
                    onCancel={() => setIsAddingAddress(false)}
                    isSubmitting={createAddressMutation.isPending}
                />
            ) : (
                <div className="space-y-3">
                    {addresses.length === 0 ? (
                        <div className="text-center py-8 text-slate-400 border border-dashed border-slate-700 rounded-lg">
                            No addresses found. Please add one.
                        </div>
                    ) : (
                        addresses.map(addr => (
                            <button 
                                key={addr.id} 
                                type="button"
                                onClick={() => setSelectedAddressId(addr.id)}
                                className={`
                                    w-full text-left
                                    relative p-4 rounded-lg border cursor-pointer transition-all focus:outline-none focus:ring-2 focus:ring-cyan-500
                                    ${selectedAddressId === addr.id 
                                        ? 'border-cyan-500 bg-cyan-500/10' 
                                        : 'border-slate-700 hover:border-slate-600 bg-slate-800/30'}
                                `}
                            >
                                <div className="flex justify-between items-start">
                                    <div>
                                        <div className="font-semibold text-white flex items-center gap-2">
                                            {addr.label}
                                            {selectedAddressId === addr.id && <Check className="w-4 h-4 text-cyan-400" />}
                                        </div>
                                        <div className="text-sm text-slate-400 mt-1">
                                            {addr.line1}<br/>
                                            {addr.line2 && <>{addr.line2}<br/></>}
                                            {addr.postalCode} {addr.city}, {addr.country}
                                        </div>
                                    </div>
                                    {addr.isDefault && <span className="text-xs bg-slate-700 text-slate-300 px-2 py-1 rounded">Default</span>}
                                </div>
                            </button>
                        ))
                    )}
                </div>
            )}
          </div>
        </div>

        {/* Right Column: Order Summary */}
        <div className="space-y-6">
            <div className="bg-slate-800 border border-slate-700 rounded-xl p-6 sticky top-6">
                <h2 className="text-xl font-semibold text-white mb-4">Order Summary</h2>
                <div className="space-y-4 mb-6">
                    {cartItems.map(item => (
                        <div key={item.itemId} className="flex justify-between text-sm">
                            <span className="text-slate-300">{item.name} x {item.quantity}</span>
                            <span className="text-white font-medium">
                                €{cart?.billingCycle === 'yearly' 
                                    ? ((item.price * (item.itemType === 'plan' ? 10 : 12)) * item.quantity).toFixed(2)
                                    : (item.price * item.quantity).toFixed(2)
                                }
                            </span>
                        </div>
                    ))}
                    <div className="border-t border-slate-700 my-4"></div>
                    <div className="flex justify-between items-center text-lg font-bold text-white">
                        <span>Total ({cart?.billingCycle})</span>
                        <span>€{currentTotal.toFixed(2)}</span>
                    </div>
                </div>

                <Button 
                    onClick={() => checkoutMutation.mutate()} 
                    disabled={!selectedAddressId || checkoutMutation.isPending || isAddingAddress || cartItems.length === 0}
                    className="w-full bg-gradient-to-r from-cyan-500 to-blue-600 hover:from-cyan-400 hover:to-blue-500 text-white font-bold py-3 rounded-lg shadow-lg hover:shadow-cyan-500/25 transition-all"
                >
                    {checkoutMutation.isPending ? <Loader2 className="w-5 h-5 animate-spin" /> : 'Proceed to Payment'}
                </Button>
                
                <div className="text-center mt-4">
                    <Link to="/cart" className="text-sm text-slate-400 hover:text-white underline">
                        Back to Cart
                    </Link>
                </div>
            </div>
        </div>

      </div>
    </div>
  )
}
