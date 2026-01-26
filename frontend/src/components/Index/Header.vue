<template>
    <div id="header">
        <router-link to="/" class="header-logo-link">
            <img v-if="core.logo" id="logo" style="max-height: 200px;" class="col-12 text-center pt-4 mb-3 header-title font-6" :src="core.logo" :alt="core.name" :title="core.name">
            <h1 v-if="!core.logo" id="title" class="col-12 text-center pt-4 mt-4 mb-3 header-title font-6">{{core.name}}</h1>
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
