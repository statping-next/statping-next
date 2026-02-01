<template>
  <nav class="dashboard-nav" :class="{ 'dashboard-nav--overlay': overlay, 'dashboard-nav--fixed': !overlay }">
    <router-link to="/dashboard/status" class="nav-link" :class="{ 'active': isActiveRoute('/dashboard/status') }" :data-text="$t('status')">{{ $t('status') }}</router-link>
    <router-link to="/dashboard/incidents" class="nav-link" :class="{ 'active': isActiveRoute('/dashboard/incidents') }" :data-text="$t('incidents')">{{ $t('incidents') }}</router-link>
    <router-link v-if="admin" to="/dashboard/logs" class="nav-link" :class="{ 'active': isActiveRoute('/dashboard/logs') }" :data-text="$t('logs')">{{ $t('logs') }}</router-link>
    <router-link to="/dashboard/services" class="nav-link" :class="{ 'active': isActiveRoute('/dashboard/services') }" :data-text="$t('services')">{{ $t('services') }}</router-link>
    <router-link v-if="admin" to="/dashboard/users" class="nav-link" :class="{ 'active': isActiveRoute('/dashboard/users') }" :data-text="$t('users')">{{ $t('users') }}</router-link>
    <router-link v-if="admin" to="/dashboard/settings" class="nav-link" :class="{ 'active': isActiveRoute('/dashboard/settings') }" :data-text="$t('settings')">{{ $t('settings') }}</router-link>
  </nav>
</template>

<script>
export default {
  name: 'DashboardNav',
  props: {
    overlay: {
      type: Boolean,
      default: false
    }
  },
  computed: {
    admin() {
      return this.$store.getters.admin
    }
  },
  methods: {
    isActiveRoute(path, exact = false) {
      if (exact) {
        return this.$route.path === path
      }
      return this.$route.path === path || this.$route.path.startsWith(path + '/')
    }
  }
}
</script>

<style scoped>
.dashboard-nav {
  display: flex;
  align-items: center;
  gap: 2px;
  flex-wrap: wrap;
  justify-content: center;
}

/* Overlay: fixed and centered in viewport (matches inner body width ~1020px) */
.dashboard-nav--overlay {
  position: fixed;
  left: 0;
  right: 0;
  top: 0;
  height: 72px;
  margin: 0;
  width: 100%;
  pointer-events: none;
  z-index: 1001;
  display: flex;
  align-items: center;
  justify-content: center;
}

.dashboard-nav--overlay .nav-link {
  pointer-events: auto;
}

/* Mobile: fixed below sticky header, same z-index as header */
.dashboard-nav--fixed {
  position: fixed;
  top: 72px;
  left: 0;
  right: 0;
  z-index: 1000;
  margin: 0;
  padding: 8px 1em 10px;
}

.dashboard-nav .nav-link {
  text-decoration: none;
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 14px;
  transition: all 0.3s ease;
  white-space: nowrap;
  position: relative;
  display: inline-block;
  text-align: center;
}

.dashboard-nav .nav-link::before {
  content: attr(data-text);
  font-weight: bold;
  display: block;
  height: 0;
  overflow: hidden;
  visibility: hidden;
  user-select: none;
  pointer-events: none;
  white-space: nowrap;
}

.dashboard-nav .nav-link.active {
  font-weight: bold;
}

/* Narrow (<1020px): fixed bar styling (width, link size) */
@media (max-width: 1019px) {
  .dashboard-nav--fixed {
    width: 100%;
    box-sizing: border-box;
  }

  .dashboard-nav .nav-link {
    font-size: 14px;
    padding: 4px 8px;
  }
}
</style>
