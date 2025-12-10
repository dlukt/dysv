import {
  HeadContent,
  Scripts,
  createRootRouteWithContext,
  useRouterState,
} from '@tanstack/react-router'
import { TanStackRouterDevtoolsPanel } from '@tanstack/react-router-devtools'
import { TanStackDevtools } from '@tanstack/react-devtools'

import Header from '../components/Header'

import TanStackQueryDevtools from '../integrations/tanstack-query/devtools'

import StoreDevtools from '../lib/demo-store-devtools'
import { getLocaleFromPath } from '../lib/locale'

import appCss from '../styles.css?url'

import type { QueryClient } from '@tanstack/react-query'

interface MyRouterContext {
  queryClient: QueryClient
}

import { NotFound } from '../components/NotFound'
import { useEffect } from 'react'
import { initPersistence } from '../lib/cart-store'

export const Route = createRootRouteWithContext<MyRouterContext>()({
  head: () => ({
    meta: [
      {
        charSet: 'utf-8',
      },
      {
        name: 'viewport',
        content: 'width=device-width, initial-scale=1',
      },
      {
        title: 'dysv.de – Redundant Kubernetes Hosting | Serverstandort Deutschland',
      },
      {
        name: 'description',
        content: 'Unkillable uptime for your Next.js, Nuxt, and React apps. 3-node redundant hosting with flat pricing. German datacenter, ISO 27001 certified.',
      },
    ],
    links: [
      {
        rel: 'stylesheet',
        href: appCss,
      },
      // Hreflang tags for multilingual SEO
      {
        rel: 'alternate',
        hrefLang: 'de',
        href: 'https://dysv.de/de',
      },
      {
        rel: 'alternate',
        hrefLang: 'en',
        href: 'https://dysv.de/en',
      },
      {
        rel: 'alternate',
        hrefLang: 'hr',
        href: 'https://dysv.de/hr',
      },
      {
        rel: 'alternate',
        hrefLang: 'x-default',
        href: 'https://dysv.de/',
      },
    ],
  }),

  shellComponent: RootDocument,
  notFoundComponent: NotFound,
})

// ... existing imports

function RootDocument({ children }: { children: React.ReactNode }) {
  const { pathname } = useRouterState({ select: (state) => state.location })
  const locale = getLocaleFromPath(pathname)

  useEffect(() => {
    initPersistence()
  }, [])


  return (
    <html lang={locale}>
      <head>
        <HeadContent />
      </head>
      <body
        suppressHydrationWarning
      >
        <Header />
        {children}
        {import.meta.env.DEV ? (
          <TanStackDevtools
            config={{
              position: 'bottom-right',
            }}
            plugins={[
              {
                name: 'Tanstack Router',
                render: <TanStackRouterDevtoolsPanel />,
              },
              TanStackQueryDevtools,
              StoreDevtools,
            ]}
          />
        ) : null}
        <Scripts />
      </body>
    </html>
  )
}
