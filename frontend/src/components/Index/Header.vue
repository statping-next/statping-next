<template>
    <div id="header">
        <router-link to="/" class="header-logo-link">
            <img id="logo" class="header-logo-img col-12 pt-4 mb-3 header-title font-6" :src="displayLogo" :alt="core.name || 'Statping'" :title="core.name || 'Statping'">
        </router-link>
        <h5 id="description" class="col-12 text-center mb-5 header-desc font-3">{{core.description}}</h5>
    </div>
</template>

<script>
import { DEFAULT_LOGO, logoUrlIfSet } from '@/constants/branding'

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
    }
}
</script>

<style scoped>
.header-logo-link {
  text-decoration: none;
  cursor: pointer;
  display: block;
  text-align: center;
}

.header-logo-img {
  display: block;
  max-width: 60%;
  max-height: 200px;
  margin-left: auto;
  margin-right: auto;
}

</style>
