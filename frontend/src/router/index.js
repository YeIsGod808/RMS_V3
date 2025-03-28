import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/store/modules/user'
import { message } from 'ant-design-vue'
import NProgress from '@/utils/progress'

const routes = [
  {
    path: '/login', // URL路径
    name: 'Login', // 路由名称
    component: () => import('@/views/login/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/register', // 注册路由
    name: 'Register',
    component: () => import('@/views/login/Register.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    component: () => import('@/components/layout/Layout.vue'),
    children: [
      {
        path: '',
        redirect: '/graph'
      },
      {
        path: 'graph',
        name: 'Graph',
        component: () => import('@/views/graph/Graph.vue'),
        meta: { requiresAuth: true }
      },
      // {
      //   path: 'learning-path',
      //   name: 'LearningPath',
      //   component: () => import('@/views/learning-path/LearningPath.vue'),
      //   meta: { requiresAuth: true }
      // },
      {
        path: 'route',
        name: 'Route',
        component: () => import('@/views/route/Route.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'admin',
        name: 'Admin',
        component: () => import('@/views/admin/Admin.vue'),
        meta: { requiresAuth: true, requiresAdmin: true }
      },
      {
        path: 'external',
        name: 'ExternalLinks',
        component: () => import('@/views/outside/Outside.vue'),
        meta: { requiresAuth: true }
      }
    ]
  },
  {
    path: '/403',
    name: '403',
    component: () => import('@/views/error/403.vue'),
    meta: {
      title: '403 - 禁止访问',
      requiresAuth: false
    }
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/'
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

// 全局前置守卫
router.beforeEach(async (to, from, next) => {
  NProgress.start()
  
  const userStore = useUserStore()
  
  // 判断是否需要登录权限
  if (to.matched.some(record => record.meta.requiresAuth)) {
    // 用户未登录，重定向到登录页
    if (!userStore.isLoggedIn) {
      next({ path: '/login', query: { redirect: to.fullPath } })
      NProgress.done()
      return
    }
    
    // 验证token有效性
    try {
      const isValid = await userStore.verifyToken()
      if (!isValid) {
        next({ path: '/login', query: { redirect: to.fullPath } })
        NProgress.done()
        return
      }
    } catch (error) {
      console.error('Token验证失败:', error)
      next({ path: '/login', query: { redirect: to.fullPath } })
      NProgress.done()
      return
    }
    
    // 检查角色权限
    if (to.matched.some(record => record.meta.requiresAdmin) && !userStore.isAdmin) {
      next({ path: '/403' })
      NProgress.done()
      return
    }
    
    if (to.matched.some(record => record.meta.requiresTeacher) && 
        !(userStore.isTeacher || userStore.isAdmin)) {
      next({ path: '/403' })
      NProgress.done()
      return
    }
  }
  
  // 已登录且要跳转到登录页，则重定向到首页
  if (to.path === '/login' && userStore.isLoggedIn) {
    next({ path: '/' })
    NProgress.done()
    return
  }
  
  next()
})

router.afterEach(() => {
  NProgress.done()
})

export default router