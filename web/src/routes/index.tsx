import { createFileRoute, Link } from '@tanstack/react-router'
import {
  Server,
  Shield,
  Zap,
  Globe,
  Clock,
  ArrowRight,
} from 'lucide-react'

import { Button } from '@/components/ui/button'
import { useTranslation } from '@/lib/i18n'

export const Route = createFileRoute('/')({
  component: LandingPage,
})

export function LandingPage() {
  const { t, locale } = useTranslation()
  const pricingPath = locale === 'de' ? '/pricing' : `/${locale}/pricing`

  const features = [
    {
      icon: <Server className="w-10 h-10 text-cyan-400" />,
      title: t.landing.features.server_title,
      description: t.landing.features.server_desc,
    },
    {
      icon: <Shield className="w-10 h-10 text-cyan-400" />,
      title: t.landing.features.k8s_title,
      description: t.landing.features.k8s_desc,
    },
    {
      icon: <Zap className="w-10 h-10 text-cyan-400" />,
      title: t.landing.features.nvme_title,
      description: t.landing.features.nvme_desc,
    },
    {
      icon: <Globe className="w-10 h-10 text-cyan-400" />,
      title: t.landing.features.geo_title,
      description: t.landing.features.geo_desc,
    },
    {
      icon: <Clock className="w-10 h-10 text-cyan-400" />,
      title: t.landing.features.pricing_title,
      description: t.landing.features.pricing_desc,
    },
    {
      icon: <Shield className="w-10 h-10 text-cyan-400" />,
      title: t.landing.features.ssl_title,
      description: t.landing.features.ssl_desc,
    },
  ]

  return (
    <div className="min-h-screen bg-gradient-to-b from-slate-900 via-slate-800 to-slate-900">
      {/* Hero */}
      <section className="relative py-24 px-6 text-center overflow-hidden">
        <div className="absolute inset-0 bg-gradient-to-r from-cyan-500/5 via-blue-500/5 to-purple-500/5" />
        <div className="relative max-w-4xl mx-auto">
          <p className="text-cyan-400 font-mono text-sm mb-4 tracking-widest uppercase">
            {t.landing.hero.surtitle}
          </p>
          <h1 className="text-5xl md:text-7xl font-black text-white mb-6 leading-tight">
            <span className="bg-gradient-to-r from-cyan-400 to-blue-400 bg-clip-text text-transparent">
              {t.landing.hero.title_line1}
            </span>
            <br />
            {t.landing.hero.title_line2}
          </h1>
          <p className="text-xl md:text-2xl text-slate-300 mb-4 max-w-2xl mx-auto">
            {t.landing.hero.subtitle}
          </p>
          <p className="text-lg text-slate-400 mb-6 max-w-xl mx-auto">
            {t.landing.hero.description}
          </p>

          {/* Framework Logos/Names - Keeping English/Technical terms */}
          <div className="flex flex-wrap items-center justify-center gap-3 mb-10 text-sm">
             {/* ... */}
            <span className="px-3 py-1.5 bg-slate-800/80 border border-slate-700 rounded-full text-slate-300">Next.js</span>
            <span className="px-3 py-1.5 bg-slate-800/80 border border-slate-700 rounded-full text-slate-300">Nuxt</span>
            <span className="px-3 py-1.5 bg-slate-800/80 border border-slate-700 rounded-full text-slate-300">Remix</span>
            <span className="px-3 py-1.5 bg-slate-800/80 border border-slate-700 rounded-full text-slate-300">Astro</span>
            <span className="px-3 py-1.5 bg-slate-800/80 border border-slate-700 rounded-full text-slate-300">SvelteKit</span>
            <span className="px-3 py-1.5 bg-slate-800/80 border border-slate-700 rounded-full text-slate-300">React</span>
            <span className="px-3 py-1.5 bg-slate-800/80 border border-slate-700 rounded-full text-slate-300">Vue</span>
            <span className="px-3 py-1.5 bg-slate-800/80 border border-slate-700 rounded-full text-slate-300">Angular</span>
          </div>

          <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
            <Link to={pricingPath}>
              <Button
                size="lg"
                className="bg-cyan-500 hover:bg-cyan-600 text-white shadow-lg shadow-cyan-500/30 px-8"
              >
                {t.landing.hero.cta}
                <ArrowRight className="w-5 h-5 ml-2" />
              </Button>
            </Link>
            <p className="text-slate-500 text-sm">{t.landing.hero.starting_at}</p>
          </div>
        </div>
      </section>

      {/* Trust Bar */}
      <section className="py-8 px-6 border-y border-slate-700/50">
        <div className="max-w-5xl mx-auto flex flex-wrap items-center justify-center gap-8 md:gap-16 text-center">
          <div>
            <p className="text-2xl font-bold text-white">99.9%</p>
            <p className="text-xs text-slate-500 uppercase tracking-wide">{t.landing.trust.sla}</p>
          </div>
          <div>
            <p className="text-2xl font-bold text-white">🇩🇪</p>
            <p className="text-xs text-slate-500 uppercase tracking-wide">{t.landing.trust.location}</p>
          </div>
          <div>
            <p className="text-2xl font-bold text-white">GDPR</p>
            <p className="text-xs text-slate-500 uppercase tracking-wide">{t.landing.trust.compliant}</p>
          </div>
        </div>
      </section>

      {/* Features */}
      <section className="py-20 px-6 max-w-7xl mx-auto">
        <div className="text-center mb-16">
          <h2 className="text-3xl md:text-4xl font-bold text-white mb-4">
            {t.landing.dev_section.title}
          </h2>
          <p className="text-lg text-slate-400 max-w-2xl mx-auto">
            {t.landing.dev_section.description}
          </p>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {features.map((feature) => (
            <div
              key={feature.title}
              className="bg-slate-800/50 backdrop-blur-sm border border-slate-700 rounded-xl p-6 hover:border-cyan-500/50 transition-all duration-300"
            >
              <div className="mb-4">{feature.icon}</div>
              <h3 className="text-xl font-semibold text-white mb-2">
                {feature.title}
              </h3>
              <p className="text-slate-400">{feature.description}</p>
            </div>
          ))}
        </div>
      </section>

      {/* CTA */}
      <section className="py-20 px-6">
        <div className="max-w-4xl mx-auto text-center bg-gradient-to-r from-cyan-500/10 to-blue-500/10 border border-cyan-500/20 rounded-2xl p-12">
          <h2 className="text-3xl md:text-4xl font-bold text-white mb-4">
            {t.landing.cta_section.title}
          </h2>
          <p className="text-lg text-slate-300 mb-8 max-w-xl mx-auto">
            {t.landing.cta_section.description}
          </p>
          <Link to={pricingPath}>
            <Button
              size="lg"
              className="bg-cyan-500 hover:bg-cyan-600 text-white shadow-lg shadow-cyan-500/30 px-10"
            >
              {t.landing.cta_section.button}
            </Button>
          </Link>
        </div>
      </section>

      {/* Footer */}
      <footer className="py-12 px-6 border-t border-slate-700/50">
        <div className="max-w-7xl mx-auto flex flex-col md:flex-row items-center justify-between gap-6">
          <div className="text-slate-400 text-sm">
            © 2025 dysv.de. {t.landing.footer.rights}
          </div>
          <nav className="flex items-center gap-6 text-sm">
            <Link to="/impressum" className="text-slate-400 hover:text-white transition-colors">
              Impressum
            </Link>
            <Link to="/datenschutz" className="text-slate-400 hover:text-white transition-colors">
              Datenschutz
            </Link>
            <Link to="/agb" className="text-slate-400 hover:text-white transition-colors">
              AGB
            </Link>
          </nav>
        </div>
      </footer>
    </div>
  )
}
