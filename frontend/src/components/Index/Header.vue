<template>
    <div id="header" style="position: relative;">
        <div class="theme-toggle-container">
            <button @click="toggleTheme" class="btn btn-sm theme-toggle-btn" :title="darkTheme ? 'Switch to Light Mode' : 'Switch to Dark Mode'">
                <span v-if="darkTheme">☀️</span>
                <span v-else>🌙</span>
            </button>
        </div>
        <img v-if="core.logo" id="logo" class="col-12 text-center pt-4 mt-4 mb-3 header-title font-6" :src="core.logo" :alt="core.name" :title="core.name">
        <h1 v-if="!core.logo" id="title" class="col-12 text-center pt-4 mt-4 mb-3 header-title font-6">{{core.name}}</h1>
        <h5 id="description" class="col-12 text-center mb-5 header-desc font-3">{{core.description}}</h5>
    </div>
</template>

<script>
export default {
  name: 'Header',
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
.theme-toggle-container {
  position: absolute;
  top: 20px;
  right: 20px;
  z-index: 100;
}

.theme-toggle-btn {
  background: transparent;
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-radius: 8px;
  padding: 8px 12px;
  font-size: 18px;
  cursor: pointer;
  transition: all 0.3s ease;
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
</style>
