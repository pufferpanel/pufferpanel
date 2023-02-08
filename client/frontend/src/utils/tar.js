export function extract(bytes) {
  const bytesToString = bytes => {
    return new TextDecoder('utf-8').decode(new Uint8Array(bytes))
  }

  const CHUNK_SIZE = 512
  const FILE_NAME_SIZE = 100
  const FILE_LENGTH_START = 124
  const FILE_LENGTH_SIZE = 11
  const BASE_OCTAL = 8
  const files = []
  const data = new Uint8Array(bytes)

  let currentFileMeta = 0
  let hasNext = true

  while (hasNext) {
    const meta = data.slice(currentFileMeta, currentFileMeta + CHUNK_SIZE)
    const nameBytes = []
    for (let i = 0; i < FILE_NAME_SIZE; i++) {
      if (meta[i] === 0) break
      nameBytes.push(meta[i])
    }
    const fileName = bytesToString(nameBytes)
    const fileLength = parseInt(bytesToString(meta.slice(FILE_LENGTH_START, FILE_LENGTH_START + FILE_LENGTH_SIZE)), BASE_OCTAL) || 0
    const paddedFileLength = (Math.floor(fileLength / CHUNK_SIZE) + (fileLength === 0 ? 0 : 1)) * CHUNK_SIZE
    const fileContent = data.slice(currentFileMeta + CHUNK_SIZE, currentFileMeta + CHUNK_SIZE + fileLength)
    files.push({ name: fileName, content: bytesToString(fileContent), blob: new Blob([fileContent]) })

    // seek to next file meta
    let j = currentFileMeta + CHUNK_SIZE + paddedFileLength
    while (j < data.length) {
      if (data[j] !== 0) {
        currentFileMeta = j
        break
      }
      j += CHUNK_SIZE
    }

    if (j >= data.length) hasNext = false
  }

  return files.filter(f => f.name.substr(-1) !== '/')
}
