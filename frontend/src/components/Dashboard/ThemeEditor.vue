<template>
    <div>
        <div class="card mb-5">
            <div class="card-header">{{ $t('default_theme') }}</div>
            <div class="card-body">
                <div class="form-group">
                    <label>{{ $t('default_theme_for_new_users') }}</label>
                    <select v-model="core.default_theme" @change="saveDefaultTheme" class="form-control">
                        <option value="">{{ $t('system_default_light') }}</option>
                        <option value="light">{{ $t('light') }}</option>
                        <option value="dark">{{ $t('dark') }}</option>
                    </select>
                    <small class="form-text text-muted">{{ $t('default_theme_desc') }}</small>
                </div>
            </div>
        </div>

        <div class="card mb-5">
            <div class="card-header">{{ $t('logo_configuration') }}</div>
            <div class="card-body">
                <div class="form-group">
                    <label>{{ $t('light_theme_logo') }}</label>
                    <input v-model="core.logo_light" @change="saveLogo" type="text" class="form-control" :placeholder="$t('logo_url_placeholder')" id="logo_light">
                    <small class="form-text text-muted">{{ $t('logo_light_desc') }}</small>
                </div>
                <div class="form-group">
                    <label>{{ $t('dark_theme_logo_optional') }}</label>
                    <input v-model="core.logo_dark" @change="saveLogo" type="text" class="form-control" :placeholder="$t('logo_dark_placeholder')" id="logo_dark">
                    <small class="form-text text-muted">{{ $t('logo_dark_optional_desc') }}</small>
                </div>
                <div class="form-group">
                    <label>{{ $t('favicon_url_optional') }}</label>
                    <input v-model="core.favicon" @change="saveLogo" type="text" class="form-control" :placeholder="$t('favicon_placeholder')" id="favicon">
                    <small class="form-text text-muted">{{ $t('favicon_desc') }}</small>
                </div>
            </div>
        </div>

        <div class="card mb-5">
        <div class="card-header">{{$t('theme_editor')}}</div>
        <div class="card-body">
        <div v-if="error" class="alert alert-danger mt-3" style="white-space: pre-line;">
            {{error}}
        </div>

            <h6 v-if="directory" id="assets_dir" class="text-muted text-monospace text-sm-center font-1 mb-4">
                {{$t('assets_dir')}}: {{directory}}
            </h6>

            <div v-if="loaded && !directory" class="jumbotron jumbotron-fluid">
                <div class="text-center col-12">
                    <h1 class="display-5">{{$t('enable_assets')}}</h1>
                    <span class="lead">{{$t('assets_desc')}}<p>
                        <button id="enable_assets" @click.prevent="createAssets" :disabled="pending" href="#" class="btn btn-primary mt-3">
                            <font-awesome-icon v-if="pending" icon="circle-notch" class="mr-2" spin/>{{pending ? $t('assets_loading') : $t('assets_btn')}}
                        </button>
                    </p></span>
                </div>
        </div>

    <form v-observe-visibility="visible" v-if="loaded && directory" @submit.prevent="saveAssets" :disabled="pending">
        <h3>{{ $t('theme_editor_variables') }}</h3>
        <codemirror v-show="loaded" v-model="vars" ref="vars" :options="cmOptions" class="codemirrorInput"/>

        <h3 class="mt-3">{{ $t('base_theme') }}</h3>
        <codemirror v-show="loaded" v-model="base" ref="base" :options="cmOptions" class="codemirrorInput"/>

        <h3 class="mt-3">{{ $t('layout_theme') }}</h3>
        <codemirror v-show="loaded" v-model="layout" ref="layout" :options="cmOptions" class="codemirrorInput"/>

        <h3 class="mt-3">{{ $t('forms_theme') }}</h3>
        <codemirror v-show="loaded" v-model="forms" ref="forms" :options="cmOptions" class="codemirrorInput"/>

        <h3 class="mt-3">{{ $t('mixins') }}</h3>
        <codemirror v-show="loaded" v-model="mixins" ref="mixins" :options="cmOptions" class="codemirrorInput"/>

        <h3 class="mt-3">{{ $t('mobile_overwrites') }}</h3>
        <codemirror v-show="loaded" v-model="mobile" ref="mobile" :options="cmOptions" class="codemirrorInput"/>

    </form>
    </div>

        <div v-if="directory" class="card-footer">
            <div class="row">
                <div class="col-6">
                    <button id="save_assets" @click.prevent="saveAssets" type="submit" class="btn btn-primary btn-block" :disabled="pending">{{pending ? $t('saving') : $t('save_styles') }}</button>
                </div>
                <div class="col-6">
                    <button id="delete_assets" @click.prevent="deleteAssets" class="btn btn-danger btn-block confirm-btn" :disabled="pending">{{ $t('delete_local_assets') }}</button>
                </div>
            </div>
        </div>
        </div>
    </div>
</template>

<script>
import Api from "../../API";

// require component
import {codemirror} from 'vue-codemirror'

import('codemirror/mode/css/css.js')

  import('codemirror/lib/codemirror.css')
  import('codemirror-colorpicker/dist/codemirror-colorpicker.css')
  import('codemirror-colorpicker')

  export default {
      name: 'ThemeEditor',
      components: {
          codemirror
      },
      computed: {
          core() {
              return this.$store.getters.core
          }
      },
      data () {
          return {
              base: null,
              layout: null,
              forms: null,
              mixins: null,
              vars: null,
              mobile: null,
              error: null,
              directory: null,
              tab: "vars",
              loaded: true,
              pending: false,
              cmOptions: {
                  height: 700,
                  tabSize: 4,
                  lineNumbers: true,
                  matchBrackets: true,
                  mode: "text/x-scss",
                  line: true,
                  colorpicker: true
              }
          }
      },
      async mounted () {
          await this.fetchTheme()
          this.changeTab('vars')
      },
      methods: {
        async saveDefaultTheme() {
          const c = this.core
          await Api.core_save(c)
          this.$store.commit('setCore', c)
        },
        async saveLogo() {
          const c = this.core
          await Api.core_save(c)
          this.$store.commit('setCore', c)
        },
        visible(isVisible, entry) {
          if (isVisible) {
            this.$refs.vars.codemirror.refresh()
            this.$refs.base.codemirror.refresh()
            this.$refs.mobile.codemirror.refresh()
            this.$refs.layout.codemirror.refresh()
            this.$refs.forms.codemirror.refresh()
            this.$refs.mixins.codemirror.refresh()
          }
        },
          async fetchTheme() {
              this.loaded = true
              this.pending = true
              const theme = await Api.theme()
              this.directory = theme.directory
              if (this.directory) {
                  this.base = theme.base
                  this.vars = theme.variables
                  this.mobile = theme.mobile
                  this.layout = theme.layout
                  this.forms = theme.forms
                  this.mixins = theme.mixins
              }
              this.pending = false
              this.loaded = true
          },
          async createAssets() {
              this.pending = true
            let resp
            try {
              resp = await Api.theme_generate(true)
            } catch(e) {
                this.error = e.response.data.error
            }
              this.pending = false
            await this.fetchTheme()
          },
        async delete() {
          this.pending = true
          const resp = await Api.theme_generate(false)
          await this.fetchTheme()
          this.pending = false
        },
          async deleteAssets() {
            const modal = {
              visible: true,
              title: this.$t('delete_local_assets'),
              body: this.$t('confirm_delete_assets'),
              btnColor: "btn-danger",
              btnText: this.$t('delete'),
              func: () => this.delete(),
            }
            this.$store.commit("setModal", modal)
          },
          async saveAssets() {
              this.pending = true
              const data = {
                base: this.base,
                layout: this.layout,
                forms: this.forms,
                mixins: this.mixins,
                variables: this.vars,
                mobile: this.mobile
              }
            let resp
            try {
              resp = await Api.theme_save(data)
            } catch(e) {
              resp = {status: 'error', error: e.response.data.error}
            }
              if (resp.error) {
                  this.error = resp.error
                  this.pending = false
                  return
              } else {
                  this.error = null
              }
              this.pending = false
              await this.fetchTheme()
          },
          changeTab (v) {
              this.tab = v
              // if (v === 'base') {
              //     this.$refs.base.codemirror.refresh();
              // } else if (v === 'vars') {
              //     this.$refs.vars.codemirror.refresh();
              // } else if (v === 'mobile') {
              //     this.$refs.mobile.codemirror.refresh();
              // }
          }
      }
  }
</script>
