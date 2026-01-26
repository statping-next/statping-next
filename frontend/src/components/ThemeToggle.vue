<template>
    <div class="global-controls">
        <button @click="toggleTheme" class="btn btn-sm theme-toggle-btn" :title="darkTheme ? 'Switch to Light Mode' : 'Switch to Dark Mode'">
            <span v-if="darkTheme">☀️</span>
            <span v-else>🌙</span>
        </button>
    </div>
</template>

<script>
export default {
  name: 'ThemeToggle',
    computed: {
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
.global-controls {
  position: fixed;
  top: 20px;
  right: 0;
  left: 0;
  max-width: 1012px;
  margin: 0 auto;
  padding: 0 20px;
  z-index: 1000;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  pointer-events: none;
}

.global-controls > * {
  pointer-events: auto;
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

.dark-theme .theme-toggle-btn {
  border-color: rgba(255, 255, 255, 0.2);
}

.dark-theme .theme-toggle-btn:hover {
  background: rgba(255, 255, 255, 0.1);
}

/* Mobile adjustments */
@media (max-width: 767px) {
  .global-controls {
    padding: 0 15px;
  }

  .theme-toggle-btn {
    padding: 8px 12px;
    font-size: 18px;
  }
}
</style>
