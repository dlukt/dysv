import { useState, useEffect } from 'react'
import { createFileRoute, Link, useNavigate } from '@tanstack/react-router'
import { useStore } from '@tanstack/react-store'
import { ArrowLeft, Trash2, CreditCard } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { useAuth } from '@/hooks/use-auth'

import { Button } from '@/components/ui/button'
import { plans, addons } from '@/data/pricing-data'
import {
  cartStore,
  hasItems,
  updateQuantity,
  removeFromCart,
} from '@/lib/cart-store'

export const Route = createFileRoute('/cart')({
  component: CartPage,
})

export function CartPage() {
  const cart = useStore(cartStore)
  const isEmpty = !hasItems(cart)
  const isYearly = cart.billingCycle === 'yearly'
  const { user } = useAuth()
  const navigate = useNavigate()

  const { t, i18n } = useTranslation()
  const locale = i18n.language
  const pricingPath = locale === 'de' ? '/pricing' : `/${locale}/pricing`

  // Derived state for display
  const cartItems = cart.items.map((item) => {
    let details: { name: string; monthlyPrice: number; limits?: unknown } | undefined
    if (item.type === 'plan') {
      details = plans.find((p) => p.id === item.id)
    } else {
      details = addons.find((a) => a.id === item.id)
    }
    return {
      ...item,
      name: details?.name || item.id,
      description: item.type === 'plan' ? (details as { limits: string } | undefined)?.limits : t('cart.domain_reg'),
      unitPrice: details?.monthlyPrice || 0,
    }
  })

  // Calculate totals
  const monthlyTotal = cartItems.reduce(
    (sum, item) => sum + item.unitPrice * item.quantity,
    0
  )

  // Yearly: Plans get 2 months free (pay 10), Addons pay 12
  const yearlyTotal = cartItems.reduce((sum, item) => {
    if (item.type === 'plan') {
      return sum + item.unitPrice * 10 * item.quantity
    }
    return sum + item.unitPrice * 12 * item.quantity
  }, 0)

  const getSessionId = () => {
    let sessionId = localStorage.getItem('dysv_session_id')
    if (!sessionId) {
      sessionId = crypto.randomUUID()
      localStorage.setItem('dysv_session_id', sessionId)
    }
    return sessionId
  }

  // Sync cart: Local is Truth. Make backend match local.
  const syncCartToBackend = async (sessionId: string) => {
    const headers = {
      'Content-Type': 'application/json',
      'X-Session-ID': sessionId,
    }

    const ensureOk = async (response: Response, context: string) => {
      if (response.ok) return
      let detail = ''
      try {
        const body = await response.json()
        detail = body.error
      } catch {} // eslint-disable-line no-empty
      const message = detail || response.statusText || 'unexpected error'
      throw new Error(`${context}: ${message}`)
    }

    // 1. Get backend state
    const cartResponse = await fetch('/api/cart', { headers })
    await ensureOk(cartResponse, 'Failed to load cart')
    const data = await cartResponse.json()
    const backendItems: { itemId: string; quantity: number }[] =
      data.cart?.items ?? []

    // 2. Set Billing Cycle
    await ensureOk(
      await fetch('/api/cart/billing-cycle', {
        method: 'POST',
        headers,
        body: JSON.stringify({ billingCycle: cart.billingCycle }),
      }),
      'Failed to set billing cycle'
    )

    // 3. Sync Items
    // Track what we've handled on backend to know what to delete
    const handledBackendItemIds = new Set<string>()

    for (const localItem of cart.items) {
      const backendItem = backendItems.find((bi) => bi.itemId === localItem.id)

      if (backendItem) {
        handledBackendItemIds.add(backendItem.itemId)
        // If quantity mismatch, update
        if (backendItem.quantity !== localItem.quantity) {
          await ensureOk(
            await fetch(`/api/cart/item/${localItem.id}`, {
              method: 'PUT',
              headers,
              body: JSON.stringify({ quantity: localItem.quantity }),
            }),
            `Failed to update quantity for ${localItem.id}`
          )
        }
      } else {
        // Not in backend, add it
        const endpoint =
          localItem.type === 'plan' ? '/api/cart/plan' : '/api/cart/addon'
        const bodyKey = localItem.type === 'plan' ? 'planId' : 'addonId'
        // For 'add', simplistic API might add to existing?
        // Our 'AddPlan' adds quantity. If we just call it with full quantity, it adds that much more?
        // No, AddPlan increments. But here we know it's 0 on backend.
        // So calling AddPlan(qty) is correct.
        await ensureOk(
          await fetch(endpoint, {
            method: 'POST',
            headers,
            body: JSON.stringify({
              [bodyKey]: localItem.id,
              quantity: localItem.quantity,
            }),
          }),
          `Failed to add ${localItem.id}`
        )
      }
    }

    // 4. Remove backend items not in local
    for (const backendItem of backendItems) {
      if (!handledBackendItemIds.has(backendItem.itemId)) {
        await ensureOk(
          await fetch(`/api/cart/item/${backendItem.itemId}`, {
            method: 'DELETE',
            headers,
          }),
          `Failed to remove ${backendItem.itemId}`
        )
      }
    }
  }

  const handleCheckout = async () => {
    if (!user) {
      navigate({ to: '/login', search: { redirect: '/cart' } })
      return
    }

    try {
      const sessionId = getSessionId()
      await syncCartToBackend(sessionId)

      // Navigate to Checkout Page for Address Selection
      navigate({ to: '/checkout' })
      
    } catch (err: unknown) {
      const message =
        err instanceof Error ? err.message : 'Checkout failed. Please try again.'
      console.error('Checkout error:', err)
      alert(message)
    }
  }

  const [mounted, setMounted] = useState(false)
  useEffect(() => {
    setMounted(true)
  }, [])

  if (!mounted || isEmpty) {
    return (
      <div className="min-h-screen bg-gradient-to-b from-slate-900 via-slate-800 to-slate-900 flex items-center justify-center">
        <div className="text-center px-6">
          <h1 className="text-3xl font-bold text-white mb-4">
            {t('cart.empty.title')}
          </h1>
          <p className="text-slate-400 mb-8">{t('cart.empty.description')}</p>
          <Link to={pricingPath}>
            <Button className="bg-cyan-500 hover:bg-cyan-600 text-white">
              {t('cart.empty.button')}
            </Button>
          </Link>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gradient-to-b from-slate-900 via-slate-800 to-slate-900">
      <main className="px-6 max-w-4xl mx-auto">
        <div className="py-6">
          <Link
            to={pricingPath}
            className="inline-flex items-center gap-2 text-slate-400 hover:text-white transition-colors"
          >
            <ArrowLeft className="w-4 h-4" />
            {t('cart.back')}
          </Link>
        </div>
        <h1 className="text-3xl font-bold text-white mb-8">{t('cart.title')}</h1>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Cart Items List */}
          <div className="lg:col-span-2 space-y-4">
            {cartItems.map((item) => (
              <div
                key={item.id}
                className="bg-slate-800/50 backdrop-blur-sm border border-slate-700 rounded-xl p-6 flex flex-col sm:flex-row items-center justify-between gap-4"
              >
                <div className="flex-grow text-center sm:text-left">
                  <h3 className="text-lg font-semibold text-white">
                    {item.name}
                  </h3>
                  <p className="text-sm text-slate-400">{item.description}</p>
                  {isYearly && item.type === 'plan' && (
                    <p className="text-sm text-cyan-400 mt-1">
                      {t('cart.billed_yearly')}
                    </p>
                  )}
                </div>

                <div className="flex items-center gap-4">
                  {/* Quantity Controls */}
                  <div className="flex items-center gap-2 bg-slate-700/50 rounded-lg p-1">
                    <button
                      type="button"
                      onClick={() => updateQuantity(item.id, item.quantity - 1)}
                      className="w-8 h-8 flex items-center justify-center text-slate-400 hover:text-white transition-colors"
                      disabled={item.quantity <= 1}
                    >
                      -
                    </button>
                    <span className="w-8 text-center text-white font-medium">
                      {item.quantity}
                    </span>
                    <button
                      type="button"
                      onClick={() => updateQuantity(item.id, item.quantity + 1)}
                      className="w-8 h-8 flex items-center justify-center text-slate-400 hover:text-white transition-colors"
                    >
                      +
                    </button>
                  </div>

                  <span className="text-xl font-bold text-white min-w-[100px] text-right">
                    €{(item.unitPrice * item.quantity).toFixed(2)}/mo
                  </span>

                  <Button
                    variant="ghost"
                    size="icon-sm"
                    onClick={() => removeFromCart(item.id)}
                    className="text-slate-400 hover:text-red-400"
                  >
                    <Trash2 className="w-4 h-4" />
                  </Button>
                </div>
              </div>
            ))}
          </div>

          {/* Order Summary */}
          <div className="lg:col-span-1">
            <div className="bg-slate-800/50 backdrop-blur-sm border border-slate-700 rounded-xl p-6 sticky top-6">
              <h2 className="text-lg font-semibold text-white mb-4">
                {t('cart.summary.title')}
              </h2>

              <div className="space-y-3 mb-6">
                <div className="flex justify-between text-slate-400">
                  <span>{t('cart.summary.subtotal')}</span>
                  <span>€{monthlyTotal.toFixed(2)}/mo</span>
                </div>
                {isYearly && (
                  <div className="flex justify-between text-cyan-400 text-sm">
                    <span>{t('cart.summary.yearly_discount')}</span>
                    <span>
                      -€
                      {(monthlyTotal * 12 - yearlyTotal).toFixed(2)}
                      /yrSaved
                    </span>
                  </div>
                )}
                <div className="border-t border-slate-700 pt-3 flex justify-between text-white font-semibold">
                  <span>{t('cart.summary.total')}</span>
                  <div className="text-right">
                    <p>€{monthlyTotal.toFixed(2)}/mo</p>
                    {isYearly && (
                      <p className="text-sm text-slate-400">
                        €{yearlyTotal.toFixed(2)}/year
                      </p>
                    )}
                  </div>
                </div>
              </div>

              <Button
                size="lg"
                onClick={handleCheckout}
                className="w-full bg-cyan-500 hover:bg-cyan-600 text-white shadow-lg shadow-cyan-500/30"
              >
                <CreditCard className="w-5 h-5 mr-2" />
                {t('cart.summary.checkout')}
              </Button>

              <p className="text-xs text-slate-500 text-center mt-4">
                {t('cart.summary.redirect')}
              </p>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
