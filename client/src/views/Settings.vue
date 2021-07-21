<!--
  - Copyright 2019 Padduck, LLC
  -  Licensed under the Apache License, Version 2.0 (the "License");
  -  you may not use this file except in compliance with the License.
  -  You may obtain a copy of the License at
  -          http://www.apache.org/licenses/LICENSE-2.0
  -  Unless required by applicable law or agreed to in writing, software
  -  distributed under the License is distributed on an "AS IS" BASIS,
  -  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  -  See the License for the specific language governing permissions and
  -  limitations under the License.
  -->

<template>
  <v-container>
    <v-row>
      <v-col cols="12">
        <v-card
          :loading="panelSettingsLoading"
          :disabled="panelSettingsLoading"
        >
          <v-card-title v-text="$t('settings.PanelSettings')" />
          <v-card-text>
            <v-row>
              <v-col cols="12">
                <ui-input
                  v-model="masterUrl"
                  end-icon="mdi-auto-fix"
                  :label="$t('settings.MasterUrl')"
                  :hint="$t('settings.MasterUrlHint')"
                  @click:append="autofillMasterUrl"
                />
              </v-col>
              <v-col cols="12">
                <ui-input
                  v-model="panelTitle"
                  :label="$t('settings.CompanyName')"
                />
              </v-col>
              <v-col cols="12">
                <ui-select
                  v-model="defaultTheme"
                  :label="$t('settings.DefaultTheme')"
                  :items="themes"
                />
              </v-col>
              <v-col cols="12">
                <v-btn
                  large
                  block
                  color="success"
                  @click="savePanelConfig()"
                  v-text="$t('common.Save')"
                />
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12">
        <v-card
          :loading="emailSettingsLoading"
          :disabled="emailSettingsLoading"
        >
          <v-card-title v-text="$t('settings.EmailSettings')" />
          <v-card-text>
            <v-row>
              <v-col cols="12">
                <ui-select
                  v-model="emailProvider"
                  :label="$t('settings.EmailProvider')"
                  :items="emailProviders"
                />
              </v-col>
              <v-col
                v-for="key in emailFields"
                :key="key"
                cols="12"
              >
                <ui-input
                  v-model="email[key]"
                  :label="$t('settings.emailFields.' + key)"
                />
              </v-col>
              <v-col cols="12">
                <v-btn
                  large
                  block
                  color="success"
                  @click="saveEmailConfig()"
                  v-text="$t('common.Save')"
                />
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
const emailProviderConfigs = {
  none: [],
  smtp: ['from', 'host', 'username', 'password'],
  mailgun: ['domain', 'from', 'key'],
  mailjet: ['domain', 'from', 'key']
}

export default {
  data () {
    return {
      panelSettingsLoading: true,
      masterUrl: '',
      panelTitle: '',
      defaultTheme: '',
      themes: ['PufferPanel'],
      emailSettingsLoading: true,
      emailProvider: 'none',
      emailProviders: [],
      emailFields: [],
      email: {
        from: '',
        domain: '',
        key: '',
        host: '',
        username: '',
        password: ''
      }
    }
  },
  watch: {
    emailProvider (newVal) {
      this.emailFields = emailProviderConfigs[newVal]
    }
  },
  mounted () {
    Object.keys(emailProviderConfigs).map(provider => {
      this.emailProviders.push({
        text: this.$t('settings.emailProviders.' + provider),
        value: provider
      })
    })

    if (process.env.NODE_ENV !== 'production') {
      emailProviderConfigs.debug = []
      this.emailProviders.push({ text: 'Debug', value: 'debug' })
    }

    this.loadData()
  },
  beforeDestroy () {
    if (process.env.NODE_ENV !== 'production') {
      delete emailProviderConfigs.debug
    }
  },
  methods: {
    async loadData () {
      this.masterUrl = await this.$api.getSetting('panel.settings.masterUrl')
      this.panelTitle = await this.$api.getSetting('panel.settings.companyName')
      this.defaultTheme = await this.$api.getSetting('panel.settings.defaultTheme')
      this.themes = (await this.$api.getConfig()).themes.available
      this.panelSettingsLoading = false
      this.emailProvider = await this.$api.getSetting('panel.email.provider')
      Object.keys(this.email).map(async key => {
        this.email[key] = await this.$api.getSetting('panel.email.' + key)
      })
      this.emailSettingsLoading = false
    },
    autofillMasterUrl () {
      this.masterUrl = window.location.origin
    },
    async savePanelConfig () {
      this.panelSettingsLoading = true
      try {
        await this.$api.setSettings({
          'panel.settings.masterUrl': this.masterUrl,
          'panel.settings.companyName': this.panelTitle,
          'panel.settings.defaultTheme': this.defaultTheme
        })
        this.$toast.success(this.$t('common.Saved'))
      } finally {
        this.panelSettingsLoading = false
      }

      this.loadData()
    },
    async saveEmailConfig () {
      this.emailSettingsLoading = true
      try {
        const data = { 'panel.email.provider': this.emailProvider }
        Object.keys(this.email).map(key => {
          data['panel.email.' + key] = this.email[key]
        })
        await this.$api.setSettings(data)
        this.$toast.success(this.$t('common.Saved'))
      } finally {
        this.emailSettingsLoading = false
      }

      this.loadData()
    }
  }
}
</script>
