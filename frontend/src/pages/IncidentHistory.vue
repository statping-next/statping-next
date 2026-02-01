<template>
  <div class="container col-md-7 col-sm-12 sm-container">
    <StickyHeader :visible="true"/>

    <div class="incident-history-content">
      <h1 class="incident-history-title mb-4">{{ $t('incident_history') }}</h1>

      <div v-if="loading" class="text-center py-5">
        <font-awesome-icon icon="circle-notch" class="text-dim" size="2x" spin/>
        <p class="mt-3 text-dim">{{ $t('loading') || 'Loading...' }}</p>
      </div>

      <div v-else-if="archivedIncidents.length === 0" class="incident-history-empty text-center py-5 text-dim">
        <p>{{ $t('incident_history_empty') || 'No archived incidents.' }}</p>
      </div>

      <div v-else class="incident-history-timeline">
        <IncidentHistoryBlock
          v-for="incident in archivedIncidents"
          :key="incident.id"
          :incident="incident"
          :service-name="serviceNameFor(incident)"
        />
      </div>
    </div>
  </div>
</template>

<script>
import Api from '@/API'

const StickyHeader = () => import(/* webpackChunkName: "index" */ '@/components/StickyHeader')
const IncidentHistoryBlock = () => import(/* webpackChunkName: "index" */ '@/components/Index/IncidentHistoryBlock')

export default {
  name: 'IncidentHistory',
  components: {
    StickyHeader,
    IncidentHistoryBlock
  },
  data() {
    return {
      loading: true,
      archivedIncidents: []
    }
  },
  computed: {
    services() {
      return this.$store.getters.services || []
    }
  },
  mounted() {
    this.loadArchived()
  },
  methods: {
    async loadArchived() {
      this.loading = true
      try {
        const data = await Api.incidents_archived()
        const list = Array.isArray(data) ? data : []
        this.archivedIncidents = [...list].sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
      } catch (e) {
        console.error(e)
        this.archivedIncidents = []
      }
      this.loading = false
    },
    serviceNameFor(incident) {
      const sid = Number(incident.service)
      if (!sid) return ''
      const service = this.services.find(s => s.id === sid)
      return service ? service.name : ''
    }
  }
}
</script>

<style scoped>
.incident-history-content {
  padding-top: 70px;
}
.incident-history-title {
  font-size: 1.75rem;
  font-weight: bold;
}
.incident-history-timeline {
  max-width: 100%;
}
</style>
