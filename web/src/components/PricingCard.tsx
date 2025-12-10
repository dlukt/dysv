// Glassmorphism pricing card component

import { Check } from 'lucide-react'
import { Button } from '@/components/ui/button'
import type { Plan } from '@/data/pricing-data'
import {
  calculateYearlyPrice,
  getYearlyMonthlyEquivalent,
} from '@/data/pricing-data'
import { useTranslation } from '@/lib/i18n'

interface PricingCardProps {
  plan: Plan
  billingCycle: 'monthly' | 'yearly'
  isSelected?: boolean
  onSelect: (planId: string) => void
}

export function PricingCard({
  plan,
  billingCycle,
  isSelected = false,
  onSelect,
}: PricingCardProps) {
  const { t } = useTranslation()
  const isYearly = billingCycle === 'yearly'
  const displayPrice = isYearly
    ? getYearlyMonthlyEquivalent(plan.monthlyPrice)
    : plan.monthlyPrice
  const totalYearly = calculateYearlyPrice(plan.monthlyPrice)
  const isPopular = plan.id === 'node-starter'

  return (
    <div
      className={`
        relative flex flex-col rounded-2xl p-6
        backdrop-blur-xl bg-slate-800/40 border
        transition-all duration-300 hover:scale-[1.02]
        ${isSelected ? 'border-cyan-400 shadow-lg shadow-cyan-500/20' : 'border-slate-700 hover:border-cyan-500/50'}
        ${isPopular ? 'ring-2 ring-cyan-400' : ''}
      `}
    >
      {isPopular && (
        <div className="absolute -top-3 left-1/2 -translate-x-1/2 px-4 py-1 bg-gradient-to-r from-cyan-500 to-blue-500 text-white text-xs font-semibold rounded-full">
          {t.pricing.card.popular}
        </div>
      )}

      <div className="mb-4">
        <h3 className="text-xl font-bold text-white mb-1">{plan.name}</h3>
        <p className="text-sm text-slate-400">{plan.targetAudience}</p>
      </div>

      <div className="mb-6">
        <div className="flex items-baseline gap-1">
          <span className="text-4xl font-black text-white">
            €{displayPrice.toFixed(2)}
          </span>
          <span className="text-slate-400">{t.pricing.card.per_month}</span>
        </div>
        {isYearly && (
          <p className="text-sm text-cyan-400 mt-1">
            €{totalYearly.toFixed(2)}{t.pricing.card.per_year} ({t.pricing.card.free_months})
          </p>
        )}
      </div>

      <div className="flex-grow mb-6">
        <p className="text-xs text-slate-500 uppercase tracking-wide mb-3">
          {plan.limits}
        </p>
        <ul className="space-y-2">
          {plan.features.map((feature) => (
            <li key={feature} className="flex items-center gap-2 text-sm text-slate-300">
              <Check className="w-4 h-4 text-cyan-400 shrink-0" />
              {feature}
            </li>
          ))}
        </ul>
      </div>

      <Button
        onClick={() => onSelect(plan.id)}
        variant={isSelected ? 'default' : 'outline'}
        className={`
          w-full
          ${isSelected
            ? 'bg-cyan-500 hover:bg-cyan-600 text-white shadow-lg shadow-cyan-500/30'
            : 'border-slate-600 hover:border-cyan-500 hover:text-cyan-400'
          }
        `}
      >
        {isSelected ? t.pricing.card.selected : t.pricing.card.add_to_cart}
      </Button>
    </div>
  )
}
