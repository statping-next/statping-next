<template>
  <div class="card shadow mb-4 global-incident-block" role="alert">
    <div class="card-body pb-2">
      <h3 class="mb-3 font-weight-bold global-incident-title">
        <span class="global-incident-emoji">⚠️</span> {{ incident.title }}
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
  name: 'GlobalIncidentBlock',
  components: {
    IncidentUpdate
  },
  props: {
    incident: {
      type: Object,
      required: true
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
.global-incident-block {
  /* Same card size/prominence as MessageBlock (global announcements) */
}

.global-incident-title {
  font-size: 1.25rem;
}

.global-incident-emoji {
  font-size: 0.9em;
  vertical-align: middle;
}

.global-incident-description {
  font-size: 1rem;
  margin-bottom: 0.75rem;
}
.global-incident-description >>> p:last-child {
  margin-bottom: 0;
}

.global-incident-dates {
  font-size: 0.875rem;
  color: var(--secondary, #6c6c6c);
  margin-bottom: 0.75rem;
}
.dark-theme .global-incident-dates {
  color: #aaaaaa;
}

.global-incident-date-sep {
  margin: 0 0.35rem;
}

/* Incident-style updates list (colored types) - same as IncidentsBlock */
.global-incident-updates {
  border-top: 1px solid rgba(0, 0, 0, 0.08);
  padding-top: 0.75rem;
  margin-top: 0.25rem;
}
.dark-theme .global-incident-updates {
  border-top-color: rgba(255, 255, 255, 0.2);
}
.global-incident-updates >>> .incident-update-row:last-child {
  border-bottom: none !important;
}
</style>
