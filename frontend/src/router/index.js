import Vue from 'vue'
import Router from 'vue-router'
import Home from '@/views/Home.vue'
import Login from '@/views/Login.vue'
import PersonInfo from '@/views/PersonInfo.vue'

Vue.use(Router)

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/person-info',
    name: 'PersonInfo',
    component: PersonInfo,
    meta: { requiresAuth: true }
  }
]

const router = new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

router.beforeEach((to, from, next) => {
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
  const isAuthenticated = false // 这里应该是你实际的认证逻辑

  if (requiresAuth && !isAuthenticated) {
    next({ name: 'Login' })
  } else {
    next()
  }
})

export default router