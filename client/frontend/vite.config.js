import path from 'path'
import { defineConfig } from "vite"
import vue from "@vitejs/plugin-vue"
import eslint from "vite-plugin-eslint"
import vueI18n from '@intlify/vite-plugin-vue-i18n'
import fs from 'fs'

export default defineConfig({
  build: {
    chunkSizeWarningLimit: 550,
    rollupOptions: {
      output: {
        entryFileNames: 'js/[name]-[hash].js',
        chunkFileNames: 'js/[name]-[hash].js',
        assetFileNames: '[ext]/[name]-[hash][extname]'
      }
    }
  },
  resolve: {
    alias: [
      {find: "@", replacement: path.resolve(__dirname, 'src')}
    ]
  },
  define: {
    localeList: fs.readdirSync('src/lang')
  },
  plugins: [
    vue(),
    vueI18n({
      runtimeOnly: false,
      include: path.resolve(__dirname, '@/lang/**')
    }),
    eslint()
  ]
})
