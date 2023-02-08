<script setup>
import { ref, onMounted, inject } from 'vue'
import { useI18n } from 'vue-i18n'
import Btn from '@/components/ui/Btn.vue'
import Icon from '@/components/ui/Icon.vue'

const { t } = useI18n()
const api = inject('api')

const props = defineProps({
  server: { type: Object, required: true }
})

const host = ref('')
const user = ref('')
const hostField = ref(null)
const userField = ref(null)
const hostCopied = ref(false)
const userCopied = ref(false)

onMounted(async () => {
  host.value = props.server.node.publicHost !== '127.0.0.1' ? props.server.node.publicHost : window.location.hostname
  host.value = host.value + ':' + props.server.node.sftpPort
  const u = await api.self.get()
  user.value = `${u.email}|${props.server.id}`
})

function copyHost() {
  hostField.value.select()
  document.execCommand('copy')
  userCopied.value = false
  hostCopied.value = true
  setTimeout(() => {
    hostCopied.value = false
  }, 6000)
}

function copyUser() {
  userField.value.select()
  document.execCommand('copy')
  hostCopied.value = false
  userCopied.value = true
  setTimeout(() => {
    userCopied.value = false
  }, 6000)
}
</script>

<template>
  <div>
    <h2 v-text="t('servers.SFTPInfo')" />
    <div>
      <b>{{t('common.Host')}}/{{t('common.Port')}}: </b>
      <span>{{host}}</span>
      <btn variant="icon" @click="copyHost()"><icon name="copy" /></btn>
      <span v-if="hostCopied" class="copied" v-text="t('common.Copied')" />
    </div>
    <input ref="hostField" :value="host" style="width:1px;height:1px;position:fixed;left:-100vw;" />
    <div>
      <b>{{t('users.Username')}}: </b>
      <span>{{user}}</span>
      <btn variant="icon" @click="copyUser()"><icon name="copy" /></btn>
      <span v-if="userCopied" class="copied" v-text="t('common.Copied')" />
    </div>
    <input ref="userField" :value="user" style="width:1px;height:1px;position:fixed;left:-100vw;" />
    <div>
      <b>{{t('users.Password')}}: </b>
      <span>{{t('users.AccountPassword')}}</span>
    </div>
    <a :href="`sftp://${user}@${host}`"><btn color="primary" v-text="t('servers.SftpConnection')" /></a>
  </div>
</template>
