<template>
    <div class="container col-md-7 col-sm-12 sm-container">

      <StickyHeader :visible="showStickyHeader"/>
      <Header/>

      <div v-if="!loaded" class="row mt-5 mb-5">
        <div class="col-12 mt-5 mb-2 text-center">
          <font-awesome-icon icon="circle-notch" class="text-dim" size="2x" spin/>
        </div>
        <div class="col-12 text-center mt-3 mb-3">
          <span class="text-dim">{{loading_text}}</span>
        </div>
      </div>

      <div class="col-12 full-col-12">
          <MessageBlock v-for="message in messages" v-bind:key="message.id" :message="message" />
          <GlobalIncidentBlock v-for="incident in globalIncidents" v-bind:key="incident.id" :incident="incident" />
      </div>

      <div class="col-12 full-col-12">
          <div v-for="service in services_no_group" v-bind:key="service.id" class="list-group online_list service-entry">
              <div class="list-group-item list-group-item-action">
                  <router-link class="no-decoration font-3" :to="serviceLink(service)">
                    {{service.name}}
                    <MessagesIcon :messages="service.messages"/>
                  </router-link>
                  <span class="badge float-right" :class="{'bg-success': service.online, 'bg-danger': !service.online }">{{service.online ? "ONLINE" : "OFFLINE"}}</span>
                  <GroupServiceFailures :service="service"/>
                  <IncidentsBlock :service="service"/>
              </div>
          </div>
      </div>

      <Group v-for="group in groups" v-bind:key="group.id" :group=group />

    </div>
</template>

<script>
import Api from "@/API";

const Group = () => import(/* webpackChunkName: "index" */ '@/components/Index/Group')
const Header = () => import(/* webpackChunkName: "index" */ '@/components/Index/Header')
const StickyHeader = () => import(/* webpackChunkName: "index" */ '@/components/StickyHeader')
const MessageBlock = () => import(/* webpackChunkName: "index" */ '@/components/Index/MessageBlock')
const GlobalIncidentBlock = () => import(/* webpackChunkName: "index" */ '@/components/Index/GlobalIncidentBlock')
const ServiceBlock = () => import(/* webpackChunkName: "index" */ '@/components/Service/ServiceBlock')
const GroupServiceFailures = () => import(/* webpackChunkName: "index" */ '@/components/Index/GroupServiceFailures')
const IncidentsBlock = () => import(/* webpackChunkName: "index" */ '@/components/Index/IncidentsBlock')
const MessagesIcon = () => import(/* webpackChunkName: "index" */ '@/components/Index/MessagesIcon')

export default {
    name: 'Index',
    components: {
      IncidentsBlock,
      GroupServiceFailures,
      ServiceBlock,
      MessageBlock,
      GlobalIncidentBlock,
      MessagesIcon,
      Group,
      Header,
      StickyHeader
    },
    data() {
        return {
            logged_in: false,
            showStickyHeader: false,
        }
    },
    computed: {
      loading_text() {
        if (this.$store.getters.groups.length === 0) {
          return "Loading Groups"
        } else if (this.$store.getters.services.length === 0) {
          return "Loading Services"
        } else if (this.$store.getters.messages == null) {
          return "Loading Announcements"
        }
      },
      loaded() {
        return this.$store.getters.services.length !== 0
      },
        core() {
          return this.$store.getters.core
        },
        messages() {
            const msgs = this.$store.getters.messages
            if (!msgs || !Array.isArray(msgs)) {
                return []
            }
            return msgs.filter(m => this.inRange(m) && m.service === 0)
        },
        globalIncidents() {
            return this.$store.getters.globalIncidents || []
        },
        groups() {
            return this.$store.getters.groupsInOrder
        },
        services() {
            return this.$store.getters.servicesInOrder
        },
        services_no_group() {
            return this.$store.getters.servicesNoGroup
        }
    },
    methods: {
        async checkLogin() {
          const token = this.$cookies.get('statping_auth')
          if (!token) {
            this.$store.commit('setLoggedIn', false)
            return
          }
          try {
            const jwt = await Api.check_token(token)
            this.$store.commit('setAdmin', jwt.admin)
            if (jwt.username) {
              this.$store.commit('setLoggedIn', true)
            }
          } catch (e) {
            console.error(e)
          }
        },
        inRange(message) {
            return this.isBetween(this.now(), message.start_on, message.start_on === message.end_on ? this.maxDate().toISOString() : message.end_on)
        },
        handleScroll() {
            // Show sticky header when logo would be fully covered (around 140px)
            this.showStickyHeader = window.scrollY > 140
        }
    },
    mounted() {
        window.addEventListener('scroll', this.handleScroll)
    },
    beforeDestroy() {
        window.removeEventListener('scroll', this.handleScroll)
    }
}
</script>

<style scoped>
.service-entry {
  margin-bottom: 6px;
}

/* Treat service header + incident as a single hoverable row */
.service-entry:hover .list-group-item {
  background-color: #ffffff !important;
}
.dark-theme .service-entry:hover .list-group-item {
  background-color: #202328 !important;
}

/* On hover, brighten the dark outer border but keep inner incident background */
.service-entry:hover >>> .incident-nested {
  border-color: #ffffff;
}
.dark-theme .service-entry:hover >>> .incident-nested {
  border-color: #202328;
}
</style>
