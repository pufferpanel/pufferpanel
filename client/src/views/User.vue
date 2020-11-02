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
              <v-text-field
                v-model="user.username"
                prepend-inner-icon="mdi-account"
                :label="$t('common.Name')"
                outlined
                hide-details
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-text-field
                v-model="user.email"
                prepend-inner-icon="mdi-email"
                :label="$t('users.Email')"
                type="email"
                outlined
                hide-details
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-text-field
                v-model="user.password"
                prepend-inner-icon="mdi-lock"
                :label="$t(create ? 'users.Password' : 'users.NewPassword')"
                :append-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
                :type="!showPassword ? 'password' : 'text'"
                :error-messages="passwordErrors"
                outlined
                :hide-details="passwordErrors === ''"
                @click:append="showPassword = !showPassword"
                @blur="validatePassword"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col class="py-0">
              <v-switch
                v-model="user.admin"
                hide-details
                :label="$t('scopes.Admin')"
                class="mt-2"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col class="py-0">
              <v-switch
                v-model="user.viewServers"
                hide-details
                :label="$t('scopes.ViewServers')"
                :disabled="user.admin"
                class="mt-2"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col class="py-0">
              <v-switch
                v-model="user.createServers"
                hide-details
                :label="$t('scopes.CreateServers')"
                :disabled="user.admin"
                class="mt-2"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col class="py-0">
              <v-switch
                v-model="user.deleteServers"
                hide-details
                :label="$t('scopes.DeleteServers')"
                :disabled="user.admin"
                class="mt-2"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col class="py-0">
              <v-switch
                v-model="user.editServerAdmin"
                hide-details
                :label="$t('scopes.EditServerAdmin')"
                :disabled="user.admin"
                class="mt-2"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col class="py-0">
              <v-switch
                v-model="user.viewNodes"
                hide-details
                :label="$t('scopes.ViewNodes')"
                :disabled="user.admin"
                class="mt-2"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col class="py-0">
              <v-switch
                v-model="user.editNodes"
                hide-details
                :label="$t('scopes.EditNodes')"
                :disabled="user.admin"
                class="mt-2"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col class="py-0">
              <v-switch
                v-model="user.deployNodes"
                hide-details
                :label="$t('scopes.DeployNodes')"
                :disabled="user.admin"
                class="mt-2"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col class="py-0">
              <v-switch
                v-model="user.viewTemplates"
                hide-details
                :label="$t('scopes.ViewTemplates')"
                :disabled="user.admin"
                class="mt-2"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col class="py-0">
              <v-switch
                v-model="user.editTemplates"
                hide-details
                :label="$t('scopes.EditTemplates')"
                :disabled="user.admin"
                class="mt-2"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col class="py-0">
              <v-switch
                v-model="user.viewUsers"
                hide-details
                :label="$t('scopes.ViewUsers')"
                :disabled="user.admin"
                class="mt-2"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col class="pt-0">
              <v-switch
                v-model="user.editUsers"
                hide-details
                :label="$t('scopes.EditUsers')"
                :disabled="user.admin"
                class="mt-2"
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
import { handleError } from '@/utils/api'

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
    loadData () {
      const ctx = this
      ctx.loading = true
      ctx.$http.get(`/api/users/${ctx.$route.params.id}/perms`).then(response => {
        ctx.user = { ...response.data }
        ctx.loading = false
      }).catch(handleError(ctx))
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
    save () {
      const ctx = this
      const url = ctx.$route.params.id ? '/api/users/' + ctx.$route.params.id : '/api/users'
      const user = ctx.user
      if (!user.password || user.password === '') delete user.password
      ctx.$http.post(url, user).then(response => {
        const id = ctx.$route.params.id || response.data.id
        ctx.$http.put(`/api/users/${id}/perms`, user).then(response => {
          ctx.$toast.success(ctx.$t(this.create ? 'users.CreateSuccess' : 'users.UpdateSuccess'))
          if (this.create) ctx.$router.push({ name: 'User', params: { id } })
        }).catch(error => {
          // eslint-disable-next-line no-console
          console.log(error)
          ctx.$toast.error(ctx.$t('users.PermsUpdateError'))
        })
      }).catch(error => {
        // eslint-disable-next-line no-console
        console.log(error)
        ctx.$toast.error(ctx.$t(this.create ? 'users.CreateError' : 'users.UpdateError'))
      })
    },
    deleteUser () {
      const ctx = this
      ctx.$http.delete(`/api/users/${ctx.$route.params.id}`).then(response => {
        ctx.$router.push({ name: 'Users' })
      }).catch(handleError(ctx))
    }
  }
}
</script>
