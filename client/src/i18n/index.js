/**
 * Vue i18n
 *
 * @library
 *
 * http://kazupon.github.io/vue-i18n/en/
 */

// Lib imports
import Vue from 'vue'
import VueI18n from 'vue-i18n'
import messages from '@/lang'

Vue.use(VueI18n)

const getLocale = () => {
  const stored = localStorage.getItem('locale')
  if (stored && messages[stored]) return stored
  const userLang = (navigator.language || navigator.userLanguage).replace('-', '_').toLowerCase()
  const test = lang => elem => elem.toLowerCase().indexOf(lang) !== -1
  const fromBrowser =
    Object.keys(messages).filter(test(userLang))[0] ||
    Object.keys(messages).filter(test(userLang.split('_')[0]))[0]
  if (fromBrowser) return fromBrowser
  return 'en_US'
}

const i18n = new VueI18n({
  locale: getLocale(),
  fallbackLocale: 'en_US',
  messages
})

export default i18n
