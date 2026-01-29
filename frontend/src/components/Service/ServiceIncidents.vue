<template>
  <div v-if="incidents && incidents.length > 0" class="card mt-3 service-incidents-card">
    <div class="card-body pt-3 pb-3">
      <div
        v-for="incident in incidents"
        :key="incident.id"
        class="service-incident mb-3"
      >
        <h5 class="service-incident-title font-weight-bold mb-2">
          <span class="service-incident-emoji mr-2">⚠️</span>
          {{ incident.title }}
        </h5>

        <div
          v-if="incident.description"
          class="service-incident-description mb-2"
          v-html="incident.description"
        />

        <div class="service-incident-dates mb-3">
          <span>Started {{ niceDate(incident.created_at) }}</span>
          <span class="mx-2">·</span>
          <span>Last updated {{ niceDate(lastUpdatedTime(incident)) }}</span>
        </div>

        <div
          v-if="incident.updates && incident.updates.length > 0"
          class="service-incident-updates"
        >
          <div
            v-for="(update, ui) in incident.updates"
            :key="update.id || ui"
            class="service-incident-update pt-2 mt-2"
          >
            <span
              class="font-weight-bold text-capitalize"
              :class="{
                'text-success': update.type.toLowerCase() === 'resolved',
                'text-danger': update.type.toLowerCase() === 'investigating',
                'text-warning': update.type.toLowerCase() === 'update'
              }"
            >
              {{ update.type }}
            </span>
            <span class="text-muted"> - {{ update.message }}</span>
            <div class="small text-muted mt-1">
              {{ niceDate(update.created_at) }} (about {{ ago(update.created_at) }} ago)
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Api from '@/API'

export default {
  name: 'ServiceIncidents',
  props: {
    service: {
      type: Object,
      required: true,
    },
  },
  data () {
    return {
      incidents: null,
    }
  },
  async mounted () {
    await this.fetchIncidents()
  },
  methods: {
    async fetchIncidents () {
      if (!this.service || !this.service.id) return
      this.incidents = await Api.incidents_service(this.service.id)
    },
    lastUpdatedTime (incident) {
      if (incident.updates && incident.updates.length > 0) {
        const sortedUpdates = [...incident.updates].sort(
          (a, b) => new Date(b.created_at) - new Date(a.created_at)
        )
        return sortedUpdates[0].created_at
      }
      return incident.updated_at
    },
  },
}
</script>

<style scoped>
.service-incidents-card {
  border-radius: 10px;
}

.service-incident-title {
  font-size: 1.6rem; /* bumped ~8px over previous size */
}

.service-incident-emoji {
  font-size: 0.9em;
  vertical-align: middle;
}

.service-incident-description {
  font-size: 1.2rem; /* +~4px over previous size */
}

.service-incident-dates {
  font-size: 1.1rem; /* +~4px over previous size */
  color: #b0b0b0;
}

.service-incident-update {
  border-top: 1px solid rgba(255, 255, 255, 0.15);
}

.service-incident-update:last-child {
  border-bottom: none;
}
</style>

