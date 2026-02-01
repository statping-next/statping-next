<template>
  <div class="card shadow mb-4 global-incident-block incident-history-block" role="alert">
    <div class="card-body">
      <h3 class="mb-3 font-weight-bold global-incident-title">
        <span class="global-incident-emoji">⚠️</span> {{ incident.title }}<span class="incident-title-service"> (<template v-if="isServiceLink"><router-link :to="serviceLink" class="incident-title-service-link">{{ serviceLabel }}</router-link></template><template v-else>{{ serviceLabel }}</template>)</span>
      </h3>
      <div v-if="incident.description" class="global-incident-description" v-html="incident.description"></div>
      <div class="global-incident-dates">
        <span class="global-incident-date-item">Started {{ niceDate(incident.created_at) }} ({{ aboutAgo(incident.created_at) }})</span>
        <span class="global-incident-date-sep">·</span>
        <span class="global-incident-date-item">Last updated {{ niceDate(lastUpdatedTime(incident)) }} ({{ aboutAgo(lastUpdatedTime(incident)) }})</span>
      </div>
      <div v-if="incident.updates && incident.updates.length > 0" class="global-incident-updates">
        <IncidentUpdate v-for="(update, ui) in incident.updates" :key="update.id || ui" :update="update" :admin="false"/>
      </div>
    </div>
  </div>
</template>

<script>
import IncidentUpdate from '@/components/Elements/IncidentUpdate'

export default {
  name: 'IncidentHistoryBlock',
  components: {
    IncidentUpdate
  },
  props: {
    incident: {
      type: Object,
      required: true
    },
    serviceName: {
      type: String,
      default: ''
    }
  },
  computed: {
    serviceLabel() {
      return this.serviceName || 'Global'
    },
    isServiceLink() {
      const sid = Number(this.incident.service)
      return sid > 0
    },
    serviceLink() {
      return { path: '/service/' + this.incident.service }
    }
  },
  methods: {
    lastUpdatedTime(incident) {
      if (incident.updates && incident.updates.length > 0) {
        const sorted = [...incident.updates].sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
        return sorted[0].created_at
      }
      return incident.updated_at
    }
  }
}
</script>

<style scoped>
/* Match GlobalIncidentBlock layout and styling */
.incident-history-block.global-incident-block {
  width: 100%;
  max-width: 100%;
  box-sizing: border-box;
  transition: background-color 0.2s ease, box-shadow 0.2s ease;
}
.incident-history-block.global-incident-block .card-body {
  padding: 1rem 1rem 0.25rem;
  box-sizing: border-box;
}
.incident-history-block.global-incident-block:hover {
  background-color: #f8f8f8;
  box-shadow: 0 4px 12px 1px rgba(0, 0, 0, 0.12);
}
.dark-theme .incident-history-block.global-incident-block:hover {
  background-color: #1a1d21 !important;
  box-shadow: 0 4px 12px 1px rgba(0, 0, 0, 0.35);
}
.incident-history-block .global-incident-title {
  font-size: calc(1.25rem - 0.5px);
}
.incident-history-block .global-incident-emoji {
  font-size: 0.9em;
  vertical-align: middle;
}
/* Match admin incident header service name: normal weight, smaller size, secondary color */
.incident-history-block .incident-title-service {
  font-weight: normal;
  font-size: 0.9em;
  color: #5a6268;
  margin-left: 0.35rem;
}
.dark-theme .incident-history-block .incident-title-service {
  color: #aaaaaa;
}
.incident-history-block .incident-title-service-link {
  text-decoration: none;
  color: #5a6268;
}
.dark-theme .incident-history-block .incident-title-service-link {
  color: #aaaaaa;
}
.incident-history-block .incident-title-service-link:hover {
  text-decoration: underline;
}
.incident-history-block .global-incident-description {
  font-size: calc(1rem - 0.5px);
  margin-bottom: 0.75rem;
}
.incident-history-block .global-incident-description >>> p:last-child {
  margin-bottom: 0;
}
.incident-history-block .global-incident-dates {
  font-size: calc(0.875rem + 0.5px);
  color: var(--secondary, #6c6c6c);
  margin-bottom: 0.75rem;
}
.dark-theme .incident-history-block .global-incident-dates {
  color: #aaaaaa;
}
.incident-history-block .global-incident-date-sep {
  margin: 0 0.35rem;
}
.incident-history-block .global-incident-updates {
  border-top: 1px solid rgba(0, 0, 0, 0.08);
  padding-top: 0.75rem;
  margin-top: 0.25rem;
}
.dark-theme .incident-history-block .global-incident-updates {
  border-top-color: rgba(255, 255, 255, 0.2);
}
.incident-history-block .global-incident-updates >>> .incident-update-row:last-child {
  border-bottom: none !important;
}
</style>
