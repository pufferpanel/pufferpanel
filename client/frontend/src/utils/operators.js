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
    case "archive":
      c = Array.isArray(o.source) ? o.source.length : 1
      if (c === 1) params.file = Array.isArray(o.source) ? o.source[0] : o.source
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
  archive: [
    {
      name: 'source',
      type: 'list',
      default: []
    },
    {
      name: 'destination',
      type: 'text',
      label: 'templates.Filename',
      default: ''
    }
  ],
  extract: [
    {
      name: 'source',
      type: 'text',
      label: 'templates.Filename',
      default: ''
    },
    {
      name: 'destination',
      type: 'text',
      default: ''
    }
  ],
  console: [
    {
      name: 'message',
      type: 'text',
      default: ''
    }
  ],
  sleep: [
    {
      name: 'duration',
      type: 'text',
      default: '5s'
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
      name: 'minecraftVersion',
      type: 'text',
      label: 'templates.MinecraftVersion',
      default: 'latest'
    },
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
    },
    {
      name: 'outputVariable',
      type: 'text',
      default: ''
    }
  ],
  neoforgedl: [
    {
      name: 'minecraftVersion',
      type: 'text',
      label: 'templates.MinecraftVersion',
      default: 'latest'
    },
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
    },
    {
      name: 'outputVariable',
      type: 'text',
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
  ],
  paperdl: [
    {
      name: 'project',
      type: 'text',
      default: 'paper'
    },
    {
      name: 'minecraftVersion',
      type: 'text',
      label: 'templates.MinecraftVersion',
      default: 'latest'
    },
    {
      name: 'build',
      type: 'text',
      default: 'latest'
    },
    {
      name: 'target',
      type: 'text',
      label: 'templates.Filename',
      default: 'server.jar'
    }
  ],
  curseforge: [
    {
      name: 'projectId',
      type: 'text',
      default: ''
    },
    {
      name: 'fileId',
      type: 'text',
      default: ''
    },
    {
      name: 'java',
      type: 'text',
      default: 'java'
    }
  ],
  nodejsdl: [
    {
      name: 'version',
      type: 'text',
      label: 'templates.Version',
      default: ''
    }
  ]
}
