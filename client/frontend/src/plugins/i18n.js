import { ref, provide } from 'vue'
import { createI18n } from 'vue-i18n'

let i18n = null
const ltr = ref(true)
const fallback = 'en_US'

const getLocale = () => {
  const stored = localStorage.getItem('locale')
  if (stored && localeList.indexOf(stored) !== -1) return stored
  const userLang = (navigator.language || navigator.userLanguage).replace('-', '_').toLowerCase()
  const test = lang => elem => elem.toLowerCase().indexOf(lang) !== -1
  const fromBrowser =
    localeList.filter(test(userLang))[0] ||
    localeList.filter(test(userLang.split('_')[0]))[0]
  if (fromBrowser) return fromBrowser
  return fallback
}

export default async () => {
  const locale = getLocale()
  i18n = createI18n({
    legacy: false,
    locale,
    fallbackLocale: fallback
  })
  await updateLocale(locale, false)

  const i18nInstall = i18n.install
  i18n.install = (app) => {
    app.provide('ltr', ltr)
    i18nInstall(app)
  }
  return i18n
}

const rtl = ['ar_SA', 'he_IL']
const files = ['common', 'env', 'errors', 'files', 'hotkeys', 'nodes', 'oauth', 'operators', 'scopes', 'servers', 'settings', 'templates', 'users']
export async function updateLocale(locale, save = true) {
  if (save) localStorage.setItem('locale', locale)
  const messages = {}
  for (let i = 0; i < files.length; i++) {
    const file = files[i]
    try {
      messages[file] = (await import(`../lang/${locale}/${file}.json`)).default
    } catch (e) {
      messages[file] = (await import(`../lang/${fallback}/${file}.json`)).default
    }
  }

  i18n.global.setLocaleMessage(locale, messages)
  i18n.global.locale.value = locale

  document.querySelector('html').setAttribute('lang', locale)
  document.querySelector('html').setAttribute('dir', rtl.indexOf(locale) === -1 ? 'ltr' : 'rtl')
  ltr.value = rtl.indexOf(locale) === -1 ? true : false
}

export const locales = localeList.map(locale => {
  let [lang, region] = locale.split('_')

  if (locale === 'sr_SP') {
    // crodwin uses the wrong country code for serbia, so we need to manually fix it
    region = 'RS'
  }

  const f = new Intl.DisplayNames(lang, { type: 'language', languageDisplay: 'standard' })
  return { value: locale, label: f.of(`${lang}-${region}`) }
})
