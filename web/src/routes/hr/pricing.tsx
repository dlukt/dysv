import { createFileRoute } from '@tanstack/react-router'
import { PricingPage } from '../pricing'

export const Route = createFileRoute('/hr/pricing')({
  component: PricingPage,
})
