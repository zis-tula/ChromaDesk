import { authFetch } from './auth.js'

const BASE_URL = '/api/v1/admin'

export async function getDevices(params = {}) {
  const query = new URLSearchParams()
  if (params.page) query.set('page', params.page)
  if (params.size) query.set('size', params.size)
  if (params.sort) query.set('sort', params.sort)
  if (params.order) query.set('order', params.order)
  if (params.search) query.set('search', params.search)
  if (params.os) query.set('os', params.os)
  if (params.online !== undefined && params.online !== '') query.set('online', params.online)

  const qs = query.toString()
  const res = await authFetch(`${BASE_URL}/devices${qs ? '?' + qs : ''}`)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

export async function getDeviceDetail(deviceId) {
  const res = await authFetch(`${BASE_URL}/devices/${deviceId}`)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

export async function getDeviceStatus(deviceId) {
  const res = await authFetch(`/api/v1/devices/${deviceId}/status`)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

export async function batchDevices(action, ids, groupId) {
  const res = await authFetch(`${BASE_URL}/devices/batch`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ action, ids, group_id: groupId })
  })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}
