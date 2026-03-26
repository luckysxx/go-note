<template>
  <div class="auth-page">
    <div class="auth-hero">
      <p class="eyebrow">GoNote · 安全空间</p>
      <h1>登录，开始记录代码</h1>
      <p class="lede">
        集中管理你的代码片段，随时随地查看与分享。
      </p>
    </div>

    <el-card class="auth-card" shadow="always">
      <div class="card-grid">
        <section class="welcome">
          <p class="badge">快速上手</p>
          <h2>保持专注的工作流</h2>
          <ul>
            <li>多语言高亮，代码片段整洁呈现</li>
            <li>可分享链接，协作沟通更高效</li>
            <li>轻量化设计，打开即用</li>
          </ul>
          <div class="cta">立即体验全新的粘贴板体验</div>
        </section>

        <section class="forms">
          <h3 class="form-title">账号登录</h3>

          <el-form :model="loginForm" label-position="top" size="large" class="form-panel">
            <el-form-item label="用户名">
              <el-input v-model="loginForm.account" placeholder="请输入用户名" />
            </el-form-item>
            <el-form-item label="密码">
              <el-input v-model="loginForm.password" type="password" placeholder="••••••••" show-password />
            </el-form-item>
            <el-button type="primary" class="full" size="large" :loading="loginLoading" @click="handleDirectLogin">
              <el-icon><Right /></el-icon>
              <span>登录</span>
            </el-button>
          </el-form>

          <div class="divider">
            <span>或</span>
          </div>

          <el-button plain class="sso-btn" size="large" @click="handleSsoLogin">
            <el-icon><Link /></el-icon>
            <span>前往统一身份登录（SSO）</span>
          </el-button>

          <div class="register-hint">
            <p>还没有账号？</p>
            <el-button plain class="register-btn" @click="goToRegister">
              <el-icon><UserFilled /></el-icon>
              <span>前往注册</span>
            </el-button>
          </div>
        </section>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Right, UserFilled, Link } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { buildSsoLoginUrl } from '@/router'
import request from '@/utils/request'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const loginLoading = ref(false)

const loginForm = reactive({
  account: '',
  password: '',
})

// user-platform 统一注册页面地址
const USER_PLATFORM_URL = import.meta.env.VITE_USER_PLATFORM_URL || 'http://localhost:5173'

// ========== 方式一：直接在 go-note 输入账号密码 ==========
const handleDirectLogin = async () => {
  if (!loginForm.account || !loginForm.password) {
    ElMessage.warning('请填写完整的登录信息')
    return
  }

  // 生成设备指纹
  let deviceId = localStorage.getItem('device_id')
  if (!deviceId) {
    deviceId = crypto.randomUUID()
    localStorage.setItem('device_id', deviceId)
  }

  loginLoading.value = true
  try {
    // 通过网关直接调 user-platform 的 login API
    const res = await request.post<unknown, {
      access_token?: string
      token?: string
      refresh_token?: string
      user_id: number
      username: string
    }>('/api/v1/users/login', {
      username: loginForm.account,
      password: loginForm.password,
      app_code: 'go-note',
      device_id: deviceId,
    })

    const token = res.access_token ?? res.token ?? ''
    authStore.setAuth(token, res.refresh_token ?? '', {
      id: res.user_id,
      username: res.username,
      email: '',
    })

    ElMessage.success('登录成功！')
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/'
    router.push(redirect)
  } catch (error: unknown) {
    const msg = error instanceof Error ? error.message : '登录失败，请检查用户名和密码'
    ElMessage.error(msg)
  } finally {
    loginLoading.value = false
  }
}

// ========== 方式二：跳转 SSO 统一登录 ==========
const handleSsoLogin = () => {
  const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/'
  window.location.href = buildSsoLoginUrl(redirect)
}

const goToRegister = () => {
  window.open(`${USER_PLATFORM_URL}/register?app_code=go-note`, '_blank')
}

onMounted(() => {
  authStore.initFromStorage()
  if (authStore.isAuthenticated) {
    router.replace('/')
  }
})
</script>

<style scoped lang="scss">
.auth-page {
  min-height: calc(100vh - 160px);
  border-radius: var(--radius-xl);
  padding: 32px 20px 48px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.auth-hero {
  max-width: 860px;
  margin: 0 auto;
  text-align: center;
  display: flex;
  flex-direction: column;
  gap: 10px;

  .eyebrow {
    display: inline-flex;
    align-self: center;
    align-items: center;
    gap: 8px;
    padding: 6px 14px;
    background: var(--accent-glow);
    color: var(--accent);
    border-radius: 999px;
    font-weight: 700;
    letter-spacing: 0.2px;
  }

  h1 {
    margin: 0;
    font-size: 2.6rem;
    font-weight: 800;
    color: var(--text-primary);
  }

  .lede {
    margin: 0;
    color: var(--text-secondary);
    font-size: 1.05rem;
  }
}

.auth-card {
  border-radius: var(--radius-xl) !important;
  box-shadow: var(--shadow-lg) !important;
  padding: 6px;
  border: 1px solid var(--border) !important;

  :deep(.el-card__body) {
    padding: 0;
  }
}

.card-grid {
  display: grid;
  grid-template-columns: 1fr 1.2fr;
  gap: 0;
  overflow: hidden;
  border-radius: 14px;
  background: var(--bg-card);
}

.welcome {
  background: linear-gradient(180deg, #4338ca 0%, #6366f1 60%, #4338ca 100%);
  color: #fff;
  padding: 32px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  position: relative;

  &::after {
    content: '';
    position: absolute;
    inset: 0;
    background: radial-gradient(circle at 30% 20%, rgba(255, 255, 255, 0.12), transparent 40%);
    pointer-events: none;
  }

  .badge {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    background: rgba(255, 255, 255, 0.2);
    padding: 6px 12px;
    border-radius: 999px;
    width: fit-content;
    font-weight: 700;
    font-size: 13px;
  }

  h2 {
    margin: 0;
    font-size: 1.8rem;
    line-height: 1.2;
  }

  ul {
    padding-left: 20px;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: 10px;
    color: rgba(255, 255, 255, 0.85);
    line-height: 1.5;
  }

  .cta {
    margin-top: auto;
    padding: 12px 14px;
    border-radius: 12px;
    background: rgba(255, 255, 255, 0.12);
    border: 1px solid rgba(255, 255, 255, 0.2);
    font-weight: 700;
    text-align: center;
  }
}

.forms {
  padding: 28px 32px;
  display: flex;
  flex-direction: column;
  gap: 18px;

  .form-title {
    margin: 0;
    font-size: 1.4rem;
    font-weight: 700;
    color: var(--text-primary);
  }

  .form-panel {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .full {
    width: 100%;
    border-radius: 10px;
    font-weight: 700;
    background: linear-gradient(135deg, var(--accent), #7c3aed);
    border: none;
  }
}

.sso-btn {
  width: 100%;
  border-radius: 10px;
  font-weight: 600;
  color: var(--accent);
  border-color: var(--accent) !important;

  &:hover {
    background: var(--accent-glow) !important;
  }
}

.divider {
  display: flex;
  align-items: center;
  gap: 12px;
  color: var(--text-muted);
  font-size: 0.95rem;

  &::before,
  &::after {
    content: '';
    flex: 1;
    height: 1px;
    background: linear-gradient(90deg, transparent, var(--border) 50%, transparent);
  }
}

.register-hint {
  text-align: center;
  display: flex;
  flex-direction: column;
  gap: 8px;
  align-items: center;

  p {
    margin: 0;
    color: var(--text-muted);
    font-size: 14px;
  }

  .register-btn {
    border-radius: 10px;
    font-weight: 600;
    color: var(--accent);
    border-color: var(--accent) !important;
    padding: 10px 24px;

    &:hover {
      background: var(--accent-glow) !important;
    }

    .el-icon {
      margin-right: 8px;
    }
  }
}

@media (max-width: 960px) {
  .card-grid {
    grid-template-columns: 1fr;
  }
}
</style>
