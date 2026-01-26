<template>
    <div v-if="services.length > 0" class="col-12 full-col-12">
        <h4 v-if="group.name !== 'Empty Group'" class="group_header mb-3 mt-4">{{group.name}}</h4>
        <div class="list-group online_list mb-4">

            <div v-for="(service, index) in services" v-bind:key="index" class="list-group-item list-group-item-action">
                <router-link class="no-decoration font-3" :to="serviceLink(service)">
                  {{service.name}}
                  <MessagesIcon :messages="service.messages"/>
                </router-link>
                <span class="badge text-uppercase float-right" :class="{'bg-success': service.online, 'bg-danger': !service.online }">
                    {{service.online ? $t('online') : $t('offline')}}
                </span>
                <span class="text-center font-2 float-right" style="margin-right: 2em" :class="{'text-muted': service.online, 'text-danger': !service.online}">{{smallText[service.id]}}</span>

            </div>

        </div>
    </div>
</template>

<script>
    const GroupServiceFailures = () => import(/* webpackChunkName: "index" */ './GroupServiceFailures');
    const IncidentsBlock = () => import(/* webpackChunkName: "index" */ './IncidentsBlock');
    const MessagesIcon = () => import(/* webpackChunkName: "index" */ '@/components/Index/MessagesIcon')

export default {
  name: 'Group',
  components: {
      IncidentsBlock,
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
        // Use last_check (actual last check time) instead of last_success (last time it was online)
        const lastCheck = service.last_check ? new Date(service.last_check).getTime() : null;
        if (lastCheck) {
          const diffInSeconds = Math.floor((now - lastCheck) / 1000);
          acc[service.id] = `last checked ${diffInSeconds} seconds ago`;
        } else {
          // Fallback if last_check is not available
          const lastSuccess = service.last_success ? new Date(service.last_success).getTime() : null;
          if (lastSuccess) {
            const diffInSeconds = Math.floor((now - lastSuccess) / 1000);
            acc[service.id] = `last checked ${diffInSeconds} seconds ago`;
          } else {
            acc[service.id] = 'never checked';
          }
        }
        return acc;
      }, {});
    },
  }
}
</script>
