<template>
    <font-awesome-layers v-if="activeMessages(messages) > 0">
        <font-awesome-icon icon="calendar" style="color:dodgerblue"/>
        <font-awesome-layers-text class="icon-text" :value="activeMessages(messages)" />
    </font-awesome-layers>
</template>

<script>
export default {
  name: 'MessagesIcon',
    props: {
        messages: {
            type: [Array, Object],
            required: false,
            default: () => []
        }
    },
    methods: {
        activeMessages(msgs) {
            if (!msgs || !Array.isArray(msgs)) {
                return 0
            }
            return msgs.filter(m => this.isAfter(this.now(), m.start_on) && this.isBefore(this.now(), m.end_on)).length
        }
    }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>