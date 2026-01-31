<template>
  <div class="card shadow mb-4 global-incident-block" role="alert">
    <div class="card-body">
      <h3 class="mb-3 font-weight-bold global-incident-title">
        <span class="global-incident-emoji">📢</span> {{ message.title }}
      </h3>
      <div class="global-incident-description">{{ message.description }}</div>
      <div class="global-incident-dates">
        <span class="global-incident-date-item">Started {{ niceDate(message.start_on) }} ({{ aboutAgo(message.start_on) }})</span>
        <span class="global-incident-date-sep">·</span>
        <span class="global-incident-date-item">Ends {{ niceDate(message.end_on) }} ({{ timeRelative(message.end_on) }})</span>
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
/* Match service incidents layout: full width, same padding, hover */
.global-incident-block {
  width: 100%;
  max-width: 100%;
  box-sizing: border-box;
  transition: background-color 0.2s ease, box-shadow 0.2s ease;
}
/* 1rem padding on all sides (pb-2 removed so Bootstrap doesn't override bottom) */
.global-incident-block .card-body {
  padding: 1rem;
  box-sizing: border-box;
}
.global-incident-block:hover {
  background-color: #f8f8f8;
  box-shadow: 0 4px 12px 1px rgba(0, 0, 0, 0.12);
}
.dark-theme .global-incident-block:hover {
  background-color: #1a1d21 !important;
  box-shadow: 0 4px 12px 1px rgba(0, 0, 0, 0.35);
}

.global-incident-title {
  font-size: calc(1.25rem - 0.5px);
}

.global-incident-emoji {
  font-size: 0.9em;
  vertical-align: middle;
}

.global-incident-description {
  font-size: calc(1rem - 0.5px);
  margin-bottom: 0.75rem;
}

.global-incident-dates {
  font-size: calc(0.875rem + 0.5px);
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
