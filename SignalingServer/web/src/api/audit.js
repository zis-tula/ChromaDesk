import { authFetch } from './auth.js'

const BASE_URL = '/api/v1/admin'

export async function getAuditLogs(params = {}) {
  const query = new URLSearchParams()
  if (params.page) query.set('page', params.page)
  if (params.size) query.set('size', params.size)
  if (params.action) query.set('action', params.action)
  if (params.admin) query.set('admin', params.admin)
  if (params.dateFrom) query.set('dateFrom', params.dateFrom)
  if (params.dateTo) query.set('dateTo', params.dateTo)

  const qs = query.toString()
  const res = await authFetch(`${BASE_URL}/audit-logs${qs ? '?' + qs : ''}`)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}
