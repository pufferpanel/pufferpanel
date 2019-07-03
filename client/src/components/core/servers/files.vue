<!--
  - Copyright 2019 Padduck, LLC
  -  Licensed under the Apache License, Version 2.0 (the "License");
  -  you may not use this file except in compliance with the License.
  -  You may obtain a copy of the License at
  -  	http://www.apache.org/licenses/LICENSE-2.0
  -  Unless required by applicable law or agreed to in writing, software
  -  distributed under the License is distributed on an "AS IS" BASIS,
  -  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  -  See the License for the specific language governing permissions and
  -  limitations under the License.
  -->

<template>
  <b-card
    header-tag="header">
    <h6 slot="header" class="mb-0">
      <span v-text="$t('common.FileManager') + ' - ' + currentPath + '          '"></span>
      <a v-if="!createFolder" @click="createFolder = true"><font-awesome-icon icon="plus"></font-awesome-icon></a>
      <div v-if="createFolder">
        <b-form-input size="sm" class="input-small" v-model="newFolderName" v-bind:placeholder="$t('common.NewFolder')"></b-form-input>
        <b-btn varient="primary" @click="submitNewFolder" v-text="$t('common.Create')"></b-btn>
        <b-btn variant="warning" @click="cancelFolderCreate" v-text="$t('common.Cancel')"></b-btn>
      </div>
    </h6>
    <b-table hover :items="files" :fields="fields" :busy="loading">
      <div slot="table-busy" class="text-center text-danger my-2">
        <b-spinner class="align-middle"/>
        <strong :text="$t('common.Loading')">Loading...</strong>
      </div>
      <template slot="name" slot-scope="data">
        <a v-if="data.item.isFile" v-bind:href="createDownloadLink(data.item)" v-bind:download="data.item.name"><strong v-text="data.value"></strong></a>
        <a v-else><strong v-on:dblclick="itemClicked(data.item)" v-text="data.value"></strong></a>
      </template>
      <template slot="size" slot-scope="data">
        <span v-if="data.value" v-text="toFileSize(data.value)"></span>
      </template>
      <template slot="modifyTime" slot-scope="data">
        <span v-if="data.value" v-text="toDate(data.value)"></span>
      </template>
      <template slot="isFile" slot-scope="data">
        <a v-bind:href="createDownloadLink(data.item)" v-bind:download="data.item.name" v-if="data.value"><font-awesome-icon v-b-tooltip.hover v-bind:title="$t('common.Download')" icon="download"></font-awesome-icon></a>
        <span class="p-1"></span>
        <a v-on:click="editButton(data.item)" v-if="data.value && data.item.size < maxEditSize"><font-awesome-icon v-b-tooltip.hover v-bind:title="$t('common.Edit')" icon="edit"></font-awesome-icon></a>
        <span class="p-1"></span>
        <a v-on:click="deleteButton(data.item)"><font-awesome-icon v-b-tooltip.hover v-bind:title="$t('common.Delete')" icon="trash"></font-awesome-icon></a>
      </template>
    </b-table>
  </b-card>
</template>

<script>
import filesize from 'filesize'

export default {
  prop: {
    server: Object
  },
  data() {
    return {
      files: [],
      currentPath: '/',
      loading: true,
      fields: {
        'name': {
          sortable: true,
          label: this.$t('common.Name')
        },
        'size': {
          sortable: true,
          label: this.$t('common.Size')
        },
        'modifyTime': {
          sortable: true,
          label: this.$t('common.LastModified')
        },
        'isFile': {
          sortable: false,
          label: this.$t('common.Actions')
        }
      },
      currentFile: '',
      fileContents: '',
      toEdit: false,
      maxEditSize: 1024 * 1024 * 20,
      createFolder: false,
      newFolderName: ''
    }
  },
  methods: {
    fetchItems(path, edit=false) {
      this.loading = true
      this.$socket.sendObj({type: 'file', action: 'get', path: path, edit: edit})
    },
    itemClicked(item) {
      if (!item.isFile) {
        this.loading = true

        if (item.name === '..') {
          let parts = this.currentPath.split('/')
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

        this.$socket.sendObj({type: 'file', action: 'get', path: this.currentPath})
      }
    },
    editButton(item) {
      this.toEdit = true
      if (this.currentPath === '/') {
        this.currentPath = '/' + item.name
      } else {
        this.currentPath = this.currentPath + '/' + item.name
      }
      this.fetchItems(this.currentPath, true)
    },
    deleteButton(item) {
      this.toEdit = false
      let path = ''
      if (this.currentPath === '/') {
        path = '/' + item.name
      } else {
        path = this.currentPath + '/' + item.name
      }
      this.loading = true
      this.$socket.sendObj({type: 'file', action: 'delete', path: path})
    },

    //utility
    toFileSize(size) {
      return filesize(size)
    },
    toDate(epoch) {
      return new Date(epoch * 1000).toLocaleString()
    },
    createDownloadLink(item) {
      return this.$router.resolve('/daemon/server/' + this.$attrs.server.id + '/file' + this.currentPath + item.name).href
    },
    cancelFolderCreate() {
      this.createFolder = false
      this.newFolderName = ''
    },
    submitNewFolder() {
      let path = this.currentPath
      if (path === '/') {
        path = path + this.newFolderName
      } else {
        path = path + '/' + this.newFolderName
      }

      this.$socket.sendObj({type: 'file', action: 'create', path: path})
      this.createFolder = false
      this.newFolderName = ''
    }
  },
  mounted() {
    let vue = this
    this.$socket.addEventListener('open', function (event) {
      vue.fetchItems(vue.currentPath)
    })

    this.$socket.addEventListener('message', function (event) {
      let data = JSON.parse(event.data)
      if (data === 'undefined') {
        return
      }
      if (data.type === 'file') {
        if (data.data) {
          if (data.data.error) {
            vue.isLoading = false
            return
          }

          let files = data.data.files

          //if we have a list of files, show them
          if (files) {
            vue.files = []
            for (let i in files) {
              let file = files[i]
              vue.files.push(file)
            }
          }
          //otherwise, it's an actual file, so we need to show it
          else {
            if (vue.toEdit) {
            } else {
              let url = vue.$router.resolve('/daemon/' + vue.$attrs.server.id + '/file' + data.data.url);
              console.log(url)
              vue.download(data.data.name, url.href)
            }
          }
          vue.currentPath = data.data.path
          vue.loading = false
        }
      }
    })
  }
}
</script>

<style scoped>
  a {
    color: inherit;
  }

  .input-small {
    width: 200px
  }
</style>