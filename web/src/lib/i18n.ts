import { useRouterState } from '@tanstack/react-router'
import { getLocaleFromPath, type Locale } from './locale'

const de = {
  landing: {
    hero: {
      surtitle: '3-Node Redundantes Hosting',
      title_line1: 'Unkillable',
      title_line2: 'Uptime',
      subtitle: 'Redundantes Kubernetes-Hosting für Ihre statischen Seiten und Node.js-Anwendungen.',
      description: 'Ihre App läuft gleichzeitig auf 3 Nodes. Fällt einer aus, übernehmen die anderen.',
      cta: 'Pläne ansehen',
      starting_at: 'Ab 3,90 €/Monat',
    },
    features: {
      server_title: 'Serverstandort Deutschland',
      server_desc: 'Ihre Daten bleiben in Deutschland. ISO 27001 zertifizierte Rechenzentren mit strenger DSGVO-Compliance.',
      k8s_title: 'Kubernetes-Powered',
      k8s_desc: 'Enterprise-Infrastruktur. Auto-Scaling, Self-Healing und Zero-Downtime Deployments.',
      nvme_title: 'NVMe Speicher',
      nvme_desc: 'Rasendschnelle NVMe SSDs für extrem kurze Ladezeiten und optimale Performance.',
      geo_title: 'Geo-Redundant',
      geo_desc: 'Maximale Ausfallsicherheit. Ihre Seite läuft über mehrere physische Verfügbarkeitszonen.',
      pricing_title: 'Einfache Preise',
      pricing_desc: 'Kein Metering. Keine Überraschungen. Nur planbare, transparente monatliche Kosten.',
      ssl_title: 'SSL Inklusive',
      ssl_desc: 'Kostenlose SSL-Zertifikate für alle Domains. Sicher per Default.',
    },
    trust: {
      sla: 'Verfügbarkeits-SLA',
      location: 'Serverstandort',
      compliant: 'Konform',
    },
    dev_section: {
      title: 'Für Entwickler gebaut',
      description: 'Deployen Sie Ihre Next.js oder Nuxt App in Minuten. Wir kümmern uns um die Infrastruktur.',
    },
    cta_section: {
      title: 'Bereit zum Start?',
      description: 'Starten Sie mit dem Static Micro Plan für nur 3,90 €/Monat. Keine Kreditkarte für Setup erforderlich.',
      button: 'Plan auswählen',
    },
    footer: {
      rights: 'Alle Rechte vorbehalten.',
    },
  },
  pricing: {
    hero: {
      title: 'Einfache Pauschalpreise',
      description: 'Keine Verbrauchmessung. Keine bösen Überraschungen.',
      description_highlight: 'Einfach deutsche Ingenieurskunst.',
    },
    toggle: {
      monthly: 'Monatlich',
      yearly: 'Jährlich',
      discount: '2 Monate gratis',
    },
    domain: {
      title: 'Domain benötigt?',
      description_pre: 'Fügen Sie eine',
      description_post: 'Domain für nur',
      price: '1,00 €/Monat',
      label: '.de Domain hinzufügen',
    },
    cart_button: {
        view: 'Warenkorb ansehen',
        item: 'Artikel',
        items: 'Artikel',
    },
    card: {
      popular: 'Beliebt',
      selected: 'Ausgewählt',
      add_to_cart: 'In den Warenkorb',
      per_month: '/Monat',
      per_year: '/Jahr',
      free_months: '2 Monate gratis',
    },
    trust: {
        iso: 'Zertifiziertes Rechenzentrum',
        uptime: 'Verfügbarkeits-SLA',
        location: 'Serverstandort Deutschland',
    }
  },
  cart: {
    empty: {
      title: 'Ihr Warenkorb ist leer',
      description: 'Fügen Sie einen Plan oder ein Addon hinzu, um zu beginnen.',
      button: 'Pläne ansehen',
    },
    title: 'Ihr Warenkorb',
    back: 'Zurück zu den Preisen',
    billed_monthly: 'Monatliche Abrechnung',
    billed_yearly: 'Jährliche Abrechnung',
    summary: {
      title: 'Zusammenfassung',
      subtotal: 'Zwischensumme',
      yearly_discount: 'Jährliche Abrechnung (2 Monate gratis)',
      total: 'Gesamt',
      checkout: 'Zur Kasse',
      redirect: 'Sie werden zu Stripe weitergeleitet, um die Zahlung abzuschließen.',
    },
    domain_reg: 'Domainregistrierung',
  },
}

const en = {
  landing: {
    hero: {
      surtitle: '3-Node Redundant Hosting',
      title_line1: 'Unkillable',
      title_line2: 'Uptime',
      subtitle: 'Redundant Kubernetes hosting for your static sites and Node.js applications.',
      description: 'Your app runs across 3 nodes simultaneously. If one fails, the others keep serving.',
      cta: 'View Plans',
      starting_at: 'Starting at €3.90/month',
    },
    features: {
      server_title: 'Server Location Germany',
      server_desc: 'Your data stays in Germany. ISO 27001 certified datacenters with strict GDPR compliance.',
      k8s_title: 'Kubernetes-Powered',
      k8s_desc: 'Enterprise-grade infrastructure. Auto-scaling, self-healing, and zero-downtime deployments.',
      nvme_title: 'NVMe Storage',
      nvme_desc: 'Blazing fast NVMe SSDs for lightning-quick page loads and optimal performance.',
      geo_title: 'Geo-Redundant',
      geo_desc: 'Unkillable uptime. Your site runs across multiple physical availability zones.',
      pricing_title: 'Flat Pricing',
      pricing_desc: 'No metering. No surprise bills. Just predictable, transparent monthly costs.',
      ssl_title: 'SSL Included',
      ssl_desc: 'Free SSL certificates for all your domains. Secure by default.',
    },
    trust: {
      sla: 'Uptime SLA',
      location: 'Data Location',
      compliant: 'Compliant',
    },
    dev_section: {
      title: 'Built for Developers',
      description: 'Deploy your Next.js or Nuxt app in minutes. We handle the infrastructure.',
    },
    cta_section: {
      title: 'Ready to Deploy?',
      description: 'Get started with our Static Micro plan for just €3.90/month. No credit card required for setup.',
      button: 'Choose Your Plan',
    },
    footer: {
      rights: 'All rights reserved.',
    },
  },
  pricing: {
    hero: {
      title: 'Simple, Flat Pricing',
      description: 'No metering. No surprise bills.',
      description_highlight: 'Just German engineering.',
    },
    toggle: {
      monthly: 'Monthly',
      yearly: 'Yearly',
      discount: '2 months free',
    },
    domain: {
      title: 'Need a domain?',
      description_pre: 'Add a',
      description_post: 'domain for just',
      price: '€1.00/month',
      label: 'Add .de domain',
    },
    cart_button: {
        view: 'View Cart',
        item: 'item',
        items: 'items',
    },
    card: {
      popular: 'Most Popular',
      selected: 'Selected',
      add_to_cart: 'Add to Cart',
      per_month: '/mo',
      per_year: '/year',
      free_months: '2 months free',
    },
    trust: {
        iso: 'Certified Datacenter',
        uptime: 'Uptime SLA',
        location: 'Server Location Germany',
    }
  },
  cart: {
    empty: {
      title: 'Your cart is empty',
      description: 'Add a plan or addon to get started.',
      button: 'View Plans',
    },
    title: 'Your Cart',
    back: 'Back to Pricing',
    billed_monthly: 'Monthly billing',
    billed_yearly: 'Yearly billing',
    summary: {
      title: 'Order Summary',
      subtotal: 'Subtotal',
      yearly_discount: 'Yearly billing (2 mo free)',
      total: 'Total',
      checkout: 'Proceed to Checkout',
      redirect: 'You\'ll be redirected to Stripe to complete payment.',
    },
    domain_reg: 'Domain registration',
  },
}

const hr = {
  landing: {
    hero: {
      surtitle: '3-Node Redundantni Hosting',
      title_line1: 'Unkillable',
      title_line2: 'Uptime',
      subtitle: 'Redundantni Kubernetes hosting za vaše statičke stranice i Node.js aplikacije.',
      description: 'Vaša aplikacija se vrti na 3 node-a istovremeno. Ako jedan padne, ostali preuzimaju.',
      cta: 'Pogledaj planove',
      starting_at: 'Već od 3,90 €/mj',
    },
    features: {
      server_title: 'Serveri u Njemačkoj',
      server_desc: 'Vaši podaci ostaju u Njemačkoj. ISO 27001 certificirani podatkovni centri sa strogom GDPR usklađenošću.',
      k8s_title: 'Kubernetes',
      k8s_desc: 'Enterprise infrastruktura. Auto-scaling, self-healing i zero-downtime deploymenti.',
      nvme_title: 'NVMe Pohrana',
      nvme_desc: 'Brzi NVMe SSD-ovi za ekstremno brzo učitavanje i optimalne performanse.',
      geo_title: 'Geo-Redundantno',
      geo_desc: 'Maksimalna pouzdanost. Vaša stranica radi u više fizičkih zona dostupnosti.',
      pricing_title: 'Fiksne Cijene',
      pricing_desc: 'Bez mjerenja prometa. Bez iznenađenja. Samo predvidljivi, transparentni mjesečni troškovi.',
      ssl_title: 'SSL Uključen',
      ssl_desc: 'Besplatni SSL certifikati za sve domene. Sigurno po defaultu.',
    },
    trust: {
      sla: 'Uptime SLA',
      location: 'Lokacija',
      compliant: 'Usklađeno',
    },
    dev_section: {
      title: 'Stvoreno za Developere',
      description: 'Deployajte Next.js ili Nuxt aplikaciju u nekoliko minuta. Mi brinemo o infrastrukturi.',
    },
    cta_section: {
      title: 'Spremni za start?',
      description: 'Započnite sa Static Micro planom za samo 3,90 €/mj. Kreditna kartica nije potrebna za setup.',
      button: 'Odaberi Plan',
    },
    footer: {
      rights: 'Sva prava pridržana.',
    },
  },
  pricing: {
    hero: {
      title: 'Jednostavne Fiksne Cijene',
      description: 'Bez mjerenja. Bez loših iznenađenja.',
      description_highlight: 'Samo njemački inženjering.',
    },
    toggle: {
      monthly: 'Mjesečno',
      yearly: 'Godišnje',
      discount: '2 mjeseca gratis',
    },
    domain: {
      title: 'Trebate domenu?',
      description_pre: 'Dodajte',
      description_post: 'domenu za samo',
      price: '1,00 €/mj',
      label: 'Dodaj .de domenu',
    },
    cart_button: {
        view: 'Vidi košaricu',
        item: 'artikl',
        items: 'artikala',
    },
    card: {
      popular: 'Najpopularnije',
      selected: 'Odabrano',
      add_to_cart: 'Dodaj u košaricu',
      per_month: '/mj',
      per_year: '/god',
      free_months: '2 mjeseca gratis',
    },
    trust: {
        iso: 'Certificirani DC',
        uptime: 'Uptime SLA',
        location: 'Lokacija Njemačka',
    }
  },
  cart: {
    empty: {
      title: 'Vaša košarica je prazna',
      description: 'Dodajte plan ili dodatak da biste započeli.',
      button: 'Pogledaj planove',
    },
    title: 'Vaša Košarica',
    back: 'Natrag na cjenik',
    billed_monthly: 'Mjesečna naplata',
    billed_yearly: 'Godišnja naplata',
    summary: {
      title: 'Sažetak narudžbe',
      subtotal: 'Međuzbroj',
      yearly_discount: 'Godišnja naplata (2 mj gratis)',
      total: 'Ukupno',
      checkout: 'Na blagajnu',
      redirect: 'Bit ćete preusmjereni na Stripe za plaćanje.',
    },
    domain_reg: 'Registracija domene',
  },
}

const translations: Record<Locale, typeof en> = {
  de,
  en,
  hr,
}

export function useTranslation() {
  const { pathname } = useRouterState({ select: (state) => state.location })
  const locale = getLocaleFromPath(pathname)
  // Default to DE if locale not found (shouldn't happen with strict types but good safety)
  const t = translations[locale] || translations.de
  return { t, locale }
}
