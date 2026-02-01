import Vuex from 'vuex'
import Vue from 'vue'
import Api from "./API"
import { DEFAULT_FAVICON } from './constants/branding'

Vue.use(Vuex)

function applyBranding(core) {
  if (!core) {
    return
  }

  // Update document title based on project name
  if (core.name) {
    document.title = core.name
  }

  // Set favicon: user's override or project default (never show default if user set one)
  const href = core.favicon || DEFAULT_FAVICON
  let link = document.getElementById('app-favicon') || document.querySelector("link[rel='shortcut icon']") || document.querySelector("link[rel='icon']")
  if (!link) {
    link = document.createElement("link")
    link.rel = "shortcut icon"
    link.id = "app-favicon"
    document.head.appendChild(link)
  }
  link.href = href
}

export const HAS_ALL_DATA = 'HAS_ALL_DATA'
export const HAS_PUBLIC_DATA = 'HAS_PUBLIC_DATA'

export const GET_CORE = 'GET_CORE'
export const GET_SERVICES = 'GET_SERVICES'
export const GET_TOKEN = 'GET_TOKEN'
export const GET_GROUPS = 'GET_GROUPS'
export const GET_MESSAGES = 'GET_MESSAGES'
export const GET_NOTIFIERS = 'GET_NOTIFIERS'
export const GET_USERS = 'GET_USERS'

export default new Vuex.Store({
  state: {
    hasAllData: false,
    hasPublicData: false,
    core: {},
    oauth: {},
    token: null,
    services: [],
    service: null,
    groups: [],
    messages: [],
    incidents: [],
    users: [],
    notifiers: [],
    checkins: [],
    admin: false,
    user: false,
    loggedIn: false,
    darkTheme: localStorage.getItem('darkTheme') === 'true' || false,
    refreshInterval: parseInt(localStorage.getItem('refreshInterval')) || 15, // 0 = off, 15 = 15s, 60 = 60s
    modal: {
      visible: false,
      title: "Modal Header",
      body: "This is the content for the modal body",
      btnText: "Save Changes",
      btnColor: "btn-primary",
      func: null,
    }
  },
  getters: {
    hasAllData: state => state.hasAllData,
    hasPublicData: state => state.hasPublicData,
    admin: state => state.admin,
    core: state => state.core,
    oauth: state => state.oauth,
    token: state => state.token,
    services: state => state.services,
    service: state => state.service,
    groups: state => state.groups,
    messages: state => state.messages,
    incidents: state => state.incidents,
    users: state => state.users,
    notifiers: state => state.notifiers,
    checkins: state => state.checkins,
    loggedIn: state => state.loggedIn,
    modal: state => state.modal,
    darkTheme: state => state.darkTheme,
    refreshInterval: state => state.refreshInterval,
    // True once core has been loaded from API - use to avoid showing default logo/favicon then replacing with theme
    brandingReady: state => !!(state.core && (state.core.name !== undefined && state.core.name !== null)),

    isAdmin: state => state.admin,
    isUser: state => state.user,

    globalMessages: state => {
      if (!state.messages || !Array.isArray(state.messages)) {
        return []
      }
      return state.messages.filter(s => !s.service || s.service === 0)
    },
    // Public view: only non-archived global incidents (like announcements hidden after end_on).
    globalIncidents: state => {
      if (!state.incidents || !Array.isArray(state.incidents)) {
        return []
      }
      return state.incidents.filter(i => Number(i.service) === 0 && !i.archived)
    },
    servicesInOrder: state => state.services.sort((a, b) => a.order_id - b.order_id),
    servicesNoGroup: state => state.services.filter(g => g.group_id === 0).sort((a, b) => a.order_id - b.order_id),
    groupsInOrder: state => state.groups.sort((a, b) => a.order_id - b.order_id),
    groupsClean: state => state.groups.filter(g => g.name !== '').sort((a, b) => a.order_id - b.order_id),
    groupsCleanInOrder: state => state.groups.filter(g => g.name !== '').sort((a, b) => a.order_id - b.order_id).sort((a, b) => a.order_id - b.order_id),
    serviceCheckins: (state) => (id) => {
      return state.checkins.filter(c => c.service_id === id)
    },
    serviceByAll: (state) => (element) => {
      if (!isNaN(parseFloat(element)) && isFinite(element)) {
        return state.services.find(s => s.id === parseInt(element))
      } else {
        return state.services.find(s => s.permalink === element)
      }
    },
    serviceById: (state) => (id) => {
      return state.services.find(s => s.permalink === id || s.id === id)
    },
    servicesInGroup: (state) => (id) => {
      return state.services.filter(s => s.group_id === id).sort((a, b) => a.order_id - b.order_id)
    },
    serviceMessages: (state) => (id) => {
      if (!state.messages || !Array.isArray(state.messages)) {
        return []
      }
      return state.messages.filter(s => s.service === id)
    },
    onlineServices: (state) => (online) => {
      return state.services.filter(s => s.online === online)
    },
    groupById: (state) => (id) => {
      return state.groups.find(g => g.id === id)
    },
    cleanGroups: (state) => () => {
      return state.groups.filter(g => g.name !== '').sort((a, b) => a.order_id - b.order_id)
    },
    userById: (state) => (id) => {
      return state.users.find(u => u.id === id)
    },
    messageById: (state) => (id) => {
      return state.messages.find(m => m.id === id)
    },
  },
  mutations: {
    setHasAllData(state, bool) {
      state.hasAllData = bool
    },
    setHasPublicData(state, bool) {
      state.hasPublicData = bool
    },
    setCore(state, core) {
      state.core = core
      applyBranding(core)
    },
    setToken(state, token) {
      state.token = token
    },
    setService(state, service) {
      state.service = service
    },
    setServices(state, services) {
      state.services = services
    },
    setCheckins(state, checkins) {
      state.checkins = checkins
    },
    setGroups(state, groups) {
      state.groups = groups
    },
    setMessages(state, messages) {
      state.messages = messages
    },
    setIncidents(state, incidents) {
      state.incidents = incidents || []
    },
    setUsers(state, users) {
      state.users = users
    },
    setNotifiers(state, notifiers) {
      state.notifiers = notifiers
    },
    setAdmin(state, admin) {
      state.admin = admin
    },
    setLoggedIn(state, loggedIn) {
      state.loggedIn = loggedIn
    },
    setUser(state, user) {
      state.user = user
    },
    setDarkTheme(state, darkTheme) {
      state.darkTheme = darkTheme
      localStorage.setItem('darkTheme', darkTheme)
      // Apply theme to html and body (html prevents white strip at bottom)
      const html = document.documentElement
      if (darkTheme) {
        html.classList.remove('light-theme')
        html.classList.add('dark-theme')
        document.body.classList.remove('light-theme')
        document.body.classList.add('dark-theme')
      } else {
        html.classList.remove('dark-theme')
        html.classList.add('light-theme')
        document.body.classList.remove('dark-theme')
        document.body.classList.add('light-theme')
      }
    },
    setRefreshInterval(state, interval) {
      state.refreshInterval = interval
      localStorage.setItem('refreshInterval', interval)
    },
    setOAuth(state, oauth) {
      state.oauth = oauth
    },
    setModal(state, modal) {
      state.modal = modal
    },
  },
  actions: {
    async getAllServices(context) {
      const services = await Api.services()
      context.commit("setServices", services);
    },
    async loadCore(context) {
      const core = await Api.core()
      const token = await Api.token()
      context.commit("setCore", core);
      context.commit('setAdmin', token)
      context.commit('setCore', core)
      context.commit('setUser', token !== undefined)

      // Ensure branding is applied on initial load as well
      applyBranding(core)

      // Apply theme: use localStorage if set, otherwise use core default
      const savedTheme = localStorage.getItem('darkTheme')
      let darkTheme = false
      if (savedTheme !== null) {
        // User has a saved preference
        darkTheme = savedTheme === 'true'
      } else if (core.default_theme) {
        // Use instance default if no user preference
        darkTheme = core.default_theme === 'dark'
      }
      context.commit('setDarkTheme', darkTheme)

      // Apply refresh rate: use localStorage if set, otherwise use core default
      const savedRefresh = localStorage.getItem('refreshInterval')
      let refreshInterval = 15
      if (savedRefresh !== null) {
        // User has a saved preference
        refreshInterval = parseInt(savedRefresh) || 15
      } else if (core.default_refresh_rate !== undefined) {
        // Use instance default if no user preference
        refreshInterval = core.default_refresh_rate || 15
      }
      context.commit('setRefreshInterval', refreshInterval)
    },
    toggleTheme(context) {
      context.commit('setDarkTheme', !context.state.darkTheme)
    },
    cycleRefreshInterval(context) {
      const current = context.state.refreshInterval
      // Cycle: 0 (off) -> 15 -> 60 -> 0
      const next = current === 0 ? 15 : current === 15 ? 60 : 0
      context.commit('setRefreshInterval', next)
    },
    async loadRequired(context) {
      const groups = await Api.groups()
      context.commit("setGroups", groups);
      const services = await Api.services()
      context.commit("setServices", services);
      const messages = await Api.messages()
      context.commit("setMessages", messages)
      const incidents = await Api.incidents_all()
      context.commit("setIncidents", incidents)
      const oauth = await Api.oauth()
      context.commit("setOAuth", oauth);
      context.commit("setHasPublicData", true)
    },
    async loadAdmin(context) {
      const groups = await Api.groups()
      context.commit("setGroups", groups);
      const services = await Api.services()
      context.commit("setServices", services);
      const messages = await Api.messages()
      context.commit("setMessages", messages)
      const incidents = await Api.incidents_all()
      context.commit("setIncidents", incidents)
      context.commit("setHasPublicData", true)
      const checkins = await Api.checkins()
      context.commit("setCheckins", checkins);
      const notifiers = await Api.notifiers()
      context.commit("setNotifiers", notifiers);
      const users = await Api.users()
      context.commit("setUsers", users);
      const oauth = await Api.oauth()
      context.commit("setOAuth", oauth);
    }
  }
});
