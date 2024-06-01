<script setup>
import { ref, onMounted, onUnmounted, inject } from 'vue'
import { useI18n } from 'vue-i18n'
import Btn from '@/components/ui/Btn.vue'
import Icon from '@/components/ui/Icon.vue'
import Loader from '@/components/ui/Loader.vue'
import Overlay from '@/components/ui/Overlay.vue'
import Editor, { skipDownload } from './files/Editor.vue'
import Upload from './files/Upload.vue'
import TextField from '@/components/ui/TextField.vue'

const { t } = useI18n()
const events = inject('events')

const props = defineProps({
  server: { type: Object, required: true }
})

const allowDirectoryUpload = 'webkitdirectory' in document.createElement('input')
const canEdit = props.server.hasScope('server.files.edit')

const fileEls = ref([])
const files = ref(null)
const file = ref(null)
const fileSizeWarn = ref(false)
const fileSizeWarnSubject = ref(null)
const currentPath = ref([])
const editorOpen = ref(false)
const loading = ref(false)
const createFileOpen = ref(false)
const createFolderOpen = ref(false)
const newItemName = ref('')

let task
let unbindEvent
onMounted(() => {
  refresh()

  unbindEvent = props.server.on('status', () => {
    refresh()
  })

  task = props.server.startTask(() => {
    refresh()
  }, 5 * 60 * 1000)
})

onUnmounted(async () => {
  if (unbindEvent) unbindEvent()
  if (task) props.server.stopTask(task)
})

async function refresh(manual = false) {
  if (manual) files.value = null // cause visual feedback on manual refresh
  const res = await props.server.getFile(getCurrentPath())
  files.value = res.sort(sortFiles)
}

function sortFiles(a, b) {
  if (a.isFile && !b.isFile) return 1
  if (!a.isFile && b.isFile) return -1
  if (a.name.toLowerCase() < b.name.toLowerCase()) return -1
  return 1
}

function getCurrentPath() {
  return currentPath.value.map(e => e.name).join('/')
}

async function openFile(f, overrideWarn = false) {
  if (f.isFile) {
    if (!skipDownload(f) && !overrideWarn && f.size > 30 * Math.pow(2, 20)) {
      fileSizeWarnSubject.value = f
      fileSizeWarn.value = true
      return
    }

    fileSizeWarn.value = false
    loading.value = true
    const path = getCurrentPath() + `/${f.name}`
    const content = skipDownload(f) ? null : await props.server.getFile(path, true)
    file.value = { ...f, content, url: props.server.getFileUrl(path) }
    editorOpen.value = true
    loading.value = false
  } else {
    let path = ''
    if (f.name === '..') {
      currentPath.value.pop()
      path = getCurrentPath()
    } else {
      path = getCurrentPath() + `/${f.name}`
    }
    const res = await props.server.getFile(path)
    files.value = res.sort(sortFiles)
    if (f.name !== '..') currentPath.value.push(f)
  }
}

async function saveFile() {
  await props.server.uploadFile(`${getCurrentPath()}/${file.value.name}`, file.value.content)
  editorOpen.value = false
  file.value = null
  refresh()
}

function getIcon(file) {
  if (!file.isFile) return 'folder'
  if (!file.extension) return 'file'
  return 'file-' + file.extension.substring(1)
}

function deleteFile(file) {
  events.emit(
    'confirm',
    t('files.ConfirmDelete', { name: file.name }),
    {
      text: t('files.Delete'),
      icon: 'remove',
      color: 'error',
      action: async () => {
        await props.server.deleteFile(getCurrentPath() + '/' + file.name)
        await refresh()
      }
    },
    {
      color: 'primary'
    }
  )
}

const numFormat = new Intl.NumberFormat('en-US', { maximumFractionDigits: 2 })
function formatFileSize(size) {
  if (!size) return '0 B'
  if (size < Math.pow(2, 10)) return numFormat.format(size) + ' B'
  if (size < Math.pow(2, 20)) return numFormat.format(size / Math.pow(2, 10)) + ' KiB'
  if (size < Math.pow(2, 30)) return numFormat.format(size / Math.pow(2, 20)) + ' MiB'
  if (size < Math.pow(2, 40)) return numFormat.format(size / Math.pow(2, 30)) + ' GiB'
  return numFormat.format(size / Math.pow(2, 40)) + ' TiB'
}

function startCreateFile() {
  newItemName.value = ''
  createFileOpen.value = true
}

function startCreateFolder() {
  newItemName.value = ''
  createFolderOpen.value = true
}

async function createFile() {
  if (!newItemName.value || newItemName.value.trim() === '') return
  await props.server.uploadFile(`${getCurrentPath()}/${newItemName.value}`, '')
  const file = { name: newItemName.value, size: 0, isFile: true }
  createFileOpen.value = false
  newItemName.value = ''
  await openFile(file)
  await refresh()
}

async function createFolder() {
  if (!newItemName.value || newItemName.value.trim() === '') return
  await props.server.createFolder(`${getCurrentPath()}/${newItemName.value}`)
  const folder = { name: newItemName.value, isFile: false }
  createFolderOpen.value = false
  newItemName.value = ''
  await openFile(folder)
  await refresh()
}

const archiveExtensions = [
  '.7z',
  '.bz2',
  '.gz',
  '.lz',
  '.lzma',
  '.rar',
  '.tar',
  '.tgz',
  '.xz',
  '.zip',
  '.zipx'
]

function isArchive (file) {
  const filename = file.name.toLowerCase()
  for (let i = 0; i < archiveExtensions.length; i++) {
    if (filename.endsWith(archiveExtensions[i])) return true
  }
  return false
}

function extract(file) {
  loading.value = true
  try {
    props.server.extractFile(`${getCurrentPath()}/${file.name}`, getCurrentPath())
    refresh()
  } finally {
    loading.value = false
  }
}

async function makeArchiveName(fileName) {
  let destination = `${getCurrentPath()}/${fileName}.zip`
  for (let i = 2; await props.server.fileExists(destination); i++) {
    destination = `${getCurrentPath()}/${fileName} (${i}).zip`
  }
  return destination
}

async function archiveCurrentDirectory() {
  loading.value = true
  try {
    const item = currentPath.value[currentPath.value.length - 1];
    let lastPathEntry = props.server.id
    if (item !== undefined) {
      lastPathEntry = currentPath.value[currentPath.value.length - 1].name
    }

    await props.server.archiveFile(
      await makeArchiveName(lastPathEntry),
      `${getCurrentPath()}`
    )
  } finally {
    setTimeout(() => {
      refresh()
      loading.value = false
    }, 500)
  }
}

async function archive(file) {
  loading.value = true
  try {
    await props.server.archiveFile(
      await makeArchiveName(file.name),
      `${getCurrentPath()}/${file.name}`
    )
  } finally {
    setTimeout(() => {
      refresh()
      loading.value = false
    }, 500)
  }
}

function downloadLink(file) {
  return props.server.getFileUrl(getCurrentPath() + '/' + file.name)
}

function fileListHotkey() {
  if (fileEls.value[0]) fileEls.value[0].focus()
}

function trackFileEl(index) {
  return (el) => fileEls.value[index] = el
}
</script>

<template>
  <div class="file-manager">
    <div class="header">
      <h2 v-text="t('servers.Files')" />
      <h3 v-text="'/' + getCurrentPath()" />
      <span class="spacer" />
      <btn v-if="canEdit" v-hotkey="'f a'" variant="icon" :tooltip="t('files.ArchiveCurrent')" @click="archiveCurrentDirectory()"><icon name="archive" /></btn>
      <upload v-if="canEdit" :path="getCurrentPath()" :server="server" hotkey="f u" @uploaded="refresh()" />
      <upload v-if="canEdit && allowDirectoryUpload" :path="getCurrentPath()" :server="server" folder hotkey="f d" @uploaded="refresh()" />
      <btn v-if="canEdit" v-hotkey="'f c f'" variant="icon" :tooltip="t('files.CreateFile')" @click="startCreateFile()"><icon name="file-create" /></btn>
      <btn v-if="canEdit" v-hotkey="'f c d'" variant="icon" :tooltip="t('files.CreateFolder')" @click="startCreateFolder()"><icon name="folder-create" /></btn>
      <btn v-hotkey="'f r'" variant="icon" :tooltip="t('files.Refresh')" @click="refresh(true)"><icon name="reload" /></btn>
    </div>
    <div v-hotkey="'f l'" class="file-list" @hotkey="fileListHotkey">
      <loader v-if="!Array.isArray(files)" />
      <!-- eslint-disable-next-line vue/no-template-shadow -->
      <a v-for="(file, index) in files" v-else :key="file.name" :ref="trackFileEl(index)" tabindex="0" class="file" @click="openFile(file)" @keydown.enter="openFile(file)">
        <icon class="file-icon" :name="getIcon(file)" />
        <div class="details">
          <div class="name">{{ file.name }}</div>
          <div v-if="file.isFile" class="size">{{ formatFileSize(file.size) }}</div>
        </div>
        <btn v-if="canEdit && file.name !== '..' && !file.isFile" tabindex="-1" variant="icon" :tooltip="t('files.Archive')" @click.stop="archive(file)">
          <icon name="archive" />
        </btn>
        <btn v-if="canEdit && file.isFile && isArchive(file)" tabindex="-1" variant="icon" :tooltip="t('files.Extract')" @click.stop="extract(file)">
          <icon name="extract" />
        </btn>
        <a v-if="file.isFile" tabindex="-1" class="dl-link" :href="downloadLink(file)" target="_blank" rel="noopener">
          <btn tabindex="-1" variant="icon" :tooltip="t('files.Download')" @click.stop="">
            <icon name="download" />
          </btn>
        </a>
        <btn v-if="canEdit && file.name !== '..'" tabindex="-1" variant="icon" :tooltip="t('files.Delete')" @click.stop="deleteFile(file)">
          <icon name="remove" />
        </btn>
      </a>
    </div>
    <overlay v-model="fileSizeWarn" closable :title="t('files.OpenLargeFile')">
      <btn color="error" @click="fileSizeWarn = false"><icon name="close" />{{ t('common.Cancel') }}</btn>
      <btn color="primary" @click="openFile(fileSizeWarnSubject, true)"><icon name="check" />{{ t('files.OpenAnyways') }}</btn>
    </overlay>
    <overlay v-model="createFileOpen" closable :title="t('files.CreateFile')">
      <text-field v-model="newItemName" />
      <btn color="primary" :disabled="!newItemName || newItemName.trim() === ''" @click="createFile()"><icon name="check" />{{ t('files.CreateFile') }}</btn>
    </overlay>
    <overlay v-model="createFolderOpen" closable :title="t('files.CreateFolder')">
      <text-field v-model="newItemName" />
      <btn color="primary" :disabled="!newItemName || newItemName.trim() === ''" @click="createFolder()"><icon name="check" />{{ t('files.CreateFolder') }}</btn>
    </overlay>
    <overlay v-model="loading" class="loader-overlay">
      <loader />
    </overlay>
    <overlay v-model="editorOpen" class="editor">
      <editor v-if="file" v-model="file" :read-only="!canEdit" @save="saveFile()" @close="editorOpen = false" />
    </overlay>
  </div>
</template>
