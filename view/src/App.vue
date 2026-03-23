<template>
  <div class="common-layout">
    <el-container class="main-container">
      <el-header class="app-header">
        <div class="header-inner">
          <div class="logo" @click="$router.push('/')">
            <el-icon class="logo-icon" :size="24">
              <EditPen />
            </el-icon>
            <span class="logo-text">GoNote</span>
          </div>
          <div class="nav-actions">
            <el-button v-if="isAuthenticated" text @click="$router.push('/')">我的代码</el-button>
            <el-button v-if="isAuthenticated" text @click="$router.push('/snippets/new')">新建片段</el-button>
            <el-button v-else text @click="$router.push('/auth')">登录 / 注册</el-button>
            <el-button text @click="$router.push('/about')">关于</el-button>
            <template v-if="isAuthenticated">
              <span class="user-name">{{ authStore.user?.username }}</span>
              <el-button text type="danger" @click="handleLogout">退出登录</el-button>
            </template>
            <a href="https://github.com/luckysxx" target="_blank" class="github-link">
              <el-button text>GitHub</el-button>
            </a>
          </div>
        </div>
      </el-header>

      <el-main class="app-main">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>

      <el-footer class="app-footer">
        <div class="footer-content">
          <p>© {{ new Date().getFullYear() }} GoNote. Built with Go & Vue 3.</p>
        </div>
      </el-footer>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { EditPen } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const isAuthenticated = computed(() => authStore.isAuthenticated)

onMounted(() => {
  authStore.initFromStorage()
})

const handleLogout = () => {
  authStore.logout()
  ElMessage.success('已退出登录')
  router.push('/auth')
}
</script>

<style lang="scss">
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&display=swap');

/* =============================================
   Design System: Modern Light
   ============================================= */
:root {
  --bg-primary: #f8fafc;
  --bg-secondary: #ffffff;
  --bg-card: #ffffff;
  --bg-glass: rgba(255, 255, 255, 0.8);
  --bg-subtle: #f1f5f9;

  --text-primary: #0f172a;
  --text-secondary: #475569;
  --text-muted: #94a3b8;

  --accent: #6366f1;
  --accent-light: #818cf8;
  --accent-glow: rgba(99, 102, 241, 0.12);

  --green: #10b981;
  --green-light: #34d399;

  --border: #e2e8f0;
  --border-light: #f1f5f9;

  --shadow-sm: 0 1px 3px rgba(0, 0, 0, 0.04), 0 1px 2px rgba(0, 0, 0, 0.06);
  --shadow-md: 0 4px 16px rgba(0, 0, 0, 0.06);
  --shadow-lg: 0 12px 40px rgba(0, 0, 0, 0.08);

  --radius-sm: 8px;
  --radius-md: 12px;
  --radius-lg: 16px;
  --radius-xl: 20px;
}

* {
  box-sizing: border-box;
}

body {
  margin: 0;
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  background: var(--bg-primary);
  color: var(--text-primary);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

/* =============================================
   Layout
   ============================================= */
.common-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background:
    radial-gradient(ellipse at 20% 0%, rgba(99, 102, 241, 0.06), transparent 50%),
    radial-gradient(ellipse at 80% 100%, rgba(16, 185, 129, 0.04), transparent 50%),
    var(--bg-primary);
}

.main-container {
  min-height: 100vh;
}

.app-header {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(20px) saturate(180%);
  border-bottom: 1px solid var(--border);
  box-shadow: var(--shadow-sm);
  position: sticky;
  top: 0;
  z-index: 100;
  padding: 0;
  height: 64px;

  .header-inner {
    max-width: 1200px;
    margin: 0 auto;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 20px;
  }

  .logo {
    display: flex;
    align-items: center;
    gap: 10px;
    cursor: pointer;
    font-weight: 700;
    font-size: 20px;
    color: var(--accent);
    transition: all 0.2s ease;

    &:hover {
      color: var(--accent-light);
      transform: translateY(-1px);
    }
  }

  .nav-actions {
    display: flex;
    align-items: center;
    gap: 6px;

    .user-name {
      color: var(--text-secondary);
      font-size: 14px;
      margin-left: 4px;
    }

    .github-link {
      text-decoration: none;
    }
  }
}

.app-main {
  padding: 40px 20px;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
  box-sizing: border-box;
}

.app-footer {
  background: transparent;
  border-top: 1px solid var(--border);
  padding: 20px;
  text-align: center;
  color: var(--text-muted);
  font-size: 14px;
  margin-top: auto;
}

/* =============================================
   Transitions
   ============================================= */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
