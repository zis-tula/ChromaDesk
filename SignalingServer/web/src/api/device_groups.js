import { authFetch } from './auth.js'

const BASE_URL = '/api/v1/admin'

export async function getGroups() {
  const res = await authFetch(`${BASE_URL}/groups`)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

export async function createGroup(data) {
  const res = await authFetch(`${BASE_URL}/groups`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

export async function updateGroup(id, data) {
  const res = await authFetch(`${BASE_URL}/groups/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

export async function deleteGroup(id) {
  const res = await authFetch(`${BASE_URL}/groups/${id}`, { method: 'DELETE' })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

export async function addDevicesToGroup(groupId, deviceIds) {
  const res = await authFetch(`${BASE_URL}/groups/${groupId}/devices`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ device_ids: deviceIds })
  })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

export async function removeDevicesFromGroup(groupId, deviceIds) {
  const res = await authFetch(`${BASE_URL}/groups/${groupId}/devices`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ device_ids: deviceIds })
  })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

export async function getGroupDevices(groupId) {
  const res = await authFetch(`${BASE_URL}/groups/${groupId}/devices`)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}
