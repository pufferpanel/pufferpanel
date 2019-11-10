import Vue from 'vue'
import upperFirst from 'lodash/upperFirst'
import camelCase from 'lodash/camelCase'
import path from 'path'

const requireComponent = require.context(
  '@/components', true, /\.vue$/
)

let serverTypes = []

requireComponent.keys().forEach(fileName => {
  const componentConfig = requireComponent(fileName)

  const componentName = upperFirst(
    camelCase(fileName.replace(/^\.\//, '').replace(/\.\w+$/, ''))
  )

  if (componentName.startsWith("CoreServersType")) {
    let name = path.basename(fileName, '.vue')
    serverTypes.push(name)
  }

  Vue.component(componentName, componentConfig.default || componentConfig)
})

Vue.component('server-render', {
  render: function (createElement, context) {
    let server = this.$attrs.server
    if (server === null) {
      return
    }

    let element = 'core-servers-type-generic'
    for (let v in serverTypes) {
      if (serverTypes[v] === server.type) {
        element = 'core-servers-type-' + serverTypes[v]
        break
      }
    }

    return createElement(element, {props: {"server": server}}, [])
  }
})