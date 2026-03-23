import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import { useAuthStore } from '@/stores/auth'

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
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/AboutView.vue'),
    },
  ],
})

router.beforeEach((to) => {
  const authStore = useAuthStore()
  if (!authStore.token) {
    authStore.initFromStorage()
  }

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    return { path: '/auth', query: { redirect: to.fullPath } }
  }

  if (to.meta.guestOnly && authStore.isAuthenticated) {
    return { path: '/' }
  }

  return true
})

export default router
