<template>
  <v-row>
    <v-col
      cols="12"
      md="6"
      offset-md="3"
    >
      <v-card>
        <v-card-title v-text="$t(create ? 'users.Add' : 'users.Edit')" />
        <v-card-text>
          <v-row>
            <v-col>
              <ui-input
                v-model="user.username"
                icon="mdi-account"
                :label="$t('common.Name')"
                hide-details
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <ui-input
                v-model="user.email"
                icon="mdi-email"
                :label="$t('users.Email')"
                type="email"
                hide-details
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <ui-password-input
                v-model="user.password"
                :label="$t(create ? 'users.Password' : 'users.NewPassword')"
                :error-messages="passwordErrors"
                :hide-details="passwordErrors === ''"
                @blur="validatePassword"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col class="py-0">
              <ui-switch
                v-model="user.admin"
                :label="$t('scopes.Admin')"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col class="py-0">
              <ui-switch
                v-model="user.viewServers"
                :label="$t('scopes.ViewServers')"
                :disabled="user.admin"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col class="py-0">
              <ui-switch
                v-model="user.createServers"
                :label="$t('scopes.CreateServers')"
                :disabled="user.admin"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col class="py-0">
              <ui-switch
                v-model="user.deleteServers"
                :label="$t('scopes.DeleteServers')"
                :disabled="user.admin"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col class="py-0">
              <ui-switch
                v-model="user.editServerAdmin"
                :label="$t('scopes.EditServerAdmin')"
                :disabled="user.admin"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col class="py-0">
              <ui-switch
                v-model="user.viewNodes"
                :label="$t('scopes.ViewNodes')"
                :disabled="user.admin"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col class="py-0">
              <ui-switch
                v-model="user.editNodes"
                :label="$t('scopes.EditNodes')"
                :disabled="user.admin"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col class="py-0">
              <ui-switch
                v-model="user.deployNodes"
                :label="$t('scopes.DeployNodes')"
                :disabled="user.admin"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col class="py-0">
              <ui-switch
                v-model="user.panelSettings"
                :label="$t('scopes.PanelSettings')"
                :disabled="user.admin"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col class="py-0">
              <ui-switch
                v-model="user.viewTemplates"
                :label="$t('scopes.ViewTemplates')"
                :disabled="user.admin"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col class="py-0">
              <ui-switch
                v-model="user.editTemplates"
                :label="$t('scopes.EditTemplates')"
                :disabled="user.admin"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col class="py-0">
              <ui-switch
                v-model="user.viewUsers"
                :label="$t('scopes.ViewUsers')"
                :disabled="user.admin"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col class="pt-0">
              <ui-switch
                v-model="user.editUsers"
                :label="$t('scopes.EditUsers')"
                :disabled="user.admin"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12">
              <v-btn
                large
                block
                color="success"
                @click="save"
                v-text="$t(create ? 'users.Add' : 'users.Update')"
              />
            </v-col>
            <v-col
              v-if="!create"
              cols="12"
            >
              <v-btn
                block
                color="error"
                @click="deleteUser"
                v-text="$t('users.Delete')"
              />
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
</template>

<script>
export default {
  data () {
    return {
      loading: false,
      showPassword: false,
      passwordErrors: '',
      create: this.$route.params.id === undefined,
      user: {}
    }
  },
  mounted () {
    if (!this.create) this.loadData()
  },
  methods: {
    async loadData () {
      this.loading = true
      this.user = await this.$api.getUser(this.$route.params.id)
      this.loading = false
    },
    validatePassword () {
      if (this.create && (!this.user.password || this.user.password === '')) {
        this.passwordErrors = this.$t('errors.ErrFieldRequired', { field: this.$t('users.Password') })
        return
      }

      if (this.user.password && this.user.password !== '' && this.user.password.length < 8) {
        this.passwordErrors = this.$t('errors.ErrPasswordRequirements')
        return
      }

      this.passwordErrors = ''
    },
    async save () {
      const user = this.user
      if (!user.password || user.password === '') delete user.password
      if (this.create) {
        const id = await this.$api.createUser(user)
        this.$toast.success(this.$t('users.CreateSuccess'))
        this.create = false
        this.$router.push({ name: 'User', params: { id } })
      } else {
        await this.$api.updateUser(this.$route.params.id, user)
        this.$toast.success(this.$t('users.UpdateSuccess'))
      }
    },
    async deleteUser () {
      this.$api.deleteUser(this.$route.params.id)
      this.$router.push({ name: 'Users' })
    }
  }
}
</script>
