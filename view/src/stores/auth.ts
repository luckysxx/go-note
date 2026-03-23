import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

interface AuthUser {
  id: number
  username: string
  email: string
}

const TOKEN_KEY = 'token'
const REFRESH_TOKEN_KEY = 'refresh_token'
const USER_KEY = 'user'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>('')
  const refreshToken = ref<string>('')
  const user = ref<AuthUser | null>(null)

  const isAuthenticated = computed(() => !!token.value)

  const initFromStorage = () => {
    const savedToken = localStorage.getItem(TOKEN_KEY)
    const savedRefreshToken = localStorage.getItem(REFRESH_TOKEN_KEY)
    const savedUser = localStorage.getItem(USER_KEY)

    token.value = savedToken || ''
    refreshToken.value = savedRefreshToken || ''

    if (savedUser) {
      try {
        user.value = JSON.parse(savedUser) as AuthUser
      } catch {
        user.value = null
        localStorage.removeItem(USER_KEY)
      }
    } else {
      user.value = null
    }
  }

  const setAuth = (nextToken: string, nextRefreshToken: string, nextUser: AuthUser) => {
    token.value = nextToken
    refreshToken.value = nextRefreshToken
    user.value = nextUser
    localStorage.setItem(TOKEN_KEY, nextToken)
    localStorage.setItem(REFRESH_TOKEN_KEY, nextRefreshToken)
    localStorage.setItem(USER_KEY, JSON.stringify(nextUser))
  }

  const updateTokens = (nextToken: string, nextRefreshToken: string) => {
    token.value = nextToken
    refreshToken.value = nextRefreshToken
    localStorage.setItem(TOKEN_KEY, nextToken)
    localStorage.setItem(REFRESH_TOKEN_KEY, nextRefreshToken)
  }

  const logout = () => {
    token.value = ''
    refreshToken.value = ''
    user.value = null
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(REFRESH_TOKEN_KEY)
    localStorage.removeItem(USER_KEY)
  }

  return {
    token,
    refreshToken,
    user,
    isAuthenticated,
    initFromStorage,
    setAuth,
    updateTokens,
    logout,
  }
})
