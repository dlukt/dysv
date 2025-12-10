import { Link, createFileRoute } from '@tanstack/react-router'
import { useState, useEffect } from 'react'
import { useStore } from '@tanstack/react-store'
import { ShoppingCart } from 'lucide-react'
import { useTranslation } from '@/lib/i18n'

import { PricingCard } from '@/components/PricingCard'
import { Switch } from '@/components/ui/switch'
import { Button } from '@/components/ui/button'
import { plans } from '@/data/pricing-data'
import {
  cartStore,
  addToCart,
  setBillingCycle,
  removeFromCart,
  getItemCount,
} from '@/lib/cart-store'

export const Route = createFileRoute('/pricing')({
  component: PricingPage,
})

export function PricingPage() {
  const cart = useStore(cartStore)
  const itemCount = getItemCount(cart)
  const { t, locale } = useTranslation()
  const cartPath = locale === 'de' ? '/cart' : `/${locale}/cart`

  const handleBillingToggle = (checked: boolean) => {
    setBillingCycle(checked ? 'yearly' : 'monthly')
  }

  const handlePlanAdd = (planId: string) => {
    addToCart(planId, 'plan')
  }

  const hasDomain = cart.items.some((i) => i.id === 'de-domain')
  const handleDomainToggle = (checked: boolean) => {
    if (checked) {
      addToCart('de-domain', 'addon')
    } else {
      removeFromCart('de-domain')
    }
  }

  const [mounted, setMounted] = useState(false)
  useEffect(() => {
    setMounted(true)
  }, [])

  return (
    <div className="min-h-screen bg-gradient-to-b from-slate-900 via-slate-800 to-slate-900">
      {/* ... existing sections ... */}

      {/* CTA */}

      <section className="py-16 px-6 text-center">
        <h1 className="text-4xl md:text-5xl font-black text-white mb-4">
          {t.pricing.hero.title}
        </h1>
        <p className="text-xl text-slate-400 max-w-2xl mx-auto mb-8">
          {t.pricing.hero.description} <span className="text-cyan-400">{t.pricing.hero.description_highlight}</span>
        </p>

        {/* Billing Toggle */}
        <div className="flex items-center justify-center gap-4 mb-12">
          <span
            className={`text-sm font-medium ${cart.billingCycle === 'monthly' ? 'text-white' : 'text-slate-500'}`}
          >
            {t.pricing.toggle.monthly}
          </span>
          <Switch
            checked={cart.billingCycle === 'yearly'}
            onCheckedChange={handleBillingToggle}
          />
          <span
            className={`text-sm font-medium ${cart.billingCycle === 'yearly' ? 'text-white' : 'text-slate-500'}`}
          >
            {t.pricing.toggle.yearly}
            <span className="ml-2 px-2 py-0.5 bg-cyan-500/20 text-cyan-400 text-xs rounded-full">
              {t.pricing.toggle.discount}
            </span>
          </span>
        </div>
      </section>

      {/* Pricing Cards */}
      <section className="px-6 max-w-6xl mx-auto">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {plans.map((plan) => (
            <PricingCard
              key={plan.id}
              plan={plan}
              billingCycle={cart.billingCycle}
              onAdd={handlePlanAdd}
            />
          ))}
        </div>
      </section>

      {/* Domain Addon */}
      <section className="py-12 px-6 max-w-6xl mx-auto">
        <div className="bg-slate-800/50 backdrop-blur-sm border border-slate-700 rounded-xl p-6 flex flex-col md:flex-row items-center justify-between gap-4">
          <div>
            <h3 className="text-lg font-semibold text-white mb-1">
              {t.pricing.domain.title}
            </h3>
            <p className="text-slate-400">
              {t.pricing.domain.description_pre} <span className="text-cyan-400 font-mono">.de</span> {t.pricing.domain.description_post}{' '}
              <span className="text-white font-semibold">{t.pricing.domain.price}</span>
            </p>
          </div>
          <div className="flex items-center gap-3">
            <span className="text-sm text-slate-400">{t.pricing.domain.label}</span>
            <Switch
              checked={hasDomain}
              onCheckedChange={handleDomainToggle}
            />
          </div>
        </div>
      </section>

      {/* CTA */}
      {/* Hydration fix: Only render CTA if we sure we are on client and state matches */}
      {mounted && itemCount > 0 && (
        <section className="py-8 px-6 text-center sticky bottom-0 z-50 pointer-events-none">
          <div className="pointer-events-auto inline-block">
             <Link to={cartPath}>
            <Button
              size="lg"
              className="bg-cyan-500 hover:bg-cyan-600 text-white shadow-lg shadow-cyan-500/30 px-8 animate-in fade-in slide-in-from-bottom-4 duration-300"
            >
              <ShoppingCart className="w-5 h-5 mr-2" />
              {t.pricing.cart_button.view} ({itemCount} {itemCount === 1 ? t.pricing.cart_button.item : t.pricing.cart_button.items})
            </Button>
          </Link>
          </div>
        </section>
      )}

      {/* Trust Badges */}
      <section className="py-16 px-6 max-w-4xl mx-auto text-center">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div>
            <p className="text-2xl font-bold text-cyan-400 mb-1">ISO 27001</p>
            <p className="text-sm text-slate-400">{t.pricing.trust.iso}</p>
          </div>
          <div>
            <p className="text-2xl font-bold text-cyan-400 mb-1">99.9%</p>
            <p className="text-sm text-slate-400">{t.pricing.trust.uptime}</p>
          </div>
          <div>
            <p className="text-2xl font-bold text-cyan-400 mb-1">🇩🇪</p>
            <p className="text-sm text-slate-400">{t.pricing.trust.location}</p>
          </div>
        </div>
      </section>
    </div>
  )
}
