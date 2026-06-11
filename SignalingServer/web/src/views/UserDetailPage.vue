<template>
  <div class="user-detail-page" v-loading="loading">
    <div class="page-header">
      <el-button :icon="ArrowLeft" @click="router.back()">{{ t('common.back') }}</el-button>
      <h2>{{ t('userMgmt.detail') }}</h2>
    </div>

    <!-- User Info -->
    <el-card shadow="never" class="info-card">
      <template #header>
        <div class="card-header">
          <el-icon><User /></el-icon>
          <span>{{ t('userMgmt.basicInfo') }}</span>
          <el-tag :type="user.status ? 'success' : 'info'" size="small" style="margin-left: auto">
            {{ user.status ? t('common.enabled') : t('common.disabled') }}
          </el-tag>
        </div>
      </template>
      <el-descriptions :column="2" border size="small">
        <el-descriptions-item label="ID">{{ user.id }}</el-descriptions-item>
        <el-descriptions-item :label="t('userMgmt.username')">{{ user.username }}</el-descriptions-item>
        <el-descriptions-item :label="t('userMgmt.phone')">{{ user.phone || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="t('userMgmt.email')">{{ user.email || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="t('userMgmt.level')">
          <el-tag size="small">{{ user.level }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="t('userMgmt.channelType')">
          <el-tag :type="user.channelType === '全球' ? 'success' : 'warning'" size="small">
            {{ user.channelType === '全球' ? t('userMgmt.channelGlobal') : t('userMgmt.channelChina') }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="t('userMgmt.deviceCount')">{{ user.deviceCount }}</el-descriptions-item>
        <el-descriptions-item :label="t('userMgmt.createdAt')">{{ formatDate(user.created_at) }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- Bound Devices -->
    <el-card shadow="never" class="info-card">
      <template #header>
        <div class="card-header">
          <el-icon><Monitor /></el-icon>
          <span>{{ t('userMgmt.boundDevices') }}</span>
        </div>
      </template>
      <el-table :data="devices" stripe style="width: 100%" size="small">
        <el-table-column prop="device_id" :label="t('devices.deviceId')" width="110" />
        <el-table-column prop="device_uuid" label="UUID" min-width="180" show-overflow-tooltip />
        <el-table-column :label="t('devices.os')" min-width="120">
          <template #default="{ row }">
            {{ row.os }}{{ row.os_version ? ' ' + row.os_version : '' }}
          </template>
        </el-table-column>
        <el-table-column :label="t('devices.status')" width="90">
          <template #default="{ row }">
            <el-tag :type="row.online ? 'success' : 'info'" size="small">
              {{ row.online ? t('devices.online') : t('devices.offline') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('devices.lastSeen')" width="170">
          <template #default="{ row }">
            {{ formatDate(row.last_seen) }}
          </template>
        </el-table-column>
      </el-table>
      <el-empty v-if="devices.length === 0" :description="t('common.noData')" />
    </el-card>

    <!-- Connection History -->
    <el-card shadow="never" class="info-card">
      <template #header>
        <div class="card-header">
          <el-icon><Connection /></el-icon>
          <span>{{ t('userMgmt.connectionHistory') }}</span>
        </div>
      </template>
      <el-table :data="connectionHistory" stripe style="width: 100%" size="small">
        <el-table-column prop="time" :label="t('dashboard.time')" width="180" />
        <el-table-column prop="deviceId" :label="t('dashboard.deviceId')" width="120" />
        <el-table-column prop="action" :label="t('dashboard.activity')" width="150" />
        <el-table-column prop="details" :label="t('dashboard.details')" show-overflow-tooltip />
        <el-table-column prop="status" :label="t('common.status')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'success' ? 'success' : 'warning'" size="small">
              {{ row.status === 'success' ? t('common.success') : t('common.failed') }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
      <el-empty v-if="connectionHistory.length === 0" :description="t('dashboard.noActivity')" />
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft, User, Monitor, Connection } from '@element-plus/icons-vue'
import { getUserDetail } from '../api/users.js'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const loading = ref(false)

const user = ref({})
const devices = ref([])
const connectionHistory = ref([])

function formatDate(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}

async function loadDetail() {
  const id = route.params.id
  if (!id) return

  loading.value = true
  try {
    const data = await getUserDetail(id)
    user.value = data.user || data
    devices.value = data.devices || []
    connectionHistory.value = data.connectionHistory || []
  } catch (e) {
    ElMessage.error(t('userMgmt.loadFailed') + ': ' + e.message)
  } finally {
    loading.value = false
  }
}

onMounted(loadDetail)
</script>

<style scoped>
.user-detail-page {
  width: 100%;
  padding: 20px;
  box-sizing: border-box;
}

.page-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
  font-size: 22px;
  font-weight: 600;
  color: #303133;
}

.info-card {
  margin-bottom: 16px;
  border-radius: 8px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}
</style>
