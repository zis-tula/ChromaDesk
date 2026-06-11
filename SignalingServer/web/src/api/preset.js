import { authFetch } from './auth.js'

const BASE_URL = '/api/v1/admin/preset'

export async function getPreset() {
  const res = await authFetch(BASE_URL)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

export async function updatePreset(data) {
  const res = await authFetch(BASE_URL, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}
