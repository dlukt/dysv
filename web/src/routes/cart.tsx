import { createFileRoute, Link } from '@tanstack/react-router'
import { useStore } from '@tanstack/react-store'
import { ArrowLeft, Trash2, CreditCard } from 'lucide-react'
import { useTranslation } from '@/lib/i18n'

import { Button } from '@/components/ui/button'
import { plans, addons } from '@/data/pricing-data'
import {
  cartStore,
  clearPlan,
  removeAddon,
  hasItems,
} from '@/lib/cart-store'

export const Route = createFileRoute('/cart')({
  component: CartPage,
})

export function CartPage() {
  const cart = useStore(cartStore)
  const isEmpty = !hasItems(cart)
  const isYearly = cart.billingCycle === 'yearly'

  const { t, locale } = useTranslation()
  const pricingPath = locale === 'de' ? '/pricing' : `/${locale}/pricing`

  const selectedPlan = plans.find((p) => p.id === cart.planId)
  const selectedAddons = addons.filter((a) => cart.addons.includes(a.id))

  // Calculate totals - always use raw monthly prices for display
  const planMonthlyPrice = selectedPlan?.monthlyPrice ?? 0
  const addonsMonthlyPrice = selectedAddons.reduce((sum, addon) => sum + addon.monthlyPrice, 0)

  // Monthly display: raw monthly price
  const monthlyTotal = planMonthlyPrice + addonsMonthlyPrice

  // Yearly display: plan gets 2 months free (10 months), addons pay full 12 months
  const planYearlyPrice = planMonthlyPrice * 10 // 2 months free
  const addonsYearlyPrice = addonsMonthlyPrice * 12 // No discount
  const yearlyTotal = planYearlyPrice + addonsYearlyPrice

  // Helper to get or create session ID
  const getSessionId = () => {
    let sessionId = localStorage.getItem('dysv_session_id')
    if (!sessionId) {
      sessionId = crypto.randomUUID()
      localStorage.setItem('dysv_session_id', sessionId)
    }
    return sessionId
  }

  // Sync cart to backend (handles additions AND removals)
  const syncCartToBackend = async (sessionId: string) => {
    // ... existing sync logic ...
    const headers = {
      'Content-Type': 'application/json',
      'X-Session-ID': sessionId,
    }

    const ensureOk = async (response: Response, context: string) => {
      if (response.ok) return

      let detail = ''
      try {
        const body = await response.json()
        if (body && typeof body.error === 'string') {
          detail = body.error
        }
      } catch {
        // Ignore JSON parse errors; fall back to status text
      }

      const message = detail || response.statusText || 'unexpected error'
      throw new Error(`${context}: ${message}`)
    }

    // First, get current backend cart to find items to remove
    const cartResponse = await fetch('/api/cart', { headers })
    await ensureOk(cartResponse, 'Failed to load cart')

    let backendAddons: string[] = []
    let backendPlans: string[] = []
    const data = await cartResponse.json()
    const backendItems = data.cart?.items ?? []
    backendAddons = backendItems
      .filter((item: { itemType: string }) => item.itemType === 'addon')
      .map((item: { itemId: string }) => item.itemId)
    backendPlans = backendItems
      .filter((item: { itemType: string }) => item.itemType === 'plan')
      .map((item: { itemId: string }) => item.itemId)

    // Remove addons that are on backend but not in local cart
    for (const addonId of backendAddons) {
      if (!cart.addons.includes(addonId)) {
        const resp = await fetch(`/api/cart/item/${addonId}`, {
          method: 'DELETE',
          headers,
        })
        await ensureOk(resp, 'Failed to remove addon from cart')
      }
    }

    // Remove backend plan if it was cleared locally
    if (!cart.planId) {
      for (const planId of backendPlans) {
        const resp = await fetch(`/api/cart/item/${planId}`, {
          method: 'DELETE',
          headers,
        })
        await ensureOk(resp, 'Failed to remove plan from cart')
      }
    }

    // Set billing cycle
    const billingResp = await fetch('/api/cart/billing-cycle', {
      method: 'POST',
      headers,
      body: JSON.stringify({ billingCycle: cart.billingCycle }),
    })
    await ensureOk(billingResp, 'Failed to set billing cycle')

    // Set plan if selected (this replaces any existing plan)
    if (cart.planId) {
      const planResp = await fetch('/api/cart/plan', {
        method: 'POST',
        headers,
        body: JSON.stringify({ planId: cart.planId }),
      })
      await ensureOk(planResp, 'Failed to set plan')
    }

    // Add addons that are in local cart but not on backend
    for (const addonId of cart.addons) {
      if (!backendAddons.includes(addonId)) {
        const resp = await fetch('/api/cart/addon', {
          method: 'POST',
          headers,
          body: JSON.stringify({ addonId }),
        })
        await ensureOk(resp, 'Failed to add addon')
      }
    }
  }

  const handleCheckout = async () => {
    try {
      const sessionId = getSessionId()

      // Sync local cart to backend before checkout
      await syncCartToBackend(sessionId)

      const response = await fetch('/api/checkout', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Session-ID': sessionId,
        },
      })

      if (!response.ok) {
        const error = await response.json()
        alert(`Checkout failed: ${error.error || 'Unknown error'}`)
        return
      }

      const data = await response.json()
      if (data.url) {
        window.location.href = data.url
      }
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Checkout failed. Please try again.'
      console.error('Checkout error:', err)
      alert(message)
    }
  }

  if (isEmpty) {
    return (
      <div className="min-h-screen bg-gradient-to-b from-slate-900 via-slate-800 to-slate-900 flex items-center justify-center">
        <div className="text-center px-6">
          <h1 className="text-3xl font-bold text-white mb-4">{t.cart.empty.title}</h1>
          <p className="text-slate-400 mb-8">{t.cart.empty.description}</p>
          <Link to={pricingPath}>
            <Button className="bg-cyan-500 hover:bg-cyan-600 text-white">
              {t.cart.empty.button}
            </Button>
          </Link>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gradient-to-b from-slate-900 via-slate-800 to-slate-900">
      {/* Cart Content */}
      <main className="px-6 max-w-4xl mx-auto">
        <div className="py-6">
          <Link to={pricingPath} className="inline-flex items-center gap-2 text-slate-400 hover:text-white transition-colors">
            <ArrowLeft className="w-4 h-4" />
            {t.cart.back}
          </Link>
        </div>
        <h1 className="text-3xl font-bold text-white mb-8">{t.cart.title}</h1>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Cart Items */}
          <div className="lg:col-span-2 space-y-4">
            {/* Plan */}
            {selectedPlan && (
              <div className="bg-slate-800/50 backdrop-blur-sm border border-slate-700 rounded-xl p-6 flex items-center justify-between">
                <div>
                  <h3 className="text-lg font-semibold text-white">{selectedPlan.name}</h3>
                  <p className="text-sm text-slate-400">{selectedPlan.limits}</p>
                  <p className="text-sm text-cyan-400 mt-1">
                    {isYearly ? t.cart.billed_yearly : t.cart.billed_monthly}
                  </p>
                </div>
                <div className="flex items-center gap-4">
                  <span className="text-xl font-bold text-white">
                    €{planMonthlyPrice.toFixed(2)}/mo
                  </span>
                  <Button
                    variant="ghost"
                    size="icon-sm"
                    onClick={clearPlan}
                    className="text-slate-400 hover:text-red-400"
                  >
                    <Trash2 className="w-4 h-4" />
                  </Button>
                </div>
              </div>
            )}

            {/* Addons */}
            {selectedAddons.map((addon) => (
              <div
                key={addon.id}
                className="bg-slate-800/50 backdrop-blur-sm border border-slate-700 rounded-xl p-6 flex items-center justify-between"
              >
                <div>
                  <h3 className="text-lg font-semibold text-white">{addon.name}</h3>
                  <p className="text-sm text-slate-400">{t.cart.domain_reg}</p>
                </div>
                <div className="flex items-center gap-4">
                  <span className="text-xl font-bold text-white">
                    €{addon.monthlyPrice.toFixed(2)}/mo
                  </span>
                  <Button
                    variant="ghost"
                    size="icon-sm"
                    onClick={() => removeAddon(addon.id)}
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
              <h2 className="text-lg font-semibold text-white mb-4">{t.cart.summary.title}</h2>

              <div className="space-y-3 mb-6">
                <div className="flex justify-between text-slate-400">
                  <span>{t.cart.summary.subtotal}</span>
                  <span>€{monthlyTotal.toFixed(2)}/mo</span>
                </div>
                {isYearly && planMonthlyPrice > 0 && (
                  <div className="flex justify-between text-cyan-400 text-sm">
                    <span>{t.cart.summary.yearly_discount}</span>
                    <span>-€{(planMonthlyPrice * 2).toFixed(2)}</span>
                  </div>
                )}
                <div className="border-t border-slate-700 pt-3 flex justify-between text-white font-semibold">
                  <span>{t.cart.summary.total}</span>
                  <div className="text-right">
                    <p>€{monthlyTotal.toFixed(2)}/mo</p>
                    {isYearly && (
                      <p className="text-sm text-slate-400">€{yearlyTotal.toFixed(2)}/year</p> // /year is hardcoded?
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
                {t.cart.summary.checkout}
              </Button>

              <p className="text-xs text-slate-500 text-center mt-4">
                {t.cart.summary.redirect}
              </p>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
