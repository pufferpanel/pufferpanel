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

  if (componentName.startsWith('CoreServersType')) {
    const name = path.basename(fileName, '.vue')
    serverTypes.push(name)
  }

  Vue.component(componentName, componentConfig.default || componentConfig)
})

Vue.component('server-render', {
  render: function (createElement, context) {
    const server = this.$attrs.server
    if (server === null) {
      return
    }

    let element = 'core-servers-type-generic'
    for (const v in serverTypes) {
      if (serverTypes[v] === server.type) {
        element = 'core-servers-type-' + serverTypes[v]
        break
      }
    }

    return createElement(element, { props: { server: server } }, [])
  }
})
