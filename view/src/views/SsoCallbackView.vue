<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

onMounted(() => {
  // 从 URL hash 中提取 SSO 回调的 token
  // 格式: /auth/callback#access_token=xxx&refresh_token=xxx&token_type=Bearer
  // query params: ?result=logged_in&user_id=1&username=xxx
  const hashStr = window.location.hash.substring(1) // 去掉 #
  const hashParams = new URLSearchParams(hashStr)

  const accessToken = hashParams.get('access_token')
  const refreshToken = hashParams.get('refresh_token')
  const userId = route.query.user_id as string
  const username = route.query.username as string

  if (!accessToken || !userId || !username) {
    // token 缺失，跳回登录
    router.replace('/auth')
    return
  }

  // 存储认证信息
  authStore.setAuth(accessToken, refreshToken ?? '', {
    id: Number(userId),
    username,
    email: '',
  })

  // 跳转到原始目标页面或首页
  const redirect = (route.query.state as string) || '/'
  router.replace(redirect)
})
</script>

<template>
  <div class="callback-page">
    <p>正在登录，请稍候...</p>
  </div>
</template>

<style scoped>
.callback-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  color: var(--text-secondary, #666);
  font-size: 1.1rem;
}
</style>
