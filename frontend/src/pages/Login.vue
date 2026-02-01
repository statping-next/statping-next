<template>
    <div class="offset-md-3 offset-lg-4 offset-0 col-lg-4 col-md-6 mt-5">

      <div class="offset-1 offset-lg-2 col-lg-8 col-10 mb-4 mb-md-3 login-logo-wrap">
        <img alt="Statping-ng Login" class="embed-responsive" :src="displayLogo">
      </div>

      <div class="login_container col-12 p-4">
            <FormLogin/>
        </div>
    </div>
</template>

<script>
  const FormLogin = () => import(/* webpackChunkName: "index" */ '@/forms/Login')
  import { DEFAULT_LOGO, logoUrlIfSet } from '@/constants/branding'

  export default {
  name: 'Login',
  components: {
      FormLogin
  },
  data () {
    return {

    }
  },
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
    loginLogo () {
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
      return this.brandingReady ? (this.loginLogo || DEFAULT_LOGO) : DEFAULT_LOGO
    }
  },
  mounted() {

  },
  methods: {

  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.login-logo-wrap {
  min-height: 48px;
}
</style>
