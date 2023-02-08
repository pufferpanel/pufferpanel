export class NodeApi {
  _api = null

  constructor(api) {
    this._api = api
  }

  fixNode(node) {
    return {
      ...node,
      publicPort: Number(node.publicPort),
      privatePort: Number(node.privatePort),
      sftpPort: Number(node.sftpPort)
    }
  }

  async list() {
    const res = await this._api.get('/api/nodes')
    return res.data
  }

  async get(id) {
    const res = await this._api.get(`/api/nodes/${id}`)
    return res.data
  }

  async deployment(id) {
    const res = await this._api.get(`/api/nodes/${id}/deployment`)
    return res.data
  }

  async create(node) {
    await this._api.post('/api/nodes/', this.fixNode(node))
    try {
      const nodes = await this.list()
      return nodes.find(n => n.name === node.name).id
    } catch (e) {
      return -1
    }
  }

  async update(id, node) {
    await this._api.post(`/api/nodes/${id}`, this.fixNode(node))
    return true
  }

  async delete(id) {
    await this._api.delete(`/api/nodes/${id}`)
    return true
  }
}
