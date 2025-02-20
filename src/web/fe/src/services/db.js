class DatabaseService {
  constructor() {
    this.baseurl = 'http://localhost:3000'
  }

  async getClients() {
    const clients = await fetch(`${this.baseurl}/clients`);
    const data = await clients.json()
    return data
  }

  async getFiles(clientId) {
    const files = await fetch(`${this.baseurl}/clients/${clientId}/files`)
    const data = await files.json()
    return data
  }

  async getAudits(clientId, fileId) {
    const audits = await fetch(`${this.baseurl}/clients/${clientId}/files/${fileId}/audits`)
    const data = await audits.json()
    return data
  }

  async getStats(clientId, fileId) {
    const stats = await fetch(`${this.baseurl}/clients/${clientId}/files/${fileId}/stats`)
    const data = await stats.json()
    return data
  }

  async getJobs(clientId, fileId) {
    const jobs = await fetch(`${this.baseurl}/clients/${clientId}/files/${fileId}/jobs`)
    const data = await jobs.json()
    return data
  }

  async getEdits(clientId, fileId, jobId) {
    const edits = await fetch(`${this.baseurl}/clients/${clientId}/files/${fileId}/jobs/${jobId}/edits`)
    const data = await edits.json()
    return data
  }
}

export default new DatabaseService();