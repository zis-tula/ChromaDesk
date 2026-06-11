<template>
  <div class="audit-page" v-loading="loading">
    <div class="page-header">
      <h2>{{ t('auditLog.title') }}</h2>
      <el-button type="primary" :icon="Refresh" size="small" @click="loadLogs">{{ t('common.refresh') }}</el-button>
    </div>

    <el-card shadow="never" class="filter-card">
      <div class="filter-bar">
        <el-select v-model="filters.action" :placeholder="t('auditLog.filterAction')" clearable style="width: 160px" @change="handleFilter">
          <el-option label="login" value="login" />
          <el-option label="create_user" value="create_user" />
          <el-option label="update_user" value="update_user" />
          <el-option label="delete_user" value="delete_user" />
          <el-option label="create_admin" value="create_admin" />
          <el-option label="update_admin" value="update_admin" />
          <el-option label="delete_admin" value="delete_admin" />
          <el-option label="update_settings" value="update_settings" />
        </el-select>
        <el-input
          v-model="filters.admin"
          :placeholder="t('auditLog.filterAdmin')"
          clearable
          style="width: 160px"
          @clear="handleFilter"
          @keyup.enter="handleFilter"
        />
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          :start-placeholder="t('auditLog.startDate')"
          :end-placeholder="t('auditLog.endDate')"
          size="default"
          style="width: 280px"
          @change="handleDateChange"
        />
        <el-button :icon="Search" @click="handleFilter">{{ t('common.search') }}</el-button>
      </div>
    </el-card>

    <el-card shadow="never" style="margin-top: 16px">
      <el-table :data="logs" stripe style="width: 100%" size="small">
        <el-table-column prop="created_at" :label="t('auditLog.time')" width="170">
          <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
        </el-table-column>
        <el-table-column prop="admin_username" :label="t('auditLog.admin')" width="120" />
        <el-table-column prop="action" :label="t('auditLog.action')" width="140">
          <template #default="{ row }">
            <el-tag size="small">{{ row.action }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="resource_type" :label="t('auditLog.resourceType')" width="120" />
        <el-table-column prop="resource_id" :label="t('auditLog.resourceId')" width="120" />
        <el-table-column prop="details" :label="t('auditLog.details')" show-overflow-tooltip />
        <el-table-column prop="ip" label="IP" width="140" />
      </el-table>

      <div class="pagination-bar">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.size"
          :page-sizes="[20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadLogs"
          @current-change="loadLogs"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Refresh, Search } from '@element-plus/icons-vue'
import { getAuditLogs } from '../api/audit.js'

const { t } = useI18n()
const loading = ref(false)
const logs = ref([])
const dateRange = ref(null)

const filters = reactive({ action: '', admin: '', dateFrom: '', dateTo: '' })
const pagination = reactive({ page: 1, size: 20, total: 0 })

function formatDate(d) {
  if (!d) return '-'
  return new Date(d).toLocaleString('zh-CN')
}

async function loadLogs() {
  loading.value = true
  try {
    const data = await getAuditLogs({
      page: pagination.page,
      size: pagination.size,
      action: filters.action,
      admin: filters.admin,
      dateFrom: filters.dateFrom,
      dateTo: filters.dateTo
    })
    logs.value = data.items || []
    pagination.total = data.total || 0
  } catch (e) {
    ElMessage.error(t('common.loadFailed') + ': ' + e.message)
  } finally {
    loading.value = false
  }
}

function handleFilter() {
  pagination.page = 1
  loadLogs()
}

function handleDateChange(val) {
  if (val) {
    filters.dateFrom = val[0].toISOString()
    filters.dateTo = val[1].toISOString()
  } else {
    filters.dateFrom = ''
    filters.dateTo = ''
  }
  handleFilter()
}

onMounted(loadLogs)
</script>

<style scoped>
.audit-page { width: 100%; padding: 20px; box-sizing: border-box; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.page-header h2 { margin: 0; font-size: 24px; font-weight: 600; color: #303133; }
.filter-card { border-radius: 8px; }
.filter-bar { display: flex; gap: 12px; flex-wrap: wrap; align-items: center; }
.pagination-bar { display: flex; justify-content: flex-end; margin-top: 16px; }
</style>
