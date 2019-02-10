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
            :items="items"
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
              slot="servers"
              slot-scope="{ item }"
            >
              <td>{{ item.name }}</td>
              <td>{{ item.node }}</td>
              <td>{{ item.address }}</td>
              <td class="text-xs-right">{{ item.online }}</td>
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
      this.createRequest().post('/api/servers').then(function (response) {
        let responseData = response.data
        if (responseData.success) {
          console.log(responseData)
        } else {
          data.error = responseData.msg
        }
      }).catch(function (error) {
        let msg = error.response.data.msg
        if (!msg) {
          msg = error
        }
        data.error = msg
      })
    }
  },
  mounted() {
    this.loadData();
  }
}
</script>
