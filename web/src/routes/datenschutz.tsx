import { createFileRoute, Link } from '@tanstack/react-router'
import { ArrowLeft } from 'lucide-react'

export const Route = createFileRoute('/datenschutz')({
  component: DatenschutzPage,
})

function DatenschutzPage() {
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
        <h1 className="text-4xl font-bold text-white mb-8">
          Datenschutzerklärung
        </h1>

        <div className="prose prose-invert prose-slate max-w-none space-y-8">
          <section>
            <h2 className="text-xl font-semibold text-white mb-4">
              1. Datenschutz auf einen Blick
            </h2>
            <div className="bg-slate-800/50 border border-slate-700 rounded-xl p-6 text-slate-300 space-y-4">
              <h3 className="text-lg font-medium text-white">
                Allgemeine Hinweise
              </h3>
              <p>
                Die folgenden Hinweise geben einen einfachen Überblick darüber,
                was mit Ihren personenbezogenen Daten passiert, wenn Sie diese
                Website besuchen. Personenbezogene Daten sind alle Daten, mit
                denen Sie persönlich identifiziert werden können.
              </p>
            </div>
          </section>

          <section>
            <h2 className="text-xl font-semibold text-white mb-4">
              2. Verantwortliche Stelle
            </h2>
            <div className="bg-slate-800/50 border border-slate-700 rounded-xl p-6 text-slate-300">
              <p className="mb-2">[Name / Firma]</p>
              <p className="mb-2">[Adresse]</p>
              <p>
                E-Mail:{' '}
                <a
                  href="mailto:datenschutz@dysv.de"
                  className="text-cyan-400 hover:text-cyan-300"
                >
                  datenschutz@dysv.de
                </a>
              </p>
            </div>
          </section>

          <section>
            <h2 className="text-xl font-semibold text-white mb-4">
              3. Datenerfassung auf dieser Website
            </h2>
            <div className="bg-slate-800/50 border border-slate-700 rounded-xl p-6 text-slate-300 space-y-4">
              <h3 className="text-lg font-medium text-white">Server-Log-Dateien</h3>
              <p>
                Der Provider der Seiten erhebt und speichert automatisch
                Informationen in so genannten Server-Log-Dateien, die Ihr
                Browser automatisch an uns übermittelt. Dies sind:
              </p>
              <ul className="list-disc list-inside text-slate-400">
                <li>Browsertyp und Browserversion</li>
                <li>verwendetes Betriebssystem</li>
                <li>Referrer URL</li>
                <li>Hostname des zugreifenden Rechners</li>
                <li>Uhrzeit der Serveranfrage</li>
                <li>IP-Adresse</li>
              </ul>
            </div>
          </section>

          <section>
            <h2 className="text-xl font-semibold text-white mb-4">
              4. Hosting
            </h2>
            <div className="bg-slate-800/50 border border-slate-700 rounded-xl p-6 text-slate-300 space-y-4">
              <p>
                Diese Website wird in deutschen Rechenzentren gehostet.
                Serverstandort ist Deutschland. Die Datenverarbeitung erfolgt
                gemäß DSGVO.
              </p>
              <p className="text-cyan-400">
                Zertifizierungen: ISO 27001
              </p>
            </div>
          </section>

          <section>
            <h2 className="text-xl font-semibold text-white mb-4">
              5. Zahlungsanbieter
            </h2>
            <div className="bg-slate-800/50 border border-slate-700 rounded-xl p-6 text-slate-300 space-y-4">
              <h3 className="text-lg font-medium text-white">Stripe</h3>
              <p>
                Für Zahlungen nutzen wir Stripe. Anbieter ist die Stripe, Inc.,
                510 Townsend Street, San Francisco, CA 94103, USA.
              </p>
              <p>
                Datenschutzerklärung:{' '}
                <a
                  href="https://stripe.com/de/privacy"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-cyan-400 hover:text-cyan-300"
                >
                  https://stripe.com/de/privacy
                </a>
              </p>
            </div>
          </section>

          <section>
            <h2 className="text-xl font-semibold text-white mb-4">
              6. Ihre Rechte
            </h2>
            <div className="bg-slate-800/50 border border-slate-700 rounded-xl p-6 text-slate-300 space-y-2">
              <p>Sie haben jederzeit das Recht:</p>
              <ul className="list-disc list-inside text-slate-400">
                <li>Auskunft über Ihre gespeicherten Daten zu erhalten</li>
                <li>Berichtigung unrichtiger Daten zu verlangen</li>
                <li>Löschung Ihrer Daten zu verlangen</li>
                <li>Einschränkung der Verarbeitung zu verlangen</li>
                <li>Der Verarbeitung zu widersprechen</li>
                <li>Datenübertragbarkeit zu verlangen</li>
              </ul>
            </div>
          </section>
        </div>
      </main>
    </div>
  )
}
