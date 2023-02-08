<script setup>
import { ref, inject, onMounted, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import Multiselect from '@vueform/multiselect'
import Btn from '@/components/ui/Btn.vue'
import Icon from '@/components/ui/Icon.vue'

const { t } = useI18n()
const api = inject('api')
const emit = defineEmits(['back', 'confirm'])

const users = ref([])
const ms = ref(null)

onMounted(async () => {
  const self = await api.self.get()
  nextTick(() => {
    ms.value.select({ value: self.username, label: self.username })
  })
})

async function searchUsers(query) {
  const res = await api.user.search(query)
  return res.map(u => {
    return {
      value: u.username,
      label: u.username
    }
  })
}

function remove(user) {
  users.value = users.value.filter(u => u.name !== user)
}

function confirm() {
  if (users.value.length === 0) return
  emit('confirm', users.value)
}
</script>

<template>
  <div class="users">
    <multiselect
      ref="ms"
      v-model="users"
      mode="tags"
      placeholder="t('server.SearchUsers')"
      :close-on-select="false"
      :can-clear="false"
      :filter-results="false"
      :min-chars="1"
      :resolve-on-load="false"
      :delay="500"
      :searchable="true"
      :options="searchUsers"
    />
    <btn color="error" @click="emit('back')"><icon name="back" />{{ t('common.Back') }}</btn>
    <btn color="primary" :disabled="users.length === 0" @click="confirm()"><icon name="check" />{{ t('common.Next') }}</btn>
  </div>
</template>
