<template>
  <v-row>
    <v-col cols="12" md="6" offset-md="3">
      <v-card>
        <v-card-title v-text="$t('common.EditUser')" />
        <v-card-text class="mt-6">
          <v-row>
            <v-col>
              <v-text-field :label="$t('common.Name')" v-model="user.username" outlined />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-text-field :label="$t('common.Email')" v-model="user.email" type="email" outlined />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12">
              <v-btn v-text="$t('common.UpdateUser')" large block color="primary" @click="updateUser" />
            </v-col>
            <v-col cols="12">
              <v-btn v-text="$t('common.DeleteUser')" block color="error" @click="deleteUser" />
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
      ctx.$http.get(`/api/users/${ctx.$route.params.id}`).then(function (response) {
        ctx.user = response.data
        console.log('USER', ctx.user)
        ctx.loading = false
      }).catch(function (error) {
        let msg = 'errors.ErrUnknownError'
        if (error && error.response && error.response.data.error) {
          if (error.response.data.error.code) {
            msg = 'errors.' + error.response.data.error.code
          } else {
            msg = error.response.data.error.msg
          }
        }

        ctx.$toast.error(ctx.$t(msg))
      })
    },
    updateUser () {
      const ctx = this
      ctx.$http.post(`/api/users/${ctx.$route.params.id}`, ctx.user).then(function (response) {
        ctx.$toast.success(ctx.$t('common.UserUpdateSuccess'))
      }).catch(function () {
        ctx.$toast.error(ctx.$t('common.UserUpdateError'))
      })
    },
    deleteUser () {
      const ctx = this
      ctx.$http.delete(`/api/users/${ctx.$route.params.id}`).then(function (response) {
        ctx.$router.push({ name: 'Users' })
      }).catch(function (error) {
        let msg = 'errors.ErrUnknownError'
        if (error && error.response && error.response.data.error) {
          if (error.response.data.error.code) {
            msg = 'errors.' + error.response.data.error.code
          } else {
            msg = error.response.data.error.msg
          }
        }

        ctx.$toast.error(ctx.$t(msg))
      })
    }
  }
}
</script>
