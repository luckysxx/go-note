import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import { useAuthStore } from '@/stores/auth'

// user-platform SSO 登录页地址
const SSO_LOGIN_URL = import.meta.env.VITE_SSO_LOGIN_URL || 'http://localhost:5173/login'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
      meta: { requiresAuth: true },
    },
    {
      path: '/auth',
      name: 'auth',
      component: () => import('../views/AuthView.vue'),
      meta: { guestOnly: true },
    },
    {
      path: '/auth/callback',
      name: 'sso-callback',
      component: () => import('../views/SsoCallbackView.vue'),
    },
    {
      path: '/snippets/new',
      name: 'snippet-new',
      component: () => import('../views/SnippetEditorView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/snippets/:id',
      name: 'snippet-detail',
      component: () => import('../views/PasteView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/snippets/:id/edit',
      name: 'snippet-edit',
      component: () => import('../views/SnippetEditorView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/paste/:id',
      redirect: (to) => ({ path: `/snippets/${to.params.id as string}` }),
    },
    {
      path: '/about',
      name: 'about',
      component: () => import('../views/AboutView.vue'),
    },
  ],
})

/**
 * 构建 SSO 登录跳转 URL
 * 携带 app_code 和 redirect_uri，登录成功后 user-platform 会带 token 跳回 callback
 */
function buildSsoLoginUrl(redirectAfterLogin: string): string {
  const callbackUrl = new URL('/auth/callback', window.location.origin)
  // 把用户原始目标页面存到 state 参数，callback 页面会用它做最终跳转
  callbackUrl.searchParams.set('state', redirectAfterLogin)

  const ssoUrl = new URL(SSO_LOGIN_URL)
  ssoUrl.searchParams.set('app_code', 'go-note')
  ssoUrl.searchParams.set('redirect_uri', callbackUrl.toString())

  return ssoUrl.toString()
}

router.beforeEach((to) => {
  const authStore = useAuthStore()
  if (!authStore.token) {
    authStore.initFromStorage()
  }

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    // 未登录 → 跳转 go-note 登录入口页，由用户主动点击跳 SSO
    return { path: '/auth', query: { redirect: to.fullPath } }
  }

  if (to.meta.guestOnly && authStore.isAuthenticated) {
    return { path: '/' }
  }

  return true
})

export { buildSsoLoginUrl }
export default router
