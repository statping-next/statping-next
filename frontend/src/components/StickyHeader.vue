<template>
    <div class="sticky-header" :class="{ 'active': visible }">
        <div class="sticky-header-content">
            <router-link to="/" class="sticky-logo-link" :class="{ 'visible': visible }">
                <img v-if="core.logo" :src="core.logo" :alt="core.name" class="sticky-logo">
                <span v-else class="sticky-title">{{core.name}}</span>
            </router-link>
            <div class="sticky-controls">
                <button @click="toggleTheme" class="btn btn-sm theme-toggle-btn" :title="darkTheme ? 'Switch to Light Mode' : 'Switch to Dark Mode'">
                    <span v-if="darkTheme">☀️</span>
                    <span v-else>🌙</span>
                </button>
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
  computed: {
    core() {
      return this.$store.getters.core
    },
    darkTheme() {
      return this.$store.getters.darkTheme
    }
  },
  methods: {
    toggleTheme() {
      this.$store.dispatch('toggleTheme')
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
  width: 58.333333%; /* Match col-md-7 */
  max-width: 1012px;
  z-index: 1000;
  background-color: transparent;
  box-shadow: none;
  transition: background-color 0.3s ease, box-shadow 0.3s ease;
}

@media (max-width: 767px) {
  .sticky-header {
    width: 100%;
    left: 0;
    transform: none;
    max-width: 100%;
  }
}

.sticky-header.active {
  /* Background and shadow set by theme classes */
}

.sticky-header-content {
  max-width: 1012px;
  margin: 0 auto;
  padding: 1em 1em;
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
  /* Border and hover colors set by theme classes */
}

.theme-toggle-btn:hover {
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
}
</style>
