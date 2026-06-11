<template>
  <div>
    <h2 class="page-title">{{ $t('devices.title') }}</h2>

    <div v-if="!authState.isLoggedIn" class="card">
      <div class="empty-state">
        <div class="empty-icon">🔒</div>
        <p>{{ $t('devices.loginRequired') }}</p>
        <button class="btn btn-primary btn-sm" style="margin-top:12px;" @click="showLogin()">{{ $t('user.login') }}</button>
      </div>
    </div>

    <template v-else>
      <div style="display:flex;justify-content:flex-end;margin-bottom:12px;">
        <button class="btn btn-secondary btn-sm" @click="refresh">🔄 {{ $t('devices.refresh') }}</button>
      </div>

      <!-- My Devices -->
      <div class="card">
        <div class="card-title">{{ $t('devices.myDevices') }}</div>
        <div v-if="myDevices.length === 0" class="empty-state">
          <div class="empty-icon">📱</div>
          <p>{{ $t('devices.noDevices') }}</p>
        </div>
        <div v-for="d in myDevices" :key="d.device_id" class="device-item">
          <div :class="['device-status', (d.online && d.logged_in) ? 'online' : 'offline']"></div>
          <div class="device-info">
            <div class="device-name">{{ d.remark || d.device_name || d.device_id }}</div>
            <div class="device-id">{{ d.device_id }}</div>
          </div>
          <div class="device-actions">
            <button v-if="d.online && d.logged_in" class="btn btn-primary btn-sm" @click="connectDevice(d.device_id, d.access_code)">{{ $t('devices.connect') }}</button>
            <button class="icon-btn" :title="$t('devices.setRemark')" @click="setRemark(d)">✏️</button>
          </div>
        </div>
      </div>

      <!-- My Favorites -->
      <div class="card">
        <div class="card-title">{{ $t('devices.myFavorites') }}</div>
        <div v-if="myFavorites.length === 0" class="empty-state">
          <div class="empty-icon">⭐</div>
          <p>{{ $t('devices.noFavorites') }}</p>
        </div>
        <div v-for="f in myFavorites" :key="f.device_id" class="device-item">
          <span style="font-size:16px;flex-shrink:0;">⭐</span>
          <div class="device-info">
            <div class="device-name">{{ f.device_name || f.device_id }}</div>
            <div class="device-id">{{ f.device_id }}</div>
          </div>
          <div class="device-actions">
            <button class="btn btn-primary btn-sm" @click="connectDevice(f.device_id, f.access_password)">{{ $t('devices.connect') }}</button>
            <button class="icon-btn danger" :title="$t('devices.removeFavorite')" @click="removeFav(f.device_id)">✕</button>
          </div>
        </div>
      </div>

      <!-- Connection Logs -->
      <div class="card">
        <div class="card-title">{{ $t('devices.connectionLogs') }}</div>
        <div v-if="logs.length === 0" class="empty-state">
          <div class="empty-icon">📋</div>
          <p>{{ $t('devices.noLogs') }}</p>
        </div>
        <div v-for="log in logs.slice(0, 20)" :key="log.id" class="device-item">
          <div :class="['device-status', log.status === 'success' ? 'online' : 'offline']"></div>
          <div class="device-info">
            <div class="device-name">{{ log.device_id }}</div>
            <div class="device-id">{{ formatTime(log.created_at) }}{{ log.duration ? ` · ${Math.floor(log.duration/60)}m${log.duration%60}s` : '' }}</div>
          </div>
        </div>
      </div>
    </template>

    <!-- Remark Dialog -->
    <div v-if="remarkDialog.show" class="dialog-overlay" @click.self="remarkDialog.show = false">
      <div class="dialog">
        <div class="card-title">{{ $t('devices.setRemark') }}</div>
        <div class="form-group">
          <input v-model="remarkDialog.value" class="form-input" type="text" @keyup.enter="saveRemark" />
        </div>
        <div style="display:flex;gap:8px;">
          <button class="btn btn-primary" style="flex:1;" @click="saveRemark">{{ $t('account.save') }}</button>
          <button class="btn btn-secondary" style="flex:1;" @click="remarkDialog.show = false">{{ $t('user.cancel') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, inject, onMounted, onBeforeUnmount } from 'vue'
import { useI18n } from 'vue-i18n'
import { userApi } from '../api/userApi'
import { userSync } from '../api/userSync'

const { t } = useI18n()
const showToast = inject('showToast')
const showLogin = inject('showLogin')
const authState = inject('authState')

const myDevices = ref([])
const myFavorites = ref([])
const logs = ref([])
const remarkDialog = ref({ show: false, device: null, value: '' })

async function refresh() {
  await Promise.all([loadDevices(), loadFavorites(), loadLogs()])
}

async function loadDevices() {
  const r = await userApi.fetchMyDevices()
  if (r.ok && r.data) myDevices.value = r.data.devices || []
}

async function loadFavorites() {
  const r = await userApi.fetchFavorites()
  if (r.ok && r.data) myFavorites.value = r.data.favorites || []
}

async function loadLogs() {
  const r = await userApi.fetchConnectionLogs()
  if (r.ok && r.data) logs.value = r.data.logs || []
}

function connectDevice(deviceId, accessCode) {
  if (!accessCode) { showToast(t('devices.noAccessCode'), 'error'); return }
  const serverUrl = localStorage.getItem('chromadesk_signaling_url') || 'ws://qdsignaling.quickcoder.cc:8000'
  const codec = localStorage.getItem('chromadesk_video_codec') || 'AV1'
  const params = new URLSearchParams({ server: serverUrl, device: deviceId, code: accessCode, codec })
  window.open(`remote.html?${params}`, `chromadesk_${deviceId}`)
}

function setRemark(device) {
  remarkDialog.value = { show: true, device, value: device.remark || device.device_name || '' }
}

async function saveRemark() {
  const { device, value } = remarkDialog.value
  await userApi.setDeviceRemark(device.device_id, value)
  remarkDialog.value.show = false
  await loadDevices()
}

async function removeFav(deviceId) {
  await userApi.removeFavorite(deviceId)
  await loadFavorites()
}

function formatTime(ts) {
  if (!ts) return ''
  return new Date(ts).toLocaleString()
}

onMounted(() => {
  if (authState.isLoggedIn) refresh()
  // Refresh the device list whenever the server pushes a relevant change
  // (device online/offline, logged in/out, access code / remark updates).
  userSync.addEventListener('devices-changed', onDevicesChanged)
  userSync.addEventListener('favorites-changed', onFavoritesChanged)
})

onBeforeUnmount(() => {
  userSync.removeEventListener('devices-changed', onDevicesChanged)
  userSync.removeEventListener('favorites-changed', onFavoritesChanged)
})

function onDevicesChanged() {
  if (authState.isLoggedIn) loadDevices()
}

function onFavoritesChanged() {
  if (authState.isLoggedIn) loadFavorites()
}
</script>

<style scoped>
.device-item {
  display: flex; align-items: center; gap: 12px;
  padding: 12px; border-radius: var(--radius);
  transition: background 0.15s;
}
.device-item:hover { background: rgba(255,255,255,0.05); }
.device-status { width: 10px; height: 10px; border-radius: 50%; flex-shrink: 0; }
.device-status.online { background: var(--success); }
.device-status.offline { background: var(--text-disabled); }
.device-info { flex: 1; min-width: 0; }
.device-name { font-size: 14px; font-weight: 600; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.device-id { font-size: 12px; color: var(--text-secondary); font-family: 'Consolas', monospace; margin-top: 2px; }
.device-actions { display: flex; gap: 4px; flex-shrink: 0; }
.icon-btn { width: 28px; height: 28px; border: none; border-radius: 6px; background: transparent; color: var(--text-secondary); cursor: pointer; display: flex; align-items: center; justify-content: center; font-size: 14px; transition: all 0.15s; }
.icon-btn:hover { background: rgba(255,255,255,0.1); color: var(--text-primary); }
.icon-btn.danger:hover { background: rgba(244,67,54,0.2); color: var(--error); }

.dialog-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.5); z-index: 1000; display: flex; align-items: center; justify-content: center; }
.dialog { background: var(--bg-card); border: 1px solid var(--border); border-radius: 12px; padding: 24px; width: 320px; max-width: 90vw; }
</style>
