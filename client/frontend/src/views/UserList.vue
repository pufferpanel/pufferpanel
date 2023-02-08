<script setup>
import { ref, inject, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'
import md5 from 'js-md5'
import Icon from '@/components/ui/Icon.vue'
import Btn from '@/components/ui/Btn.vue'
import Loader from '@/components/ui/Loader.vue'

const api = inject('api')
const { t } = useI18n()

const users = ref([])
let lastPage = 0
let loadingPage = false
const allUsersLoaded = ref(false)
const loaderRef = ref(null)
const firstEntry = ref(null)

function addUsers(newUsers) {
  newUsers.map(user => users.value.push(user))
}

function isLoaderVisible() {
  if (!loaderRef.value) return false
  const vw = window.innerWidth || document.documentElement.clientWidth
  const vh = window.innerHeight || document.documentElement.clientHeight
  const rect = loaderRef.value.$el.getBoundingClientRect()
  return rect.top >= 0 && rect.left >= 0 && rect.bottom <= vh && rect.right <= vw
}

async function loadPage(page = 1) {
  loadingPage = true
  const data = await api.user.list(page)
  addUsers(data.users)
  lastPage = data.paging.page
  allUsersLoaded.value = data.paging.page * data.paging.pageSize >= (data.paging.total || 0)
  nextTick(() => {
    loadingPage = false
    if (!allUsersLoaded.value && isLoaderVisible()) loadPage(lastPage + 1)
  })
}

function onScroll() {
  if (!loadingPage && isLoaderVisible()) loadPage(lastPage + 1)
}

onMounted(() => {
  nextTick(() => {
    loadPage()
    window.addEventListener('scroll', onScroll)
  })
})

onUnmounted(() => {
  window.removeEventListener('scroll', onScroll)
})

function setFirstEntry(ref) {
  if (!firstEntry.value) firstEntry.value = ref
}

function focusList() {
  firstEntry.value.$el.focus()
}
</script>

<template>
  <div class="userlist">
    <h1 v-text="t('users.Users')" />
    <div v-hotkey="'l'" class="list" @hotkey="focusList()">
      <div v-for="user in users" :key="user.id" class="list-item">
        <router-link :ref="setFirstEntry" :to="{ name: 'UserView', params: { id: user.id } }">
          <div class="user">
            <img :src="'https://www.gravatar.com/avatar/' + md5(user.email) + '?d=mp'" class="avatar" />
            <div>
              <span class="title">{{user.username}}</span>
              <span class="subline">{{user.email}}</span>
            </div>
          </div>
        </router-link>
      </div>
      <div v-if="!allUsersLoaded" ref="loaderRef" class="list-item">
        <loader small />
      </div>
      <div v-if="$api.auth.hasScope('users.create')" class="list-item">
        <router-link v-hotkey="'c'" :to="{ name: 'UserCreate' }">
          <div class="createLink"><icon name="plus" />{{ t('users.Add') }}</div>
        </router-link>
      </div>
    </div>
  </div>
</template>
