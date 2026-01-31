<template>
  <div class="card shadow mb-4 global-incident-block" role="alert">
    <div class="card-body pb-2">
      <h3 class="mb-3 font-weight-bold global-incident-title">
        <span class="global-incident-emoji">📢</span> {{ message.title }}
      </h3>
      <div class="global-incident-description">{{ message.description }}</div>
      <div class="global-incident-dates">
        <span class="global-incident-date-item">Started {{ niceDate(message.start_on) }}</span>
        <span class="global-incident-date-sep">·</span>
        <span class="global-incident-date-item">Ends {{ niceDate(message.end_on) }}</span>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'MessageBlock',
  props: {
    message: {
      type: Object,
      required: true,
    }
  }
}
</script>

<style scoped>
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

.global-incident-dates {
  font-size: 0.875rem;
  color: var(--secondary, #6c6c6c);
  margin-bottom: 0;
}
.dark-theme .global-incident-dates {
  color: #aaaaaa;
}

.global-incident-date-sep {
  margin: 0 0.35rem;
}
</style>
