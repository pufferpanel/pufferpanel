<template>
  <common-overlay
    v-model="show"
    card
    closable
    :title="$t('common.SelectLanguage')"
  >
    <v-row>
      <v-col
        v-for="lang in Object.keys($i18n.messages)"
        :key="lang"
        cols="12"
        sm="6"
        md="3"
      >
        <v-btn
          text
          @click="setLocale(lang)"
          v-text="getText(lang)"
        />
      </v-col>
    </v-row>
    <v-row>
      <v-col class="d-flex justify-center">
        <a
          target="_blank"
          rel="noopener noreferrer"
          href="https://crowdin.com/project/pufferpanel"
        >
          {{ $t('common.HelpTranslate') }}
          <v-icon class="caption">mdi-launch</v-icon>
        </a>
      </v-col>
    </v-row>
  </common-overlay>
</template>

<script>
export default {
  props: {
    value: { type: Boolean, default: () => false }
  },
  computed: {
    show: {
      // double negation to effectively make a copy to prevent mutating props
      get: function () { return !!this.value },
      set: function (newValue) { this.$emit('input', newValue) }
    }
  },
  methods: {
    getText (locale) {
      const langName = this.$i18n.messages[locale].common.LanguageName
      let flag = ''
      if (locale.indexOf('_') >= 0) {
        const parts = locale.split('_')
        const last = parts[parts.length - 1]
        if (last.length === 2 && last.toUpperCase() !== 'SP') {
          flag = this.getFlag(parts[parts.length - 1])
        } else {
          flag = this.getFlag(parts[0])
        }
      } else {
        flag = this.getFlag(locale)
      }
      return (!flag || flag === '') ? langName : `${flag} ${langName}`
    },
    getFlag (cc) {
      return cc.toUpperCase().replace(/./g, char => String.fromCodePoint(char.charCodeAt(0) + 127397))
    },
    setLocale (locale) {
      this.$i18n.locale = locale
      localStorage.setItem('locale', locale)
      this.show = false
    }
  }
}
</script>
