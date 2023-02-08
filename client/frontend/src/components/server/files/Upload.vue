<script setup>
import { ref, inject } from 'vue'
import { useI18n } from 'vue-i18n'
import Btn from '@/components/ui/Btn.vue'
import Icon from '@/components/ui/Icon.vue'
import Overlay from '@/components/ui/Overlay.vue'

const { t } = useI18n()
const toast = inject('toast')

const props = defineProps({
  server: { type: Object, required: true },
  folder: { type: Boolean, default: () => false },
  hotkey: { type: String, default: () => undefined },
  path: { type: String, required: true }
})

const emit = defineEmits(['uploaded'])

const input = ref(null)
const uploading = ref(false)
const uploadState = ref(null)

function trigger() {
  input.value.click()
}

async function prepareFolders(files) {
  uploadState.value = { state: 'create folders' }
  const toCreate = new Set()
  const exist = new Set()
  for (let i = 0; i < files.length; i++) {
    let filePath = files[i].webkitRelativePath.split('/')
    filePath.pop()
    filePath = [props.path, ...filePath].join('/')
    if (toCreate.has(filePath) || exist.has(filePath)) continue
    const type = await props.server.fileExists(filePath)

    if (!type) toCreate.add(filePath)
    if (type === 'folder') exist.add(filePath)
    if (type === 'file') throw `${filePath} exists as file`
  }

  for (let path of toCreate) {
    await props.server.createFolder(path)
  }
}

function onUploadProgress(event) {
  const i = uploadState.value.current
  const file = uploadState.value.files[i]
  file.progress = Math.min(event.loaded, file.size)
}

async function uploadFiles(event) {
  uploadState.value = { state: 'preparing' }
  uploading.value = true
  try {
    if (event.target.webkitdirectory) await prepareFolders(event.target.files)
    uploadState.value = {
      state: 'upload files',
      current: 0,
      total: event.target.files.length,
      files: []
    }
    for (let i = 0; i < event.target.files.length; i++) {
      const file = event.target.files[i]
      uploadState.value.files[i + 1] = {
        name: file.webkitRelativePath || file.name,
        size: file.size,
        progress: 0
      }
    }
    for (let i = 0; i < event.target.files.length; i++) {
      uploadState.value.current = i + 1
      const file = event.target.files[i]
      let path = props.path + '/'
      if (file.webkitRelativePath) {
        path = path + file.webkitRelativePath
      } else {
        path = path + file.name
      }
      await props.server.uploadFile(path, file, onUploadProgress)
    }
  } catch(e) {
    console.error('file upload failed', e)
    toast.error(t('files.UploadFailed'))
  } finally {
    toast.success(t('files.UploadSuccess'))
    emit('uploaded')
    uploading.value = false
    uploadState.value = null
  }
}
</script>

<template>
  <div>
    <btn v-hotkey="hotkey" variant="icon" @click="trigger()">
      <icon :name="folder ? 'folder-upload' : 'file-upload'" />
    </btn>
    <input ref="input" type="file" multiple :webkitdirectory="folder" @change="uploadFiles" />
    <overlay v-model="uploading" :title="t('files.UploadProgress')">
      <div v-if="uploadState.state === 'preparing'">
        <div class="progress">
          <span v-text="t('files.PreparingUpload')" />
          <progress />
        </div>
      </div>
      <div v-if="uploadState.state === 'create folders'">
        <div class="progress">
          <span v-text="t('files.CreatingFolders')" />
          <progress />
        </div>
      </div>
      <div v-if="uploadState.state === 'upload files'">
        <div class="upload-file-count" v-text="t('files.CurrentlyUploading', { current: uploadState.current, total: uploadState.total })" />
        <div class="upload-file-name" v-text="uploadState.files[uploadState.current].name" />
        <div class="progress">
          <span v-text="t('files.Current')" />
          <progress
            :max="uploadState.files[uploadState.current].size"
            :value="uploadState.files[uploadState.current].progress"
          />
        </div>
        <div class="progress">
          <span v-text="t('files.Total')" />
          <progress
            :max="uploadState.files.reduce((a, b) => a + b.size, 0)"
            :value="uploadState.files.reduce((a, b) => a + b.progress, 0)"
          />
        </div>
      </div>
    </overlay>
  </div>
</template>
