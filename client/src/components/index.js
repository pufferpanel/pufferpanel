import Vue from 'vue'
import upperFirst from 'lodash/upperFirst'
import camelCase from 'lodash/camelCase'
import path from 'path'

const requireComponent = require.context('@/components', true, /\.vue$/)

const serverTypes = []

requireComponent.keys().forEach(fileName => {
  const componentConfig = requireComponent(fileName)

  const componentName = upperFirst(
    camelCase(fileName.replace(/^\.\//, '').replace(/\.\w+$/, ''))
  )

  if (componentName.startsWith('ServerType')) {
    const name = path.basename(fileName, '.vue')
    serverTypes.push(name)
  }

  Vue.component(componentName, componentConfig.default || componentConfig)
})

Vue.component('server-render', {
  render (createElement, context) {
    const server = this.$attrs.server
    if (server === null) {
      return
    }

    let element = 'server-type-generic'
    for (const v in serverTypes) {
      if (serverTypes[v] === server.type) {
        element = 'server-type-' + serverTypes[v]
        break
      }
    }

    return createElement(element, { props: { server: server } }, [])
  }
})
