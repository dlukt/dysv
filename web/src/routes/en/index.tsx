import { createFileRoute } from '@tanstack/react-router'

import { LandingPage } from '../index'

export const Route = createFileRoute('/en/')({
  component: EnglishLayout,
})

function EnglishLayout() {
  return <LandingPage />
}
