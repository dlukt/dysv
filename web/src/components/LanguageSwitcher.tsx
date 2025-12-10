import { useNavigate, useRouterState } from '@tanstack/react-router'
import { Globe } from 'lucide-react'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { type Locale, switchLocale, supportedLocales, getLocaleFromPath } from '@/lib/locale'

const localeLabels: Record<Locale, string> = {
  de: 'Deutsch',
  en: 'English',
  hr: 'Hrvatski',
}

export function LanguageSwitcher() {
  const navigate = useNavigate()
  const { pathname } = useRouterState({ select: (state) => state.location })
  const currentLocale = getLocaleFromPath(pathname)

  const handleValueChange = (value: string) => {
    const newPath = switchLocale(pathname, value as Locale)
    navigate({ to: newPath })
  }

  return (
    <Select value={currentLocale} onValueChange={handleValueChange}>
      <SelectTrigger className="w-[130px] h-9 bg-transparent border-slate-700 text-slate-300 hover:text-white hover:border-slate-600 focus:ring-offset-0 focus:ring-0 data-[placeholder]:text-slate-300">
        <div className="flex items-center gap-2">
            <Globe className="w-4 h-4 shrink-0" />
            <SelectValue />
        </div>
      </SelectTrigger>
      <SelectContent className="bg-slate-800 border-slate-700 text-slate-300">
        {supportedLocales.map((loc) => (
          <SelectItem 
            key={loc} 
            value={loc}
            className="focus:bg-slate-700 focus:text-white cursor-pointer"
          >
            {localeLabels[loc]}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  )
}
