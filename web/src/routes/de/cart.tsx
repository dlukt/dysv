import { createFileRoute } from '@tanstack/react-router'
import { CartPage } from '../cart'

export const Route = createFileRoute('/de/cart')({
  component: CartPage,
})
