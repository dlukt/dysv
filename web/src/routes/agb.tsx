import { createFileRoute, Link } from '@tanstack/react-router'
import { ArrowLeft } from 'lucide-react'

export const Route = createFileRoute('/agb')({
  component: AGBPage,
})

function AGBPage() {
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
          Allgemeine Geschäftsbedingungen (AGB)
        </h1>

        <div className="prose prose-invert prose-slate max-w-none space-y-8">
          <section>
            <h2 className="text-xl font-semibold text-white mb-4">
              § 1 Geltungsbereich
            </h2>
            <div className="bg-slate-800/50 border border-slate-700 rounded-xl p-6 text-slate-300 space-y-4">
              <p>
                Diese Allgemeinen Geschäftsbedingungen gelten für alle
                Verträge, die zwischen dysv.de und dem Kunden über die auf
                dieser Website angebotenen Hosting-Dienste geschlossen werden.
              </p>
            </div>
          </section>

          <section>
            <h2 className="text-xl font-semibold text-white mb-4">
              § 2 Vertragsschluss
            </h2>
            <div className="bg-slate-800/50 border border-slate-700 rounded-xl p-6 text-slate-300 space-y-4">
              <p>
                Die Darstellung der Produkte auf der Website stellt kein
                rechtlich bindendes Angebot, sondern eine Aufforderung zur
                Bestellung dar.
              </p>
              <p>
                Der Vertrag kommt durch die Annahme der Bestellung durch uns
                zustande, die durch eine Bestätigungs-E-Mail erfolgt.
              </p>
            </div>
          </section>

          <section>
            <h2 className="text-xl font-semibold text-white mb-4">
              § 3 Preise und Zahlungsbedingungen
            </h2>
            <div className="bg-slate-800/50 border border-slate-700 rounded-xl p-6 text-slate-300 space-y-4">
              <p>
                Alle angegebenen Preise sind Endpreise und enthalten die
                gesetzliche Mehrwertsteuer.
              </p>
              <ul className="list-disc list-inside text-slate-400">
                <li>
                  <strong className="text-white">Monatliche Abrechnung:</strong>{' '}
                  Zahlung monatlich im Voraus
                </li>
                <li>
                  <strong className="text-white">Jährliche Abrechnung:</strong>{' '}
                  Zahlung jährlich im Voraus (2 Monate kostenlos)
                </li>
              </ul>
              <p>Die Zahlung erfolgt über Stripe.</p>
            </div>
          </section>

          <section>
            <h2 className="text-xl font-semibold text-white mb-4">
              § 4 Domain-Registrierung
            </h2>
            <div className="bg-slate-800/50 border border-slate-700 rounded-xl p-6 text-slate-300 space-y-4">
              <div className="bg-yellow-500/10 border border-yellow-500/30 rounded-lg p-4">
                <p className="text-yellow-400 font-medium">
                  Wichtiger Hinweis zur Kündigung:
                </p>
                <p className="text-slate-300 mt-2">
                  Bei Kündigung des Hosting-Vertrags innerhalb der ersten 12
                  Monate wird für die .de Domain eine einmalige Gebühr von{' '}
                  <strong className="text-white">€10,00</strong> berechnet, um
                  die Registrierungskosten zu decken.
                </p>
              </div>
            </div>
          </section>

          <section>
            <h2 className="text-xl font-semibold text-white mb-4">
              § 5 Leistungsumfang
            </h2>
            <div className="bg-slate-800/50 border border-slate-700 rounded-xl p-6 text-slate-300 space-y-4">
              <p>Der Leistungsumfang ergibt sich aus dem gewählten Tarif:</p>
              <ul className="list-disc list-inside text-slate-400">
                <li>Static Micro: Shared RAM, 1GB Storage</li>
                <li>Node Starter: 1 vCPU (Shared, Burstable), 512MB RAM, 5GB Storage</li>
                <li>Node Pro: 2 vCPU (Dedicated), 4GB RAM, 20GB Storage</li>
              </ul>
              <p>
                Die angegebenen Ressourcen werden durch Kubernetes-Quotas
                garantiert.
              </p>
            </div>
          </section>

          <section>
            <h2 className="text-xl font-semibold text-white mb-4">
              § 6 Verfügbarkeit
            </h2>
            <div className="bg-slate-800/50 border border-slate-700 rounded-xl p-6 text-slate-300 space-y-4">
              <p>
                Wir garantieren eine Verfügbarkeit von{' '}
                <strong className="text-cyan-400">99,9%</strong> im
                Jahresdurchschnitt.
              </p>
              <p>
                Geplante Wartungsarbeiten werden mindestens 24 Stunden im Voraus
                angekündigt und zählen nicht als Ausfallzeit.
              </p>
            </div>
          </section>

          <section>
            <h2 className="text-xl font-semibold text-white mb-4">
              § 7 Kündigung
            </h2>
            <div className="bg-slate-800/50 border border-slate-700 rounded-xl p-6 text-slate-300 space-y-4">
              <p>
                <strong className="text-white">Monatliche Verträge:</strong>{' '}
                Kündigung jederzeit zum Monatsende.
              </p>
              <p>
                <strong className="text-white">Jährliche Verträge:</strong>{' '}
                Kündigung mit einer Frist von 30 Tagen zum Vertragsende.
              </p>
            </div>
          </section>

          <section>
            <h2 className="text-xl font-semibold text-white mb-4">
              § 8 Schlussbestimmungen
            </h2>
            <div className="bg-slate-800/50 border border-slate-700 rounded-xl p-6 text-slate-300 space-y-4">
              <p>Es gilt deutsches Recht.</p>
              <p>
                Gerichtsstand ist, soweit gesetzlich zulässig, der Sitz des
                Anbieters.
              </p>
            </div>
          </section>
        </div>
      </main>
    </div>
  )
}
