const TOKEN_KEY = 'chromadesk_admin_token'

export function getToken() {
  return localStorage.getItem(TOKEN_KEY)
}

export function setToken(token) {
  localStorage.setItem(TOKEN_KEY, token)
}

export function removeToken() {
  localStorage.removeItem(TOKEN_KEY)
}

export async function login(user, password, totp_code) {
  const body = { user, password }
  if (totp_code) body.totp_code = totp_code
  const res = await fetch('/api/v1/admin/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body)
  })
  const data = await res.json()
  if (data.error === '2fa_required') throw new Error('2fa_required')
  if (!res.ok) throw new Error(data.error || 'Login failed')
  setToken(data.token)
  return data
}

export function logout() {
  removeToken()
}

// Authenticated fetch wrapper
export async function authFetch(url, options = {}) {
  const token = getToken()
  const headers = { ...options.headers }
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const res = await fetch(url, { ...options, headers })

  if (res.status === 401) {
    removeToken()
    window.location.hash = '#/login'
    throw new Error('unauthorized')
  }

  return res
}
