const path = require('path')
const { defineConfig } = require('vite')

module.exports = defineConfig({
  build: {
    lib: {
      entry: path.resolve(__dirname, 'src/index.js'),
      name: 'PufferPanel',
      fileName: (format) => {
        if (format === 'es') {
          return 'pufferpanel.mjs'
        } else {
          return 'pufferpanel.cjs'
        }
      }
    },
    minify: false,
    rollupOptions: {
      external: ['axios'],
      output: {
        globals: {
          axios: 'axios'
        }
      }
    }
  }
})
