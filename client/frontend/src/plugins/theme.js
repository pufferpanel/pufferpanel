import { ref, inject } from 'vue'
import { contrast, deriveColor } from '@/utils/theme'
import { extract } from '@/utils/tar'

import defaultStyles from '@/themes/default/theme.scss?inline'
import defaultManifest from '@/themes/default/manifest.json'

let api = null
let config = null
let userSettings = null

const generatedStyles = document.createElement('style')
generatedStyles.id = 'generated-styles'
document.head.append(generatedStyles)

const styles = document.createElement('style')
styles.id = 'styles'
document.head.append(styles)

const activeTheme = ref('')
const themeSettings = ref({})
const rootClasses = ref([])
const sidebarClosedBelow = ref(1200)
const availableThemes = ['PufferPanel']

function appendStyle(style) {
  generatedStyles.textContent = generatedStyles.textContent + '\n' + style
}

function getDefaultValue(name, definition) {
  switch (definition.type) {
    case 'class':
      const defaults =  definition.options.filter(e => e.default)
      if (defaults.length < 1) {
        if (!definition.options[0].value) console.error(`no default found for setting '${name}'`)
        return definition.options[0].value
      }
      return defaults[0].value
    case 'color':
      return definition.default
    default:
      return null
  }
}

function handleSetting(name, definition, value, extra = {}) {
  if (value == undefined) value = getDefaultValue(name, definition)
  switch (definition.type) {
    case 'class':
      rootClasses.value.push(value)
      break
    case 'color':
      try {
        const derive = Array.isArray(definition.derive) ? definition.derive : []
        const derived = derive.map(d => deriveColor(value, d)).join('\n')
        appendStyle(`#app {\n${definition.var}: ${value};\n${derived}}`)
      } catch (e) {
        if (!extra.isRetry) {
          console.warn('Applying requested color setting failed', e)
          handleSetting(name, definition, undefined, { isRetry: true })
        } else {
          console.error('Applying color setting failed', e)
        }
      }
      break
    default:
      throw new Error(`type '${definition.type}' for setting ${name} is not a valid setting type`)
  }
  themeSettings.value[name] = value
}

let blobs = []

const themeApi = {
  getThemes: () => availableThemes,
  getActiveTheme: () => activeTheme.value,
  setTheme: async (newTheme, settings = {}, save = true) => {
    blobs.map(b => URL.revokeObjectURL(b))
    blobs = []
    if (availableThemes.indexOf(newTheme) === -1) {
      console.error('invalid theme selection, falling back to default')
      newTheme = 'PufferPanel'
      settings = {}
    }
    activeTheme.value = newTheme
    themeSettings.value = {}
    rootClasses.value = []
    generatedStyles.textContent = ''
    const themeData = { files: [] }
    if (newTheme !== 'PufferPanel') {
      const themeFiles = extract(await api.getTheme(newTheme))
      let manifestSeen = false
      themeFiles.map(file => {
        if (file.name === 'manifest.json') {
          themeData.manifest = JSON.parse(file.content)
          manifestSeen = true
        } else if (file.name === 'theme.css') {
          themeData.css = file.content
        } else {
          themeData.files.push(file)
        }
      })
      if (!manifestSeen) throw 'no theme manifest found'
      if (themeData.css) {
        themeData.files.map(file => {
          const url = URL.createObjectURL(file.blob)
          console.log('replace', file.name, url)
          themeData.css = themeData.css
            .split(`url('${file.name}')`).join(`url('${url}')`)
            .split(`url("${file.name}")`).join(`url("${url}")`)
          blobs.push(url)
        })
      }
    }
    const manifest = newTheme === 'PufferPanel' ? defaultManifest : themeData.manifest
    styles.textContent = newTheme === 'PufferPanel' ? defaultStyles : (manifest.keepDefaultCss ? defaultStyles + '\n' + themeData.css : themeData.css)
    sidebarClosedBelow.value = manifest.sidebarClosedBelow || 1200
    Object.keys(manifest.settings).map(key => {
      handleSetting(key, manifest.settings[key], settings[key])
    })
    if (save) {
      await userSettings.set('theme', newTheme)
      await userSettings.set('themeSettings', JSON.stringify(settings))
    }
  },
  getThemeSettings: async (theme = activeTheme.value) => {
    let definition
    if (theme === 'PufferPanel') {
      definition = {...defaultManifest.settings}
    } else {
      const themeData = extract(await api.getTheme(theme))
      const manifest = themeData.find(t => t.name === 'manifest.json')
      if (!manifest) throw 'no theme manifest found'
      definition = JSON.parse(manifest.content).settings
    }
    if (theme === activeTheme.value) {
      for (const setting in themeSettings.value) {
        definition[setting].current = themeSettings.value[setting]
      }
    } else {
      for (const setting in definition) {
        definition[setting].current = getDefaultValue(setting, definition[setting])
      }
    }
    return definition
  },
  setThemeSettings: (settings, save = true) => {
    themeApi.setTheme(activeTheme.value, settings, save)
  },
  serializeThemeSettings: (settings) => {
    const res = {}
    Object.keys(settings).map(key => {
      if (settings[key].current) res[key] = settings[key].current
    })
    return JSON.stringify(res)
  },
  deserializeThemeSettings: (settings, serialized) => {
    try {
      const deser = JSON.parse(serialized)
      const res = {...settings}
      if (deser && typeof deser === 'object' && !Array.isArray(deser)) {
        Object.keys(deser).map(key => {
          if (res[key]) res[key].current = deser[key]
        })
        return res
      } else {
        console.warn('Settings could not be deserialized, invalid value', serialized)
      }
    } catch (e) {
      console.error('Settings could not be deserialized', serialized, e)
    }

    return settings
  },
  // copy to prevent uncontrolled 3rd party mutation
  getThemeClasses: () => [...rootClasses.value],
  getThemeStyleAttributes: () => { return { ...rootAttributes.value } }
}

function initTheme() {
  try {
    themeApi.setTheme(
      userSettings.get('theme') || config.themes.active,
      JSON.parse(
        userSettings.get('themeSettings') || config.themes.settings
      ),
      false
    )
  } catch (e) {
    console.error('Default settings could not be applied', e)
    try {
      themeApi.setTheme(
        config.themes.active,
        JSON.parse(
          config.themes.settings
        ),
        false
      )
    } catch (e) {
      console.error('Falling back to default theme setting failed', e)
    }
  }
}

export default (conf) => {
  return {
    install: (app) => {
      config = conf
      app.config.globalProperties.$theme = themeApi
      app.provide('theme', themeApi)
      app.provide('themeClasses', rootClasses)
      app.provide('sidebarClosedBelow', sidebarClosedBelow)
      config.themes.available.filter(e => e !== 'PufferPanel').map(e => availableThemes.push(e))
      availableThemes.sort()
      api = app.config.globalProperties.$api
      userSettings = app.config.globalProperties.$userSettings
      initTheme()

      // unhide app after first theme load to prevent rendering unstyled elements
      document.getElementById('app').style.opacity = 1
      setTimeout(() => {
        document.getElementById('hideApp').remove()
        document.getElementById('app').removeAttribute('style')
      }, 200)

      app.config.globalProperties.$events.on('login', initTheme)
      app.config.globalProperties.$events.on('userSettingsReloaded', initTheme)
      app.config.globalProperties.$events.on('logout', initTheme)
    }
  }
}
