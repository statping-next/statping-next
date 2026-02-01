<template>
    <div class="sticky-header" :class="{ 'active': visible || adminMode, 'admin-mode': adminMode }" :style="headerStyle">
        <div class="sticky-header-content">
            <!-- Left side: Logo (always on left) -->
            <div class="sticky-header-left">
                <router-link to="/" class="sticky-logo-link" :class="{ 'visible': visible || adminMode }">
                    <img :src="displayLogo" :alt="core.name || 'Statping'" class="sticky-logo">
                </router-link>
            </div>

            <!-- Middle: Navigation items (admin-only dashboard) -->
            <nav v-if="adminMode" class="sticky-nav">
                <router-link to="/dashboard/status" class="nav-link" :class="{ 'active': isActiveRoute('/dashboard/status') }" :data-text="$t('status')">{{ $t('status') }}</router-link>
                <router-link to="/dashboard/incidents" class="nav-link" :class="{ 'active': isActiveRoute('/dashboard/incidents') }" :data-text="$t('incidents')">{{ $t('incidents') }}</router-link>
                <router-link v-if="admin" to="/dashboard/logs" class="nav-link" :class="{ 'active': isActiveRoute('/dashboard/logs') }" :data-text="$t('logs')">{{ $t('logs') }}</router-link>
                <router-link to="/dashboard/services" class="nav-link" :class="{ 'active': isActiveRoute('/dashboard/services') }" :data-text="$t('services')">{{ $t('services') }}</router-link>
                <router-link v-if="admin" to="/dashboard/users" class="nav-link" :class="{ 'active': isActiveRoute('/dashboard/users') }" :data-text="$t('users')">{{ $t('users') }}</router-link>
                <router-link v-if="admin" to="/dashboard/settings" class="nav-link" :class="{ 'active': isActiveRoute('/dashboard/settings') }" :data-text="$t('settings')">{{ $t('settings') }}</router-link>
                <router-link v-if="admin" to="/dashboard/help" class="nav-link" :class="{ 'active': isActiveRoute('/dashboard/help') }" :data-text="$t('help')">{{ $t('help') }}</router-link>
            </nav>

            <!-- Right side: Controls and buttons -->
            <div class="sticky-header-right">
                <div v-if="!adminMode" class="sticky-controls">
                    <button @click="cycleRefresh" class="btn btn-sm refresh-btn" :title="refreshTitle">
                        <font-awesome-icon icon="sync-alt" />
                        <span v-if="refreshInterval > 0" class="refresh-badge">{{refreshInterval}}s</span>
                        <span v-else class="refresh-badge off">OFF</span>
                    </button>
                    <button @click="toggleTheme" class="btn btn-sm theme-toggle-btn" :title="darkTheme ? 'Switch to Light Mode' : 'Switch to Dark Mode'">
                        <font-awesome-icon v-if="darkTheme" icon="sun" />
                        <font-awesome-icon v-else icon="moon" />
                    </button>
                    <router-link to="/dashboard/status" class="btn btn-sm admin-btn" title="Admin Dashboard">
                        <font-awesome-icon icon="cog" />
                    </router-link>
                </div>
                <div v-if="adminMode" class="sticky-controls">
                    <button @click="cycleRefresh" class="btn btn-sm refresh-btn" :title="refreshTitle">
                        <font-awesome-icon icon="sync-alt" />
                        <span v-if="refreshInterval > 0" class="refresh-badge">{{refreshInterval}}s</span>
                        <span v-else class="refresh-badge off">OFF</span>
                    </button>
                    <button @click="toggleTheme" class="btn btn-sm theme-toggle-btn" :title="darkTheme ? 'Switch to Light Mode' : 'Switch to Dark Mode'">
                        <font-awesome-icon v-if="darkTheme" icon="sun" />
                        <font-awesome-icon v-else icon="moon" />
                    </button>
                    <button @click="logout" class="btn btn-sm logout-btn" :title="$t('logout')">
                        <font-awesome-icon icon="sign-out-alt" />
                    </button>
                    <router-link to="/" class="btn btn-sm back-btn" title="Back to Home">
                        <font-awesome-icon icon="arrow-left" />
                    </router-link>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { DEFAULT_LOGO, logoUrlIfSet } from '@/constants/branding'

export default {
  name: 'StickyHeader',
  props: {
    visible: {
      type: Boolean,
      default: true
    },
    adminMode: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      containerWidth: null
    }
  },
  computed: {
    brandingReady() {
      return this.$store.getters.brandingReady
    },
    core() {
      return this.$store.getters.core
    },
    darkTheme() {
      return this.$store.getters.darkTheme
    },
    refreshInterval() {
      return this.$store.getters.refreshInterval
    },
    admin() {
      return this.$store.getters.admin
    },
    refreshTitle() {
      if (this.refreshInterval === 0) {
        return 'Auto-refresh: OFF (Click to enable)'
      } else {
        return `Auto-refresh: ${this.refreshInterval}s (Click to change)`
      }
    },
    headerStyle() {
      if (this.containerWidth && window.innerWidth >= 768) {
        return {
          width: `${this.containerWidth}px`
        }
      }
      return {}
    },
    currentLogo() {
      // Only use admin-set logos (non-empty); otherwise null → template uses DEFAULT_LOGO
      if (this.darkTheme) {
        const url = logoUrlIfSet(this.core.logo_dark) || logoUrlIfSet(this.core.logo_light)
        if (url) return url
      } else {
        const url = logoUrlIfSet(this.core.logo_light)
        if (url) return url
      }
      const url = logoUrlIfSet(this.core.logo)
      return url || null
    },
    displayLogo() {
      return this.brandingReady ? (this.currentLogo || DEFAULT_LOGO) : DEFAULT_LOGO
    }
  },
  methods: {
    toggleTheme() {
      this.$store.dispatch('toggleTheme')
    },
    cycleRefresh() {
      this.$store.dispatch('cycleRefreshInterval')
    },
    async logout() {
      const Api = (await import('../API')).default
      await Api.logout()
      this.$store.commit('setHasAllData', false)
      this.$store.commit('setToken', null)
      this.$store.commit('setAdmin', false)
      this.$store.commit('setUser', false)
      // Remove the authentication cookie
      this.$cookies.remove('statping_auth')
      // Force a full page reload to ensure all state is cleared
      window.location.href = '/'
    },
    isActiveRoute(path, exact = false) {
      if (exact) {
        // For Status, match /dashboard/status
        return this.$route.path === path
      } else {
        // For other routes, match the path or any sub-path
        return this.$route.path === path || this.$route.path.startsWith(path + '/')
      }
    },
    updateWidth() {
      if (window.innerWidth < 768) {
        this.containerWidth = null
        return
      }
      const container = document.querySelector('.container.col-md-7')
      if (container) {
        this.containerWidth = container.offsetWidth
      }
    }
  },
  mounted() {
    this.updateWidth()
    window.addEventListener('resize', this.updateWidth)
    this.$nextTick(() => {
      this.updateWidth()
    })
  },
  beforeDestroy() {
    window.removeEventListener('resize', this.updateWidth)
  },
  watch: {
    '$route'() {
      this.$nextTick(() => {
        this.updateWidth()
      })
    }
  }
}
</script>

<style scoped>
.sticky-header {
  position: fixed;
  top: 0;
  left: 50%;
  transform: translateX(-50%);
  z-index: 1000;
  background-color: transparent;
  box-shadow: none;
  transition: background-color 0.3s ease, box-shadow 0.3s ease;
}

/* Width is set dynamically via JavaScript to match container.col-md-7 exactly */
@media (min-width: 768px) {
  .sticky-header {
    /* Width set via :style binding */
  }
}

@media (max-width: 767px) {
  .sticky-header {
    width: 100%;
    left: 0;
    transform: none;
    max-width: 100%;
  }

  .sticky-header-content {
    width: 100%;
    max-width: 100%;
  }
}

.sticky-header.active,
.sticky-header.admin-mode {
  /* Background and shadow set by theme classes */
  /* Always visible in admin mode */
}

.sticky-header-content {
  max-width: 1012px;
  margin: 0 auto;
  padding: 1em;
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.sticky-header-left,
.sticky-header-right {
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

.sticky-nav {
  display: flex;
  align-items: center;
  gap: 2px;
  flex: 1;
  justify-content: center;
  margin: 0 20px;
}

.sticky-nav .nav-link {
  text-decoration: none;
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 14px;
  transition: all 0.3s ease;
  white-space: nowrap;
  position: relative;
  display: inline-block;
  text-align: center;
  /* Colors set by theme classes */
}

/* Use a pseudo-element to reserve space for bold text to prevent width shift */
.sticky-nav .nav-link::before {
  content: attr(data-text);
  font-weight: bold;
  display: block;
  height: 0;
  overflow: hidden;
  visibility: hidden;
  user-select: none;
  pointer-events: none;
  white-space: nowrap;
}

.sticky-nav .nav-link:hover {
  /* Background set by theme classes */
}

.sticky-nav .nav-link.active {
  font-weight: bold;
  /* Background and color set by theme classes */
}

.logout-btn,
.back-btn {
  background: transparent;
  border-radius: 8px;
  padding: 6px 10px;
  font-size: 16px;
  cursor: pointer;
  transition: all 0.3s ease;
  line-height: 1;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 40px;
  /* Border and hover colors set by theme classes */
}

.logout-btn:hover,
.back-btn:hover {
  transform: scale(1.1);
  /* Background set by theme classes */
}

.back-btn {
  /* Same styling as logout-btn, set by theme classes */
}

.sticky-logo-link.admin-logo {
  opacity: 1;
  transform: none;
  pointer-events: auto;
}

/* Ensure admin mode header is always visible and styled */
.sticky-header.admin-mode {
  background-color: transparent;
}

.sticky-header.admin-mode.active {
  /* Background set by theme classes */
}

.sticky-logo-link {
  display: flex;
  align-items: center;
  justify-content: center;
  text-decoration: none;
  color: inherit;
  margin-left: 0;
  opacity: 0;
  transform: translate(20px, 20px);
  pointer-events: none;
  transition: opacity 0.4s ease, transform 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
  flex-shrink: 0;
  /* Fixed 140x40px container to prevent layout shifts */
  width: 140px;
  height: 40px;
  min-width: 140px;
  min-height: 40px;
  max-width: 140px;
  max-height: 40px;
}

.sticky-logo-link.visible {
  opacity: 1;
  transform: translate(0, 0);
  pointer-events: auto;
}

.sticky-logo {
  height: 40px;
  width: auto;
  max-width: 140px;
  object-fit: contain;
}

.sticky-title {
  font-weight: bold;
  white-space: normal; /* Allow wrapping to two lines */
  line-height: 1.2;
  word-wrap: break-word;
  overflow-wrap: break-word;
  hyphens: auto;
  /* Auto-scale font size to fit within 140x40px container */
  font-size: clamp(0.7rem, 2vw, 1rem);
  /* Fit within fixed container */
  width: 100%;
  max-width: 140px;
  max-height: 40px;
  display: block;
  text-align: center;
  overflow: hidden;
  /* Color set by theme classes */
}

/* On smaller screens, adjust font size but keep container size */
@media (max-width: 767px) {
  .sticky-logo-link {
    width: 120px;
    min-width: 120px;
    max-width: 120px;
  }
  .sticky-title {
    max-width: 120px;
    font-size: clamp(0.6rem, 3vw, 0.9rem);
  }
  .sticky-logo {
    max-width: 120px;
  }
}

.sticky-controls {
  display: flex;
  gap: 10px;
  margin-right: 0;
}

.theme-toggle-btn {
  background: transparent;
  border-radius: 8px;
  padding: 6px 10px;
  font-size: 16px;
  cursor: pointer;
  transition: all 0.3s ease;
  line-height: 1;
  min-width: 40px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  /* Border and hover colors set by theme classes */
}

.theme-toggle-btn:hover {
  transform: scale(1.1);
  /* Background set by theme classes */
}

.refresh-btn {
  background: transparent;
  border-radius: 8px;
  padding: 6px 10px;
  font-size: 16px;
  cursor: pointer;
  transition: all 0.3s ease;
  line-height: 1;
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 40px;
  /* Border and hover colors set by theme classes */
}

.refresh-btn:hover {
  transform: scale(1.1);
  /* Background set by theme classes */
}

.refresh-badge {
  position: absolute;
  bottom: 0px;
  right: 0px;
  font-size: 10px;
  font-weight: bold;
  line-height: 1;
  pointer-events: none;
  padding-right: 4px;
  padding-bottom: 4px;
}

.refresh-badge.off {
  font-size: 9px;
}

.admin-btn {
  background: transparent;
  border-radius: 8px;
  padding: 6px 10px;
  font-size: 16px;
  cursor: pointer;
  transition: all 0.3s ease;
  line-height: 1;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 40px;
  /* Border and hover colors set by theme classes */
}

.admin-btn:hover {
  transform: scale(1.1);
  /* Background set by theme classes */
}

/* Mobile adjustments */
@media (max-width: 767px) {
  .sticky-header-content {
    width: 100%;
    padding: 8px 15px;
    max-width: 100%;
    flex-wrap: wrap;
  }

  .sticky-logo {
    height: 35px;
  }

  .sticky-controls {
    margin-right: 0;
  }

  .theme-toggle-btn {
    padding: 8px 12px;
    font-size: 18px;
  }

  .refresh-btn {
    padding: 8px 12px;
    font-size: 18px;
  }

  .admin-btn {
    padding: 8px 12px;
    font-size: 18px;
  }

  .logout-btn,
  .back-btn {
    padding: 8px 12px;
    font-size: 18px;
  }

  .sticky-nav {
    order: 3;
    width: 100%;
    margin: 10px 0 0 0;
    flex-wrap: wrap;
    justify-content: flex-start;
    gap: 2px;
  }

  .sticky-nav .nav-link {
    font-size: 12px;
    padding: 4px 8px;
  }

  .sticky-header-left,
  .sticky-header-right {
    width: auto;
  }
}
</style>
