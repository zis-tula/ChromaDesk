// User-data sync over WebSocket. Subscribes to the same /user/sync stream
// used by the Qt client so the device list reacts in real time when a host
// goes online/offline or a user logs in/out on a device.

import { userApi } from './userApi'

const DEVICE_EVENTS = new Set([
  'device_online',
  'device_offline',
  'device_logged_in',
  'device_logged_out',
  'device_access_code_changed',
  'device_remark_changed',
])

const FAVORITE_EVENTS = new Set([
  'favorite_added',
  'favorite_updated',
  'favorite_removed',
])

const RECONNECT_DELAY_MS = 5000

class UserSync extends EventTarget {
  constructor() {
    super()
    this._ws = null
    this._reconnectTimer = null
    this._stopped = true
    this._token = ''
  }

  /** Start (or restart) the sync WebSocket using the current token. */
  start() {
    const token = userApi.getToken()
    if (!token) return
    this._token = token
    this._stopped = false
    this._open()
  }

  /** Stop the sync WebSocket and prevent reconnection. */
  stop() {
    this._stopped = true
    clearTimeout(this._reconnectTimer)
    this._reconnectTimer = null
    if (this._ws) {
      this._ws.onclose = null
      try { this._ws.close() } catch { /* noop */ }
      this._ws = null
    }
  }

  // --- internal ------------------------------------------------------------

  _wsUrl() {
    let base = userApi.getServerUrl().replace(/\/+$/, '')
    if (base.startsWith('http://')) base = base.replace(/^http:\/\//, 'ws://')
    else if (base.startsWith('https://')) base = base.replace(/^https:\/\//, 'wss://')
    return `${base}/api/v1/user/sync?token=${encodeURIComponent(this._token)}`
  }

  _open() {
    try {
      this._ws = new WebSocket(this._wsUrl())
    } catch {
      this._scheduleReconnect()
      return
    }

    this._ws.onopen = () => {
      this.dispatchEvent(new CustomEvent('connected'))
    }
    this._ws.onmessage = (e) => this._onMessage(e)
    this._ws.onclose = () => {
      this.dispatchEvent(new CustomEvent('disconnected'))
      this._scheduleReconnect()
    }
    this._ws.onerror = () => { /* onclose will fire next */ }
  }

  _scheduleReconnect() {
    if (this._stopped) return
    clearTimeout(this._reconnectTimer)
    this._reconnectTimer = setTimeout(() => {
      // Token may have been refreshed/cleared while disconnected.
      const token = userApi.getToken()
      if (!token) return
      this._token = token
      this._open()
    }, RECONNECT_DELAY_MS)
  }

  _onMessage(event) {
    let msg
    try { msg = JSON.parse(event.data) } catch { return }
    const type = msg && msg.type
    if (!type) return
    if (DEVICE_EVENTS.has(type)) {
      this.dispatchEvent(new CustomEvent('devices-changed', { detail: msg }))
    } else if (FAVORITE_EVENTS.has(type)) {
      this.dispatchEvent(new CustomEvent('favorites-changed', { detail: msg }))
    }
  }
}

export const userSync = new UserSync()
