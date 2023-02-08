<script>
import { useI18n } from 'vue-i18n'
import Ace from '@/components/ui/Ace.vue'
import Btn from '@/components/ui/Btn.vue'
import Icon from '@/components/ui/Icon.vue'

const extensions = {}
// image formats
extensions['.jpg']  = { type: 'image', disableSave: true }
extensions['.jpeg'] = { type: 'image', disableSave: true }
extensions['.png']  = { type: 'image', disableSave: true }
extensions['.gif']  = { type: 'image', disableSave: true }
// audio formats
extensions['.mp3']  = { type: 'audio', disableSave: true }
extensions['.wav']  = { type: 'audio', disableSave: true }
extensions['.ogg']  = { type: 'audio', disableSave: true }
extensions['.flac'] = { type: 'audio', disableSave: true }
extensions['.aac']  = { type: 'audio', disableSave: true }
extensions['.alac'] = { type: 'audio', disableSave: true }
// video formats
extensions['.mp4']  = { type: 'video', disableSave: true }
extensions['.webm'] = { type: 'video', disableSave: true }
extensions['.avi']  = { type: 'video', disableSave: true }
extensions['.mkv']  = { type: 'video', disableSave: true }
extensions['.m4a']  = { type: 'video', disableSave: true }

function getType(file) {
  return (extensions[file.extension] || {}).type
}

export function skipDownload(file) {
  return ['image', 'audio', 'video'].indexOf(getType(file)) !== -1
}

export default {
  components: {
    Ace,
    Btn,
    Icon
  },
  props: {
    modelValue: { type: Object, required: true }
  },
  emits: ['update:modelValue', 'save', 'close'],
  setup(props, { emit }) {
const { t } = useI18n()

function emitUpdate(event) {
  emit('update:modelValue', { ...props.modelValue, content: event })
}

    return { t, emit, emitUpdate, extensions, getType, skipDownload }
  }
}
</script>

<template>
  <div>
    <div class="overlay-header">
      <h1 class="title" v-text="modelValue.name" />
      <btn v-if="!((extensions[modelValue.extension] || {}).disableSave)" variant="text" @click="emit('save')"><icon name="save" /> {{ t('common.Save') }}</btn>
      <btn v-hotkey="'Escape'" variant="icon" @click="emit('close')"><icon name="close" /></btn>
    </div>
    <img v-if="getType(modelValue) === 'image'" class="file-viewer" :src="modelValue.url" />
    <video v-else-if="getType(modelValue) === 'video'" class="file-viewer" controls>
      <source :src="modelValue.url" />
      <div class="warning unsupported" v-text="t('errors.VideoUnsupported')" />
    </video>
    <audio v-else-if="getType(modelValue) === 'audio'" class="file-viewer" controls>
      <source :src="modelValue.url" />
      <div class="warning unsupported" v-text="t('errors.AudioUnsupported')" />
    </audio>
    <ace v-else id="file-editor" :model-value="modelValue.content" class="file-editor" :file="modelValue.name" theme="monokai" @update:modelValue="emitUpdate" />
  </div>
</template>
