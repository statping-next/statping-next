import Vue from 'vue'
import VueI18n from 'vue-i18n'
import language from './languages'

Vue.use(VueI18n)

export default new VueI18n({
  fallbackLocale: 'en',
  messages: language
})
