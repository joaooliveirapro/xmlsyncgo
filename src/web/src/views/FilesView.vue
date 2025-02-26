<template>
  <ReusableTable :fetchData="fetchFiles">
    <template #header>
      <h1>{{ client.name }} Files</h1>
    </template>
    <template #table-header="slotProps">
      <th @click="slotProps.sort('id')">
        Id
        <span v-if="slotProps.sortColumn === 'id'">
          {{ slotProps.sortDirection === 'asc' ? '▲' : '▼' }}
        </span>
      </th>
      <th @click="slotProps.sort('hostname')">
        Host
        <span v-if="slotProps.sortColumn === 'hostname'">
          {{ slotProps.sortDirection === 'asc' ? '▲' : '▼' }}
        </span>
      </th>
      <th @click="slotProps.sort('remoteFilename')">
        Filename
        <span v-if="slotProps.sortColumn === 'remoteFilename'">
          {{ slotProps.sortDirection === 'asc' ? '▲' : '▼' }}
        </span>
      </th>
      <th @click="slotProps.sort('rootKey')">
        Root Key
        <span v-if="slotProps.sortColumn === 'rootKey'">
          {{ slotProps.sortDirection === 'asc' ? '▲' : '▼' }}
        </span>
      </th>
      <th @click="slotProps.sort('jobNodeKey')">
        Job Node Key
        <span v-if="slotProps.sortColumn === 'jobNodeKey'">
          {{ slotProps.sortDirection === 'asc' ? '▲' : '▼' }}
        </span>
      </th>
      <th @click="slotProps.sort('externalReferenceKey')">
        Ext Ref Key
        <span v-if="slotProps.sortColumn === 'externalReferenceKey'">
          {{ slotProps.sortDirection === 'asc' ? '▲' : '▼' }}
        </span>
      </th>
      <th @click="slotProps.sort('createdAt')">
        Created at
        <span v-if="slotProps.sortColumn === 'createdAt'">
          {{ slotProps.sortDirection === 'asc' ? '▲' : '▼' }}
        </span>
      </th>
      <th @click="slotProps.sort('updatedAt')">
        Updated at
        <span v-if="slotProps.sortColumn === 'updatedAt'">
          {{ slotProps.sortDirection === 'asc' ? '▲' : '▼' }}
        </span>
      </th>
    </template>
    <template #table-rows="slotProps">
      <tr v-for="f in slotProps.items" :key="f.id">
        <td>{{ f.id }}</td>
        <td>{{ f.hostname }}</td>
        <td>
          <router-link :to="{ name: 'file', params: { clientId: f.id, fileId: f.id } }">{{ f.remoteFilename }}</router-link>
        </td>
        <td>{{ f.rootKey }}</td>
        <td>{{ f.jobNodeKey }}</td>
        <td>{{ f.externalReferenceKey }}</td>
        <td>{{ new Date(f.createdAt).toUTCString() }}</td>
        <td>{{ new Date(f.updatedAt).toUTCString() }}</td>
      </tr>
    </template>
  </ReusableTable>
</template>


<script>
import feather from "feather-icons";
import dbservice from "../services/db";
import { mapState } from 'vuex'
import ReusableTable from '../components/tableIterator.vue'
import store from '../store/index';
export default {
  name: "FilesView",
  mounted() {
    feather.replace();
  },
  components: {
    ReusableTable
  },
  created() {
    store.dispatch('updateClient', store.state.clients.find(c => c.id == this.$route.params.clientId))
  },
  computed: {
    ...mapState(['client'])
  },
  methods: {
    async fetchFiles(params) {
      return await dbservice.getFiles(this.$route.params.clientId, params);
    },
  }
};
</script>

<style scoped>
tr {
  cursor: pointer;
}
</style>