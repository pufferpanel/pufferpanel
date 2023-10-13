<script setup>
import { ref, inject, nextTick } from 'vue'
import { useRouter, RouterLink } from 'vue-router';
import { useI18n } from 'vue-i18n'
import Loader from './Loader.vue'
import TextField from './TextField.vue'

const api = inject('api')
const router = useRouter()
const { t } = useI18n()

let timeout = null

const active = ref(false)
const loading = ref(false)
const query = ref('')
const maxIndex = ref(0)
const currIndex = ref(0)
const currRef = ref(null)
const containerRef = ref(null)
const links = ref([])

const servers = ref([])
const templates = ref([])
const users = ref([])
const nodes = ref([])

function cancel() {
  if (timeout) {
    clearTimeout(timeout)
    timeout = null
  }
  setTimeout(() => {
    //active.value = false
  }, 500)
}

function close() {
  active.value = false
  query.value = ''
  document.activeElement.blur()
}

function input() {
  if (timeout) {
    clearTimeout(timeout)
    timeout = null
  }
  timeout = setTimeout(search, 500)
}

async function search() {
  loading.value = true
  active.value = true
  const q = query.value.toLowerCase()
  reset()
  await Promise.all([
    findServers(q),
    api.auth.hasScope('users.info.search') ? findUsers(q) : Promise.resolve(),
    api.auth.hasScope('nodes.view') ? findNodes(q) : Promise.resolve(),
    api.auth.hasScope('templates.view') ? findTemplates(q) : Promise.resolve()
  ])
  let mi = servers.value.length + users.value.length + nodes.value.length
  templates.value.map(repo => mi += repo.templates.length)
  maxIndex.value = mi - 1
  loading.value = false
}

async function findServers(query) {
  servers.value = (await api.server.list(1, 5, query)).servers
}

async function findTemplates(query) {
  templates.value = (await api.template.listAllTemplates()).map(repo => {
    repo.templates = repo.templates.filter(template => {
      if (template.name.toLowerCase().indexOf(query) > -1) return true
      return template.display.toLowerCase().indexOf(query) > -1
    }).slice(0, 5)
    return repo
  }).filter(repo => {
    return repo.templates.length > 0
  })
}

async function findUsers(query) {
  const byName = await api.user.search(query, 5)
  const byEmail = (await api.user.searchEmail(query, 5)).filter(u => {
    return byName.filter(n => {
      return n.id === u.id
    }).length === 0
  })
  users.value = byName.concat(byEmail).slice(0, 5)
}

async function findNodes(query) {
  nodes.value = (await api.node.list()).filter(node => node.name.toLowerCase().indexOf(query) > -1).slice(0, 5)
}

function reset() {
  maxIndex.value = 0
  currIndex.value = 0
  currRef.value = null
  containerRef.value = null
  links.value = []
  servers.value = []
  templates.value = []
  users.value = []
  nodes.value = []
}

function getServerAddress(server) {
  let ip = server.node.publicHost
  if (server.ip && server.ip !== '0.0.0.0') {
    ip = server.ip
  }
  return ip + (server.port ? ':' + server.port : '')
}

function serverIndex(i) {
  return i
}

function userIndex(i) {
  return servers.value.length + i
}

function nodeIndex(i) {
  return servers.value.length + users.value.length + i
}

function templateIndex(repo, i) {
  let offset = servers.value.length + users.value.length + nodes.value.length
  let foundRepo = false
  templates.value.map(r => {
    if (repo.id === r.id) foundRepo = true
    if (!foundRepo) offset += r.templates.length
  })
  return offset + i
}

function setRef(i) {
  return ref => {
    if (currIndex.value === i)
      currRef.value = ref
  }
}

function scrollIfNeeded(complete) {
  nextTick(() => {
    const container = containerRef.value.getBoundingClientRect()
    const curr = currRef.value.getBoundingClientRect()
    if (container.top > curr.top) {
      complete ? containerRef.value.scrollTo({ behavior: 'smooth', top: 0 }) : containerRef.value.scrollBy({ behavior: 'smooth', top: curr.top - container.top })
    }
    if (container.bottom < curr.bottom) {
      complete ? containerRef.value.scrollTo({ behavior: 'smooth', top: 9999 }) : containerRef.value.scrollBy({ behavior: 'smooth', top: curr.bottom - container.bottom })
    }
  })
}

function up() {
  if (currIndex.value === 0) {
    currIndex.value = maxIndex.value
    scrollIfNeeded(true)
  } else {
    currIndex.value = currIndex.value - 1
    scrollIfNeeded(false)
  }
}

function down() {
  if (currIndex.value === maxIndex.value) {
    currIndex.value = 0
    scrollIfNeeded(true)
  } else {
    currIndex.value = currIndex.value + 1
    scrollIfNeeded(false)
  }
}

function go() {
  router.push(links[currIndex.value])
  close()
}

function link(index, to) {
  links[index] = to
  return to
}
</script>

<template>
  <span class="global-search">
    <text-field v-model="query" v-hotkey="'/'" icon="search" @change="input()" @blur="cancel()" @keyup.up="up()" @keyup.down="down()" @keyup.enter="go()" @keyup.esc="close()" />
    <div v-if="active && loading" class="results">
      <loader />
    </div>
    <div v-if="active && !loading" ref="containerRef" class="results">
      <div v-if="servers.length > 0" class="server-results">
        <h3 v-text="t('servers.Servers')" />
        <div v-for="(server, i) in servers" :key="server.id" :ref="setRef(serverIndex(i))" :class="['result', currIndex === serverIndex(i) ? 'selected' : '' ]">
          <router-link :to="link(serverIndex(i), { name: 'ServerView', params: { id: server.id } })">
            <div class="title">{{ server.name }}</div>
            <div class="subline">{{ getServerAddress(server) }} @ {{ server.node.name }}</div>
          </router-link>
        </div>
      </div>
      <div v-if="users.length > 0" class="user-results">
        <h3 v-text="t('users.Users')" />
        <div v-for="(user, i) in users" :key="user.id" :ref="setRef(userIndex(i))" :class="['result', currIndex === userIndex(i) ? 'selected' : '' ]">
          <router-link :to="link(userIndex(i), { name: 'UserView', params: { id: user.id } })">
            <div class="title">{{ user.username }}</div>
            <div class="subline">{{ user.email }}</div>
          </router-link>
        </div>
      </div>
      <div v-if="nodes.length > 0" class="node-results">
        <h3 v-text="t('nodes.Nodes')" />
        <div v-for="(node, i) in nodes" :key="node.id" :ref="setRef(nodeIndex(i))" :class="['result', currIndex === nodeIndex(i) ? 'selected' : '' ]">
          <router-link :to="link(nodeIndex(i), { name: 'NodeView', params: { id: node.id } })">
            <div class="title">{{ node.name }}</div>
            <div class="subline">{{ node.publicHost + ':' + node.publicPort }}</div>
          </router-link>
        </div>
      </div>
      <div v-if="templates.length > 0" class="template-results">
        <h3 v-text="t('templates.Templates')" />
        <div v-for="repo in templates" :key="repo.id">
          <h4 v-text="repo.name" />
          <div v-for="(template, i) in repo.templates" :key="template.name" :ref="setRef(templateIndex(repo, i))" :class="['result', currIndex === templateIndex(repo, i) ? 'selected' : '' ]">
            <router-link :to="link(templateIndex(repo, i), { name: 'TemplateView', params: { repo: repo.id, id: template.name } })">
              <div class="title">{{ template.display }}</div>
              <div class="subline">{{ template.name }}</div>
            </router-link>
          </div>
        </div>
      </div>
      <div v-if="maxIndex === -1" class="no-results" v-text="t('common.NoResults')" />
    </div>
  </span>
</template>
