import { createStore } from 'vuex'

export default createStore({
  state: {
    client: {},
    clients: [],
    files: [],
    jobs: [],
    audits: []
  },
  getters: {
  },
  mutations: {
    client(state, newClient) {
      state.client = newClient
    },
    clients(state, newClients) {
      state.clients = newClients
    },
    files(state, newFiles) {
      state.files = newFiles
    },
    jobs(state, newJobs) {
      state.jobs = newJobs
    },
    audits(state, newAudits) {
      state.audits = newAudits
    }
  },
  actions: {
    updateClient({commit}, client) {
      commit('client', client)
    },
    updateClients({commit}, clients) {
      commit('clients', clients)
    },
    updateFiles({commit}, files) {
      commit('files', files)
    },
    updateJobs({commit}, jobs) {
      commit('jobs', jobs)
    },
    updateAudits({commit}, jobs) {
      commit('audits', jobs)
    }
  },
  modules: {
  }
})
