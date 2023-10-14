module.exports = {
  globals: {
    defineEmits: 'readonly',
    defineExpose: 'readonly',
    defineProps: 'readonly',
    localeList: 'readonly',
    window: 'readonly'
  },
  extends: [
    "eslint:recommended",
    "plugin:vue/vue3-recommended",
    "prettier",
  ],
};
