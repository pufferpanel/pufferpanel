<!--
  - Copyright 2019 Padduck, LLC
  -  Licensed under the Apache License, Version 2.0 (the "License");
  -  you may not use this file except in compliance with the License.
  -  You may obtain a copy of the License at
  -          http://www.apache.org/licenses/LICENSE-2.0
  -  Unless required by applicable law or agreed to in writing, software
  -  distributed under the License is distributed on an "AS IS" BASIS,
  -  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  -  See the License for the specific language governing permissions and
  -  limitations under the License.
  -->

<template>
  <v-card>
    <v-card-title>
      <span v-text="$t('files.FileManager') + ' - ' + currentPath" />
      <v-btn
        v-if="(server.permissions.putServerFiles || isAdmin()) && !createFolder"
        icon
        @click="createFolder = true"
      >
        <v-icon>mdi-folder-plus</v-icon>
      </v-btn>
      <div v-if="createFolder">
        <v-text-field
          v-model="newFolderName"
          hide-details
          class="input-small ml-2 mt-0 pt-0"
          :placeholder="$t('files.NewFolder')"
        />
        <v-btn
          icon
          color="success"
          @click="submitNewFolder"
        >
          <v-icon>mdi-check</v-icon>
        </v-btn>
        <v-btn
          icon
          color="error"
          @click="cancelFolderCreate"
        >
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </div>
    </v-card-title>
    <v-card-text>
      <v-dialog v-model="confirmDeleteOpen" max-width="600">
        <v-card>
          <v-card-title v-text="$t('files.ConfirmDelete')" />
          <v-card-actions>
            <v-spacer />
            <v-btn v-text="$t('common.Cancel')" @click="deleteCancelled()" color="error" />
            <v-btn v-text="$t('common.Confirm')" @click="deleteConfirmed()" color="success" />
          </v-card-actions>
        </v-card>
      </v-dialog>
      <v-data-table
        :items-per-page="-1"
        :headers="headers"
        :no-data-text="$t('common.NoFiles')"
        hide-default-footer
        :items="files"
        @click:row="itemClicked"
      >
        <template v-slot:item.size="{ item }">
          {{ item.size ? toFileSize(item.size) : '' }}
        </template>
        <template v-slot:item.modifyTime="{ item }">
          {{ item.modifyTime ? toDate(item.modifyTime) : '' }}
        </template>
        <template v-slot:item.isFile="{ item }">
          <v-tooltip
            v-if="item.isFile && !(item.size > maxEditSize)"
            bottom
          >
            <template v-slot:activator="{ on }">
              <!-- letting the click propagate through to the table row, does the same anyways -->
              <v-btn
                icon
                v-on="on"
              >
                <v-icon>mdi-pencil</v-icon>
              </v-btn>
            </template>
            <span v-text="$t('common.Edit')" />
          </v-tooltip>
          <v-tooltip
            v-if="item.isFile"
            bottom
          >
            <template v-slot:activator="{ on }">
              <v-btn
                icon
                :href="createDownloadLink(item)"
                target="_blank"
                v-on="on"
              >
                <v-icon>mdi-download</v-icon>
              </v-btn>
            </template>
            <span v-text="$t('common.Download')" />
          </v-tooltip>
          <v-tooltip
            v-if="item.name !== '..'"
            bottom
          >
            <template v-slot:activator="{ on }">
              <v-btn
                icon
                v-on="on"
                @click.stop="deleteRequest(item)"
              >
                <v-icon>mdi-trash-can</v-icon>
              </v-btn>
            </template>
            <span v-text="$t('common.Delete')" />
          </v-tooltip>
        </template>
      </v-data-table>

      <v-overlay :value="editOpen">
        <v-card
          :dark="isDark()"
          :light="!isDark()"
          class="d-flex flex-column"
          height="90vh"
          width="90vw"
        >
          <v-card-title v-text="currentFile" />
          <v-card-text
            id="editor"
            class="flex-grow-1"
          >
            <ace
              v-model="fileContents"
              editor-id="fileEditor"
              :theme="isDark() ? 'monokai' : 'github'"
	      :file="currentFile"
            />
          </v-card-text>
          <v-card-actions class="px-4 pb-4">
            <div class="flex-grow-1" />
            <v-btn v-text="$t('common.Cancel')" color="error" @click="cancelEdit()" />
            <v-btn v-text="$t('common.Save')" color="success" @click="saveEdit()" />
          </v-card-actions>
        </v-card>
      </v-overlay>

      <div v-if="server.permissions.putServerFiles || isAdmin()">
        <v-file-input
          v-model="uploadFiles"
          multiple
          chips
          show-size
          counter
          :placeholder="$t('files.Upload')"
          class="mb-4"
        />
        <div v-if="uploading">
          <v-progress-linear
            :value="(uploadCurrent / uploadSize) * 100"
            buffer-value="0"
            stream
            class="mb-4"
          />
        </div>
        <v-btn
          block
          color="primary"
          :disabled="!(uploadFiles.length > 0) || uploading"
          @click="transmitFiles"
          v-text="$t('files.Upload')"
        />
      </div>
    </v-card-text>
  </v-card>
</template>

<script>
import filesize from 'filesize'
import { isDark } from '@/utils/dark'

export default {
  props: {
    server: { type: Object, default: function () { return {} } }
  },
  data () {
    return {
      files: [],
      currentPath: '/',
      loading: true,
      headers: [
        {
          value: 'name',
          text: this.$t('common.Name'),
          sortable: true
        },
        {
          value: 'size',
          text: this.$t('common.Size'),
          sortable: true
        },
        {
          value: 'modifyTime',
          text: this.$t('files.LastModified'),
          sortable: true
        },
        {
          value: 'isFile',
          text: this.$t('common.Actions'),
          sortable: false
        }
      ],
      currentFile: '',
      fileContents: '',
      editOpen: false,
      maxEditSize: 1024 * 1024 * 20,
      createFolder: false,
      newFolderName: '',
      uploadFiles: [],
      uploading: false,
      uploadCurrent: 0,
      uploadSize: 0,
      toDelete: null,
      confirmDeleteOpen: false
    }
  },
  mounted () {
    const vue = this
    this.$socket.addEventListener('open', function (event) {
      vue.fetchItems(vue.currentPath)
    })

    this.$socket.addEventListener('message', function (event) {
      const data = JSON.parse(event.data)
      if (data === 'undefined') {
        return
      }
      if (data.type === 'file') {
        if (data.data) {
          if (data.data.error) {
            vue.isLoading = false
            return
          }

          vue.files = (data.data.files || []).sort(function (a, b) {
            if (!a.size && !b.size) return 0
            if (a.size && b.size) return 0
            if (a.size && !b.size) return 1
            return -1
          })
          if (data.data.path !== '') {
            vue.currentPath = data.data.path
          }
          vue.loading = false
        }
      }
    })
  },
  methods: {
    fetchItems (path, edit = false) {
      this.loading = true
      this.$socket.sendObj({ type: 'file', action: 'get', path: path, edit: edit })
    },
    itemClicked (item) {
      if (!item.isFile) {
        this.loading = true

        if (item.name === '..') {
          const parts = this.currentPath.split('/')
          parts.pop()
          this.currentPath = parts.join('/')
          if (this.currentPath === '') {
            this.currentPath = '/'
          }
        } else {
          let path = this.currentPath
          if (path === '/') {
            path += item.name
          } else {
            path += '/' + item.name
          }
          this.currentPath = path
        }

        this.$socket.sendObj({ type: 'file', action: 'get', path: this.currentPath })
      } else {
        if (item.size > this.maxEditSize) return
        let path = this.currentPath
        if (path === '/') {
          path += item.name
        } else {
          path += '/' + item.name
        }
        const ctx = this
        this.$http.get(`/daemon/server/${this.server.id}/file/${path}`).then(function (response) {
          ctx.currentFile = item.name
          ctx.fileContents = response.data
          ctx.editOpen = true
        }).catch(function () {
          ctx.$toast.error(ctx.$t('common.FileLoadFailed'))
        })
      }
    },
    cancelEdit () {
      this.editOpen = false
      this.currentFile = ''
      this.fileContents = ''
    },
    saveEdit () {
      let path = this.currentPath
      if (path === '/') {
        path += this.currentFile
      } else {
        path += '/' + this.currentFile
      }
      const file = new Blob([this.fileContents])
      const formData = new FormData()
      formData.append('file', file)
      const ctx = this
      this.$http.put(`/daemon/server/${this.server.id}/file/${path}`, formData, { headers: { 'Content-Type': 'multipart/form-data' } }).then(function (response) {
        ctx.editOpen = false
        ctx.currentFile = ''
        ctx.fileContents = ''
        ctx.$toast.success(ctx.$t('common.Saved'))
      }).catch(function () {
        ctx.$toast.error(ctx.$t('common.SaveFailed'))
      })
    },
    deleteRequest (item) {
      this.toDelete = item
      this.confirmDeleteOpen = true
    },
    deleteConfirmed () {
      let path = ''
      if (this.currentPath === '/') {
        path = '/' + this.toDelete.name
      } else {
        path = this.currentPath + '/' + this.toDelete.name
      }
      this.loading = true
      this.$socket.sendObj({ type: 'file', action: 'delete', path: path })
      this.toDelete = null
      this.confirmDeleteOpen = false
    },
    deleteCancelled () {
      this.toDelete = null
      this.confirmDeleteOpen = null
    },

    // utility
    toFileSize (size) {
      return filesize(size)
    },
    toDate (epoch) {
      return new Date(epoch * 1000).toLocaleString()
    },
    createDownloadLink (item) {
      let path = this.currentPath
      if (path === '/') {
        path += item.name
      } else {
        path += '/' + item.name
      }
      return '/daemon/server/' + this.server.id + '/file' + path
    },
    cancelFolderCreate () {
      this.createFolder = false
      this.newFolderName = ''
    },
    submitNewFolder () {
      let path = this.currentPath
      if (path === '/') {
        path = path + this.newFolderName
      } else {
        path = path + '/' + this.newFolderName
      }

      this.$socket.sendObj({ type: 'file', action: 'create', path: path })
      this.createFolder = false
      this.newFolderName = ''
    },
    transmitFiles () {
      this.uploading = true
      this.uploadNextItem(this)
    },
    uploadNextItem (vue) {
      this.uploadSingleFile(vue.uploadFiles[0]).then(function () {
        vue.uploadFiles.shift()
        if (vue.uploadFiles.length === 0) {
          vue.uploading = false
          vue.fetchItems(vue.currentPath)
          return
        }
        vue.uploadNextItem(vue)
      })
    },
    uploadSingleFile (item) {
      let path = this.currentPath
      if (path === '/') {
        path += item.name
      } else {
        path += '/' + item.name
      }
      this.uploadCurrent = 0
      this.uploadSize = item.size

      const vue = this
      return this.$http({
        method: 'put',
        url: '/daemon/server/' + this.server.id + '/file' + path,
        data: item,
        onUploadProgress: function (event) {
          vue.uploadCurrent = event.loaded
          vue.uploadSize = event.total
        }
      })
    },
    isDark
  }
}
</script>

<style>
  #editor .v-input__control, #editor .v-input__slot, #editor .v-text-field__slot, #editor textarea {
    height: 100%;
  }
</style>
