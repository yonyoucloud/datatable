<template>
  <div>
    <vxe-form>
      <vxe-form-item title="当前表">
        <vxe-select v-model="table" placeholder="请选择表名" :options="tables" clearable filterable :loading="loadingTables" @change="changeTable"></vxe-select>
      </vxe-form-item>
    </vxe-form>
    <vxe-grid v-bind="gridOptions"></vxe-grid>
  </div>
</template>

<script setup lang="jsx">
import { ref, onMounted, reactive } from 'vue'
import VXETable, { config } from 'vxe-table'
import FilterComplex from './components/FilterComplex.vue'
import axios from 'axios'

const baseURL = window.CONFIG && window.CONFIG.apiHost || 'http://localhost:8889/'

const openMessage = (options) => {
  VXETable.modal.message(options)
}
const closeMessage = (id) => {
  VXETable.modal.close(id)
}

// 创建一个条件的渲染器
VXETable.renderer.add('FilterComplex', {
  // 不显示底部按钮，使用自定义的按钮
  showFilterFooter: false,
  // 筛选模板
  renderFilter (renderOpts, params) {
    return <FilterComplex params={ params }></FilterComplex>
  },
  // 重置数据方法
  filterResetMethod (params) {
    const { options } = params
    options.forEach((option) => {
      option.data = { type: '', name: '' }
    })
  },
  // 筛选数据方法
  filterMethod (params) {
    const { option, row, column } = params
    const cellValue = row[column.field]
    const { name, type } = option.data
    if (cellValue) {
      return cellValue.indexOf(name) > -1
    }
    return false
  }
})

const loadingTables = ref(false)
const table = ref()
const tables = []

const tableContents = ref([])
const loadingContents = ref(false)
const pagerConfig = ref({ 
  enabled: true,
  background: true,
  layouts: ['PrevJump', 'PrevPage', 'JumpNumber', 'NextPage', 'NextJump', 'Sizes', 'FullJump', 'Total'],
  pageSize: 20,
  currentPage: 1,
  total: 0
})
const tableColumns = ref([])
const sortField = ref(""), sortOrder = ref("")
let filterOptions = new Map()

// 模拟分页接口
const getTableContents = async () => {
  try {
    loadingContents.value = true
    const response = await axios.post(baseURL + 'api/v1/list/table/contents/' + table.value, {
      page_size: pagerConfig.value.pageSize,
      current_page: pagerConfig.value.currentPage,
      sort_field: sortField.value,
      sort_order: sortOrder.value,
      filter: JSON.stringify(Object.fromEntries(filterOptions))
    })
    const data = response.data
    tableContents.value = data.data.list
    pagerConfig.value.total = data.data.total
  } catch (error) {
    openMessage({ content: error, status: 'error' })
  } finally {
    loadingContents.value = false
  }
}

// 前端排序
// const handleSortChange = ({ sortList }) => {
//   sortList.map((item) => `${item.field},${item.order}`).join('; ')
// }

// 后端排序
const handleSortChange = ({ field, order }) => {
  sortField.value = field
  sortOrder.value = order
  getTableContents()
}

const handlePagerChange = ({ currentPage, pageSize }) => {
  pagerConfig.value.currentPage = currentPage
  pagerConfig.value.pageSize = pageSize
  getTableContents()
}

const handleFilterChange = ({ column, property, values, datas, filterList, $event }) => {
  if (datas.length > 0) {
    const { type, name } = datas[0]
    filterOptions.set(property, { type: type, name: name })
  } else {
    filterOptions.delete(property)
  }
  getTableContents()
}

const gridOptions = reactive({
  round: true,
  border: true,
  // height: 530,
  size: 'mini',
  loading: loadingContents,
  rowConfig: {
    keyField: 'id'
  },
  columnConfig: {
    resizable: true,
    isHover: true
  },
  checkboxConfig: {
    reserve: true
  },
  pagerConfig: pagerConfig,
  sortConfig: {
    // multiple: true,
    remote: true,
    // trigger: 'cell'
  },
  columns: tableColumns,
  data: tableContents,
  onSortChange: handleSortChange,
  onPageChange: handlePagerChange,
  filterConfig: {
    remote: true,
  },
  onFilterChange: handleFilterChange
})

const getTables = async () => {
  try {
    loadingTables.value = true
    const response = await axios.get(baseURL + 'api/v1/list/tables')
    const data = response.data.data
    for (let index = 0; index < data.length; index++) {
      tables.push({
        value: data[index],
        label: data[index]
      })
    }
  } catch (error) {
    openMessage({ content: error, status: 'error' })
  } finally {
    loadingTables.value = false
  }
}

const changeTable = () => {
  if (table.value == null || table.value == "") {
    return
  }
  getFields()
}

const getFields = async () => {
  try {
    loadingTables.value = true
    const response = await axios.get(baseURL + 'api/v1/list/table/fields/' + table.value)
    const data = response.data.data
    tableColumns.value = []
    tableColumns.value.push({ type: 'seq', title: '序号', width: '60', fixed: 'left' })
    filterOptions.clear()
    for (let index = 0, length = data.length; index < length; index++) {
      let field = data[index].Field
      tableColumns.value.push({ 
        field: field,
        title: field,
        width: '200',
        sortable: true,
        showOverflow: true,
        filters: [{ data: { type: '', name: '' } }],
        filterRender: {name: 'FilterComplex'}
      })
    }
    getTableContents()
  } catch (error) {
    openMessage({ content: error, status: 'error' })
  } finally {
    loadingTables.value = false
  }
}

onMounted(() => {
  getTables()
})
</script>

<style scoped>
.vxe-select {
    width: 350px;
}
</style>
