import { authFetch } from './auth.js'

const BASE_URL = '/api/v1/admin/user-list'

export async function getUsers(params = {}) {
  const query = new URLSearchParams()
  if (params.page) query.set('page', params.page)
  if (params.size) query.set('size', params.size)
  if (params.sort) query.set('sort', params.sort)
  if (params.order) query.set('order', params.order)
  if (params.search) query.set('search', params.search)
  if (params.level) query.set('level', params.level)
  if (params.status !== undefined && params.status !== '') query.set('status', params.status)
  if (params.channelType) query.set('channelType', params.channelType)

  const qs = query.toString()
  const res = await authFetch(`${BASE_URL}${qs ? '?' + qs : ''}`)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

export async function getUser(id) {
  const res = await authFetch(`${BASE_URL}/${id}`)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

export async function getUserDetail(id) {
  const res = await authFetch(`${BASE_URL}/${id}/details`)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

export async function createUser(data) {
  const res = await authFetch(BASE_URL, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({}))
    throw new Error(err.error || `HTTP ${res.status}`)
  }
  return res.json()
}

export async function updateUser(id, data) {
  const res = await authFetch(`${BASE_URL}/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({}))
    throw new Error(err.error || `HTTP ${res.status}`)
  }
  return res.json()
}

export async function deleteUser(id) {
  const res = await authFetch(`${BASE_URL}/${id}`, {
    method: 'DELETE'
  })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

export async function batchUsers(action, ids, level) {
  const res = await authFetch('/api/v1/admin/user-list/batch', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ action, ids, level })
  })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

export async function updateUserDeviceCount(id, deviceCount) {
  const res = await authFetch(`${BASE_URL}/${id}/device-count`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ deviceCount })
  })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}
