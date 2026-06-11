<template>
  <div class="settings-page" v-loading="loading">
    <div class="page-header">
      <h2>{{ t('settings.title') }}</h2>
      <el-button type="primary" size="small" @click="handleSave" :loading="saving" :icon="Check">
        {{ t('settings.saveSettings') }}
      </el-button>
    </div>

    <div class="settings-container">
      <!-- 基础配置 -->
      <el-card class="settings-card" shadow="never">
        <template #header>
          <div class="card-header">
            <el-icon><Setting /></el-icon>
            <span>{{ t('settings.basic') }}</span>
          </div>
        </template>
        <el-form label-width="120px" label-position="top">
          <el-form-item :label="t('settings.siteEnabled')">
            <el-switch v-model="form.siteEnabled" :active-text="t('common.on')" :inactive-text="t('common.off')" />
          </el-form-item>
          <el-form-item :label="t('settings.siteName')">
            <el-input v-model="form.siteName" :placeholder="t('settings.siteNamePlaceholder')" style="max-width:400px" />
          </el-form-item>
        </el-form>
      </el-card>

      <!-- ICE / TURN 配置 -->
      <el-card class="settings-card" shadow="never">
        <template #header>
          <div class="card-header">
            <el-icon><Connection /></el-icon>
            <span>{{ t('settings.ice') }}</span>
          </div>
        </template>
        <el-form label-width="180px" label-position="top">
          <el-form-item :label="t('settings.turnUrls')">
            <div class="form-tip">{{ t('settings.turnUrlsTip') }}</div>
            <ListEditor
              v-model="turnUrlList"
              :placeholder="t('settings.turnUrlsPlaceholder')"
              :add-text="t('settings.addTurnServer')"
            />
          </el-form-item>

          <el-form-item :label="t('settings.turnAuthSecret')">
            <el-input
              v-model="form.turnAuthSecret"
              :placeholder="t('settings.turnAuthSecretPlaceholder')"
              style="max-width:500px"
              show-password
              type="password"
            />
            <div class="form-tip">{{ t('settings.turnAuthSecretTip') }}</div>
          </el-form-item>

          <el-form-item :label="t('settings.turnTtl')">
            <el-input-number v-model="form.turnCredentialTtl" :min="300" :max="604800" :step="3600" />
            <div class="form-tip">{{ t('settings.turnTtlTip') }}</div>
          </el-form-item>

          <el-form-item :label="t('settings.stunUrls')">
            <div class="form-tip">{{ t('settings.stunUrlsTip') }}</div>
            <ListEditor
              v-model="stunUrlList"
              :placeholder="t('settings.stunUrlsPlaceholder')"
              :add-text="t('settings.addStunServer')"
            />
          </el-form-item>
        </el-form>
      </el-card>

      <!-- 安全配置 -->
      <el-card class="settings-card" shadow="never">
        <template #header>
          <div class="card-header">
            <el-icon><Lock /></el-icon>
            <span>{{ t('settings.security') }}</span>
          </div>
        </template>
        <el-form label-width="180px" label-position="top">
          <el-form-item :label="t('settings.apiKey')">
            <el-input
              v-model="form.apiKey"
              :placeholder="t('settings.apiKeyPlaceholder')"
              style="max-width:500px"
              show-password
              type="password"
            />
            <div class="form-tip">
              {{ t('settings.apiKeyTip') }}
            </div>
          </el-form-item>

          <el-form-item :label="t('settings.allowedOrigins')">
            <div class="form-tip">
              {{ t('settings.allowedOriginsTip') }}
            </div>
            <ListEditor
              v-model="allowedOriginList"
              :placeholder="t('settings.allowedOriginsPlaceholder')"
              :add-text="t('settings.addOrigin')"
            />
          </el-form-item>
        </el-form>
      </el-card>

      <!-- IP 白名单 -->
      <el-card class="settings-card" shadow="never">
        <template #header>
          <div class="card-header">
            <el-icon><Lock /></el-icon>
            <span>{{ t('settings.ipWhitelist') }}</span>
          </div>
        </template>
        <el-form label-width="180px" label-position="top">
          <el-form-item :label="t('settings.ipWhitelist')">
            <div class="form-tip">{{ t('settings.ipWhitelistTip') }}</div>
            <ListEditor
              v-model="ipWhitelistList"
              :placeholder="t('settings.ipWhitelistPlaceholder')"
              :add-text="t('settings.addIp')"
            />
          </el-form-item>
        </el-form>
      </el-card>

      <!-- 阿里云短信配置 -->
      <el-card class="settings-card" shadow="never">
        <template #header>
          <div class="card-header">
            <el-icon><ChatLineRound /></el-icon>
            <span>{{ t('settings.sms') }}</span>
            <el-tag v-if="isSmsEnabled" type="success" size="small" style="margin-left:auto">{{ t('common.enabled') }}</el-tag>
            <el-tag v-else type="info" size="small" style="margin-left:auto">{{ t('common.disabled') }}</el-tag>
          </div>
        </template>
        <el-form label-width="180px" label-position="top">
          <div class="form-tip" style="margin-bottom:12px">
            {{ t('settings.smsTip') }}
          </div>

          <el-form-item :label="t('settings.smsKeyId')">
            <el-input
              v-model="form.smsAccessKeyId"
              :placeholder="t('settings.smsKeyIdPlaceholder')"
              style="max-width:500px"
            />
          </el-form-item>

          <el-form-item :label="t('settings.smsKeySecret')">
            <el-input
              v-model="form.smsAccessKeySecret"
              :placeholder="t('settings.smsKeySecretPlaceholder')"
              style="max-width:500px"
              show-password
              type="password"
            />
          </el-form-item>

          <el-form-item :label="t('settings.smsSignName')">
            <el-input
              v-model="form.smsSignName"
              :placeholder="t('settings.smsSignNamePlaceholder')"
              style="max-width:400px"
            />
            <div class="form-tip">{{ t('settings.smsSignNameTip') }}</div>
          </el-form-item>

          <el-form-item :label="t('settings.smsTemplateCode')">
            <el-input
              v-model="form.smsTemplateCode"
              :placeholder="t('settings.smsTemplateCodePlaceholder')"
              style="max-width:400px"
            />
            <div class="form-tip">{{ t('settings.smsTemplateCodeTip') }}</div>
          </el-form-item>
        </el-form>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Setting, Check, Connection, Lock, ChatLineRound } from '@element-plus/icons-vue'
import { useSettingsStore } from '../stores/settings.js'
import { authFetch } from '../api/auth.js'
import ListEditor from './ListEditor.vue'

const { t } = useI18n()

const settingsStore = useSettingsStore()
const loading = ref(false)
const saving = ref(false)

const form = reactive({
  siteEnabled: true,
  siteName: 'ChromaDesk',
  loginLogo: '',
  smallLogo: '',
  favicon: '',
  turnUrls: '',
  turnAuthSecret: '',
  turnCredentialTtl: 86400,
  stunUrls: '',
  apiKey: '',
  allowedOrigins: '',
  smsAccessKeyId: '',
  smsAccessKeySecret: '',
  smsSignName: '',
  smsTemplateCode: '',
  adminIpWhitelist: ''
})

const isSmsEnabled = computed(() =>
  form.smsAccessKeyId && form.smsAccessKeySecret && form.smsSignName && form.smsTemplateCode
)

const turnUrlList = ref([])
const stunUrlList = ref([])
const allowedOriginList = ref([])
const ipWhitelistList = ref([])

function textToList(text) {
  if (!text) return []
  return text.split('\n').map(s => s.trim()).filter(Boolean)
}

function listToText(list) {
  return list.filter(s => s.trim() !== '').join('\n')
}

function syncListsFromForm() {
  turnUrlList.value = textToList(form.turnUrls)
  stunUrlList.value = textToList(form.stunUrls)
  allowedOriginList.value = textToList(form.allowedOrigins)
  ipWhitelistList.value = textToList(form.adminIpWhitelist)
}

function syncFormFromLists() {
  form.turnUrls = listToText(turnUrlList.value)
  form.stunUrls = listToText(stunUrlList.value)
  form.allowedOrigins = listToText(allowedOriginList.value)
  form.adminIpWhitelist = listToText(ipWhitelistList.value)
}

async function loadSettings() {
  loading.value = true
  try {
    const res = await authFetch('/api/v1/admin/settings')
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const data = await res.json()
    Object.assign(form, data)
    if (!form.turnCredentialTtl) form.turnCredentialTtl = 86400
    syncListsFromForm()
  } catch (e) {
    ElMessage.error(t('settings.loadFailed') + ': ' + e.message)
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  syncFormFromLists()
  saving.value = true
  try {
    const res = await authFetch('/api/v1/admin/settings', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(form)
    })
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    settingsStore.updateSettings(form)
    ElMessage.success(t('settings.saved'))
  } catch (e) {
    ElMessage.error(t('settings.saveFailed') + ': ' + e.message)
  } finally {
    saving.value = false
  }
}

onMounted(loadSettings)
</script>

<style scoped>
.settings-page {
  width: 100%;
  padding: 20px;
  box-sizing: border-box;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #303133;
}

.settings-container {
  display: flex;
  flex-direction: column;
  gap: 20px;
  width: 100%;
}

.settings-card {
  border-radius: 8px;
  width: 100%;
  overflow: hidden;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}

.form-tip {
  margin-top: 4px;
  margin-bottom: 8px;
  color: #909399;
  font-size: 13px;
  line-height: 1.5;
}

@media (max-width: 768px) {
  .settings-page { padding: 10px; }
  .page-header { flex-direction: column; align-items: flex-start; gap: 10px; }
}
</style>
