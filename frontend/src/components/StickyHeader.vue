<template>
    <div class="sticky-header" :class="{ 'visible': visible }">
        <div class="sticky-header-content">
            <router-link to="/" class="sticky-logo-link">
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
  left: 0;
  right: 0;
  z-index: 1000;
  background-color: #ffffff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transform: translateY(-100%);
  transition: transform 0.3s ease;
}

.sticky-header.visible {
  transform: translateY(0);
}

.sticky-header-content {
  max-width: 1012px;
  margin: 0 auto;
  padding: 10px 15px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 58.333333%; /* Match col-md-7 */
}

.sticky-logo-link {
  display: flex;
  align-items: center;
  text-decoration: none;
  color: inherit;
  margin-left: 0;
}

.sticky-logo {
  height: 40px;
  width: auto;
}

.sticky-title {
  font-size: 1.2rem;
  font-weight: bold;
  color: #4e4e4e;
}

.sticky-controls {
  display: flex;
  gap: 10px;
  margin-right: 1em;
}

.theme-toggle-btn {
  background: transparent;
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-radius: 8px;
  padding: 6px 10px;
  font-size: 16px;
  cursor: pointer;
  transition: all 0.3s ease;
  line-height: 1;
}

.theme-toggle-btn:hover {
  background: rgba(0, 0, 0, 0.05);
  transform: scale(1.1);
}

/* Dark theme */
.dark-theme .sticky-header {
  background-color: #05070a;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.4);
}

.dark-theme .sticky-title {
  color: #ffffff;
}

.dark-theme .theme-toggle-btn {
  border-color: rgba(255, 255, 255, 0.2);
}

.dark-theme .theme-toggle-btn:hover {
  background: rgba(255, 255, 255, 0.1);
}

/* Mobile adjustments */
@media (max-width: 767px) {
  .sticky-header-content {
    width: 100%;
    padding: 8px 15px;
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
