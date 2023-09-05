import { createApp } from 'vue'
import VXETable from 'vxe-table'
import 'vxe-table/lib/style.css'
import App from './App.vue'

function useTable(app : any) {
  app.use(VXETable);
}

createApp(App).use(useTable).mount('#app')