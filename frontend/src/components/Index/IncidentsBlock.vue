<template>
    <div v-if="incidents && incidents.length > 0" class="incidents-block">
        <div v-for="(incident, i) in incidents" :key="incident.id" class="incident-nested" role="alert">
            <h3 class="incident-title font-weight-bold">
                <span class="incident-emoji">⚠️</span> {{ incident.title }}
            </h3>
            <div v-if="incident.description" class="incident-description" v-html="incident.description"></div>
            <div class="incident-dates">
                <span class="incident-date-item">Started {{ niceDate(incident.created_at) }} ({{ aboutAgo(incident.created_at) }})</span>
                <span class="incident-date-sep">·</span>
                <span class="incident-date-item">Last updated {{ niceDate(lastUpdatedTime(incident)) }} ({{ aboutAgo(lastUpdatedTime(incident)) }})</span>
            </div>
            <div v-if="incident.updates && incident.updates.length > 0" class="incident-updates">
                <IncidentUpdate v-for="(update, ui) in incident.updates" :key="update.id || ui" :update="update" :admin="false"/>
            </div>
        </div>
    </div>
</template>

<script>
import Api from '../../API';
import IncidentUpdate from "@/components/Elements/IncidentUpdate";

export default {
  name: 'IncidentsBlock',
  components: {
    IncidentUpdate
  },
  props: {
        service: {
            type: Object,
            required: true
        }
    },
    data() {
        return {
            incidents: null,
        }
    },
    mounted () {
        this.getIncidents()
    },
    methods: {
        badgeClass(val) {
          switch (val.toLowerCase()) {
            case "resolved":
              return "badge-success"
            case "update":
              return "badge-info"
            case "investigating":
              return "badge-danger"
          }
        },
      async getIncidents() {
        this.incidents = await Api.incidents_service(this.service.id)
      },
      async incident_updates(incident) {
        return await Api.incident_updates(incident)
      },
      lastUpdatedTime(incident) {
        // If there are updates, use the most recent update's created_at
        if (incident.updates && incident.updates.length > 0) {
          // Updates are typically ordered by ID DESC, so first one is most recent
          // But to be safe, find the most recent by comparing dates
          const sortedUpdates = [...incident.updates].sort((a, b) => {
            return new Date(b.created_at) - new Date(a.created_at)
          })
          return sortedUpdates[0].created_at
        }
        // Otherwise use the incident's updated_at
        return incident.updated_at
      }
    }
}
</script>

<style scoped>
.incidents-block {
  margin-top: -16px;
  padding-top: 0;
}

.incident-nested {
  padding: 16px 16px 0;
  margin-bottom: 0;
  border-radius: 0 0 6px 6px;
  /* Lighter background so it reads as "coming out of" the service row */
  background: #f0f0f0;
  /* 12px border matching service entry background color (doubled from 6px) */
  border: 12px solid #f8f8f8;
}
/* When inside a list-group-item (Index), extend to match service entry width exactly */
.list-group-item .incident-nested {
  margin-left: -1.25rem;
  margin-right: -1.25rem;
  /* 12px border + 16px internal padding */
  padding-left: calc(12px + 16px);
  padding-right: calc(12px + 16px);
}
/* When inside a service-entry-wrapper (Group), extend to match service entry width exactly */
.service-entry-wrapper .incident-nested {
  margin-left: 0;
  margin-right: 0;
}
.dark-theme .incident-nested {
  background: #1a1d21;
  border-color: #0f1114;
}

/* Match global announcement/incident font sizes and spacing */
.incident-title {
  font-size: calc(1.25rem - 0.5px);
  margin-bottom: 1rem;
}
.incident-emoji {
  font-size: 0.9em;
  vertical-align: middle;
}

.incident-description {
  font-size: calc(1rem - 0.5px);
  margin-bottom: 0.5rem;
}
.incident-description >>> p:last-child {
  margin-bottom: 0;
}

.incident-dates {
  font-size: calc(0.875rem + 0.5px);
  color: #6c6c6c;
  margin-bottom: 0.75rem;
}
.dark-theme .incident-dates {
  color: #aaaaaa;
}
.incident-date-sep {
  margin: 0 0.35rem;
}

.incident-updates {
  border-top: 1px solid rgba(0, 0, 0, 0.08);
  padding-top: 0.5rem;
  margin-top: 0.25rem;
}
.dark-theme .incident-updates {
  border-top-color: rgba(255, 255, 255, 0.2);
}
/* Remove bottom border from last update row */
.incident-updates >>> .incident-update-row:last-child {
  border-bottom: none !important;
}
</style>
