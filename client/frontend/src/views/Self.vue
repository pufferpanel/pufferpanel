<script setup>
import { ref, inject, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { updateLocale, locales } from '@/plugins/i18n'
import Overlay from '@/components/ui/Overlay.vue'
import TextField from '@/components/ui/TextField.vue'
import Dropdown from '@/components/ui/Dropdown.vue'
import Btn from '@/components/ui/Btn.vue'
import Icon from '@/components/ui/Icon.vue'
import Loader from '@/components/ui/Loader.vue'
import OAuth from '@/components/ui/OAuth.vue'
import Tab from '@/components/ui/Tab.vue'
import Tabs from '@/components/ui/Tabs.vue'
import ThemeSetting from '@/components/ui/ThemeSetting.vue'

const { t, locale, fallbackLocale } = useI18n()
const api = inject('api')
const toast = inject('toast')
const themeApi = inject('theme')
const theme = ref(themeApi.getActiveTheme())
const themeSettings = ref({})
const user = ref(undefined)
const acc = ref(undefined)
const newPass = ref({ old: '', new: '', confirm: ''})
const otpEnabled = ref(false)
const otpEnrolling = ref(false)
const otpQrCode = ref(false)
const otpSecret = ref(false)
const otpDisabling = ref(false)
const token = ref('')
const selectedLocale = ref(locale.value)

onMounted(async () => {
  themeSettings.value = await themeApi.getThemeSettings()
  const data = await api.self.get()
  acc.value = { username: data.username, email: data.email, password: '' }
  user.value = data
  otpEnabled.value = await api.self.isOtpEnabled()
})

async function themeChanged() {
  themeSettings.value = await themeApi.getThemeSettings(theme.value)
}

function savePreferences() {
  if (locale.value !== selectedLocale.value) {
    updateLocale(selectedLocale.value)
  }

  const settings = {}
  Object.keys(themeSettings.value).map(key => {
    settings[key] = themeSettings.value[key].current
  })

  themeApi.setTheme(theme.value, settings)
  toast.success(t('self.PreferencesUpdated'))
}

async function startOtpEnroll() {
  const data = await api.self.startOtpEnroll()
  otpEnrolling.value = true
  otpQrCode.value = data.img
  otpSecret.value = data.secret
}

function resetOtpEnroll() {
  otpEnrolling.value = false
  otpQrCode.value = false
  otpSecret.value = false
  token.value = ''
}

async function confirmOtpEnroll() {
  await api.self.validateOtpEnroll(token.value)
  resetOtpEnroll()
  otpEnabled.value = await api.self.isOtpEnabled()
  toast.success(t('users.UpdateSuccess'))
}

function startOtpDeactivation() {
  otpDisabling.value = true
}

function resetOtpDeactivation() {
  otpDisabling.value = false
  token.value = ''
}

async function confirmOtpDeactivation() {
  await api.self.disableOtp(token.value)
  resetOtpDeactivation()
  otpEnabled.value = await api.self.isOtpEnabled()
  toast.success(t('users.UpdateSuccess'))
}

function isValidUsername(u) {
  return u.length >= 5
}

function isValidEmail(e) {
  return e.match(/.+@.+\..{2,}/)
}

function isValidPassword(p) {
  return p.length >= 8
}

function canSubmitDetailsChange() {
  return isValidUsername(acc.value.username) && isValidEmail(acc.value.email) && isValidPassword(acc.value.password)
}

function canSubmitPasswordChange() {
  return isValidPassword(newPass.value.old) && isValidPassword(newPass.value.new) && newPass.value.new === newPass.value.confirm
}

async function submitDetailsChange() {
  if (!canSubmitDetailsChange()) return
  await api.self.updateDetails(acc.value.username, acc.value.email, acc.value.password)
  toast.success(t('users.InfoChanged'))
}

async function submitPasswordChange() {
  if (!canSubmitPasswordChange()) return
  await api.self.changePassword(newPass.value.old, newPass.value.new)
  toast.success(t('users.PasswordChanged'))
}

function updateThemeSetting(name, newSetting) {
  themeSettings.value[name] = newSetting
}
</script>

<template>
  <div v-if="!user" class="self loading">
    <div class="loader"><loader /></div>
  </div>
  <div v-else class="self">
    <tabs anchors>
      <tab id="preferences" :title="t('users.Preferences')" icon="settings" hotkey="t s">
        <div class="preferences">
          <h1 v-text="t('users.Preferences')" />
          <dropdown v-model="selectedLocale" class="locale-select" :options="locales" :label="t('common.Language')" :hint="`[${t('common.HelpTranslate')}](https://translate.pufferpanel.com)`">
            <template #singlelabel="{ value }">
              <div class="multiselect-single-label">
                <span :data-locale="value.value" /> {{ value.label }}
              </div>
            </template>

            <template #option="{ option }">
              <span :data-locale="option.value" /> {{ option.label }}
            </template>
          </dropdown>
          <dropdown v-model="theme" :options="$theme.getThemes()" :label="t('common.theme.Theme')" @change="themeChanged()" />
          <theme-setting v-for="(setting, name) in themeSettings" :key="name" :model-value="setting" @update:modelValue="updateThemeSetting(name, $event)" />
          <btn color="primary" @click="savePreferences()"><icon name="save" />{{ t('users.SavePreferences') }}</btn>
        </div>
      </tab>
      <tab id="account" :title="t('users.ChangeInfo')" icon="account" hotkey="t a">
        <div class="accountdetails">
          <h1 v-text="t('users.ChangeInfo')" />
          <form>
            <text-field v-model="acc.username" icon="account" :label="t('users.Username')" />
            <text-field v-model="acc.email" icon="email" type="email" :label="t('users.Email')" />
            <text-field v-model="acc.password" icon="lock" type="password" :label="t('users.ConfirmPassword')" />
            <btn :disabled="!canSubmitDetailsChange()" color="primary" @click="submitDetailsChange()"><icon name="save" />{{ t('users.ChangeInfo') }} </btn>
          </form>
        </div>
      </tab>
      <tab id="changepassword" :title="t('users.ChangePassword')" icon="lock" hotkey="t p">
        <div class="changepassword">
          <h1 v-text="t('users.ChangePassword')" />
          <form>
            <text-field v-model="newPass.old" icon="lock" type="password" :label="t('users.OldPassword')" />
            <text-field v-model="newPass.new" icon="lock" type="password" :label="t('users.NewPassword')" />
            <text-field v-model="newPass.confirm" icon="lock" type="password" :label="t('users.ConfirmPassword')" />
            <btn :disabled="!canSubmitPasswordChange()" color="primary" @click="submitPasswordChange()"><icon name="save" />{{ t('users.ChangePassword') }}</btn>
          </form>
        </div>
      </tab>
      <tab id="otp" :title="t('users.Otp')" icon="2fa" hotkey="t 2">
        <div class="mfa">
          <h1 v-text="t('users.Otp')" />
          <span class="description">{{ t('users.OtpHint') }}</span>
          <btn v-if="otpEnabled" class="otp-toggle" color="error" @click="startOtpDeactivation()"><icon name="lock-off" />{{ t('users.OtpDisable') }}</btn>
          <btn v-else class="otp-toggle" color="primary" @click="startOtpEnroll()"><icon name="lock" />{{ t('users.OtpEnable') }}</btn>
          <overlay v-model="otpEnrolling" class="otp-enroll" :title="t('users.OtpEnable')" closable @close="resetOtpEnroll()">
            <div class="otp-enroll-content">
              <div class="otp-enroll-qr">
                <img :src="otpQrCode" />
              </div>
              <div class="otp-enroll-info">
                <h2 v-text="t('users.OtpSetupHint')" />
                <div><b class="otp-enroll-secret" v-text="t('users.OtpSecret')" /></div>
                <span class="otp-enroll-secret" v-text="otpSecret" />
                <h3 class="otp-enroll-confirm" v-text="t('users.OtpConfirm')" />
                <text-field v-model="token"/>
              </div>
            </div>
            <div class="otp-enroll-actions">
              <btn color="error" @click="resetOtpEnroll()" v-text="t('common.Cancel')" />
              <btn color="primary" @click="confirmOtpEnroll()" v-text="t('users.OtpEnable')" />
            </div>
          </overlay>
          <overlay v-model="otpDisabling" class="otp-deactivation" :title="t('users.OtpDisable')" closable @close="resetOtpDeactivation()">
            <div class="otp-deactivation-content">
              <text-field v-model="token" :label="t('users.OtpConfirm')" />
            </div>
            <div class="otp-deactivation-actions">
              <btn color="error" @click="resetOtpDeactivation()" v-text="t('common.Cancel')" />
              <btn color="primary" @click="confirmOtpDeactivation()" v-text="t('users.OtpDisable')" />
            </div>
          </overlay>
        </div>
      </tab>
      <tab id="oauth" :title="t('oauth.Clients')" icon="api" hotkey="t o">
        <div class="oauth">
          <h1 v-text="t('oauth.Clients')" />
          <o-auth />
        </div>
      </tab>
    </tabs>
  </div>
</template>
