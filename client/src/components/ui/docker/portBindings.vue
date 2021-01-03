<template>
  <v-container>
    <v-row
      v-for="(entry, i) in value"
      :key="i"
    >
      <v-col cols="12">
        <v-divider />
      </v-col>
      <v-col
        cols="12"
        class="d-flex align-center"
      >
        <div>{{ entry }}</div>
        <div class="flex-grow-1" />
        <v-btn
          icon
          @click="startEdit(i)"
        >
          <v-icon>mdi-pencil</v-icon>
        </v-btn>
        <v-btn
          icon
          @click="remove(i)"
        >
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12">
        <v-btn
          text
          block
          @click="add()"
          v-text="$t('common.Add')"
        />
      </v-col>
    </v-row>
    <ui-overlay
      v-model="edit"
      card
      closable
      @close="reset()"
    >
      <v-row>
        <v-col
          cols="12"
          lg="3"
        >
          <ui-input
            v-model="host"
            :label="$t('common.Host')"
          />
        </v-col>
        <v-col
          cols="12"
          lg="3"
        >
          <ui-input
            v-model="outsidePort"
            type="number"
            :label="$t('env.docker.OutsidePort')"
          />
        </v-col>
        <v-col
          cols="12"
          lg="3"
        >
          <ui-input
            v-model="insidePort"
            type="number"
            :label="$t('env.docker.InsidePort')"
          />
        </v-col>
        <v-col
          cols="12"
          lg="3"
        >
          <ui-input-suggestions
            v-model="protocol"
            :label="$t('common.Protocol')"
            :items="['tcp', 'udp']"
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="12">
          <v-btn
            color="success"
            block
            @click="save()"
            v-text="$t('common.Save')"
          />
        </v-col>
      </v-row>
    </ui-overlay>
  </v-container>
</template>

<script>
export default {
  props: {
    value: { type: Array, required: true }
  },
  data () {
    return {
      edit: false,
      new: true,
      editIndex: 0,
      host: '0.0.0.0',
      outsidePort: '',
      insidePort: '',
      protocol: 'tcp'
    }
  },
  methods: {
    onInput (index, event) {
      const changed = [...this.value]
      changed[index] = event
      this.$emit('input', changed)
    },
    remove (index) {
      const changed = [...this.value]
      changed.splice(index, 1)
      this.$emit('input', changed)
    },
    reset () {
      this.host = '0.0.0.0'
      this.outsidePort = ''
      this.insidePort = ''
      this.protocol = 'tcp'
      this.edit = false
      this.new = true
      this.editIndex = 0
    },
    add () {
      this.new = true
      this.edit = true
    },
    startEdit (i) {
      this.editIndex = i
      this.new = false

      const binding = this.value[i]
      this.host = binding.split(':')[0]
      this.outsidePort = binding.split(':')[1]
      this.insidePort = binding.split(':')[2].split('/')[0]
      this.protocol = binding.split('/')[1]

      this.edit = true
    },
    save () {
      const changed = [...this.value]
      const binding = `${this.host}:${this.outsidePort}:${this.insidePort}/${this.protocol}`
      if (this.new) {
        changed.push(binding)
      } else {
        changed[this.editIndex] = binding
      }
      this.reset()
      this.$emit('input', changed)
    }
  }
}
</script>
