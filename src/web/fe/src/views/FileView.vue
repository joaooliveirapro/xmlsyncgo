<template>
  <div>
    <h5 class="card-title my-4">{{file.remoteFilename}}</h5>
    <ul class="nav nav-pills mb-3" id="pills-tab" role="tablist">
      <li class="nav-item" role="presentation">
        <button class="nav-link" id="pills-home-tab" data-bs-toggle="pill" data-bs-target="#pills-home" type="button" role="tab" aria-controls="pills-home" aria-selected="true">Overview</button>
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
          </template>
          <template #table-header>
            <th>Created at</th>
            <th>Jobs added</th>
            <th>Jobs edited</th>
            <th>Jobs removed</th>
            <th>Keys added</th>
            <th>Keys edited</th>
            <th>Keys removed</th>
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
          <template #table-header>
            <th>Id</th>
            <th>External Id</th>
            <th>Edits</th>
            <th>Last Edit</th>
            <th>Created</th>
            <th>Updated</th>
          </template>
          <template #table-rows="slotProps">
            <tr v-for="j in slotProps.items" :key="j.id">
              <td>{{ j.id }}</td>
              <td>{{ j.externalReference }}</td>
              <td>{{ j.edits?.length }}</td>
              <td>{{ j.edits?.[j.edits?.length - 1].type }}</td>
              <td>{{ new Date(j.createdAt).toUTCString() }}</td>
              <td>{{ new Date(j.updatedAt).toUTCString() }}</td>
            </tr>
          </template>
        </ReusableTable>
      </div>
      
      <div class="tab-pane fade" id="pills-contact" role="tabpanel" aria-labelledby="pills-contact-tab">...</div>
    </div>
  </div>
  
  <!-- <div>
    <div class="card mt-4">
      <div class="card-body">
        <h5 class="card-title">{{file.remote_filename}}</h5>
        <h6 class="card-subtitle mb-2 text-muted">Last update: {{file.updated_at}}</h6>
        <h2>Overview</h2>
        <p>Last 24 hours</p>
        <div class="d-flex justify-content-around">
          <div class="card " style="max-width: 18rem;">
            <div class="card-body">
              <h5 class="card-title">Total</h5>
              <p class="card-text text-center h2">{{jobs.length}}</p>
            </div>
          </div>
          <div class="card " style="max-width: 18rem;">
            <div class="card-body">
              <h5 class="card-title">Added</h5>
              <p class="card-text text-center h2">{{jobs.length}}</p>
            </div>
          </div>
          <div class="card " style="max-width: 18rem;">
            <div class="card-body">
              <h5 class="card-title">Edited</h5>
              <p class="card-text text-center h2">{{jobs.length}}</p>
            </div>
          </div>
          <div class="card " style="max-width: 18rem;">
            <div class="card-body">
              <h5 class="card-title">Removed</h5>
              <p class="card-text text-center h2">{{jobs.length}}</p>
            </div>
          </div>
        </div>
        
        <hr>
        <h2>Properties</h2>
        <ul>
          <li v-for="k in Object.keys(file)" :key="k" class="d-flex justify-content-between">
            <span>{{k}}</span> 
            <span>{{file[k]}}</span>
          </li>
        </ul>

  </div> -->
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
    ...mapState(['file'])
  },
  methods: {
    async getJobs(params) {
      const {clientId, fileId} = this.$route.params; 
      const response = await dbservice.getJobs(clientId, fileId, params);
      return response
    },
    async getStats(params) {
      const {clientId, fileId} = this.$route.params; 
      const response = await dbservice.getStats(clientId, fileId, params);
      return response
    }
  }
};
</script>

<style>
</style>