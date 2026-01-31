<template>
    <div class="incident-updates-body">

        <div v-if="updates.length === 0" class="incident-updates-empty">
            No updates yet. Add one below.
        </div>

        <div class="incident-updates-list">
            <IncidentUpdate v-for="update in updates" :key="update.id" :update="update" :onUpdate="onUpdateCallback" :admin="true"/>
        </div>

        <form class="incident-updates-add row" @submit.prevent="createIncidentUpdate">
            <div class="col-12 col-md-3 mb-2 mb-md-0">
                <select v-model="incident_update.type" class="form-control form-control-sm">
                    <option value="Investigating">Investigating</option>
                    <option value="Update">Update</option>
                    <option value="Unknown">Unknown</option>
                    <option value="Resolved">Resolved</option>
                </select>
            </div>
            <div class="col-12 col-md-7 mb-2 mb-md-0">
                <input v-model="incident_update.message" name="description" class="form-control form-control-sm" placeholder="Update message..." required>
            </div>
            <div class="col-12 col-md-2">
                <button @click.prevent="createIncidentUpdate"
                        :disabled="!incident_update.message"
                        type="submit" class="btn btn-sm btn-primary btn-block">
                    Add
                </button>
            </div>
        </form>

    </div>
</template>

<script>
    import Api from "../API";
    const IncidentUpdate = () => import(/* webpackChunkName: "index" */ "@/components/Elements/IncidentUpdate");

    export default {
        name: 'FormIncidentUpdates',
        components: {IncidentUpdate},
        props: {
            incident: {
                type: Object,
                required: true
            }
        },
        data () {
            return {
                updates: [],
                incident_update: {
                    incident: this.incident.id,
                    message: "",
                    type: "Investigating" // TODO: default to something.. theres is no error checking for blank submission...
                }
            }
        },

        async mounted() {
            await this.loadUpdates()
        },

        methods: {
            async createIncidentUpdate() {
                this.res = await Api.incident_update_create(this.incident_update)
                if (this.res.status === "success") {
                    this.updates.push(this.res.output)
                    this.$emit('updated')
                }

                // reset the form data
                this.incident_update = {
                    incident: this.incident.id,
                    message: "",
                    type: "Investigating"
                }

            },

            async loadUpdates() {
                this.updates = await Api.incident_updates(this.incident)
            },
            async onUpdateCallback() {
                await this.loadUpdates()
                this.$emit('updated')
            }
        }
    }
</script>

<style scoped>
  .incident-updates-body {
    padding: 1rem 1.25rem;
  }

  .incident-updates-empty {
    font-size: 0.875rem;
    color: var(--secondary, #6c757d);
    margin-bottom: 1rem;
  }

  .incident-updates-list {
    margin-bottom: 1.25rem;
  }

  .incident-updates-add {
    padding-top: 1rem;
    border-top: 1px solid rgba(0, 0, 0, 0.08);
    align-items: center;
  }

  .incident-updates-add .form-control,
  .incident-updates-add .btn {
    height: 2.25rem;
  }
</style>
