<script setup lang="ts">
import { reactive, ref } from 'vue'
import RightOperate from '@/components/RightOperate.vue'
import { useArcoMessage } from '@/hooks/message'
import { constant } from '@/assets/script/constant'
import type { DictValue, DictValueQueryForm } from '@/api/system/dict/types'
import { dictApi } from '@/api/system/dict'
import DictValueForm from '@/views/system/dict/dict-value-form.vue'

const { successMessage, loading } = useArcoMessage();

const modalShow = ref<boolean>(false);
const show = (dictId: number, dictKey: string) => {
  queryForm.dictId = dictId;
  queryForm.dictKey = dictKey;
  queryTableData();
  modalShow.value = true;
}
const onClose = () => {
  tableData.value = [];
  recordTotal.value = 0;
  Object.assign(queryForm, defaultQueryForm);
}

const tableLoading = ref<boolean>(false);
const recordTotal = ref<number>(0);
const tableData = ref<DictValue[]>([]);
const defaultQueryForm: DictValueQueryForm = {
  page: 1,
  size: 10,
  dictId: undefined,
  label: undefined,
  createTimeBegin: undefined,
  createTimeEnd: undefined
}
const dateRange = ref<string[]>([]);
const queryForm = reactive<DictValueQueryForm>({ ...defaultQueryForm })
const queryTableData = async () => {
  tableLoading.value = true;
  try {
    const result = await dictApi.pageDictValue(queryForm);
    if (result && result.code === 200) {
      const { total, records } = result.data;
      recordTotal.value = total;
      tableData.value = records;
    }
  } finally {
    tableLoading.value = false;
  }
}

const queryLoading = ref<boolean>(false);
const handleQuery = async () => {
  queryForm.page = 1;
  queryLoading.value = true;
  try {
    await queryTableData();
  } finally {
    queryLoading.value = false;
  }
}

const handleReset = () => {
  const { dictId, dictKey } = queryForm;
  Object.assign(queryForm, defaultQueryForm);
  queryForm.dictId = dictId;
  queryForm.dictKey = dictKey;
  dateRange.value = [];
}

const handleDateChange = () => {
  if (!dateRange.value) {
    queryForm.createTimeBegin = undefined;
    queryForm.createTimeEnd = undefined;
  } else {
    const [ begin, end ] = dateRange.value;
    queryForm.createTimeBegin = begin;
    queryForm.createTimeEnd = end;
  }
}

const handleDelete = async (record: DictValue) => {
  const loadingMsg = loading('数据删除中')
  try {
    const result = await dictApi.deleteDictValue(record.id);
    if (result.code === 200) {
      successMessage('删除成功');
      await queryTableData();
    }
  } finally {
    loadingMsg.close();
  }
}

const handleQuickChange = async (newValue: number | string | boolean, id: number) => {
  console.log(newValue, id)
  const newStatus = Number(newValue);
  const result = await dictApi.updateDictValueStatus({ id: id, status: newStatus })
  if (result.code === 200) {
    successMessage('更新成功');
    return true;
  }
  return false;
}

const formRef = ref();
const showForm = (record?: DictValue) => {
  formRef.value.show(queryForm.dictId, queryForm.dictKey, record);
}

defineExpose({
  show
})

</script>

<template>
  <a-drawer v-model:visible="modalShow" @close="onClose" width="1200px" title="字典详情" :footer="false">
    <div class="table-card">
      <div class="search-div">
        <div class="search-item">
          <label>标签</label>
          <a-input v-model="queryForm.label" placeholder="请输入字典数据标签" />
        </div>
        <div class="search-item">
          <label>创建时间</label>
          <a-range-picker allow-clear v-model="dateRange" @change="handleDateChange" />
        </div>
        <div class="search-buttons">
          <a-button type="primary" @click="handleQuery" :loading="queryLoading">
            <template #icon><icon-search /></template>
            搜索
          </a-button>
          <a-button @click="handleReset">
            <template #icon><icon-refresh /></template>
            重置
          </a-button>
        </div>
      </div>
      <div class="flex justify-between">
        <div class="flex" style="column-gap: 12px">
          <a-button type="primary" @click="showForm(undefined)">
            <template #icon><icon-plus /></template>
            新增
          </a-button>
        </div>
        <right-operate @refresh="queryTableData"/>
      </div>
      <a-table :data="tableData" :loading="tableLoading" :pagination="false">
        <template #columns>
          <a-table-column title="标签" data-index="label" />
          <a-table-column title="数据值" data-index="value" />
          <a-table-column title="排序" data-index="sort" />
          <a-table-column title="创建时间" data-index="createTime" />
          <a-table-column title="备注" data-index="remark" />
          <a-table-column title="状态">
            <template #cell="{ record }">
              <template v-if="record.status < 2">
                <a-switch :checked-value="0" :unchecked-value="1" v-model="record.status"
                          :before-change="newValue => handleQuickChange(newValue, record.dictId)"
                />
              </template>
              <template v-else>
                <span style="color: red">字典禁用</span>
              </template>
            </template>
          </a-table-column>
          <a-table-column title="操作" align="center">
            <template #cell="{ record }">
              <a-button type="text" shape="circle" @click="showForm(record)">
                <template #icon><icon-edit /></template>
              </a-button>
              <a-popconfirm content="是否确认删除数据？" type="error" position="lt"
                            :ok-button-props="{ status: 'danger' }"
                            @ok="handleDelete(record)"
              >
                <a-button type="text" shape="circle" status="danger">
                  <template #icon><icon-delete /></template>
                </a-button>
              </a-popconfirm>
            </template>
          </a-table-column>
        </template>
      </a-table>
      <div class="flex justify-end">
        <a-pagination :total="recordTotal" size="medium"
                      v-model:current="queryForm.page"
                      v-model:page-size="queryForm.size"
                      :page-size-options="constant.pageSizeOption"
                      show-page-size
                      show-total
                      @page-size-change="handleQuery"
                      @change="queryTableData"
        />
      </div>
      <dict-value-form ref="formRef" @reload="queryTableData"/>
    </div>
  </a-drawer>
</template>

<style scoped lang="scss">

</style>