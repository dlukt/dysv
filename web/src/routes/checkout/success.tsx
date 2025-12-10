import { createFileRoute, Link } from '@tanstack/react-router'
import { CheckCircle } from 'lucide-react'
import { z } from 'zod'
import { useEffect } from 'react'

import { Button } from '@/components/ui/button'
import { clearCart } from '@/lib/cart-store'

const searchSchema = z.object({
  session_id: z.string().optional(),
})

export const Route = createFileRoute('/checkout/success')({
  component: CheckoutSuccessPage,
  validateSearch: (search) => searchSchema.parse(search),
})

function CheckoutSuccessPage() {
  const { session_id } = Route.useSearch()

  // Clear cart on successful checkout
  useEffect(() => {
    clearCart()
    // Also clear the backend session ID to ensure a fresh cart for next time
    localStorage.removeItem('dysv_session_id')
  }, [])

  return (
    <div className="min-h-screen bg-slate-900 flex items-center justify-center p-4">
      <div className="bg-slate-800 border border-slate-700 rounded-xl p-8 max-w-md w-full text-center shadow-2xl">
        <div className="mb-6 flex justify-center">
             <div className="h-24 w-24 bg-cyan-500/10 rounded-full flex items-center justify-center">
                <CheckCircle className="w-12 h-12 text-cyan-400" />
            </div>
        </div>
        
        <h1 className="text-3xl font-bold text-white mb-2">
          Payment Successful!
        </h1>
        <p className="text-slate-400 mb-6">
          Thank you for your order. Your account has been updated.
        </p>

        {session_id && (
            <div className="bg-slate-900/50 rounded-lg p-4 mb-8 text-sm">
                <span className="text-slate-500 block mb-1">Order Reference</span>
                <code className="text-cyan-400 font-mono break-all">{session_id.slice(-8).toUpperCase()}</code>
            </div>
        )}

        <div className="space-y-3">
             <Link to="/" className="block">
                <Button className="w-full bg-cyan-500 hover:bg-cyan-600 text-white font-semibold h-12">
                    Return to Dashboard
                </Button>
            </Link>
             <Link to="/cart" className="block">
                <Button variant="ghost" className="w-full text-slate-400 hover:text-white">
                    Back to Shop
                </Button>
            </Link>
        </div>
      </div>
    </div>
  )
}
