<template>
  <div v-if="messagesInRange.length > 0" class="service-announcements-block">
    <div v-for="message in messagesInRange" :key="message.id" class="service-announcement-nested" role="alert">
      <h5 class="incident-title font-weight-bold">
        <span class="incident-emoji">📢</span> {{ message.title }}
      </h5>
      <div v-if="message.description" class="incident-description" v-html="message.description"></div>
      <div class="incident-dates">
        <span class="incident-date-item">Started {{ niceDate(message.start_on) }} ({{ aboutAgo(message.start_on) }})</span>
        <span class="incident-date-sep">·</span>
        <span class="incident-date-item">Ends {{ niceDate(message.end_on) }} ({{ timeRelative(message.end_on) }})</span>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ServiceAnnouncements',
  props: {
    service: {
      type: Object,
      required: true
    }
  },
  computed: {
    messagesInRange() {
      const list = this.$store.getters.serviceMessages(this.service.id) || []
      return list.filter(m => this.inRange(m))
    }
  },
  methods: {
    inRange(message) {
      return this.isBetween(
        this.now(),
        message.start_on,
        message.start_on === message.end_on ? this.maxDate().toISOString() : message.end_on
      )
    }
  }
}
</script>

<style scoped>
/* Block layout – same as IncidentsBlock */
.service-announcements-block {
  margin-top: -16px;
  padding-top: 0;
}

/* Vertical padding matches incident-nested exactly: 16px top, 0 bottom (incident uses 16px 16px 0) */
.service-announcement-nested {
  padding: 16px 16px 0;
  margin-bottom: 0;
  border-radius: 0 0 6px 6px;
  background: #f0f0f0;
  border: 12px solid #f8f8f8;
}
.service-announcement-nested:first-child {
  border-top-left-radius: 0;
  border-top-right-radius: 0;
  /* Slightly more top so title aligns with incident (user: "just slightly too little on the top") */
  padding-top: 18px;
}
.service-announcement-nested:last-child {
  margin-bottom: 0;
  /* No extra bottom padding – match incident-nested (incident has padding 16px 16px 0, no padding-bottom) */
  padding-bottom: 0;
}
.service-announcements-block:has(+ .incidents-block) .service-announcement-nested:last-child {
  border-radius: 0;
}

.list-group-item .service-announcement-nested {
  margin-left: -1.25rem;
  margin-right: -1.25rem;
  padding-left: calc(12px + 16px);
  padding-right: calc(12px + 16px);
}

.service-entry-wrapper .service-announcement-nested {
  margin-left: 0;
  margin-right: 0;
}

.dark-theme .service-announcement-nested {
  background: #1a1d21;
  border-color: #0f1114;
}

/* Content styles – identical to IncidentsBlock (same classes, same values) */
.incident-title {
  font-size: 1rem;
  margin-bottom: 0.35rem;
}
.incident-emoji {
  font-size: 0.85em;
  vertical-align: middle;
}

.incident-description {
  font-size: 0.9rem;
  margin-bottom: 0.5rem;
}
.incident-description >>> p:last-child {
  margin-bottom: 0;
}

.incident-dates {
  font-size: 0.8rem;
  color: #6c6c6c;
  margin-bottom: 0.75rem;
}
.dark-theme .incident-dates {
  color: #aaaaaa;
}
.incident-date-sep {
  margin: 0 0.35rem;
}
</style>
