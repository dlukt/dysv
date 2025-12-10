import { createFileRoute, Link } from '@tanstack/react-router'
import { CheckCircle } from 'lucide-react'

import { Button } from '@/components/ui/button'
import { clearCart } from '@/lib/cart-store'
import { useEffect } from 'react'

export const Route = createFileRoute('/checkout/success')({
  component: CheckoutSuccessPage,
})

function CheckoutSuccessPage() {
  // Clear cart on successful checkout
  useEffect(() => {
    clearCart()
  }, [])

  return (
    <div className="min-h-screen bg-gradient-to-b from-slate-900 via-slate-800 to-slate-900 flex items-center justify-center">
      <div className="text-center px-6 max-w-md">
        <CheckCircle className="w-20 h-20 text-cyan-400 mx-auto mb-6" />
        <h1 className="text-3xl font-bold text-white mb-4">
          Payment Successful!
        </h1>
        <p className="text-slate-400 mb-8">
          Thank you for your order. You'll receive a confirmation email shortly
          with your account details and next steps.
        </p>
        <Link to="/">
          <Button className="bg-cyan-500 hover:bg-cyan-600 text-white">
            Back to Home
          </Button>
        </Link>
      </div>
    </div>
  )
}
