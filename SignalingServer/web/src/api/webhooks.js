import { authFetch } from './auth.js'

const BASE_URL = '/api/v1/admin'

export async function getWebhooks() {
  const res = await authFetch(`${BASE_URL}/webhooks`)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

export async function createWebhook(data) {
  const res = await authFetch(`${BASE_URL}/webhooks`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

export async function updateWebhook(id, data) {
  const res = await authFetch(`${BASE_URL}/webhooks/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

export async function deleteWebhook(id) {
  const res = await authFetch(`${BASE_URL}/webhooks/${id}`, { method: 'DELETE' })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

export async function testWebhook(id) {
  const res = await authFetch(`${BASE_URL}/webhooks/${id}/test`, { method: 'POST' })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}
