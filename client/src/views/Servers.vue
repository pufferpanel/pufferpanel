<template>
  <v-container
    fill-height
    fluid
    grid-list-xl
  >
    <v-layout
      justify-center
      wrap
    >
      <v-flex
        md12
      >
        <material-notification
          v-if="error"
          color="error"
          v-text="error"
        />
        <material-card
          color="blue"
          title="Servers"
        >
          <v-data-table
            :headers="headers"
            :items="servers"
            hide-actions
          >
            <template
              slot="headerCell"
              slot-scope="{ header }"
            >
              <span
                class="subheading font-weight-light text-success text--darken-3"
                v-text="header.text"
              />
            </template>
            <template
              slot="items"
              slot-scope="{ item }"
            >
              <td>{{ item.name }}</td>
              <td></td>
              <td></td>
              <td class="text-xs-right"></td>
            </template>
          </v-data-table>
        </material-card>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
export default {
  data () {
    return {
      headers: [
        {
          sortable: false,
          text: 'Name',
          value: 'name'
        },
        {
          sortable: false,
          text: 'Node',
          value: 'node'
        },
        {
          sortable: false,
          text: 'Address',
          value: 'address'
        },
        {
          sortable: false,
          text: 'Online',
          value: 'online',
          align: 'right'
        }
      ],
      servers: [
      ],
      error: null
    }
  },
  methods: {
    loadData() {
      let vueData = this
      this.createRequest().get('/api/servers').then(function (response) {
        let responseData = response.data
        if (responseData.success) {
          vueData.servers = responseData.data
        } else {
          vueData.error = responseData.msg
        }
      }).catch(function (error) {
        let data = error.response.data
        let msg =  'unknown error'
        if (data) {
          msg = error
        } else if (msg.msg) {
          msg = msg.msg
        }
        vueData.error = msg
      })
    }
  },
  mounted() {
    this.loadData()
  }
}
</script>
