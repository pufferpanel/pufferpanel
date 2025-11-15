<script setup>
import { ref, inject, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import Btn from '@/components/ui/Btn.vue'
import Collapse from '@/components/ui/Collapse.vue'
import Icon from '@/components/ui/Icon.vue'
import OtpInput from '@/components/ui/OtpInput.vue'
import Overlay from '@/components/ui/Overlay.vue'
import TextField from '@/components/ui/TextField.vue'
import Toggle from '@/components/ui/Toggle.vue'

const { t } = useI18n()
const api = inject('api')
const toast = inject('toast')

const passkeysSupported = isSecureContext && window.navigator && navigator.credentials

const loading = ref(false)
const otpEnabled = ref(false)
const otpEnrolling = ref(false)
const otpQrCode = ref(false)
const otpSecret = ref(false)
const otpDisabling = ref(false)
const otpRecovery = ref(false)
const recoveryCodes = ref([])
const regeneratingRecoveryCodes = ref(false)
const token = ref('')
const passkeys = ref([])
const passkeyName = ref('')
const enrollingPasskey = ref(false)
const allowPasswordless = ref(false)

onMounted(async () => {
  loading.value = true
  otpEnabled.value = await api.self.isOtpEnabled()
  await loadPasskeys()
  loading.value = false
})

async function loadPasskeys() {
  passkeys.value = await api.self.getPasskeys()
  const self = await api.self.get()
  allowPasswordless.value = self.allowPasswordlessLogin && passkeys.value.length
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
  const res = await api.self.validateOtpEnroll(token.value)
  resetOtpEnroll()
  recoveryCodes.value = res.recoveryCodes
  otpEnabled.value = await api.self.isOtpEnabled()
  toast.success(t('users.UpdateSuccess'))
}

function startRegenerateRecoveryCodes() {
  regeneratingRecoveryCodes.value = true
}

function resetRegenerateRecoveryCodes() {
  regeneratingRecoveryCodes.value = false
  token.value = ''
  otpRecovery.value = false
}

async function confirmRegenerateRecoveryCodes() {
  const res = await api.self.regenerateRecoveryCodes(token.value)
  resetRegenerateRecoveryCodes()
  recoveryCodes.value = res.recoveryCodes
  toast.success(t('users.UpdateSuccess'))
}

function saveRecoveryCodes() {
  const el = document.createElement('a')
  el.setAttribute('href', 'data:text/plain;charset=utf-8,' + encodeURIComponent(recoveryCodes.value.join('\n')))
  el.setAttribute('download', 'recovery-codes.txt')
  el.click()
}

function startOtpDeactivation() {
  otpDisabling.value = true
}

function resetOtpDeactivation() {
  otpDisabling.value = false
  token.value = ''
  otpRecovery.value = false
}

async function confirmOtpDeactivation() {
  await api.self.disableOtp(token.value)
  resetOtpDeactivation()
  otpEnabled.value = await api.self.isOtpEnabled()
  toast.success(t('users.UpdateSuccess'))
}

async function enrollPasskey() {
  if (passkeyName.value.trim().length === 0) return
  await api.self.enrollPasskey(passkeyName.value.trim())
  toast.success(t('users.UpdateSuccess'))
  passkeyName.value = ''
  enrollingPasskey.value = false
  loading.value = true
  await loadPasskeys()
  loading.value = false
}

async function removePasskey(id) {
  await api.self.deletePasskey(id)
  loading.value = true
  await loadPasskeys()
  loading.value = false
}

async function setPasswordless(e) {
  await api.self.setAllowPasswordlessLogin(e)
  toast.success(t('users.UpdateSuccess'))
  await loadPasskeys()
}
</script>

<template>
  <div v-if="loading" class="mfa loading"><loader /></div>
  <div v-else class="mfa">
    <h1 v-text="t('users.2fa')" />
    <div class="otp">
      <h2 v-text="t('users.Otp')" />
      <span class="description">{{ t('users.OtpHint') }}</span>
      <btn v-if="otpEnabled" class="otp-regenerate-recovery-codes" variant="text" @click="startRegenerateRecoveryCodes()"><icon name="refresh" />{{ t('users.RegenerateRecoveryCodes') }}</btn>
      <btn v-if="otpEnabled" class="otp-toggle" color="error" @click="startOtpDeactivation()"><icon name="lock-off" />{{ t('users.OtpDisable') }}</btn>
      <btn v-else class="otp-toggle" color="primary" @click="startOtpEnroll()"><icon name="lock" />{{ t('users.OtpEnable') }}</btn>

      <!-- Confirm OTP Enrollment -->
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
            <otp-input @update:modelValue="token = $event" @complete="token = $event; confirmOtpEnroll()" />
          </div>
        </div>
        <div class="otp-enroll-actions">
          <btn color="error" @click="resetOtpEnroll()" v-text="t('common.Cancel')" />
          <btn color="primary" @click="confirmOtpEnroll()" v-text="t('users.OtpEnable')" />
        </div>
      </overlay>

      <!-- Confirm Recovery Code Regeneration -->
      <overlay v-model="regeneratingRecoveryCodes" class="recovery-code-regeneration" :title="t('users.RegenerateRecoveryCodes')" closable @close="resetRegenerateRecoveryCodes()">
        <div class="recovery-code-regeneration-content">
          <div v-text="t('users.OtpConfirm')" />
          <div v-text="t('users.RegenerateRecoveryCodesHint')" />
          <otp-input v-if="!otpRecovery" @update:modelValue="token = $event" @complete="token = $event; confirmRegenerateRecoveryCodes()" />
          <text-field v-else v-model="token" autofocus />
          <btn variant="text" @click="otpRecovery = !otpRecovery; token = ''" v-text="otpRecovery ? t('users.OtpUseAuthenticator') : t('users.OtpUseRecovery')" />
        </div>
        <div class="recovery-code-regeneration-actions">
          <btn color="error" @click="resetRegenerateRecoveryCodes()" v-text="t('common.Cancel')" />
          <btn color="primary" @click="confirmRegenerateRecoveryCodes()" v-text="t('users.RegenerateRecoveryCodes')" />
        </div>
      </overlay>

      <!-- Confirm Disabling OTP -->
      <overlay v-model="otpDisabling" class="otp-deactivation" :title="t('users.OtpDisable')" closable @close="resetOtpDeactivation()">
        <div class="otp-deactivation-content">
          <span v-text="t('users.OtpConfirm')" />
          <otp-input v-if="!otpRecovery" @update:modelValue="token = $event" @complete="token = $event; confirmOtpDeactivation()" />
          <text-field v-else v-model="token" autofocus />
          <btn variant="text" @click="otpRecovery = !otpRecovery; token = ''" v-text="otpRecovery ? t('users.OtpUseAuthenticator') : t('users.OtpUseRecovery')" />
        </div>
        <div class="otp-deactivation-actions">
          <btn color="error" @click="resetOtpDeactivation()" v-text="t('common.Cancel')" />
          <btn color="primary" @click="confirmOtpDeactivation()" v-text="t('users.OtpDisable')" />
        </div>
      </overlay>

      <!-- Display Recovery Codes -->
      <overlay :model-value="recoveryCodes.length > 0" class="recovery-codes" :title="t('users.RecoveryCodes')" closable @close="recoveryCodes = []">
        <div v-text="t('users.RecoveryCodesHint')" />
        <div class="codes">
          <span v-for="code in recoveryCodes" :key="code" class="code"><code v-text="code" /></span>
        </div>
        <btn variant="text" @click="saveRecoveryCodes()"><icon name="download" />{{ t('users.SaveRecoveryCodes') }}</btn>
      </overlay>
    </div>
    <div v-if="passkeysSupported" class="passkeys">
      <h2 v-text="t('users.Passkeys')" />
      <div class="list">
        <div v-for="passkey in passkeys" :key="passkey.id" class="list-item">
          <div class="passkey">
            <span class="name">{{passkey.name}}</span>
            <btn variant="icon" class="remove" @click="removePasskey(passkey.id)"><icon name="remove" /></btn>
          </div>
        </div>
      </div>
      <btn color="primary" @click="enrollingPasskey = true"><icon name="lock" />{{ t('users.EnrollPasskey') }}</btn>
      <collapse :title="t('users.PasskeyAdvanced')">
        <toggle
          :model-value="allowPasswordless"
          :label="t('users.AllowPasswordlessLogin')"
          :hint="t('users.AllowPasswordlessLoginHint')"
          :disabled="passkeys.length === 0"
          @update:modelValue="setPasswordless($event)"
        />
      </collapse>

      <!-- Enroll Passkey -->
      <overlay v-model="enrollingPasskey" class="passkey-enroll" :title="t('users.EnrollPasskey')" closable @close="passkeyName = ''">
        <text-field v-model="passkeyName" autofocus :label="t('common.Name')" @keydown.enter="enrollPasskey()" />
        <btn variant="primary" :disabled="passkeyName.trim().length === 0" @click="enrollPasskey()"><icon name="lock" />{{ t('users.EnrollPasskey') }}</btn>
      </overlay>
    </div>
  </div>
</template>
