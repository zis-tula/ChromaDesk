import { reactive } from 'vue'
import { userApi } from '../api/userApi'

// Simple reactive auth state shared across components
export const authState = reactive({
  isLoggedIn: userApi.isLoggedIn(),
  userInfo: userApi.getUserInfo(),
  smsEnabled: false,
})

export function updateAuthState() {
  authState.isLoggedIn = userApi.isLoggedIn()
  authState.userInfo = userApi.getUserInfo()
}
