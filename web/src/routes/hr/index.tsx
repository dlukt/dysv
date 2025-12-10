import { createFileRoute } from '@tanstack/react-router'

import { LandingPage } from '../index'

export const Route = createFileRoute('/hr/')({
  component: CroatianLayout,
})

function CroatianLayout() {
  return <LandingPage />
}
