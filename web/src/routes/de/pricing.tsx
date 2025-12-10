import { createFileRoute } from '@tanstack/react-router'
import { PricingPage } from '../pricing'

export const Route = createFileRoute('/de/pricing')({
  component: PricingPage,
})
