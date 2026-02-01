<template>
    <div class="row mt-4">
        <div class="mx-auto">
            <img alt="Statping-ng 404" class="" style="max-width:480px" :src="displayLogo">
        </div>

        <div class="col-12 mt-5 mb-5">
            <div class="text-center">
                <h2>Page Not Found</h2>
                <h5 class="text-muted">This URL doesn't seem to be available</h5>

                <router-link class="btn btn-outline-success mt-4" to="/">Back To Homepage</router-link>
            </div>
        </div>

    </div>
</template>

<script>
import { DEFAULT_LOGO, logoUrlIfSet } from '@/constants/branding'

export default {
  name: 'NotFound',
  computed: {
    brandingReady () {
      return this.$store.getters.brandingReady
    },
    core () {
      return this.$store.getters.core
    },
    darkTheme () {
      return this.$store.getters.darkTheme
    },
    currentLogo () {
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
    displayLogo () {
      return this.brandingReady ? (this.currentLogo || DEFAULT_LOGO) : DEFAULT_LOGO
    }
  }
}
</script>
