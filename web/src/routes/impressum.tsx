import { createFileRoute, Link } from '@tanstack/react-router'
import { ArrowLeft } from 'lucide-react'

export const Route = createFileRoute('/impressum')({
  component: ImpressumPage,
})

function ImpressumPage() {
  return (
    <div className="min-h-screen bg-gradient-to-b from-slate-900 via-slate-800 to-slate-900">
      {/* Header */}
      <header className="py-6 px-6 max-w-4xl mx-auto">
        <Link
          to="/"
          className="inline-flex items-center gap-2 text-slate-400 hover:text-white transition-colors"
        >
          <ArrowLeft className="w-4 h-4" />
          Zurück zur Startseite
        </Link>
      </header>

      {/* Content */}
      <main className="px-6 max-w-4xl mx-auto py-12">
        <h1 className="text-4xl font-bold text-white mb-8">Impressum</h1>

        <div className="prose prose-invert prose-slate max-w-none">
          <section className="mb-8">
            <h2 className="text-xl font-semibold text-white mb-4">
              Angaben gemäß § 5 TMG
            </h2>
            <div className="bg-slate-800/50 border border-slate-700 rounded-xl p-6 text-slate-300">
              <p className="mb-2">
                <strong className="text-white">dysv.de</strong>
              </p>
              <p className="mb-2">[Vollständiger Name / Firmenname]</p>
              <p className="mb-2">[Straße und Hausnummer]</p>
              <p className="mb-2">[PLZ Ort]</p>
              <p>Deutschland</p>
            </div>
          </section>

          <section className="mb-8">
            <h2 className="text-xl font-semibold text-white mb-4">Kontakt</h2>
            <div className="bg-slate-800/50 border border-slate-700 rounded-xl p-6 text-slate-300">
              <p className="mb-2">
                Telefon: <span className="text-slate-400">[Telefonnummer]</span>
              </p>
              <p>
                E-Mail:{' '}
                <a
                  href="mailto:info@dysv.de"
                  className="text-cyan-400 hover:text-cyan-300"
                >
                  info@dysv.de
                </a>
              </p>
            </div>
          </section>

          <section className="mb-8">
            <h2 className="text-xl font-semibold text-white mb-4">
              Umsatzsteuer-ID
            </h2>
            <div className="bg-slate-800/50 border border-slate-700 rounded-xl p-6 text-slate-300">
              <p>
                Umsatzsteuer-Identifikationsnummer gemäß § 27 a
                Umsatzsteuergesetz:
              </p>
              <p className="text-slate-400 mt-2">[USt-IdNr.]</p>
            </div>
          </section>

          <section className="mb-8">
            <h2 className="text-xl font-semibold text-white mb-4">
              Verantwortlich für den Inhalt nach § 55 Abs. 2 RStV
            </h2>
            <div className="bg-slate-800/50 border border-slate-700 rounded-xl p-6 text-slate-300">
              <p className="mb-2">[Name des Verantwortlichen]</p>
              <p className="mb-2">[Adresse]</p>
            </div>
          </section>

          <section>
            <h2 className="text-xl font-semibold text-white mb-4">
              EU-Streitschlichtung
            </h2>
            <div className="bg-slate-800/50 border border-slate-700 rounded-xl p-6 text-slate-300">
              <p>
                Die Europäische Kommission stellt eine Plattform zur
                Online-Streitbeilegung (OS) bereit:{' '}
                <a
                  href="https://ec.europa.eu/consumers/odr/"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-cyan-400 hover:text-cyan-300"
                >
                  https://ec.europa.eu/consumers/odr/
                </a>
              </p>
              <p className="mt-2">
                Unsere E-Mail-Adresse finden Sie oben im Impressum.
              </p>
            </div>
          </section>
        </div>
      </main>
    </div>
  )
}
