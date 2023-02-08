export function generateOperatorLabel(t, o) {
  let c = 0
  let params = { ...o }
  switch(o.type) {
    case "download":
      c = Array.isArray(o.files) ? o.files.length : 1
      if (c === 1) params.file = Array.isArray(o.files) ? o.files[0] : o.files
      return t(`operators.${o.type}.formatted`, params, c)
    case "command":
      c = Array.isArray(o.commands) ? o.commands.length : 1
      if (c === 1) params.command = Array.isArray(o.commands) ? o.commands[0] : o.commands
      return t(`operators.${o.type}.formatted`, params, c)
    default:
      return t(`operators.${o.type}.formatted`, params, c)
  }
}

export const operators = {
  download: [
    {
      name: 'files',
      type: 'list',
      default: []
    }
  ],
  command: [
    {
      name: 'commands',
      type: 'list',
      default: []
    }
  ],
  alterfile: [
    {
      name: 'file',
      type: 'text',
      label: 'templates.Filename',
      default: ''
    },
    {
      name: 'regex',
      type: 'boolean',
      default: true
    },
    {
      name: 'search',
      type: 'text',
      default: ''
    },
    {
      name: 'replace',
      type: 'text',
      default: ''
    }
  ],
  writefile: [
    {
      name: 'target',
      type: 'text',
      label: 'templates.Filename',
      default: ''
    },
    {
      name: 'text',
      type: 'textarea',
      modeFile: 'target',
      default: ''
    }
  ],
  move: [
    {
      name: 'source',
      type: 'text',
      default: ''
    },
    {
      name: 'target',
      type: 'text',
      default: ''
    }
  ],
  mkdir: [
    {
      name: 'target',
      type: 'text',
      label: 'common.Name',
      default: ''
    }
  ],
  steamgamedl: [
    {
      name: 'appId',
      type: 'text',
      default: ''
    }
  ],
  javadl: [
    {
      name: 'version',
      type: 'text',
      label: 'templates.Version',
      default: ''
    }
  ],
  mojangdl: [
    {
      name: 'version',
      type: 'text',
      label: 'templates.Version',
      default: ''
    },
    {
      name: 'target',
      type: 'text',
      label: 'templates.Filename',
      default: ''
    }
  ],
  forgedl: [
    {
      name: 'version',
      type: 'text',
      label: 'templates.Version',
      default: ''
    },
    {
      name: 'filename',
      type: 'text',
      label: 'templates.Filename',
      default: ''
    }
  ],
  spongeforgedl: [
    {
      name: 'releaseType',
      type: 'text',
      default: ''
    }
  ],
  fabricdl: [
    {
      name: 'targetFile',
      type: 'text',
      label: 'templates.Filename',
      default: ''
    }
  ]
}
