<template>
  <ReusableTable :fetchData="fetchClients">
    <template #header>
      <h1>Clients</h1>
    </template>
    <template #table-header="slotProps">
      <th @click="slotProps.sort('id')">
        ID
        <span v-if="slotProps.sortColumn === 'id'">
          {{ slotProps.sortDirection === 'asc' ? '▲' : '▼' }}
        </span>
      </th>
      <th @click="slotProps.sort('name')">
        Name
        <span v-if="slotProps.sortColumn === 'name'">
          {{ slotProps.sortDirection === 'asc' ? '▲' : '▼' }}
        </span>
      </th>
      <th @click="slotProps.sort('createdAt')">
        Created At
        <span v-if="slotProps.sortColumn === 'createdAt'">
          {{ slotProps.sortDirection === 'asc' ? '▲' : '▼' }}
        </span>
      </th>
      <th @click="slotProps.sort('updatedAt')">
        Updated At
        <span v-if="slotProps.sortColumn === 'updatedAt'">
          {{ slotProps.sortDirection === 'asc' ? '▲' : '▼' }}
        </span>
      </th>
    </template>
    <template #table-rows="slotProps">
      <tr v-for="c in slotProps.items" :key="c.id">
        <td>{{ c.id }}</td>
        <td>
          <router-link :to="{ name: 'files', params: { clientId: c.id } }">{{ c.name }}</router-link>
        </td>
        <td>{{ new Date(c.createdAt).toUTCString() }}</td>
        <td>{{ new Date(c.updatedAt).toUTCString() }}</td>
      </tr>
    </template>
  </ReusableTable>
</template>

<script>
import dbservice from '../services/db';
import ReusableTable from '../components/tableIterator.vue';
import feather from 'feather-icons';

export default {
  name: "ClientsView",
  components: {
    ReusableTable,
  },
  mounted() {
    feather.replace();
  },
  methods: {
    async fetchClients(params) {
      return await dbservice.getClients(params);
    },
  },
};
</script>