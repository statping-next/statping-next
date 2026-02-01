<template>
    <div class="container col-md-7 col-sm-12">
      <div class="dashboard-header-wrapper">
        <StickyHeader :visible="true" :admin-mode="true"/>
        <DashboardNav :overlay="!isNarrow"/>
      </div>
      <div v-if="modal" class="modal-backdrop"></div>
      <Modal/>
      <div class="content-below-header">
        <router-view :admin="admin"/>
      </div>
    </div>
</template>

<script>
  import Modal from "@/components/Elements/Modal";
  import DashboardNav from "@/components/Dashboard/DashboardNav";
  const StickyHeader = () => import(/* webpackChunkName: "index" */ '@/components/StickyHeader')

  export default {
  name: 'Dashboard',
  components: {
    Modal,
    StickyHeader,
    DashboardNav,
  },
  data () {
      return {
          authenticated: false,
          loaded: false,
          isNarrow: false,
      }
  },
    computed: {
      modal() {
        return this.$store.getters.modal.visible
      },
      admin() {
        return this.$store.getters.admin
      },
      user() {
        return this.$store.getters.user
      }
    },
    mounted() {
      this.updateNarrow()
      window.addEventListener('resize', this.updateNarrow)
    },
    beforeDestroy() {
      window.removeEventListener('resize', this.updateNarrow)
    },
    methods: {
      updateNarrow() {
        this.isNarrow = typeof window !== 'undefined' && window.innerWidth < 1020
      }
    }
  }
</script>

<style scoped>
.dashboard-header-wrapper {
  position: relative;
  min-height: 52px;
}
@media (max-width: 1019px) {
  .dashboard-header-wrapper {
    min-height: 116px;
  }
}
</style>

