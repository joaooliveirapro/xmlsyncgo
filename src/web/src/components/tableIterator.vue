<template>
  <div>
    <div class="d-flex justify-content-between align-items-center">
      <slot name="header"></slot>
    </div>
    <table class="table table-light table-hover table-striped mb-1">
      <thead>
        <tr>
          <slot name="table-header" :sort="sort" :sortColumn="sortColumn" :sortDirection="sortDirection"></slot>
        </tr>
      </thead>
      <tbody>
        <slot name="table-rows" :items="sortedData"></slot>
      </tbody>
    </table>
    <div class="d-flex justify-content-between">
      <span class="text-muted">Page {{ currentPage }} of {{ response.totalPages }}</span>
      <span class="m-0 text-muted">{{ response.data.length }} total results</span>
    </div>
    <nav aria-label="Page navigation example">
      <ul class="pagination justify-content-center">
        <li :class="currentPage == 1 ? 'page-item disabled' : 'page-item'">
          <a @click.prevent="goToPreviousPage" class="page-link" href="#" tabindex="-1" aria-disabled="true">Previous</a>
        </li>
        <li :class="response.totalPages > 1 ? 'page-item' : 'page-item disabled'">
          <a @click.prevent="goToNextPage" class="page-link" href="#">Next</a>
        </li>
      </ul>
    </nav>
  </div>
</template>

<script>
export default {
  name: 'ReusableTable',
  props: {
    fetchData: {
      type: Function,
      required: true,
    },
    initialPage: {
      type: Number,
      default: 1,
    },
  },
  mounted() {
    this.fetchDataAndUpdate(this.initialPage);
  },
  data() {
    return {
      currentPage: this.initialPage,
      response: {
        data: [],
        page: 0,
        total: 0,
        totalPages: 0,
      },
      sortColumn: null,
      sortDirection: 'asc',
    };
  },
  computed: {
    sortedData() {
      if (!this.sortColumn) {
        return this.response.data;
      }

      return [...this.response.data].sort((a, b) => {
        let valueA = a[this.sortColumn];
        let valueB = b[this.sortColumn];

        // Attempt to convert to numbers if possible
        const numA = Number(valueA);
        const numB = Number(valueB);

        if (!isNaN(numA) && !isNaN(numB)) {
          valueA = numA;
          valueB = numB;
        }

        let comparison = 0;
        if (valueA < valueB) {
          comparison = -1;
        } else if (valueA > valueB) {
          comparison = 1;
        }

        return this.sortDirection === 'asc' ? comparison : -comparison;
      });
    },
  },
  watch: {
    currentPage(newVal) {
      this.fetchDataAndUpdate(newVal);
    },
  },
  methods: {
    async fetchDataAndUpdate(page) {
      try {
        const result = await this.fetchData({ pageNumber: page });
        this.response = result;
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    },
    goToPreviousPage() {
      if (this.currentPage > 1) {
        this.currentPage--;
      }
    },
    goToNextPage() {
      if (this.currentPage < this.response.totalPages) {
        this.currentPage++;
      }
    },
    sort(column) {
      if (this.sortColumn === column) {
        this.sortDirection = this.sortDirection === 'asc' ? 'desc' : 'asc';
      } else {
        this.sortColumn = column;
        this.sortDirection = 'asc';
      }
    },
  },
};
</script>

<style scoped>
tr {
  cursor: pointer;
}
</style>