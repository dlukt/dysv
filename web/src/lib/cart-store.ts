// Cart store using TanStack Store
// Manages plan selection, billing cycle, and add-ons

import { Store } from '@tanstack/store'

export type BillingCycle = 'monthly' | 'yearly'

export interface CartItem {
  id: string
  type: 'plan' | 'addon'
  quantity: number
}

export interface CartState {
  items: CartItem[]
  billingCycle: BillingCycle
}

const initialState: CartState = {
  items: [],
  billingCycle: 'monthly',
}

export const cartStore = new Store<CartState>(initialState)

// Persistence
export function initPersistence() {
  if (typeof window === 'undefined') return
  console.log('Initializing persistence...')

  // Load initial state
  const saved = localStorage.getItem('dysv_cart')
  if (saved) {
    try {
      const parsed = JSON.parse(saved)
      console.log('Restoring cart from storage:', parsed)
      cartStore.setState(() => parsed)
    } catch (e) {
      console.error('Failed to parse cart', e)
    }
  } else {
    console.log('No saved cart found in storage')
  }

  // Save on change
  cartStore.subscribe(() => {
    const state = cartStore.state
    console.log('Saving cart to storage:', state)
    localStorage.setItem('dysv_cart', JSON.stringify(state))
  })
}

// Actions
export function addToCart(id: string, type: 'plan' | 'addon') {
  console.log('Action: addToCart', id, type)
  cartStore.setState((state) => {
    const existingItem = state.items.find((item) => item.id === id)
    if (existingItem) {
      return {
        ...state,
        items: state.items.map((item) =>
          item.id === id ? { ...item, quantity: item.quantity + 1 } : item
        ),
      }
    }
    return {
      ...state,
      items: [...state.items, { id, type, quantity: 1 }],
    }
  })
}

export function updateQuantity(id: string, quantity: number) {
  console.log('Action: updateQuantity', id, quantity)
  cartStore.setState((state) => {
    if (quantity <= 0) {
      return {
        ...state,
        items: state.items.filter((item) => item.id !== id),
      }
    }
    return {
      ...state,
      items: state.items.map((item) =>
        item.id === id ? { ...item, quantity } : item
      ),
    }
  })
}

export function removeFromCart(id: string) {
  cartStore.setState((state) => ({
    ...state,
    items: state.items.filter((item) => item.id !== id),
  }))
}

export function setBillingCycle(cycle: BillingCycle) {
  cartStore.setState((state) => ({
    ...state,
    billingCycle: cycle,
  }))
}

export function clearCart() {
  cartStore.setState(() => initialState)
}

// Selectors
export function hasItems(state: CartState): boolean {
  return state.items.length > 0
}

export function getItemCount(state: CartState): number {
  return state.items.reduce((sum, item) => sum + item.quantity, 0)
}

export function getCartItem(state: CartState, id: string): CartItem | undefined {
  return state.items.find((item) => item.id === id)
}
