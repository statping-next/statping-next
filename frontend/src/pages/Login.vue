<template>
    <div class="offset-md-3 offset-lg-4 offset-0 col-lg-4 col-md-6 mt-5">

      <div class="offset-1 offset-lg-2 col-lg-8 col-10 mb-4 mb-md-3 login-logo-wrap">
        <div v-if="!brandingReady" class="login-logo-placeholder" aria-hidden="true"></div>
        <img v-else alt="Statping-ng Login" class="embed-responsive" :src="loginLogo || DEFAULT_LOGO">
      </div>

      <div class="login_container col-12 p-4">
            <FormLogin/>
        </div>
    </div>
</template>

<script>
  const FormLogin = () => import(/* webpackChunkName: "index" */ '@/forms/Login')
  import { DEFAULT_LOGO } from '@/constants/branding'

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
      // Use theme-specific logo (same logic as Header / StickyHeader)
      if (this.darkTheme) {
        if (this.core.logo_dark) return this.core.logo_dark
        if (this.core.logo_light) return this.core.logo_light
      } else {
        if (this.core.logo_light) return this.core.logo_light
      }
      if (this.core.logo) return this.core.logo
      return null
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
.login-logo-placeholder {
  width: 100%;
  height: 48px;
  background: transparent;
}
</style>
