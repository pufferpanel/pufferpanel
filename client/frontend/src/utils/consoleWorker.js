import AnsiUp from 'ansi_up'

const ansiup = new AnsiUp()
ansiup.ansi_to_html('\u001b[m')

function handleCarriageReturn(line) {
  if (line.indexOf('\r') !== -1) {
    const parts = line.split('\r')
    let result = parts.shift()
    parts.map(part => {
      result = part + result.substring(part.length)
    })
    return result
  }

  return line
}

function markDaemon(line, name) {
  if (line.trim().indexOf('[DAEMON]') === 0) {
    return `<span class="daemon-marker" data-name="${name}"></span>` + line.substring(8)
  }

  return line
}

function handleLine(line, name) {
  return markDaemon(handleCarriageReturn(line), name)
}

let lastIncompleteLine = null

onmessage = function (e) {
  let newLines = (Array.isArray(e.data.logs) ? e.data.logs.join('') : e.data.logs).replaceAll('\r\n', '\n')
  const endOnNewline = newLines.endsWith('\n')
  newLines = newLines.split('\n')
  if (endOnNewline) newLines.pop()

  const updates = []

  let last = null
  newLines.map(line => {
    line = ansiup.ansi_to_html(line)
    if (lastIncompleteLine) line = lastIncompleteLine.line + line

    if (lastIncompleteLine) {
      updates.push({ op: 'update', content: handleLine(line, e.data.name) })
      last = { line }
    } else {
      updates.push({ op: 'append', content: handleLine(line, e.data.name) })
      last = { line }
    }
  })

  postMessage(updates)

  if (!endOnNewline) {
    lastIncompleteLine = last
  } else {
    lastIncompleteLine = null
  }
}
