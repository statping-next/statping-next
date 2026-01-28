<template>
    <form @submit.prevent="saveSettings">
    <div class="card">
        <div class="card-header">Statping-ng Settings</div>
        <div class="card-body">
                <div class="form-group">
                    <label>{{ $t('project_name') }}</label>
                    <input v-model="core.name" type="text" class="form-control" placeholder="Great Uptime" id="project">
                </div>

                <div class="form-group">
                    <label>{{ $t('project_logo') }}</label>
                    <div class="form-control-plaintext">
                        <small class="text-muted">Configure logos in <a href="#" @click.prevent="goToThemeSettings">Theme settings</a></small>
                    </div>
                </div>

                <div class="form-group">
                    <label>{{ $t('description') }}</label>
                    <input v-model="core.description" type="text" class="form-control" placeholder="Great Uptime" id="description">
                </div>

                <div class="form-group row">
                    <div class="col-8 col-sm-9">
                        <label>{{ $t('domain') }}</label>
                        <input v-model="core.domain" type="url" class="form-control" id="domain">
                    </div>
                    <div class="col-4 col-sm-3 mt-sm-1 mt-0">
                        <label class="d-inline d-sm-none">{{$t('enable_cdn')}}</label>
                        <label class="d-none d-sm-block">{{$t('enable_cdn')}}</label>
                        <span @click="core.using_cdn = !!core.using_cdn" class="switch" id="using_cdn">
                            <input v-model="core.using_cdn" type="checkbox" name="using_cdn" class="switch" id="switch-normal" :checked="core.using_cdn">
                            <label for="switch-normal"></label>
                          </span>
                    </div>
                </div>

                <div class="form-group">
                    <label>{{ $t('footer') }}</label>
                    <textarea v-model="core.footer" rows="4" class="form-control" id="footer">{{core.footer}}</textarea>
                    <small class="form-text text-muted">{{ $t('footer_notes') }}</small>
                </div>

                <div class="form-group">
                    <label>{{ $t('language') }}</label>
                    <select v-model="core.language" class="form-control">
                        <option value="en">English</option>
                        <option value="es">Spanish</option>
                        <option value="fr">French</option>
                        <option value="ru">Russian</option>
                        <option value="de">German</option>
                        <option value="cs">Czech</option>
                        <option value="ja">Japanese</option>
                        <option value="it">Italian</option>
                        <option value="ko">Korean</option>
                        <option value="zh">Chinese</option>
                        <option value="sv">Swedish</option>
                    </select>
                </div>

                <div class="form-group">
                    <label>Default Refresh Rate</label>
                    <select v-model="core.default_refresh_rate" class="form-control">
                        <option :value="15">15 seconds</option>
                        <option :value="60">60 seconds</option>
                        <option :value="0">Off</option>
                    </select>
                    <small class="form-text text-muted">Default auto-refresh interval for new users. Users can override this in their browser.</small>
                </div>

                <div class="form-group row mt-3">
                    <label class="col-sm-10 col-form-label">Show failure details to unauthenticated users?</label>
                    <div class="col-sm-2 float-right">
                        <span @click="core.show_failures_to_unauthenticated = !!core.show_failures_to_unauthenticated" class="switch" id="show_failures_to_unauthenticated">
                        <input v-model="core.show_failures_to_unauthenticated" type="checkbox" name="show_failures_to_unauthenticated" class="switch" id="switch-show-failures" :checked="core.show_failures_to_unauthenticated">
                        <label for="switch-show-failures"></label>
                      </span>
                    </div>
                    <div class="col-12">
                        <small class="form-text text-muted">By default failure details are only visible to authenticated users. Enable this to allow all visitors to view failure details. WARNING: Will expose check endpoints from error messages!</small>
                    </div>
                </div>

                <div class="form-group row mt-3">
                    <label class="col-sm-10 col-form-label">{{ $t('send_reports') }}</label>
                    <div class="col-sm-2 float-right">
                        <span @click="core.allow_reports = !!core.allow_reports" class="switch" id="allow_report">
                        <input v-model="core.allow_reports" type="checkbox" name="allow_report" class="switch" id="switch_allow_report" :checked="core.allow_reports">
                        <label for="switch_allow_report"></label>
                      </span>
                    </div>
                    <div class="col-12">
                        <small>{{ $t('send_reports_desc') }}</small>
                    </div>
                </div>

        </div>
        <div class="card-footer">
            <button @click.prevent="saveSettings" id="save_core" type="submit" class="btn btn-primary btn-block" v-bind:disabled="loading">
                <font-awesome-icon v-if="loading" icon="circle-notch" class="mr-2" spin/>{{ $t('save_settings') }}
            </button>
        </div>
    </div>
    </form>
</template>

<script>
  import Api from '../API'

  export default {
      name: 'CoreSettings',
    data () {
      return {
        loading: false
      }
    },
      computed: {
          core() {
              return this.$store.getters.core
          }
      },
      methods: {
          async saveSettings() {
            this.loading = true
              const c = this.core
              await Api.core_save(c)
              this.$store.commit('setCore', c)
            this.$i18n.locale = c.language || "en";
            this.loading = false
          },
          selectAll() {
              this.$refs.input.select();
          },
          goToThemeSettings() {
              // Navigate to Settings page and switch to Theme tab
              // Use $root to find the Settings component or traverse parent tree
              let parent = this.$parent
              let found = false
              while (parent && !found) {
                  if (parent.$options && parent.$options.name === 'Settings') {
                      found = true
                      if (parent.changeTab) {
                          // Create a mock event object that changeTab expects
                          const mockEvent = {
                              target: {
                                  id: 'v-pills-style-tab'
                              }
                          }
                          parent.changeTab(mockEvent)
                          // Also trigger click on the actual tab link to ensure Bootstrap updates
                          this.$nextTick(() => {
                              const tabLink = document.getElementById('v-pills-style-tab')
                              if (tabLink) {
                                  tabLink.click()
                              }
                          })
                      }
                      break
                  }
                  parent = parent.$parent
              }
          }
      }
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
