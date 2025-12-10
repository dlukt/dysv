import { createFileRoute } from '@tanstack/react-router'
import { CartPage } from '../cart'

export const Route = createFileRoute('/hr/cart')({
  component: CartPage,
})
