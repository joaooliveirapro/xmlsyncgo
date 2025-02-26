<template>
  <div>
    <h3 class="card-title my-4">{{file.remoteFilename}}</h3>
    
    <div>

    </div>

    <ul class="nav nav-pills mb-3" id="pills-tab" role="tablist">
      <li class="nav-item" role="presentation">
        <button class="nav-link" id="pills-home-tab" data-bs-toggle="pill" data-bs-target="#pills-home" type="button" role="tab" aria-controls="pills-home" aria-selected="true">Stats</button>
      </li>
      <li class="nav-item" role="presentation">
        <button class="nav-link" id="pills-profile-tab" data-bs-toggle="pill" data-bs-target="#pills-profile" type="button" role="tab" aria-controls="pills-profile" aria-selected="false">Jobs</button>
      </li>
      <li class="nav-item" role="presentation">
        <button class="nav-link" id="pills-contact-tab" data-bs-toggle="pill" data-bs-target="#pills-contact" type="button" role="tab" aria-controls="pills-contact" aria-selected="false">Properties</button>
      </li>
    </ul>
    <div class="tab-content" id="pills-tabContent">
      <div class="tab-pane fade show active" id="pills-home" role="tabpanel" aria-labelledby="pills-home-tab">
        <ReusableTable :fetchData="getStats">
          <template #header>
            <h1>Stats</h1>
            <div class="d-flex">
              <div class="card">
                <div class="card-body">
                  <h5 class="cart-title">Jobs added</h5>
                  {{lastAudit()}}
                </div>
              </div>
              <div class="card">
                <div class="card-body">
                  <h5 class="cart-title">Jobs removed</h5>
                </div>
              </div>
              <div class="card">
                <div class="card-body">
                  <h5 class="cart-title">Jobs edited</h5>
                </div>
              </div>
            </div>
          </template>
          <template #table-header="slotProps">
            <th @click="slotProps.sort('createdAt')">
              Created at
              <span v-if="slotProps.sortColumn === 'createdAt'">
                {{ slotProps.sortDirection === 'asc' ? '▲' : '▼' }}
              </span>
            </th>
            <th @click="slotProps.sort('jobsAdded')">
              Jobs added
            </th>
            <th @click="slotProps.sort('jobsEdited')">
              Jobs edited
            </th>
            <th @click="slotProps.sort('jobsRemoved')">
              Jobs removed
            </th>
            <th @click="slotProps.sort('keysAdded')">
              Keys added
            </th>
            <th @click="slotProps.sort('keysEdited')">
              Keys edited
            </th>
            <th @click="slotProps.sort('keysRemoved')">
              Keys removed
            </th>
          </template>
          <template #table-rows="slotProps">
            <tr v-for="s in slotProps.items" :key="s.id">
              <td>{{ new Date(s.createdAt).toUTCString() }}</td>
              <td>{{ JSON.parse(s.jsonStr).JOBS_ADDED }}</td>
              <td>{{ JSON.parse(s.jsonStr).JOBS_EDITED }}</td>
              <td>{{ JSON.parse(s.jsonStr).JOBS_REMOVED }}</td>
              <td>{{ JSON.parse(s.jsonStr).ADDED_KEYS }}</td>
              <td>{{ JSON.parse(s.jsonStr).EDITED_KEYS }}</td>
              <td>{{ JSON.parse(s.jsonStr).REMOVED_KEYS }}</td>
            </tr>
          </template>
        </ReusableTable>
      </div>
      
      <div class="tab-pane fade" id="pills-profile" role="tabpanel" aria-labelledby="pills-profile-tab">
        <!-- Jobs ReusableTable -->
        <ReusableTable :fetchData="getJobs">
          <template #header>
            <h1>Jobs</h1>
          </template>
          <template #table-header="slotProps">
            <th @click="slotProps.sort('id')">
              Id
            </th>
            <th @click="slotProps.sort('externalReference')">
              External Id
            </th>
            <th @click="slotProps.sort('editsLength')">
              Edits
            </th>
            <th @click="slotProps.sort('lastEditType')">
              Last Edit
            </th>
            <th @click="slotProps.sort('createdAt')">
              Created
              <span v-if="slotProps.sortColumn === 'createdAt'">
                {{ slotProps.sortDirection === 'asc' ? '▲' : '▼' }}
              </span>
            </th>
            <th @click="slotProps.sort('updatedAt')">
              Updated
              <span v-if="slotProps.sortColumn === 'updatedAt'">
                {{ slotProps.sortDirection === 'asc' ? '▲' : '▼' }}
              </span>
            </th>
          </template>
          <template #table-rows="slotProps">
            <tr v-for="j in slotProps.items" :key="j.id">
              <td>{{ j.id }}</td>
              <td>{{ j.externalReference }}</td>
              <td>{{ j.edits?.length }}</td>
              <td>
                <span v-if="lastEdit(j) == 'REMOVED_JOB'" class="text-danger">▼</span>
                <span v-if="lastEdit(j) == 'ADDED_JOB'" class="text-success">▲</span>
                <span v-if="lastEdit(j).indexOf('KEY') > -1">✏</span>
                {{ lastEdit(j) }}
              </td>
              <td>{{ new Date(j.createdAt).toUTCString() }}</td>
              <td>{{ new Date(j.updatedAt).toUTCString() }}</td>
            </tr>
          </template>
        </ReusableTable>
      </div>
      
      <div class="tab-pane fade" id="pills-contact" role="tabpanel" aria-labelledby="pills-contact-tab">...</div>
    </div>

  </div>
</template>

<script>
import store from '../store';
import { mapState } from 'vuex'
import dbservice from '../services/db'
import ReusableTable from '../components/tableIterator.vue';

export default {
  components: {ReusableTable},
  data(){
    return {
      file: {
        remoteFilename: ''
      }
    }
  },
  async created() {
    this.file = store.state.files.find(f => f.id == this.$route.params.fileId);
  },
  computed: {
    ...mapState(['file', 'audits'])
  },
  methods: {
    lastEdit(j) {
      return j.edits?.[j.edits?.length - 1]?.type
    },
    lastAudit() {
      return this.audits?.[this.audits.length - 1]
    },
    async getJobs(params) {
      const {clientId, fileId} = this.$route.params; 
      const response = await dbservice.getJobs(clientId, fileId, params);
      return response
    },
    async getStats(params) {
      const {clientId, fileId} = this.$route.params; 
      const response = await dbservice.getStats(clientId, fileId, params);
      return response
    },
    async getAudits(params) {
      const {clientId, fileId} = this.$route.params; 
      const response = await dbservice.getAudits(clientId, fileId, params);
      return response
    }
  }
};
</script>

<style>
</style>