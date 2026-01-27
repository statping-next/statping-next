<template>
    <div id="header">
        <router-link to="/" class="header-logo-link">
            <img v-if="currentLogo" id="logo" style="max-height: 200px;" class="col-12 text-center pt-4 mb-3 header-title font-6" :src="currentLogo" :alt="core.name" :title="core.name">
            <h1 v-if="!currentLogo" id="title" class="col-12 text-center pt-4 mt-4 mb-3 header-title font-6">{{core.name}}</h1>
        </router-link>
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
      },
      currentLogo() {
        // Use theme-specific logo if available
        if (this.darkTheme) {
          // Dark theme: use dark logo if set, otherwise fallback to light logo
          if (this.core.logo_dark) {
            return this.core.logo_dark
          } else if (this.core.logo_light) {
            return this.core.logo_light
          }
        } else {
          // Light theme: use light logo if set
          if (this.core.logo_light) {
            return this.core.logo_light
          }
        }
        // Fallback to legacy logo field for backward compatibility
        if (this.core.logo) {
          return this.core.logo
        }
        return null
      }
    }
}
</script>

<style scoped>
.header-logo-link {
  text-decoration: none;
  cursor: pointer;
  display: block;
}

.header-logo-link:hover {
  opacity: 0.9;
}
</style>
