<template>
    <div class="incident-update-row" role="alert">
        <div class="incident-update-main">
            <span class="incident-update-type" :class="updateTypeClass">{{ update.type }}</span>
            <span class="incident-update-message">{{ update.message }}</span>
            <button v-if="admin" @click="delete_update(update)" type="button" class="incident-update-delete" aria-label="Delete update">
                <font-awesome-icon icon="times" />
            </button>
        </div>
        <div class="incident-update-datetime">{{ niceDate(update.created_at) }} (about {{ ago(update.created_at) }} ago)</div>
    </div>
</template>

<script>
  import Api from "@/API";

  export default {
    name: "IncidentUpdate",
    props: {
      update: {
        required: true
      },
      admin: {
        required: true
      },
      onUpdate: {
        required: false
      }
    },
    computed: {
      updateTypeClass() {
        const t = (this.update.type || '').toLowerCase()
        if (t === 'resolved') return 'incident-update-type--resolved'
        if (t === 'investigating') return 'incident-update-type--investigating'
        if (t === 'update') return 'incident-update-type--update'
        return 'incident-update-type--unknown'
      }
    },
    methods: {
      async delete_update(update) {
        this.res = await Api.incident_update_delete(update)
        if (this.res.status === "success") {
          this.onUpdate()
        }
      },
    }
  }
</script>

<style scoped>
  .incident-update-row {
    padding: 0.75rem 0;
    border-bottom: 1px solid rgba(0, 0, 0, 0.06);
  }

  .incident-update-row:last-child {
    border-bottom: none;
  }

  .incident-update-main {
    display: flex;
    align-items: flex-start;
    gap: 0.5rem;
  }

  .incident-update-type {
    font-weight: 600;
    font-size: 0.875rem;
    text-transform: capitalize;
    flex-shrink: 0;
  }

  .incident-update-type--resolved {
    color: #28a745;
  }

  .incident-update-type--investigating {
    color: #dc3545;
  }

  .incident-update-type--update {
    color: #fd7e14;
  }

  .incident-update-type--unknown {
    color: #6c757d;
  }

  .incident-update-message {
    flex: 1;
    color: inherit;
    font-size: 0.9375rem;
  }

  .incident-update-delete {
    background: none;
    border: none;
    color: #6c757d;
    padding: 0.125rem;
    cursor: pointer;
    font-size: 0.875rem;
    flex-shrink: 0;
  }

  .incident-update-delete:hover {
    color: #dc3545;
  }

  .incident-update-datetime {
    font-size: 0.8125rem;
    color: #6c757d;
    margin-top: 0.35rem;
    margin-left: 0;
  }
  .dark-theme .incident-update-datetime {
    color: #aaaaaa;
  }
</style>
