<script setup>
import { ref, inject, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/ui/Icon.vue'
import Loader from '@/components/ui/Loader.vue'

const { t } = useI18n()
const api = inject('api')
const nodesLoaded = ref(false)
const nodes = ref([])
const firstEntry = ref(null)

onMounted(async () => {
  nodes.value = await api.node.list()
  nodesLoaded.value = true
})

function setFirstEntry(ref) {
  if (!firstEntry.value) firstEntry.value = ref
}

function focusList() {
  firstEntry.value.$el.focus()
}
</script>

<template>
  <div class="nodelist">
    <h1 v-text="t('nodes.Nodes')" />
    <div v-hotkey="'l'" class="list" @hotkey="focusList()">
      <div v-for="node in nodes" :key="node.name" class="list-item">
        <router-link v-if="node.id" :ref="setFirstEntry" :to="{ name: 'NodeView', params: { id: node.id } }">
          <div class="node">
            <span class="title">{{node.name}}</span>
            <span class="subline">{{node.publicHost + ':' + node.publicPort}}</span>
          </div>
        </router-link>
        <div v-else class="node disabled">
          <span class="title">{{node.name}}</span>
          <span class="subline">{{node.publicHost + ':' + node.publicPort}}</span>
        </div>
      </div>
      <div v-if="!nodesLoaded" class="list-item">
        <loader small />
      </div>
      <div v-if="$api.auth.hasScope('nodes.create')" class="list-item">
        <router-link v-hotkey="'c'" :to="{ name: 'NodeCreate' }">
          <div class="createLink"><icon name="plus" />{{ t('nodes.Add') }}</div>
        </router-link>
      </div>
    </div>
  </div>
</template>
