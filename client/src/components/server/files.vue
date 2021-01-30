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
  <v-card
    :loading="loading"
    :disabled="loading"
  >
    <v-card-title>
      <span
        class="flex-grow-1"
        v-text="$t('files.FileManager')"
      />
      <v-btn
        v-if="server.permissions.putServerFiles || isAdmin()"
        icon
        @click="createFile = true"
      >
        <v-icon>mdi-file-plus</v-icon>
      </v-btn>
      <ui-overlay
        v-model="createFile"
        :title="$t('files.NewFile')"
        card
      >
        <v-row>
          <v-col>
            <ui-input
              v-model="newFileName"
              hide-details
              autofocus
              :label="$t('common.Name')"
              @keyup.esc="cancelFileCreate()"
              @keyup.enter="confirmFileCreate()"
            />
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <v-btn
              block
              color="error"
              @click="cancelFileCreate()"
              v-text="$t('common.Cancel')"
            />
          </v-col>
          <v-col>
            <v-btn
              block
              color="success"
              :disabled="newFileName === ''"
              @click="confirmFileCreate()"
              v-text="$t('common.Create')"
            />
          </v-col>
        </v-row>
      </ui-overlay>
      <v-btn
        v-if="server.permissions.putServerFiles || isAdmin()"
        icon
        @click="createFolder = true"
      >
        <v-icon>mdi-folder-plus</v-icon>
      </v-btn>
      <ui-overlay
        v-model="createFolder"
        :title="$t('files.NewFolder')"
        card
      >
        <v-row>
          <v-col>
            <ui-input
              v-model="newFolderName"
              hide-details
              autofocus
              :label="$t('common.Name')"
              @keyup.esc="cancelFolderCreate()"
              @keyup.enter="confirmFolderCreate()"
            />
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <v-btn
              block
              color="error"
              @click="cancelFolderCreate()"
              v-text="$t('common.Cancel')"
            />
          </v-col>
          <v-col>
            <v-btn
              block
              color="success"
              :disabled="newFolderName === ''"
              @click="confirmFolderCreate()"
              v-text="$t('common.Create')"
            />
          </v-col>
        </v-row>
      </ui-overlay>
      <v-btn
        icon
        @click="fetchItems(currentPath)"
      >
        <v-icon>mdi-refresh</v-icon>
      </v-btn>
    </v-card-title>
    <v-card-subtitle v-text="currentPath" />
    <v-card-text>
      <v-dialog
        v-model="confirmDeleteOpen"
        max-width="600"
      >
        <v-card>
          <v-card-title v-text="$t('files.ConfirmDelete')" />
          <v-card-actions>
            <v-spacer />
            <v-btn
              color="error"
              @click="deleteCancelled()"
              v-text="$t('common.Cancel')"
            />
            <v-btn
              color="success"
              @click="deleteConfirmed()"
              v-text="$t('common.Confirm')"
            />
          </v-card-actions>
        </v-card>
      </v-dialog>

      <v-list>
        <v-list-item
          v-for="file in files"
          :key="file.name"
          @click="itemClicked(file)"
        >
          <v-list-item-content>
            <v-list-item-title v-text="file.name" />
            <v-list-item-subtitle
              v-if="file.isFile"
              v-text="toFileSize(file.size || 0)"
            />
            <v-list-item-subtitle
              v-if="file.modifyTime"
              v-text="$t('files.LastModified') + ': ' + toDate(file.modifyTime)"
            />
          </v-list-item-content>
          <v-list-item-action class="flex-row">
            <v-tooltip
              v-if="file.isFile && !(file.size > maxEditSize)"
              bottom
            >
              <template v-slot:activator="{ on }">
                <v-btn
                  icon
                  v-on="on"
                  @click.stop="itemClicked(file)"
                >
                  <v-icon>mdi-pencil</v-icon>
                </v-btn>
              </template>
              <span v-text="$t('common.Edit')" />
            </v-tooltip>
            <v-tooltip
              v-if="file.isFile"
              bottom
            >
              <template v-slot:activator="{ on }">
                <v-btn
                  icon
                  :href="createDownloadLink(file)"
                  target="_blank"
                  v-on="on"
                >
                  <v-icon>mdi-download</v-icon>
                </v-btn>
              </template>
              <span v-text="$t('files.Download')" />
            </v-tooltip>
            <v-tooltip
              v-if="file.name !== '..'"
              bottom
            >
              <template v-slot:activator="{ on }">
                <v-btn
                  icon
                  v-on="on"
                  @click.stop="deleteRequest(file)"
                >
                  <v-icon>mdi-trash-can</v-icon>
                </v-btn>
              </template>
              <span v-text="$t('common.Delete')" />
            </v-tooltip>
          </v-list-item-action>
        </v-list-item>
      </v-list>

      <ui-overlay
        v-model="editOpen"
        :title="currentFile"
        card
        closable
        :on-close="cancelEdit"
      >
        <ace
          v-model="fileContents"
          editor-id="fileEditor"
          height="75vh"
          :theme="isDark() ? 'monokai' : 'github'"
          :file="currentFile"
        />
        <template v-slot:actions>
          <div class="flex-grow-1" />
          <v-btn
            color="error"
            @click="cancelEdit()"
            v-text="$t('common.Cancel')"
          />
          <v-btn
            color="success"
            @click="saveEdit()"
            v-text="$t('common.Save')"
          />
        </template>
      </ui-overlay>

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
    server: { type: Object, default: () => {} }
  },
  data () {
    return {
      files: [],
      currentPath: '/',
      loading: true,
      currentFile: '',
      fileContents: '',
      editOpen: false,
      maxEditSize: 1024 * 1024 * 20,
      createFolder: false,
      newFolderName: '',
      createFile: false,
      newFileName: '',
      uploadFiles: [],
      uploading: false,
      uploadCurrent: 0,
      uploadSize: 0,
      toDelete: null,
      confirmDeleteOpen: false
    }
  },
  mounted () {
    this.fetchItems(this.currentPath)

    this.$api.addServerListener(this.server.id, 'file', event => {
      if (event.error) {
        this.isLoading = false
        return
      }

      this.files = (event.files || []).sort((a, b) => {
        if (a.isFile && !b.isFile) return 1
        if (!a.isFile && b.isFile) return -1
        if (a.name.toLowerCase() > b.name.toLowerCase()) return 1
        if (a.name.toLowerCase() < b.name.toLowerCase()) return -1
        return 0
      })

      if (event.path !== '') {
        this.currentPath = event.path
      }

      this.loading = false
    })
  },
  methods: {
    fetchItems (path) {
      this.loading = true
      this.$api.requestServerFile(this.server.id, path)
    },
    async itemClicked (item) {
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

        this.fetchItems(this.currentPath)
      } else {
        if (item.size > this.maxEditSize) return
        let path = this.currentPath
        if (path === '/') {
          path += item.name
        } else {
          path += '/' + item.name
        }

        this.fileContents = await this.$api.downloadServerFile(this.server.id, path, true)
        this.currentFile = item.name
        this.editOpen = true
      }
    },
    cancelEdit () {
      this.editOpen = false
      this.currentFile = ''
      this.fileContents = ''
    },
    async saveEdit () {
      let path = this.currentPath
      if (path === '/') {
        path += this.currentFile
      } else {
        path += '/' + this.currentFile
      }

      await this.$api.uploadServerFile(this.server.id, path, this.fileContents)
      this.editOpen = false
      this.currentFile = ''
      this.fileContents = ''
      this.$toast.success(this.$t('common.Saved'))
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
      this.$api.requestServerFileDeletion(this.server.id, path)
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

      return this.$api.getServerFileUrl(this.server.id, path)
    },
    cancelFileCreate () {
      this.createFile = false
      this.newFileName = ''
    },
    confirmFileCreate () {
      if (this.newFileName === '') return

      const ctx = this
      ctx.uploadSingleFile(new File([], this.newFileName)).then(() => {
        ctx.fetchItems(ctx.currentPath)
      })
      ctx.createFile = false
      ctx.newFileName = ''
    },
    cancelFolderCreate () {
      this.createFolder = false
      this.newFolderName = ''
    },
    confirmFolderCreate () {
      if (this.newFolderName === '') return

      let path = this.currentPath
      if (path === '/') {
        path = path + this.newFolderName
      } else {
        path = path + '/' + this.newFolderName
      }

      this.$api.requestServerFolderCreation(this.server.id, path)
      this.createFolder = false
      this.newFolderName = ''
    },
    transmitFiles () {
      this.uploading = true
      this.uploadNextItem(this)
    },
    uploadNextItem (ctx) {
      this.uploadSingleFile(ctx.uploadFiles[0]).then(() => {
        ctx.uploadFiles.shift()
        if (ctx.uploadFiles.length === 0) {
          ctx.uploading = false
          ctx.fetchItems(ctx.currentPath)
          return
        }
        ctx.uploadNextItem(ctx)
      })
    },
    async uploadSingleFile (item) {
      let path = this.currentPath
      if (path === '/') {
        path += item.name
      } else {
        path += '/' + item.name
      }
      this.uploadCurrent = 0
      this.uploadSize = item.size

      return this.$api.uploadServerFile(this.server.id, path, item, event => {
        this.uploadCurrent = event.loaded
        this.uploadSize = event.total
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
