import { createApp } from 'vue'
import { createI18n } from 'vue-i18n'
import App from './App.vue'
import router from './router'
import { zh } from './i18n/zh'
import { en } from './i18n/en'

const savedLang = localStorage.getItem('chromadesk_lang')
const browserLang = navigator.language.startsWith('zh') ? 'zh-CN' : 'en-US'
const locale = (savedLang === 'zh-CN' || savedLang === 'en-US') ? savedLang : browserLang

const i18n = createI18n({
  legacy: false,
  locale,
  fallbackLocale: 'en-US',
  messages: { 'zh-CN': zh, 'en-US': en },
})

const app = createApp(App)
app.use(router)
app.use(i18n)
app.mount('#app')
