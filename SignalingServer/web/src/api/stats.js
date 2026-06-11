import { authFetch } from './auth.js'

const BASE_URL = '/api/v1/admin'

export async function getStats() {
  const res = await authFetch(`${BASE_URL}/stats`)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

export async function getSystemStatus() {
  const res = await authFetch(`${BASE_URL}/system/status`)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

export async function getConnectionStatus() {
  const res = await authFetch(`${BASE_URL}/connections`)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

export async function getTrends(range_ = '24h') {
  const res = await authFetch(`${BASE_URL}/trends?range=${range_}`)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

export async function getActivity(params = {}) {
  const query = new URLSearchParams()
  if (params.page) query.set('page', params.page)
  if (params.size) query.set('size', params.size)
  if (params.deviceId) query.set('deviceId', params.deviceId)
  if (params.status) query.set('status', params.status)
  if (params.dateFrom) query.set('dateFrom', params.dateFrom)
  if (params.dateTo) query.set('dateTo', params.dateTo)

  const qs = query.toString()
  const res = await authFetch(`${BASE_URL}/activity${qs ? '?' + qs : ''}`)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}
