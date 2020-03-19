<template>
  <common-overlay v-model="show" card closable :title="$t('common.SelectLanguage')">
    <v-row>
      <v-col cols="12" sm="6" md="3" v-for="lang in Object.keys($i18n.messages)" :key="lang">
        <v-btn text v-text="getText(lang)" @click="setLocale(lang)" />
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
        flag = this.getFlag(parts[parts.length - 1])
      } else {
        flag = this.getFlag(locale)
      }
      return (!flag || flag === '') ? langName : `${flag} ${langName}`
    },
    getFlag (cc) {
      return cc.toUpperCase().replace(/./g, char => String.fromCodePoint(char.charCodeAt(0)+127397))
    },
    setLocale (locale) {
      this.$i18n.locale = locale
      localStorage.setItem('locale', locale)
      this.show = false
    }
  }
}
</script>
