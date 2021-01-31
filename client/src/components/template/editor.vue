<template>
  <div>
    <v-row>
      <v-col>
        <v-card>
          <v-card-title v-text="$t('templates.Variables')" />
          <v-card-text>
            <template-variables v-model="value.vars" />
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <v-card>
          <v-card-title v-text="$t('templates.Install')" />
          <v-card-text>
            <template-processors v-model="value.install" />
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <v-card>
          <v-card-title v-text="$t('templates.RunConfig')" />
          <v-card-text>
            <ui-input
              v-model="value.command"
              :label="$t('templates.Command')"
            />
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <v-card>
          <v-card-title v-text="$t('templates.Shutdown')" />
          <v-card-text>
            <template-shutdown v-model="value.stop" />
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <v-card>
          <v-card-title v-text="$t('templates.PreHook')" />
          <v-card-text>
            <template-processors v-model="value.pre" />
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <v-card>
          <v-card-title v-text="$t('templates.PostHook')" />
          <v-card-text>
            <template-processors v-model="value.post" />
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <v-card>
          <v-card-title v-text="$t('templates.EnvVars')" />
          <v-card-text>
            <ui-map-input
              v-model="value.envVars"
              @input="$forceUpdate()"
            />
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <v-card>
          <v-card-title v-text="$t(server ? 'templates.Environment' : 'templates.Environments')" />
          <v-card-text v-if="!server">
            <template-environments v-model="value.supportedEnvs" />
            <ui-select
              v-if="Object.keys(value.supportedEnvs).length > 0"
              v-model="value.defaultEnv"
              :label="$t('templates.DefaultEnvironment')"
              :items="configuredEnvironments"
            />
          </v-card-text>
          <v-card-text v-else>
            <ui-env-config
              v-model="value.defaultEnv"
              :no-fields-text="$t('env.NoEnvFields')"
            />
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script>
export default {
  props: {
    value: { type: Object, required: true },
    server: { type: Boolean, default: () => false }
  },
  computed: {
    configuredEnvironments () {
      return this.value.supportedEnvs.map(elem => { return { text: elem.type, value: elem } })
    }
  },
  watch: {
    'value.supportedEnvs' (val) {
      if (this.value.defaultEnv) {
        this.value.defaultEnv = val.filter(elem => {
          return elem.type === this.value.defaultEnv.type
        })[0]
      }
    }
  }
}
</script>
