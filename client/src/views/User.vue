<template>
  <v-row>
    <v-col
      cols="12"
      md="6"
      offset-md="3"
    >
      <v-card>
        <v-card-title v-text="$t('users.Edit')" />
        <v-card-text class="mt-6">
          <v-row>
            <v-col>
              <v-text-field
                v-model="user.username"
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
                :label="$t('users.Email')"
                type="email"
                outlined
                hide-details
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
                color="primary"
                @click="updateUser"
                v-text="$t('users.Update')"
              />
            </v-col>
            <v-col cols="12">
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
      loading: true,
      user: {}
    }
  },
  mounted () {
    this.loadData()
  },
  methods: {
    loadData () {
      const ctx = this
      ctx.$http.get(`/api/users/${ctx.$route.params.id}/perms`).then(response => {
        ctx.user = { ...response.data }
        ctx.loading = false
      }).catch(handleError(ctx))
    },
    updateUser () {
      const ctx = this
      ctx.$http.post(`/api/users/${ctx.$route.params.id}`, ctx.user).then(response => {
        ctx.$http.put(`/api/users/${ctx.$route.params.id}/perms`, ctx.user).then(response => {
          ctx.$toast.success(ctx.$t('users.UpdateSuccess'))
        }).catch(error => {
          // eslint-disable-next-line no-console
          console.log(error)
          ctx.$toast.error(ctx.$t('users.PermsUpdateError'))
        })
      }).catch(error => {
        // eslint-disable-next-line no-console
        console.log(error)
        ctx.$toast.error(ctx.$t('users.UpdateError'))
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
