import js from '@eslint/js'
import vue from 'eslint-plugin-vue'
import prettier from 'eslint-config-prettier/flat'
import globals from 'globals'

export default [
  js.configs.recommended,
  ...vue.configs['flat/recommended'],
  {
    rules: {
      'no-unused-vars': [
        'error',
        {
          caughtErrors: 'none'
        }
      ],
      'vue/multi-word-component-names': 'off',
      'vue/no-v-text-v-html-on-component': 'off',
      'vue/v-on-event-hyphenation': 'off'
    },
    languageOptions: {
      globals: {
        ...globals.browser,
        localeList: false
      },
    }
  },
  prettier
]
