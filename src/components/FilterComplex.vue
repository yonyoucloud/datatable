<template>
    <div class="my-filter-complex">
      <div class="my-fc-type">
        <vxe-radio v-model="myFilter.option.data.type" name="fType" label="IN">包含</vxe-radio>
        <vxe-radio v-model="myFilter.option.data.type" name="fType" label="=">等于</vxe-radio>
        <vxe-radio v-model="myFilter.option.data.type" name="fType" label="<>">不等于</vxe-radio>
        <vxe-radio v-model="myFilter.option.data.type" name="fType" label=">">大于</vxe-radio>
        <vxe-radio v-model="myFilter.option.data.type" name="fType" label=">=">大于或等于</vxe-radio>
        <vxe-radio v-model="myFilter.option.data.type" name="fType" label="<">小于</vxe-radio>
        <vxe-radio v-model="myFilter.option.data.type" name="fType" label="<=">小于或等于</vxe-radio>
        <vxe-radio v-model="myFilter.option.data.type" name="fType" label="RANGE">范围</vxe-radio>
      </div>
      <div class="my-fc-name">
        <vxe-form :data="myFilter.option.data">
          <vxe-form-item title="值1">
            <template #default="{ data }">
              <vxe-input v-model="data.value1" type="text" size="mini" placeholder="请输入值1" @input="changeOptionEvent()"></vxe-input>
            </template>
          </vxe-form-item>
          <vxe-form-item title="值2">
            <template #default="{ data }">
              <vxe-input v-model="data.value2" type="text" size="mini" placeholder="请输入值2" @input="changeOptionEvent()"></vxe-input>
            </template>
          </vxe-form-item>
          <vxe-form-item :title-prefix="{ icon: 'vxe-icon-warning-triangle' }">
            <template #title>
              <span style="color: red;">开启该字段统计（必须是时间字段）</span>
            </template>
            <template #default="{ data }">
              <vxe-switch v-model="data.openstat" open-label="是" close-label="否"></vxe-switch>
            </template>
          </vxe-form-item>
        </vxe-form>
      </div>
      <div class="my-fc-footer">
        <vxe-button status="primary" @click="confirmEvent">确认</vxe-button>
        <vxe-button @click="resetEvent">重置</vxe-button>
      </div>
    </div>
  </template>
  
  <script lang="ts">
  import { defineComponent, type PropType, reactive } from 'vue'
  import { type VxeGlobalRendererHandles } from 'vxe-table'
  
  export default defineComponent({
    name: 'FilterComplex',
    props: {
      params: Object as PropType<VxeGlobalRendererHandles.RenderFilterParams>
    },
    setup (props) {
      const myFilter = reactive({
        option: null as any
      })
  
      const load = () => {
        const { params } = props
        if (params) {
          const { column } = params
          const option = column.filters[0]
          myFilter.option = option
        }
      }
  
      const changeOptionEvent = () => {
        const { params } = props
        const { option } = myFilter
        if (params && option) {
          const { $panel } = params
          const checked = !!option.data.value1
          $panel.changeOption(null, checked, option)
        }
      }
  
      const confirmEvent = () => {
        const { params } = props
        if (params) {
          const { $panel } = params
          $panel.confirmFilter()
        }
      }
  
      const resetEvent = () => {
        const { params } = props
        if (params) {
          const { $panel } = params
          $panel.resetFilter()
        }
      }
  
      load()
  
      return {
        myFilter,
        changeOptionEvent,
        confirmEvent,
        resetEvent
      }
    }
  })
  </script>
  
  <style scoped>
  .my-filter-complex {
    width: 380px;
    padding: 5px 15px 10px 15px;
  }
  .my-filter-complex .my-fc-type {
    padding: 8px 0;
  }
  .my-filter-complex .my-fc-footer {
    text-align: center;
    margin-top: 10px;
  }
  .my-fc-name .vxe-input {
    float: left;
    width: 150px;
  }
  </style>