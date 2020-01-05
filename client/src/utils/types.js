export function typeNode (node) {
  return {
    name: node.name,
    publicHost: node.publicHost,
    publicPort: Number(node.publicPort),
    privateHost: node.privateHost,
    privatePort: Number(node.privatePort),
    sftpPort: Number(node.sftpPort)
  }
}
