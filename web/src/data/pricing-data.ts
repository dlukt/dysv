// Pricing data for dysv.de hosting plans
// Based on SPEC.md Section 2.A

export interface Plan {
  id: string
  name: string
  monthlyPrice: number
  targetAudience: string
  limits: string
  features: string[]
}

export interface Addon {
  id: string
  name: string
  monthlyPrice: number
}

export const plans: Plan[] = [
  {
    id: 'static-micro',
    name: 'Static Micro',
    monthlyPrice: 3.9,
    targetAudience: 'React/Vue SPAs',
    limits: 'Shared RAM, 1GB Storage',
    features: [
      'Static site hosting',
      'Shared resources',
      '1GB NVMe storage',
      'SSL included',
      'German datacenter',
    ],
  },
  {
    id: 'node-starter',
    name: 'Node Starter',
    monthlyPrice: 9.9,
    targetAudience: 'Personal Blogs',
    limits: '1 vCPU (Shared), 512MB RAM, 5GB Storage',
    features: [
      'Next.js / Nuxt support',
      'High-Performance Burstable CPU',
      '512MB RAM',
      '5GB NVMe storage',
      'SSL included',
      'German datacenter',
    ],
  },
  {
    id: 'node-pro',
    name: 'Node Pro',
    monthlyPrice: 39.9,
    targetAudience: 'E-commerce/SaaS',
    limits: '2 vCPU (Dedicated), 4GB RAM, 20GB Storage',
    features: [
      'Next.js / Nuxt support',
      'Dedicated Core Performance',
      '4GB RAM',
      '20GB NVMe storage',
      'SSL included',
      'German datacenter',
      'Priority support',
    ],
  },
]

export const addons: Addon[] = [
  {
    id: 'de-domain',
    name: '.de Domain',
    monthlyPrice: 1.0,
  },
]

// Yearly billing gives 2 months free (pay for 10, get 12)
export const YEARLY_DISCOUNT_MONTHS = 2

export function calculateYearlyPrice(monthlyPrice: number): number {
  return monthlyPrice * (12 - YEARLY_DISCOUNT_MONTHS)
}

export function getYearlyMonthlyEquivalent(monthlyPrice: number): number {
  return calculateYearlyPrice(monthlyPrice) / 12
}
