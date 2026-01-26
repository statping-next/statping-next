<template>
    <div class="container col-md-7 col-sm-12">
      <StickyHeader :visible="true" :admin-mode="true"/>
      <div v-if="modal" class="modal-backdrop"></div>
      <Modal/>
      <div style="padding-top: 70px;">
        <router-view :admin="admin"/>
      </div>
    </div>
</template>

<script>
  import Modal from "@/components/Elements/Modal";
  const StickyHeader = () => import(/* webpackChunkName: "index" */ '@/components/StickyHeader')

  export default {
  name: 'Dashboard',
  components: {
    Modal,
    StickyHeader,
  },
  data () {
      return {
          authenticated: false,
          loaded: false,
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
      // if (!this.user || !this.admin) {
      //   this.$router.push('/login')
      // }
    }
  }
</script>

