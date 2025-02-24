import { createRouter, createWebHistory } from 'vue-router'
import DashboardView from '../views/DashboardView.vue'
import ClientsView from '../views/ClientsView.vue'
import FileView from '../views/FileView.vue'
import FilesView from '../views/FilesView.vue'
const routes = [
  {
    path: '/',
    name: 'home',
    component: DashboardView
  },
  {
    path: '/clients',
    name: 'clients',
    component: ClientsView
  },
  {
    path: '/clients/:clientId/files',
    name: 'files',
    component: FilesView
  },
  {
    path: '/clients/:clientId/files/:fileId',
    name: 'file',
    component: FileView
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
