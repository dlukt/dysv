import { createFileRoute } from '@tanstack/react-router'

import { LandingPage } from '../index'

export const Route = createFileRoute('/de/')({
  component: GermanLayout,
})

function GermanLayout() {
  return <LandingPage />
}
