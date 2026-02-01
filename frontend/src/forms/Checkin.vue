<template>
    <div>
        <div v-for="(checkin, i) in checkins" class="col-12 alert alert-light" role="alert">
            <span class="badge badge-pill badge-info text-uppercase">{{checkin.name}}</span>
            <span class="float-right font-2">Last checkin {{ago(checkin.last_hit)}}</span>
            <span class="float-right font-2 mr-3">Check Every {{checkin.interval}} seconds</span>
            <span class="float-right font-2 mr-3">Grace Period {{checkin.grace}} seconds</span>
            <span class="d-block mt-2">
                <input type="text" class="form-control" :value="`${core.domain}/checkin/${checkin.api_key}`" readonly>
                <span class="small">{{ $t('send_get_every_seconds', { interval: checkin.interval }) }}
                    <button @click.prevent="deleteCheckin(checkin)" type="button" class="btn btn-danger btn-xs float-right mt-1">Delete</button>
                </span>
            </span>
        </div>

        <div class="col-12 alert alert-light">
            <form @submit.prevent="saveCheckin">
                <div class="form-group row">
                    <div class="col-12 col-md-5">
                        <label for="checkin_interval" class="col-form-label">{{ $t('checkin_name_label') }}</label>
                        <input v-model="checkin.name" type="text" name="name" class="form-control" id="checkin_name" :placeholder="$t('new_checkin_placeholder')">
                    </div>
                    <div class="col-12 col-md-5">
                        <label for="checkin_interval" class="col-form-label">{{ $t('interval_minutes_label') }}</label>
                        <input v-model.number="checkin.interval" type="number" name="interval" class="form-control" id="checkin_interval" :placeholder="$t('checkin_interval_placeholder')" min="1">
                    </div>
                    <div class="col-12 col-md-5">
                        <label class="col-form-label"></label>
                        <button @click.prevent="saveCheckin" type="submit" id="submit" class="btn btn-success d-block mt-2">{{ $t('save_checkin') }}</button>
                    </div>
                </div>
            </form>
        </div>
        </div>
</template>

<script>
  import Api from "../API";

  export default {
      name: 'Checkin',
      props: {
          service: {
              type: Object,
              required: true
          }
      },
      data() {
          return {
              checkin: {
                  name: "",
                  interval: 60,
                  service_id: this.service.id
              }
          }
      },
      mounted() {

      },
      computed: {
          checkins() {
              return this.$store.getters.serviceCheckins(this.service.id)
          },
          core() {
              return this.$store.getters.core
          },
      },
      methods: {
          fixInts() {
              const c = this.checkin
              this.checkin.interval = parseInt(c.interval)
              this.checkin.grace = parseInt(c.grace)
              return this.checkin
          },
          async saveCheckin() {
              const c = this.fixInts()
              await Api.checkin_create(c)
              await this.updateCheckins()
          },
          async deleteCheckin(checkin) {
              await Api.checkin_delete(checkin)
              await this.updateCheckins()
          },
          async updateCheckins() {
              const checkins = await Api.checkins()
              this.$store.commit('setCheckins', checkins)
          }
      }
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
