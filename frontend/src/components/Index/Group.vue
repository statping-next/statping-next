<template>
    <div v-if="services.length > 0" class="col-12 full-col-12">
        <h4 v-if="group.name !== 'Empty Group'" class="group_header mb-3 mt-4">{{group.name}}</h4>
        <div v-for="(service, index) in services" v-bind:key="service.id" class="service-entry-wrapper">
            <div class="list-group online_list">
                <router-link :to="serviceLink(service)" class="list-group-item list-group-item-action" style="display: flex; align-items: center; justify-content: space-between; text-decoration: none; height: 3.5rem; transform: none !important; margin: 0 !important; z-index: 0 !important; box-shadow: none !important; position: relative !important;">
                    <span class="no-decoration font-3" style="display: flex; align-items: center;">
                      {{service.name}}
                      <MessagesIcon :messages="service.messages"/>
                    </span>
                    <div style="display: flex; align-items: center; gap: 1em;">
                      <span v-if="smallText[service.id].downFor" class="text-center font-2" style="display: flex; align-items: baseline;">
                        <span v-if="smallText[service.id].downFor" class="text-danger">{{smallText[service.id].downFor}}</span>
                      </span>
                      <span class="text-center font-2" style="display: flex; align-items: baseline; margin-right: 1em;">
                        <span :class="{'text-muted': service.online, 'text-white': !service.online}">{{smallText[service.id].lastChecked}}</span>
                      </span>
                      <span class="badge text-uppercase" :class="{'bg-success': service.online, 'bg-danger': !service.online }" style="min-width: 70px; text-align: center;">
                          {{service.online ? $t('online') : $t('offline')}}
                      </span>
                    </div>
                </router-link>
            </div>
            <ServiceAnnouncements :service="service"/>
            <IncidentsBlock :service="service"/>
        </div>
    </div>
</template>

<script>
    const GroupServiceFailures = () => import(/* webpackChunkName: "index" */ './GroupServiceFailures');
    const IncidentsBlock = () => import(/* webpackChunkName: "index" */ './IncidentsBlock');
    const ServiceAnnouncements = () => import(/* webpackChunkName: "index" */ './ServiceAnnouncements');
    const MessagesIcon = () => import(/* webpackChunkName: "index" */ '@/components/Index/MessagesIcon')

export default {
  name: 'Group',
  components: {
      IncidentsBlock,
      ServiceAnnouncements,
      GroupServiceFailures,
      MessagesIcon
  },
  props: {
    group: {
      type: Object,
      required: true,
    }
  },
  computed: {
    services() {
      return this.$store.getters.servicesInGroup(this.group.id)
    },
    smallText() {
      const now = Date.now();
      return this.services.reduce((acc, service) => {
        // Helper function to safely parse a date string
        const parseDate = (dateStr) => {
          if (!dateStr) return null;
          // Handle zero time or invalid dates
          if (dateStr === '0001-01-01T00:00:00Z' || dateStr === '' || dateStr === null) {
            return null;
          }
          const date = new Date(dateStr);
          // Check if date is valid (not NaN and not in the year 1)
          if (isNaN(date.getTime()) || date.getFullYear() < 1970) {
            return null;
          }
          return date.getTime();
        };

        // Use last_check (actual last check time) instead of last_success (last time it was online)
        const lastCheck = parseDate(service.last_check);
        let lastCheckedText = '';

        if (lastCheck && lastCheck > 0) {
          const diffInSeconds = Math.floor((now - lastCheck) / 1000);
          // Sanity check: if the difference is more than 10 years, something is wrong
          if (diffInSeconds > 0 && diffInSeconds < 315360000) {
            lastCheckedText = `last checked ${this.formatDuration(diffInSeconds)} ago`;
          } else {
            // Fallback to last_success if last_check is invalid
            const lastSuccess = parseDate(service.last_success);
            if (lastSuccess && lastSuccess > 0) {
              const diffInSeconds = Math.floor((now - lastSuccess) / 1000);
              if (diffInSeconds > 0 && diffInSeconds < 315360000) {
                lastCheckedText = `last checked ${this.formatDuration(diffInSeconds)} ago`;
              } else {
                lastCheckedText = 'never checked';
              }
            } else {
              lastCheckedText = 'never checked';
            }
          }
        } else {
          // Fallback if last_check is not available
          const lastSuccess = parseDate(service.last_success);
          if (lastSuccess && lastSuccess > 0) {
            const diffInSeconds = Math.floor((now - lastSuccess) / 1000);
            if (diffInSeconds > 0 && diffInSeconds < 315360000) {
              lastCheckedText = `last checked ${this.formatDuration(diffInSeconds)} ago`;
            } else {
              lastCheckedText = 'never checked';
            }
          } else {
            lastCheckedText = 'never checked';
          }
        }

        // For offline services, show "down for" duration before "last checked"
        // "down for" = time since downtime started (first failure that started the downtime)
        // "last checked" = time since most recent check (regardless of status)
        // Check if service is offline (handle both boolean false and string "false")
        const isOffline = service.online === false || service.online === 'false' || service.online === 0;
        const downtimeStarted = parseDate(service.downtime_started);

        if (isOffline && downtimeStarted && downtimeStarted > 0) {
          const downForSeconds = Math.floor((now - downtimeStarted) / 1000);
          // Sanity check for downForSeconds too
          if (downForSeconds > 0 && downForSeconds < 315360000) {
            acc[service.id] = {
              downFor: `down for ${this.formatDuration(downForSeconds)}`,
              lastChecked: lastCheckedText
            };
          } else {
            acc[service.id] = {
              downFor: '',
              lastChecked: lastCheckedText
            };
          }
        } else {
          acc[service.id] = {
            downFor: '',
            lastChecked: lastCheckedText
          };
        }

        return acc;
      }, {});
    },
  }
}
</script>

<style scoped>
.service-entry-wrapper {
  margin-bottom: 6px;
}
.service-entry-wrapper:last-child {
  margin-bottom: 0;
}

/* Treat service header + incident as a single hoverable row */
.service-entry-wrapper:hover .list-group-item {
  background-color: #ffffff !important;
}
.dark-theme .service-entry-wrapper:hover .list-group-item {
  background-color: #202328 !important;
}

/* On hover, brighten the dark outer border (announcements and incidents) */
.service-entry-wrapper:hover >>> .incident-nested,
.service-entry-wrapper:hover >>> .service-announcement-nested {
  border-color: #ffffff;
}
.dark-theme .service-entry-wrapper:hover >>> .incident-nested,
.dark-theme .service-entry-wrapper:hover >>> .service-announcement-nested {
  border-color: #202328;
}
</style>
