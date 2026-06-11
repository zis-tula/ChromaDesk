import { createI18n } from 'vue-i18n'
import zhCN from './zh-CN.js'
import enUS from './en-US.js'

function getDefaultLocale() {
  const stored = localStorage.getItem('chromadesk_admin_locale')
  if (stored && ['zh-CN', 'en-US'].includes(stored)) return stored
  const nav = navigator.language || navigator.userLanguage || 'en'
  return nav.startsWith('zh') ? 'zh-CN' : 'en-US'
}

const i18n = createI18n({
  legacy: false,
  locale: getDefaultLocale(),
  fallbackLocale: 'en-US',
  messages: {
    'zh-CN': zhCN,
    'en-US': enUS
  }
})

export function setLocale(locale) {
  i18n.global.locale.value = locale
  localStorage.setItem('chromadesk_admin_locale', locale)
}

export function getLocale() {
  return i18n.global.locale.value
}

export default i18n
