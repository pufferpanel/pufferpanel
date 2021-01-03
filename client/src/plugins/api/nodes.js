// used as a mixin for the ApiClient, use of `this` refers to the ApiClient instance
export const NodesApi = {
  getNodes () {
    return this.withErrorHandling(async ctx => {
      return (await ctx.$http.get('/api/nodes')).data
    })
  },

  createNode (newNode) {
    return this.withErrorHandling(async ctx => {
      await ctx.$http.post('/api/nodes', fixNode(newNode))
      return true
    })
  },

  getNode (id) {
    return this.withErrorHandling(async ctx => {
      return (await ctx.$http.get(`/api/nodes/${id}`)).data
    })
  },

  getNodeDeployment (id) {
    return this.withErrorHandling(async ctx => {
      return (await ctx.$http.get(`/api/nodes/${id}/deployment`)).data
    })
  },

  updateNode (id, updatedNode) {
    return this.withErrorHandling(async ctx => {
      await ctx.$http.put(`/api/nodes/${id}`, fixNode(updatedNode))
      return true
    })
  },

  deleteNode (id) {
    return this.withErrorHandling(async ctx => {
      await ctx.$http.delete(`/api/nodes/${id}`)
      return true
    })
  }
}

function fixNode (node) {
  return {
    ...node,
    publicPort: Number(node.publicPort),
    privatePort: Number(node.privatePort),
    sftpPort: Number(node.sftpPort)
  }
}
