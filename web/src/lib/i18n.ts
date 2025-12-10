import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';

import en from '../locales/en.json';
import de from '../locales/de.json';
import hr from '../locales/hr.json';

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources: {
      en: { translation: en },
      de: { translation: de },
      hr: { translation: hr },
    },
    lng: 'de', // Force default to 'de' (server-side matches getLocaleFromPath default)
    fallbackLng: 'de',
    debug: import.meta.env.DEV,
    interpolation: {
      escapeValue: false, // React escapes by default
    },
    detection: {
        order: ['path', 'navigator'],
        lookupFromPathIndex: 0
    }
  });

export default i18n;
