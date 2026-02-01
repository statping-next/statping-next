<template>
    <div id="header">
        <router-link to="/" class="header-logo-link">
            <div v-if="!brandingReady" class="header-logo-placeholder" aria-hidden="true"></div>
            <img v-else id="logo" style="max-height: 200px;" class="col-12 text-center pt-4 mb-3 header-title font-6" :src="currentLogo || DEFAULT_LOGO" :alt="core.name" :title="core.name">
        </router-link>
        <h5 id="description" class="col-12 text-center mb-5 header-desc font-3">{{core.description}}</h5>
    </div>
</template>

<script>
import { DEFAULT_LOGO } from '@/constants/branding'

export default {
  name: 'Header',
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

.header-logo-placeholder {
  width: 100%;
  min-height: 80px;
  background: transparent;
}
</style>
