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
          No incidents. Create a global or service-specific incident.
        </div>

        <template v-else>
          <div v-if="activeIncidents.length === 0" class="alert alert-dark d-block mt-3 mb-0">
            {{ $t('no_active_incidents') }}
          </div>
          <div v-for="incident in activeIncidents" :key="incident.id" class="card incident-card mb-4">
            <div class="card-header incident-card-header">
              <div class="d-flex align-items-center flex-wrap incident-header-left">
                <span class="incident-title">Incident: {{ incident.title }}</span>
                <span class="incident-service-name">({{ serviceName(serviceById(incident.service)) }})</span>
                <span class="ml-2 incident-auto-archive-icon" :title="incident.auto_archive_enabled ? $t('auto_archive_after_resolution') + (incident.auto_archive_delay_minutes ? ' (' + (incident.auto_archive_delay_minutes >= 60 ? Math.floor(incident.auto_archive_delay_minutes / 60) + 'h' : incident.auto_archive_delay_minutes + 'm') + ')' : ' (immediate)') : $t('auto_archive_disabled')">
                  <font-awesome-icon icon="clock" :class="{ 'incident-auto-archive-off': !incident.auto_archive_enabled }" />
                </span>
              </div>
              <div v-if="$store.state.admin" class="btn-group incident-header-actions">
                <button @click="toggleArchiveIncident(incident)" class="btn btn-sm btn-outline-secondary" :title="$t('incident_archive')">
                  <font-awesome-icon icon="archive" />
                </button>
                <button @click="editIncident(incident)" class="btn btn-sm btn-outline-secondary" :title="$t('incident_edit')">
                  <font-awesome-icon icon="edit" />
                </button>
                <button @click="deleteIncident(incident)" class="btn btn-sm btn-danger" :title="$t('incident_delete')">
                  <font-awesome-icon icon="times" />
                </button>
              </div>
            </div>
            <FormIncidentUpdates :incident="incident" @updated="onIncidentUpdated" />
            <div class="incident-meta">Created: {{ niceDate(incident.created_at) }} | Last Update: {{ niceDate(incident.updated_at) }}</div>
          </div>

          <div v-if="archivedIncidentsList.length > 0" class="archived-incidents-drawer mt-4">
            <button type="button" class="archived-incidents-toggle btn btn-link p-0 text-left d-flex align-items-center" @click="archivedDrawerOpen = !archivedDrawerOpen" :aria-expanded="archivedDrawerOpen">
              <font-awesome-icon :icon="archivedDrawerOpen ? 'chevron-down' : 'chevron-right'" class="archived-incidents-chevron mr-2" />
              <span class="archived-incidents-label">{{ $t('archived_incidents') }}</span>
            </button>
            <div v-show="archivedDrawerOpen" class="archived-incidents-content mt-2">
              <div v-for="incident in archivedIncidentsList" :key="incident.id" class="card incident-card mb-4">
                <div class="card-header incident-card-header">
                  <div class="d-flex align-items-center flex-wrap incident-header-left">
                    <span class="incident-title">Incident: {{ incident.title }}</span>
                    <span class="incident-service-name">({{ serviceName(serviceById(incident.service)) }})</span>
                    <span class="badge badge-secondary ml-2">{{ $t('incident_archived') }}</span>
                    <span class="ml-2 incident-auto-archive-icon" :title="incident.auto_archive_enabled ? $t('auto_archive_after_resolution') + (incident.auto_archive_delay_minutes ? ' (' + (incident.auto_archive_delay_minutes >= 60 ? Math.floor(incident.auto_archive_delay_minutes / 60) + 'h' : incident.auto_archive_delay_minutes + 'm') + ')' : ' (immediate)') : $t('auto_archive_disabled')">
                      <font-awesome-icon icon="clock" :class="{ 'incident-auto-archive-off': !incident.auto_archive_enabled }" />
                    </span>
                  </div>
                  <div v-if="$store.state.admin" class="btn-group incident-header-actions">
                    <button @click="toggleArchiveIncident(incident)" class="btn btn-sm btn-outline-secondary" :title="$t('incident_unarchive')">
                      <font-awesome-icon icon="folder-open" />
                    </button>
                    <button @click="editIncident(incident)" class="btn btn-sm btn-outline-secondary" :title="$t('incident_edit')">
                      <font-awesome-icon icon="edit" />
                    </button>
                    <button @click="deleteIncident(incident)" class="btn btn-sm btn-danger" :title="$t('incident_delete')">
                      <font-awesome-icon icon="times" />
                    </button>
                  </div>
                </div>
                <FormIncidentUpdates :incident="incident" @updated="onIncidentUpdated" />
                <div class="incident-meta">Created: {{ niceDate(incident.created_at) }} | Last Update: {{ niceDate(incident.updated_at) }}</div>
              </div>
            </div>
          </div>
        </template>
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
                  <select v-model.number="newIncident.service" class="form-control" required :key="'incident-service-' + showIncidentModal">
                    <option :value="0">{{ $t('global_incident') }}</option>
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
                <div class="col-sm-8 offset-sm-4">
                  <div class="custom-control custom-checkbox">
                    <input id="new-auto-archive" v-model="newIncident.auto_archive_enabled" type="checkbox" class="custom-control-input">
                    <label class="custom-control-label" for="new-auto-archive">{{ $t('auto_archive_after_resolution') }}</label>
                  </div>
                </div>
              </div>
              <div v-if="newIncident.auto_archive_enabled" class="form-group row">
                <label class="col-sm-4 col-form-label">{{ $t('auto_archive_delay_hint') }}</label>
                <div class="col-sm-8 form-inline">
                  <input v-model.number="newIncident.auto_archive_hours" type="number" min="0" max="72" class="form-control form-control-sm mr-2" style="width: 4rem">
                  <span>h</span>
                </div>
              </div>
              <div class="form-group row">
                <div class="col-sm-12">
                  <button type="submit" class="btn btn-primary" :disabled="!newIncident.title || !newIncident.description">
                    {{ $t('incident_create') }}
                  </button>
                </div>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>

    <!-- Edit Incident modal -->
    <div v-if="showEditIncidentModal && editingIncident" class="incidents-modal-backdrop" @click.self="closeEditIncidentModal">
      <div class="incidents-modal-dialog" :style="modalDialogStyle">
        <div class="incidents-modal-content">
          <div class="incidents-modal-header">
            <h5 class="incidents-modal-title">{{ $t('incident_edit') }}</h5>
            <button type="button" class="incidents-modal-close" @click="closeEditIncidentModal" aria-label="Close">
              <font-awesome-icon icon="times" />
            </button>
          </div>
          <div class="incidents-modal-body">
            <form @submit.prevent="saveIncidentEdit">
              <div class="form-group row">
                <label class="col-sm-4 col-form-label">{{ $t('title') }}</label>
                <div class="col-sm-8">
                  <input v-model="editingIncident.title" type="text" class="form-control" required>
                </div>
              </div>
              <div class="form-group row">
                <label class="col-sm-4 col-form-label">{{ $t('description') }}</label>
                <div class="col-sm-8">
                  <textarea v-model="editingIncident.description" rows="5" class="form-control" required></textarea>
                </div>
              </div>
              <div class="form-group row">
                <div class="col-sm-8 offset-sm-4">
                  <div class="custom-control custom-checkbox">
                    <input id="edit-auto-archive" v-model="editingIncident.auto_archive_enabled" type="checkbox" class="custom-control-input">
                    <label class="custom-control-label" for="edit-auto-archive">{{ $t('auto_archive_after_resolution') }}</label>
                  </div>
                </div>
              </div>
              <div v-if="editingIncident.auto_archive_enabled" class="form-group row">
                <label class="col-sm-4 col-form-label">{{ $t('auto_archive_delay_hint') }}</label>
                <div class="col-sm-8 form-inline">
                  <input v-model.number="editingIncident.auto_archive_hours" type="number" min="0" max="72" class="form-control form-control-sm mr-2" style="width: 4rem" :key="'edit-hours-' + (editingIncident.id || '')">
                  <span>h</span>
                </div>
              </div>
              <div class="form-group row">
                <div class="col-sm-12">
                  <button type="submit" class="btn btn-primary">{{ $t('save') }}</button>
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
        showEditIncidentModal: false,
        editingIncident: null,
        archivedDrawerOpen: false,
        newIncident: {
          title: '',
          description: '',
          service: 0,
          auto_archive_enabled: false,
          auto_archive_hours: 12,
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
        const list = serviceId ? this.incidents.filter(i => i.service === serviceId) : this.incidents
        return [...list].sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
      },
      activeIncidents() {
        return this.filteredIncidents.filter(i => !i.archived)
      },
      archivedIncidentsList() {
        return this.filteredIncidents.filter(i => i.archived)
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
      } else {
        this.newIncident.service = 0
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
      showEditIncidentModal(visible) {
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
        this.newIncident = {
          title: '',
          description: '',
          service: 0,
          auto_archive_enabled: false,
          auto_archive_hours: 12,
        }
        this.showIncidentModal = true
        this.$nextTick(() => { this.newIncident.service = 0 })
      },
      closeIncidentModal() {
        this.showIncidentModal = false
      },
      editIncident(incident) {
        const minutes = incident.auto_archive_delay_minutes
        const hours = (typeof minutes === 'number' && !Number.isNaN(minutes))
          ? Math.min(72, Math.max(0, Math.floor(minutes / 60)))
          : 0
        this.editingIncident = {
          id: incident.id,
          title: incident.title,
          description: incident.description,
          auto_archive_enabled: Boolean(incident.auto_archive_enabled),
          auto_archive_hours: hours,
        }
        this.showEditIncidentModal = true
      },
      closeEditIncidentModal() {
        this.showEditIncidentModal = false
        this.editingIncident = null
      },
      async saveIncidentEdit() {
        if (!this.editingIncident) return
        const incident = this.incidents.find(i => i.id === this.editingIncident.id)
        if (!incident) return
        const delayMinutes = Math.min(72 * 60, Math.max(0, (this.editingIncident.auto_archive_hours || 0) * 60))
        const res = await Api.incident_update(incident, {
          title: this.editingIncident.title,
          description: this.editingIncident.description,
          auto_archive_enabled: this.editingIncident.auto_archive_enabled,
          auto_archive_delay_minutes: this.editingIncident.auto_archive_enabled ? delayMinutes : 0,
        })
        if (res.status === 'success') {
          await this.loadIncidents()
          this.closeEditIncidentModal()
        }
      },
      async toggleArchiveIncident(incident) {
        const res = await Api.incident_archive(incident, !incident.archived)
        if (res.status === 'success') {
          incident.archived = !incident.archived
          await this.loadIncidents()
        }
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
        const isGlobal = this.newIncident.service === 0
        const delayMinutes = this.newIncident.auto_archive_enabled
          ? Math.min(72 * 60, Math.max(0, (this.newIncident.auto_archive_hours || 0) * 60))
          : 0
        const payload = {
          title: this.newIncident.title,
          description: this.newIncident.description,
          auto_archive_enabled: this.newIncident.auto_archive_enabled || false,
          auto_archive_delay_minutes: delayMinutes,
        }
        const res = isGlobal
          ? await Api.incident_create_global({ ...payload, service: 0 })
          : await Api.incident_create(this.newIncident.service, { ...this.newIncident, ...payload })
        if (res.status === 'success') {
          const incidents = await Api.incidents_all()
          this.incidents = incidents
          this.$store.commit('setIncidents', incidents)
          this.newIncident = { title: '', description: '', service: 0, auto_archive_enabled: false, auto_archive_hours: 12 }
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
    justify-content: space-between;
    flex-wrap: wrap;
    padding: 0.75rem 1rem;
    gap: 0.5rem;
  }

  .incident-header-left {
    flex: 0 1 auto;
  }

  .incident-header-actions {
    flex: 0 0 auto;
    margin-left: auto;
  }

  .incident-title {
    font-weight: 600;
  }

  .incident-service-name {
    color: var(--secondary, #6c757d);
    font-weight: normal;
    margin-left: 0.35rem;
  }

  .incident-auto-archive-icon {
    color: var(--secondary, #6c757d);
  }

  .incident-auto-archive-icon .incident-auto-archive-off {
    opacity: 0.4;
  }

  .incident-meta {
    font-size: 0.875rem;
    color: var(--secondary, #6c757d);
    padding: 0.5rem 1rem;
    border-top: 1px solid rgba(0, 0, 0, 0.1);
    background: rgba(0, 0, 0, 0.03);
  }

  .archived-incidents-drawer {
    border-top: 1px solid rgba(0, 0, 0, 0.1);
    padding-top: 0.5rem;
    padding-bottom: 1.5rem;
  }

  .archived-incidents-toggle {
    text-decoration: none;
    color: var(--secondary, #6c757d);
    font-weight: 600;
    font-size: 0.95rem;
  }

  .archived-incidents-toggle:hover {
    text-decoration: none;
    color: var(--primary, #007bff);
  }

  .archived-incidents-chevron {
    width: 0.75rem;
    color: inherit;
  }

  .archived-incidents-label {
    user-select: none;
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
