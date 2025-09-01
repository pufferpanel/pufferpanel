<script setup>
import {ref, inject, onMounted, computed} from 'vue'
import { useI18n } from 'vue-i18n'
const events = inject('events')
import Loader from '@/components/ui/Loader.vue'
import Btn from '@/components/ui/Btn.vue'
import Icon from '@/components/ui/Icon.vue'
import TextField from '@/components/ui/TextField.vue'

const { t, locale } = useI18n()
const toast = inject('toast')

const props = defineProps({
  server: { type: Object, required: true }
})

const backups = ref(null)
const backupName = ref("")
const backupRunning = ref(false)
const loading = ref(false)
const sortedBackups = computed(() => backups.value.slice().sort((a, b) => b.createdAt.localeCompare(a.createdAt)));

onMounted(async () => {
  await loadBackups()
})

async function loadBackups() {
  backups.value = await props.server.getBackups()
}

function isBackingUp() {
  return backupRunning.value
}

function isLoading() {
  return !Array.isArray(backups.value) || loading.value
}

async function save() {
  try {
    backupRunning.value = true
    await props.server.createBackup(backupName.value)
    toast.success(t('backup.BackupStarted'))
    await loadBackups()
  }
  finally {
    backupRunning.value = false
  }
}

/*
const numFormat = new Intl.NumberFormat('en-US', { maximumFractionDigits: 2 })
function formatFileSize(size) {
  if (!size) return '0 B'
  if (size < Math.pow(2, 10)) return numFormat.format(size) + ' B'
  if (size < Math.pow(2, 20)) return numFormat.format(size / Math.pow(2, 10)) + ' KiB'
  if (size < Math.pow(2, 30)) return numFormat.format(size / Math.pow(2, 20)) + ' MiB'
  if (size < Math.pow(2, 40)) return numFormat.format(size / Math.pow(2, 30)) + ' GiB'
  return numFormat.format(size / Math.pow(2, 40)) + ' TiB'
}
*/

function promptRestore(file){
  events.emit(
      'confirm',
      {
        title: t('backup.RestorePrompt'),
        body: t('backup.RestorePromptBody'),
      },
      {
        text: t('backup.Restore'),
        icon: 'remove',
        action: () => {
          restore(file)
        }
      },
      {
        color: 'neutral'
      }
    )
}

async function restore(file) {
  try {
    loading.value = true
    await props.server.restoreBackup(file.id);
    toast.success(t('backup.RestoreStarted'))
    await loadBackups()
  }
  finally {
    loading.value = false
  }
}

function promptDelete(file){
  events.emit(
      'confirm',
      {
        title: t('backup.DeletePrompt'),
        body: t('backup.DeletePromptBody'),
      },
      {
        text: t('backup.Delete'),
        icon: 'restore',
        color: 'error',
        action: () => {
          deleteBackup(file)
        }
      },
      {
        color: 'primary'
      }
    )
}

async function deleteBackup(file) {
  try {
    loading.value = true
    await props.server.deleteBackup(file.id);
    toast.success(t('backup.Deleted'))
    await loadBackups()
  }
  finally {
    loading.value = false
  }
}

const intl = new Intl.DateTimeFormat(
  [locale.value.replace('_', '-'), 'en'],
  { day: '2-digit', month: '2-digit', year: 'numeric', hour: 'numeric', minute: 'numeric', second: 'numeric' }
)

</script>

<template>
  <div class="backup-manager">
    <h2 v-text="t('backup.Backup')" />
    <div v-if="server.hasScope('server.backup.create')">
      <text-field v-model="backupName" :label="t('backup.Name')" />
      <btn color="primary" :disabled="isBackingUp() || isLoading()" @click="save()">
        <icon v-if="!isBackingUp()" name="plus" />
        <icon v-else name="loading" spin /> {{ t('backup.Create') }}
      </btn>
    </div>

    <div class="group-header">
      <div class="title">
        <h3 v-text="t('backup.BackupsHeader')" />
      </div>
    </div>
    <div class="backup-list">
      <loader v-if="isLoading()" />
      <!-- eslint-disable-next-line vue/no-template-shadow -->
      <div v-for="backup in sortedBackups" v-else :key="backup.id" tabindex="0" class="backup">
        <icon class="file-icon" name="file" />
        <div class="details">
          <div class="name">{{ backup.name }} ({{ intl.format(new Date(backup.createdAt)) }})</div>
          <!--<div class="size">{{ formatFileSize(backup.fileSize) }}</div> -->
        </div>
        <btn
          v-if="server.hasScope('server.backup.restore')"
          tabindex="-1"
          variant="icon"
          :tooltip="t('backup.Restore')"
          :disabled="isBackingUp()"
          @click.stop="promptRestore(backup)"
        >
          <icon name="restore" />
        </btn>
        <a tabindex="-1" class="dl-link" :href="props.server.getBackupUrl(backup.id)" target="_blank" rel="noopener">
          <btn tabindex="-1" variant="icon" :tooltip="t('backup.Download')">
            <icon name="download" />
          </btn>
        </a>
        <btn
          v-if="server.hasScope('server.backup.delete')"
          tabindex="-1"
          variant="icon"
          :tooltip="t('backup.Delete')"
          :disabled="isBackingUp()"
          @click.stop="promptDelete(backup)"
        >
          <icon name="remove" />
        </btn>
      </div>
    </div>
  </div>
</template>
