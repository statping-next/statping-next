<template>
    <div class="sticky-header" :class="{ 'active': visible }" :style="headerStyle">
        <div class="sticky-header-content">
            <router-link to="/" class="sticky-logo-link" :class="{ 'visible': visible }">
                <img v-if="core.logo" :src="core.logo" :alt="core.name" class="sticky-logo">
                <span v-else class="sticky-title">{{core.name}}</span>
            </router-link>
            <div class="sticky-controls">
                <button @click="cycleRefresh" class="btn btn-sm refresh-btn" :title="refreshTitle">
                    <font-awesome-icon icon="sync-alt" />
                    <span v-if="refreshInterval > 0" class="refresh-badge">{{refreshInterval}}s</span>
                    <span v-else class="refresh-badge off">OFF</span>
                </button>
                <button @click="toggleTheme" class="btn btn-sm theme-toggle-btn" :title="darkTheme ? 'Switch to Light Mode' : 'Switch to Dark Mode'">
                    <font-awesome-icon v-if="darkTheme" icon="sun" />
                    <font-awesome-icon v-else icon="moon" />
                </button>
                <router-link to="/dashboard" class="btn btn-sm admin-btn" title="Admin Dashboard">
                    <font-awesome-icon icon="cog" />
                </router-link>
            </div>
        </div>
    </div>
</template>

<script>
export default {
  name: 'StickyHeader',
  props: {
    visible: {
      type: Boolean,
      default: true
    }
  },
  data() {
    return {
      containerWidth: null
    }
  },
  computed: {
    core() {
      return this.$store.getters.core
    },
    darkTheme() {
      return this.$store.getters.darkTheme
    },
    refreshInterval() {
      return this.$store.getters.refreshInterval
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
    }
  },
  methods: {
    toggleTheme() {
      this.$store.dispatch('toggleTheme')
    },
    cycleRefresh() {
      this.$store.dispatch('cycleRefreshInterval')
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

.sticky-header.active {
  /* Background and shadow set by theme classes */
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

.sticky-logo-link {
  display: flex;
  align-items: center;
  text-decoration: none;
  color: inherit;
  margin-left: 0;
  opacity: 0;
  transform: translate(20px, 20px);
  pointer-events: none;
  transition: opacity 0.4s ease, transform 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.sticky-logo-link.visible {
  opacity: 1;
  transform: translate(0, 0);
  pointer-events: auto;
}

.sticky-logo {
  height: 40px;
  width: auto;
}

.sticky-title {
  font-size: 1.2rem;
  font-weight: bold;
  /* Color set by theme classes */
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
}
</style>
