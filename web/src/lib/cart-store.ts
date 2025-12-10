// Cart store using TanStack Store
// Manages plan selection, billing cycle, and add-ons

import { Store } from '@tanstack/store'

export type BillingCycle = 'monthly' | 'yearly'

export interface CartState {
  planId: string | null
  billingCycle: BillingCycle
  addons: string[]
}

const initialState: CartState = {
  planId: null,
  billingCycle: 'monthly',
  addons: [],
}

export const cartStore = new Store<CartState>(initialState)

// Actions
export function setPlan(planId: string) {
  cartStore.setState((state) => ({
    ...state,
    planId,
  }))
}

export function clearPlan() {
  cartStore.setState((state) => ({
    ...state,
    planId: null,
  }))
}

export function setBillingCycle(cycle: BillingCycle) {
  cartStore.setState((state) => ({
    ...state,
    billingCycle: cycle,
  }))
}

export function addAddon(addonId: string) {
  cartStore.setState((state) => ({
    ...state,
    addons: state.addons.includes(addonId)
      ? state.addons
      : [...state.addons, addonId],
  }))
}

export function removeAddon(addonId: string) {
  cartStore.setState((state) => ({
    ...state,
    addons: state.addons.filter((id) => id !== addonId),
  }))
}

export function toggleAddon(addonId: string) {
  cartStore.setState((state) => ({
    ...state,
    addons: state.addons.includes(addonId)
      ? state.addons.filter((id) => id !== addonId)
      : [...state.addons, addonId],
  }))
}

export function clearCart() {
  cartStore.setState(() => initialState)
}

// Selectors
export function hasItems(state: CartState): boolean {
  return state.planId !== null || state.addons.length > 0
}

export function getItemCount(state: CartState): number {
  return (state.planId ? 1 : 0) + state.addons.length
}
