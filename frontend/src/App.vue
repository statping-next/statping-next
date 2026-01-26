<template>
  <div id="app">
    <router-view/>
    <Footer v-if="$route.path !== '/setup'"/>
  </div>
</template>

<script>
  const Footer = () => import(/* webpackChunkName: "index" */ "./components/Index/Footer");

  export default {
    name: 'app',
    components: {
      Footer
    },
    data() {
      return {
        loaded: false,
        version: "",
        refreshTimer: null
      }
    },
    computed: {
      core() {
        return this.$store.getters.core
      },
      refreshInterval() {
        return this.$store.getters.refreshInterval
      },
      isAdmin() {
        return this.$store.state.admin
      }
    },
    watch: {
      refreshInterval(newInterval, oldInterval) {
        this.setupRefreshTimer()
      },
      '$route'(to, from) {
        // Restart timer when route changes
        this.setupRefreshTimer()
      }
    },
    created() {
      // Apply theme immediately on app creation
      const darkTheme = localStorage.getItem('darkTheme') === 'true'
      if (darkTheme) {
        document.body.classList.add('dark-theme')
      } else {
        document.body.classList.add('light-theme')
      }
    },
    async beforeMount() {
      await this.$store.dispatch('loadCore')

      this.$i18n.locale = this.core.language || "en";
      // this.$i18n.locale = "ru";

      if (!this.core.setup) {
        this.$router.push('/setup')
      }
      if (this.$route.path !== '/setup') {
        if (this.$store.state.admin) {
          await this.$store.dispatch('loadAdmin')
        } else {
          await this.$store.dispatch('loadRequired')
        }
        this.loaded = true
      }
    },
    methods: {
      setupRefreshTimer() {
        // Clear existing timer
        if (this.refreshTimer) {
          clearInterval(this.refreshTimer)
          this.refreshTimer = null
        }

        // Set up new timer if interval is active
        if (this.refreshInterval > 0 && this.$route.path !== '/setup') {
          this.refreshTimer = setInterval(() => {
            // Refresh data based on admin status
            if (this.isAdmin) {
              this.$store.dispatch('loadAdmin')
            } else {
              this.$store.dispatch('loadRequired')
            }
          }, this.refreshInterval * 1000)
        }
      }
    },
    async mounted() {
      if (this.$route.path !== '/setup') {
        if (this.$store.state.admin) {
          this.logged_in = true
          // await this.$store.dispatch('loadAdmin')
        }
      }
      // Set up initial refresh timer
      this.setupRefreshTimer()
    },
    beforeDestroy() {
      // Clean up timer on component destroy
      if (this.refreshTimer) {
        clearInterval(this.refreshTimer)
        this.refreshTimer = null
      }
    }
  }
</script>

<style lang="scss">
    @import "./assets/css/bootstrap.min.css";
    @import "./assets/scss/index";
</style>
