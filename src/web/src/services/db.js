import store from "../store/index";

class DatabaseService {
  constructor() {
    this.baseurl = 'http://localhost:3000/v1'
  }

  async getClients({pageNumber}) {
    const clients = await fetch(`${this.baseurl}/clients?page=${pageNumber ?? 1}`);
    const data = await clients.json()
    store.dispatch('updateClients', data.data)
    return data
  }

  async getFiles(clientId, {pageNumber}) {
    const files = await fetch(`${this.baseurl}/clients/${clientId}/files?page=${pageNumber ?? 1}`)
    const data = await files.json()
    store.dispatch('updateFiles', data.data)
    return data
  }

  async getJobs(clientId, fileId, {pageNumber}) {
    const jobs = await fetch(`${this.baseurl}/clients/${clientId}/files/${fileId}/jobs?page=${pageNumber ?? 1}`)
    const data = await jobs.json()
    store.dispatch('updateJobs', data.data)
    return data
  }

  async getStats(clientId, fileId, {pageNumber}) {
    const stats = await fetch(`${this.baseurl}/clients/${clientId}/files/${fileId}/stats?page=${pageNumber ?? 1}`)
    const data = await stats.json()
    return data
  }

  async getAudits(clientId, fileId, {pageNumber}) {
    const audits = await fetch(`${this.baseurl}/clients/${clientId}/files/${fileId}/audits?page=${pageNumber ?? 1}`)
    const data = await audits.json()
    store.dispatch('updateAudits', data.data)
    return data
  }

  


  async getEdits(clientId, fileId, jobId) {
    const edits = await fetch(`${this.baseurl}/clients/${clientId}/files/${fileId}/jobs/${jobId}/edits`)
    const data = await edits.json()
    return data
  }
}

export default new DatabaseService();