export class TemplateApi {
  _api = null

  constructor(api) {
    this._api = api
  }

  async listRepos() {
    const res = await this._api.get('/api/templates')
    return res.data
  }

  async listRepoTemplates(repo) {
    const res = await this._api.get(`/api/templates/${repo}`)
    return res.data
  }

  async listAllTemplates() {
    const res = {}
    const repos = await this.listRepos()
    // treating `local` specially to have it come first in iterations over the
    // result objects properties
    if (repos.filter(r => r.name === 'local').length > 0)
      res['local'] = await this.listRepoTemplates('local')
    await Promise.all(repos.filter(r => r.name !== 'local').map(async repo => {
      res[repo.name] = await this.listRepoTemplates(repo.name)
    }))
    return res
  }

  async get(repo, name) {
    const res = await this._api.get(`/api/templates/${repo}/${name}`)
    return res.data
  }

  async exists(repo, name) {
    try {
      const res = await this._api.get(`/api/templates/${repo}/${name}`, undefined, undefined, {unhandledErrors: [404]})
      if (!res.status || res.status >= 400) throw new Error('template doesn\'t exist, behave yourself js...')
      return true
    } catch (e) {
      return false
    }
  }

  async save(name, template) {
    await this._api.put(`/api/templates/local/${name}`, template)
    return true
  }

  async delete(name) {
    await this._api.delete(`/api/templates/local/${name}`)
    return true
  }

  async getRepo(repo) {
    const res = await this._api.get(`/api/templates/${repo}`)
    return res.data
  }

  async saveRepo(repo, config) {
    await this._api.put(`/api/templates/${repo}`, config)
    return true
  }

  async deleteRepo(repo) {
    await this._api.delete(`/api/templates/${repo}`)
    return true
  }
}
