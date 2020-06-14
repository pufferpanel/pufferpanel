module.exports = {
  devServer: {
    disableHostCheck: true
  },
  transpileDependencies: ['vuetify'],
  configureWebpack: config => {
    if (process.env.NODE_ENV === 'production') {
      const SWPrecache = require('sw-precache-webpack-plugin')
      config.plugins.push(
        new SWPrecache({
          cacheId: 'pufferpanel',
          filepath: 'public/service-worker.js',
          staticFileGlobs: [
            'index.html',
            'auth/login',
            'server',
            'manifest.json',
            'js/*',
            'css/*'
          ],
          stripPrefix: '/'
        })
      )
    }
  }
}
