<template>
  <ReusableTable :fetchData="fetchClients">
    <template #header>
      <h1>Clients</h1>
    </template>
    <template #table-header>
      <th>ID</th>
      <th>Name</th>
      <th>Created At</th>
      <th>Updated At</th>
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