import { useState, useEffect } from 'react'
import { Link, useRouterState } from '@tanstack/react-router'
import { useStore } from '@tanstack/react-store'
import { ShoppingCart, Menu, X } from 'lucide-react'

import { cartStore, getItemCount } from '@/lib/cart-store'
import { getLocaleFromPath, type Locale } from '@/lib/locale'
import { LanguageSwitcher } from './LanguageSwitcher'
import { useAuth } from '@/hooks/use-auth'

const navLabels: Record<Locale, { home: string; pricing: string; cart: string; login: string; logout: string; account: string }> = {
  de: { home: 'Startseite', pricing: 'Preise', cart: 'Warenkorb', login: 'Anmelden', logout: 'Abmelden', account: 'Konto' },
  en: { home: 'Home', pricing: 'Pricing', cart: 'Cart', login: 'Login', logout: 'Logout', account: 'Account' },
  hr: { home: 'Početna', pricing: 'Cjenik', cart: 'Košarica', login: 'Prijava', logout: 'Odjava', account: 'Račun' },
}

export default function Header() {
  const [isMenuOpen, setIsMenuOpen] = useState(false)
  const cart = useStore(cartStore)
  const itemCount = getItemCount(cart)
  const { pathname } = useRouterState({ select: (state) => state.location })
  const locale = getLocaleFromPath(pathname)
  const labels = navLabels[locale] ?? navLabels.de
  const { user, logout } = useAuth()
  
  // Logic: Default locale (de) uses root paths, others get prefixed.
  const base = locale === 'de' ? '' : `/${locale}`
  const homePath = base || '/'
  const pricingPath = `${base}/pricing`
  const cartPath = `${base}/cart`
  const loginPath = `/login` // Login is global? Or should be localized? For now, global.

  // Close menu when path changes
  useEffect(() => {
    setIsMenuOpen(false)
  }, [pathname])

  const NavItems = ({ mobile = false }: { mobile?: boolean }) => (
    <>
      <Link
        to={homePath}
        className={`hover:text-white transition-colors ${mobile ? 'block py-2 text-lg' : ''}`}
        activeProps={{ className: `text-white ${mobile ? 'block py-2 text-lg font-semibold' : ''}` }}
      >
        {labels.home}
      </Link>
      <Link
        to={pricingPath}
        className={`hover:text-white transition-colors ${mobile ? 'block py-2 text-lg' : ''}`}
        activeProps={{ className: `text-white ${mobile ? 'block py-2 text-lg font-semibold' : ''}` }}
      >
        {labels.pricing}
      </Link>
      
      {/* Auth Links */}
      {user ? (
        <>
            <div className={`flex items-center gap-4 ${mobile ? 'flex-col items-start w-full' : ''}`}>
                 <Link 
                    to="/account/addresses"
                    className={`text-sm text-cyan-400 font-medium hover:text-cyan-300 transition-colors ${mobile ? 'py-1' : ''}`}
                 >
                    {user.username}
                 </Link>
                 <button
                    onClick={() => logout()}
                    className={`text-slate-300 hover:text-white transition-colors ${mobile ? 'block py-2 text-lg w-full text-left' : ''}`}
                 >
                    {labels.logout}
                 </button>
            </div>
        </>
      ) : (
        <Link
            to={loginPath}
            className={`hover:text-white transition-colors ${mobile ? 'block py-2 text-lg' : ''}`}
            activeProps={{ className: `text-white ${mobile ? 'block py-2 text-lg font-semibold' : ''}` }}
        >
            {labels.login}
        </Link>
      )}

      <Link
        to={cartPath}
        className={`relative inline-flex items-center gap-2 rounded-lg transition-colors ${
          mobile 
            ? 'flex w-full py-2 text-lg text-slate-300 hover:text-white' 
            : 'px-3 py-2 border border-slate-700 hover:border-cyan-500/60 hover:text-white'
        }`}
        activeProps={{ 
          className: mobile
             ? 'flex w-full py-2 text-lg text-white font-semibold'
             : 'relative inline-flex items-center gap-2 px-3 py-2 rounded-lg border border-cyan-500 text-white' 
        }}
      >
        <ShoppingCart className="w-5 h-5 md:w-4 md:h-4" />
        <span>{labels.cart}</span>
        {itemCount > 0 && (
          <span className={`bg-cyan-500 text-white text-xs rounded-full flex items-center justify-center ${
            mobile ? 'ml-2 w-6 h-6' : 'absolute -top-2 -right-2 w-5 h-5'
          }`}>
            {itemCount}
          </span>
        )}
      </Link>
      
      {!mobile && <div className="h-6 w-px bg-slate-800" />}
      
      <div className={mobile ? 'pt-4 border-t border-slate-800 mt-2' : ''}>
         <LanguageSwitcher />
      </div>
    </>
  )

  return (
    <header className="sticky top-0 z-40 bg-slate-900/80 backdrop-blur border-b border-slate-800">
      <div className="max-w-7xl mx-auto px-6 py-4">
        <div className="flex items-center justify-between">
          <Link to={homePath} className="text-2xl font-black text-white">
            <span className="text-slate-400">dysv</span>
            <span className="text-cyan-400">.de</span>
          </Link>

          {/* Desktop Nav */}
          <nav className="hidden md:flex items-center gap-6 text-slate-300">
            <NavItems />
          </nav>

          {/* Mobile Menu Toggle */}
          <button 
            className="md:hidden text-slate-300 hover:text-white"
            onClick={() => setIsMenuOpen(!isMenuOpen)}
            aria-label="Toggle menu"
          >
            {isMenuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
          </button>
        </div>

        {/* Mobile Nav */}
        {isMenuOpen && (
          <nav className="md:hidden flex flex-col gap-2 mt-4 pb-4 animate-in fade-in slide-in-from-top-4 duration-200">
            <NavItems mobile />
          </nav>
        )}
      </div>
    </header>
  )
}
