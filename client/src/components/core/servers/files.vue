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
    <h6 slot="header" class="mb-0" v-text="$t('common.FileManager') + ' - ' + currentPath"></h6>
    <b-table hover selectable select-mode="single" @row-selected="itemClicked" :items="files" :fields="fields"
             :busy="loading">
      <div slot="table-busy" class="text-center text-danger my-2">
        <b-spinner class="align-middle"/>
        <strong :text="$t('common.Loading')">Loading...</strong>
      </div>
      <template slot="name" slot-scope="data">
        <strong v-text="data.value"></strong>
      </template>
      <template slot="size" slot-scope="data">
        <span v-if="data.value" v-text="toFileSize(data.value)"></span>
      </template>
      <template slot="modifyTime" slot-scope="data">
        <span v-if="data.value" v-text="toDate(data.value)"></span>
      </template>
    </b-table>
  </b-card>
</template>

<script>

import filesize from 'filesize'

export default {
  data () {
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
      toEdit: false
    }
  },
  methods: {
    fetchItems (path) {
      this.loading = true
      this.$socket.sendObj({ type: 'file', action: 'get', path: path })
    },
    toFileSize (size) {
      return filesize(size)
    },
    toDate (epoch) {
      return new Date(epoch * 1000).toLocaleString()
    },
    itemClicked (rows) {
      if (rows === 'undefined' || rows.length === 0) {
        return
      }

      if (this.loading) {
        return
      }

      let item = rows[0]
      if (item === undefined) {
        return
      }

      if (item.isFile) {
        this.loading = true
        let path = this.currentPath
        if (path === '/') {
          path += item.name
        } else {
          path += '/' + item.name
        }
        this.currentFile = item.name
        this.$socket.sendObj({ type: 'file', action: 'get', path: path })
      } else {
        if (item.name === '..') {
          let parts = this.currentPath.split('/')
          parts.pop()
          this.currentPath = parts.join('/')
          if (this.currentPath === '') {
            this.currentPath = '/'
          }
        } else {
          if (this.currentPath === '/') {
            this.currentPath = '/' + item.name
          } else {
            this.currentPath = this.currentPath + '/' + item.name
          }
        }
        this.fetchItems(this.currentPath)
      }
    },
    download(filename, text) {
      let element = document.createElement('a');
      element.setAttribute('href', 'data:application/octet-stream;charset=utf-8;base64,' + text);
      element.setAttribute('download', filename);

      element.style.display = 'none';
      document.body.appendChild(element);
      element.click();
      document.body.removeChild(element);
    }
  },
  mounted () {
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
            //if we got the data in the buffer, we have it already
            if (data.data.contents) {
              if (this.toEdit) {

              } else {
                vue.download(data.data.name, data.data.contents)
              }
            }
          }

          vue.loading = false
        }
      }
    })
  }
}
</script>

<style scoped>
</style>