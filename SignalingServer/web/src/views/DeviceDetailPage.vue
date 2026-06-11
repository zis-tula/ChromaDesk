<template>
  <div class="device-detail-page" v-loading="loading">
    <div class="page-header">
      <el-button :icon="ArrowLeft" @click="router.back()">{{ t('common.back') }}</el-button>
      <h2>{{ t('devices.detail') }}</h2>
    </div>

    <!-- Device Info -->
    <el-card shadow="never" class="info-card">
      <template #header>
        <div class="card-header">
          <el-icon><Monitor /></el-icon>
          <span>{{ t('devices.basicInfo') }}</span>
          <el-tag :type="device.online ? 'success' : 'info'" size="small" style="margin-left: auto">
            {{ device.online ? t('devices.online') : t('devices.offline') }}
          </el-tag>
        </div>
      </template>
      <el-descriptions :column="2" border size="small">
        <el-descriptions-item :label="t('devices.deviceId')">{{ device.device_id }}</el-descriptions-item>
        <el-descriptions-item label="UUID">{{ device.device_uuid }}</el-descriptions-item>
        <el-descriptions-item :label="t('devices.os')">{{ device.os }}{{ device.os_version ? ' ' + device.os_version : '' }}</el-descriptions-item>
        <el-descriptions-item :label="t('devices.appVersion')">{{ device.app_version || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="t('devices.lastSeen')">{{ formatDate(device.last_seen) }}</el-descriptions-item>
        <el-descriptions-item :label="t('devices.createdAt')">{{ formatDate(device.created_at) }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- Bound User -->
    <el-card shadow="never" class="info-card" v-if="boundUser">
      <template #header>
        <div class="card-header">
          <el-icon><User /></el-icon>
          <span>{{ t('devices.boundUser') }}</span>
        </div>
      </template>
      <el-descriptions :column="2" border size="small">
        <el-descriptions-item :label="t('userMgmt.username')">{{ boundUser.username }}</el-descriptions-item>
        <el-descriptions-item :label="t('userMgmt.phone')">{{ boundUser.phone || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="t('userMgmt.email')">{{ boundUser.email || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="t('userMgmt.level')">
          <el-tag size="small">{{ boundUser.level }}</el-tag>
        </el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- Connection History -->
    <el-card shadow="never" class="info-card">
      <template #header>
        <div class="card-header">
          <el-icon><Connection /></el-icon>
          <span>{{ t('devices.connectionHistory') }}</span>
        </div>
      </template>
      <el-table :data="connectionHistory" stripe style="width: 100%" size="small">
        <el-table-column prop="time" :label="t('dashboard.time')" width="180" />
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
import { ArrowLeft, Monitor, User, Connection } from '@element-plus/icons-vue'
import { getDeviceDetail } from '../api/admin_device.js'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const loading = ref(false)

const device = ref({})
const boundUser = ref(null)
const connectionHistory = ref([])

function formatDate(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}

async function loadDetail() {
  const deviceId = route.params.deviceId
  if (!deviceId) return

  loading.value = true
  try {
    const data = await getDeviceDetail(deviceId)
    device.value = data.device || data
    boundUser.value = data.boundUser || null
    connectionHistory.value = data.connectionHistory || []
  } catch (e) {
    ElMessage.error(t('devices.loadFailed') + ': ' + e.message)
  } finally {
    loading.value = false
  }
}

onMounted(loadDetail)
</script>

<style scoped>
.device-detail-page {
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
