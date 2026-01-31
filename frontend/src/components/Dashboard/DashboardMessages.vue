<template>
  <div class="col-12">
    <!-- Announcements (messages) -->
    <div class="card contain-card mb-4">
      <div class="card-header d-flex justify-content-between align-items-center">
        <span>{{ $t('announcements') }}</span>
        <button v-if="$store.state.admin" @click="openCreateAnnouncement" class="btn btn-sm btn-outline-light incidents-create-btn" :title="$t('message_create')">
          <font-awesome-icon icon="plus" />
        </button>
      </div>
      <div class="card-body pt-0">
        <div v-if="messages.length === 0" class="alert alert-dark d-block mt-3 mb-0">
          No announcements. Create a global or service-specific announcement.
        </div>
        <table v-else class="table table-striped">
          <thead>
            <tr>
              <th scope="col">{{ $t('title') }}</th>
              <th scope="col" class="d-none d-md-table-cell">{{ $tc('service', 1) }}</th>
              <th scope="col" class="d-none d-md-table-cell">{{ $t('begins') }}</th>
              <th scope="col"></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="message in messages" :key="message.id">
              <td>{{ message.title }}</td>
              <td class="d-none d-md-table-cell">
                <router-link :to="serviceLink(service(message.service))">{{ serviceName(service(message.service)) }}</router-link>
              </td>
              <td class="d-none d-md-table-cell">{{ niceDate(message.start_on) }}</td>
              <td class="text-right">
                <div v-if="$store.state.admin" class="btn-group">
                  <button @click.prevent="editMessage(message)" class="btn btn-sm btn-outline-secondary">
                    <font-awesome-icon icon="edit" />
                  </button>
                  <button @click.prevent="deleteMessage(message)" class="btn btn-sm btn-danger">
                    <font-awesome-icon icon="times" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Announcement modal -->
    <div v-if="edit" class="incidents-modal-backdrop" :style="modalBackdropStyle" @click.self="closeAnnouncementModal">
      <div class="incidents-modal-dialog" :style="modalDialogStyle">
        <div class="incidents-modal-content">
          <div class="incidents-modal-header">
            <h5 class="incidents-modal-title">{{ messageToEdit.id ? $t('message_edit') : $t('message_create') }}</h5>
            <button type="button" class="incidents-modal-close" @click="closeAnnouncementModal" aria-label="Close">
              <font-awesome-icon icon="times" />
            </button>
          </div>
          <div class="incidents-modal-body">
            <FormMessage :edit="editChange" :in_message="messageToEdit" embedded />
          </div>
        </div>
      </div>
    </div>

    <!-- Incidents (service incidents with updates) -->
    <div class="card contain-card mb-4">
      <div class="card-header d-flex justify-content-between align-items-center">
        <span>{{ $t('incidents') }}</span>
        <button v-if="$store.state.admin && $store.getters.services && $store.getters.services.length" @click="openCreateIncidentModal" class="btn btn-sm btn-outline-light incidents-create-btn" :title="$t('incident_create')">
          <font-awesome-icon icon="plus" />
        </button>
      </div>
      <div class="card-body pt-0">
        <div v-if="incidents.length === 0" class="alert alert-dark d-block mt-3 mb-0">
          No incidents. Create a service incident.
        </div>

        <div v-for="incident in filteredIncidents" :key="incident.id" class="card incident-card mb-4">
          <div class="card-header incident-card-header">
            <span class="incident-title">Incident: {{ incident.title }}</span>
            <span class="incident-service-name">({{ serviceName(serviceById(incident.service)) }})</span>
            <button @click="deleteIncident(incident)" class="btn btn-sm btn-danger incident-delete-btn" :title="$t('incident_delete')">
              <font-awesome-icon icon="times" />
            </button>
          </div>
          <FormIncidentUpdates :incident="incident" @updated="onIncidentUpdated" />
          <div class="incident-meta">Created: {{ niceDate(incident.created_at) }} | Last Update: {{ niceDate(incident.updated_at) }}</div>
        </div>
      </div>
    </div>

    <!-- Create Incident modal -->
    <div v-if="showIncidentModal" class="incidents-modal-backdrop" @click.self="closeIncidentModal">
      <div class="incidents-modal-dialog" :style="modalDialogStyle">
        <div class="incidents-modal-content">
          <div class="incidents-modal-header">
            <h5 class="incidents-modal-title">{{ $t('incident_create') }}</h5>
            <button type="button" class="incidents-modal-close" @click="closeIncidentModal" aria-label="Close">
              <font-awesome-icon icon="times" />
            </button>
          </div>
          <div class="incidents-modal-body">
            <form @submit.prevent="createIncident">
              <div class="form-group row">
                <label class="col-sm-4 col-form-label">{{ $tc('service', 1) }}</label>
                <div class="col-sm-8">
                  <select v-model.number="newIncident.service" class="form-control" required>
                    <option v-for="s in $store.getters.services" :key="s.id" :value="s.id">{{ s.name }}</option>
                  </select>
                </div>
              </div>
              <div class="form-group row">
                <label class="col-sm-4 col-form-label">{{ $t('title') }}</label>
                <div class="col-sm-8">
                  <input v-model="newIncident.title" type="text" class="form-control" placeholder="Incident Title" required>
                </div>
              </div>
              <div class="form-group row">
                <label class="col-sm-4 col-form-label">{{ $t('description') }}</label>
                <div class="col-sm-8">
                  <textarea v-model="newIncident.description" rows="5" class="form-control" required></textarea>
                </div>
              </div>
              <div class="form-group row">
                <div class="col-sm-12">
                  <button type="submit" class="btn btn-primary" :disabled="!newIncident.title || !newIncident.description || !newIncident.service">
                    {{ $t('incident_create') }}
                  </button>
                </div>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  import Api from '../../API'
  const FormMessage = () => import(/* webpackChunkName: "dashboard" */ '../../forms/Message')
  const FormIncidentUpdates = () => import(/* webpackChunkName: "dashboard" */ '@/forms/IncidentUpdates')

  export default {
    name: 'DashboardMessages',
    components: { FormMessage, FormIncidentUpdates },
    data() {
      return {
        edit: false,
        messageToEdit: {},
        incidents: [],
        showIncidentModal: false,
        newIncident: {
          title: '',
          description: '',
          service: null,
        },
      }
    },
    computed: {
      /* Backdrop: full page (no style override; CSS keeps it full viewport) */
      modalBackdropStyle() {
        return {}
      },
      /* Dialog: match Announcements block / container width, centered */
      modalDialogStyle() {
        if (!this.modalContentWidth) {
          return { width: '100%', maxWidth: '100%' }
        }
        const w = `${this.modalContentWidth}px`
        return { width: w, minWidth: w, maxWidth: w }
      },
      messages() {
        return this.$store.getters.messages || []
      },
      filteredIncidents() {
        const serviceId = this.$route.query.service ? Number(this.$route.query.service) : null
        if (!serviceId) return this.incidents
        return this.incidents.filter(i => i.service === serviceId)
      },
    },
    async mounted() {
      this.updateModalWidth()
      window.addEventListener('resize', this.updateModalWidth)
      await this.loadMessages()
      await this.loadIncidents()
      const serviceId = this.$route.query.service ? Number(this.$route.query.service) : null
      if (serviceId) {
        this.newIncident.service = serviceId
        this.$nextTick(() => { this.showIncidentModal = true })
      } else if (this.$store.getters.services && this.$store.getters.services.length) {
        this.newIncident.service = this.$store.getters.services[0].id
      }
    },
    beforeDestroy() {
      window.removeEventListener('resize', this.updateModalWidth)
    },
    watch: {
      edit(visible) {
        if (visible) this.$nextTick(() => requestAnimationFrame(() => this.updateModalWidth()))
      },
      showIncidentModal(visible) {
        if (visible) this.$nextTick(() => requestAnimationFrame(() => this.updateModalWidth()))
      },
    },
    methods: {
      /* Exact same logic as StickyHeader: match .container.col-md-7 width (full main content area) */
      updateModalWidth() {
        if (window.innerWidth < 768) {
          this.modalContentWidth = null
          return
        }
        const container = document.querySelector('.container.col-md-7')
        if (container) {
          this.modalContentWidth = container.offsetWidth
        } else {
          this.modalContentWidth = Math.min(1012, window.innerWidth - 30)
        }
      },
      openCreateAnnouncement() {
        this.messageToEdit = {}
        this.edit = true
      },
      closeAnnouncementModal() {
        this.edit = false
        this.messageToEdit = {}
      },
      editChange(v) {
        this.edit = v
        if (!v) this.messageToEdit = {}
      },
      editMessage(m) {
        this.messageToEdit = m
        this.edit = true
      },
      openCreateIncidentModal() {
        if (this.$store.getters.services && this.$store.getters.services.length) {
          this.newIncident.service = this.$store.getters.services[0].id
        }
        this.newIncident.title = ''
        this.newIncident.description = ''
        this.showIncidentModal = true
      },
      closeIncidentModal() {
        this.showIncidentModal = false
      },
      service(id) {
        return this.$store.getters.serviceById(id) || {}
      },
      serviceById(id) {
        return this.$store.getters.serviceById(id) || {}
      },
      serviceName(svc) {
        if (!svc || !svc.name) return 'Global'
        return svc.name
      },
      async loadMessages() {
        const messages = await Api.messages()
        this.$store.commit('setMessages', messages)
      },
      async loadIncidents() {
        this.incidents = await Api.incidents_all()
      },
      async delete(m) {
        await Api.message_delete(m.id)
        await this.loadMessages()
      },
      deleteMessage(m) {
        const modal = {
          visible: true,
          title: this.$t('delete_announcement') || 'Delete Announcement',
          body: `Are you sure you want to delete "${m.title}"?`,
          btnColor: 'btn-danger',
          btnText: this.$t('delete') || 'Delete',
          func: () => this.delete(m),
        }
        this.$store.commit('setModal', modal)
      },
      async deleteIncident(incident) {
        const modal = {
          visible: true,
          title: this.$t('incident_delete') || 'Delete Incident',
          body: `Are you sure you want to delete incident "${incident.title}"?`,
          btnColor: 'btn-danger',
          btnText: this.$t('delete') || 'Delete',
          func: async () => {
            await Api.incident_delete(incident)
            await this.loadIncidents()
          },
        }
        this.$store.commit('setModal', modal)
      },
      async createIncident() {
        const res = await Api.incident_create(this.newIncident.service, this.newIncident)
        if (res.status === 'success') {
          this.incidents.push(res.output)
          this.newIncident = { title: '', description: '', service: this.$store.getters.services && this.$store.getters.services[0] ? this.$store.getters.services[0].id : null }
          this.closeIncidentModal()
        }
      },
      onIncidentUpdated() {
        this.loadIncidents()
      },
    },
  }
</script>

<style scoped>
  .incidents-create-btn {
    padding: 0.25rem 0.5rem;
    min-width: 32px;
  }

  .incident-card {
    border-radius: 8px;
    overflow: hidden;
  }

  .incident-card-header {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    padding: 0.75rem 1rem;
  }

  .incident-title {
    font-weight: 600;
  }

  .incident-service-name {
    color: var(--secondary, #6c757d);
    font-weight: normal;
    margin-left: 0.35rem;
  }

  .incident-delete-btn {
    margin-left: auto;
  }

  .incident-meta {
    font-size: 0.875rem;
    color: var(--secondary, #6c757d);
    padding: 0.5rem 1rem;
    border-top: 1px solid rgba(0, 0, 0, 0.1);
    background: rgba(0, 0, 0, 0.03);
  }

  /* Modals */
  .incidents-modal-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.7);
    z-index: 1050;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 1rem;
  }

  .incidents-modal-dialog {
    /* Width set via :style to match .container.col-md-7 (see updateModalWidth) */
    width: 100%;
    max-height: 90vh;
    overflow: auto;
  }

  .incidents-modal-content {
    background: var(--dark, #343a40);
    color: #e8e8e8;
    border-radius: 8px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
  }

  .incidents-modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 1rem 1.25rem;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  }

  .incidents-modal-title {
    margin: 0;
    font-size: 1.125rem;
  }

  .incidents-modal-close {
    background: transparent;
    border: none;
    color: rgba(255, 255, 255, 0.7);
    padding: 0.25rem;
    cursor: pointer;
    font-size: 1.25rem;
    line-height: 1;
  }

  .incidents-modal-close:hover {
    color: #fff;
  }

  .incidents-modal-body {
    padding: 1.25rem;
    background: #1e2125 !important;
  }

  /* Dark form controls inside modals so they don't inherit light theme */
  .incidents-modal-content .form-control,
  .incidents-modal-content .form-control:focus,
  .incidents-modal-content select.form-control {
    background-color: #25282e !important;
    border-color: rgba(255, 255, 255, 0.2);
    color: #e8e8e8;
  }
  .incidents-modal-content .form-control::placeholder {
    color: rgba(255, 255, 255, 0.4);
  }
  .incidents-modal-content label {
    color: rgba(255, 255, 255, 0.9);
  }
  .incidents-modal-content .switch label {
    color: rgba(255, 255, 255, 0.8);
  }
</style>
