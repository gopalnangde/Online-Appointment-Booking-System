import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import Login from '../views/Login.vue'
import Register from '../views/Register.vue'
import Profile from '../views/Profile.vue'
import ProviderDashboard from '../views/ProviderDashboard.vue'
import Providers from '../views/Providers.vue'
import ProviderProfile from '../views/ProviderProfile.vue'
import Appointments from '../views/Appointments.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home,
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { guest: true },
  },
  {
    path: '/register',
    name: 'Register',
    component: Register,
    meta: { guest: true },
  },
  {
    path: '/profile',
    name: 'Profile',
    component: Profile,
    meta: { requiresAuth: true },
  },
  {
    path: '/appointments',
    name: 'Appointments',
    component: Appointments,
    meta: { requiresAuth: true },
  },
  {
    path: '/provider/dashboard',
    name: 'ProviderDashboard',
    component: ProviderDashboard,
    meta: { requiresAuth: true, role: 'ServiceProvider' },
  },
  {
    path: '/providers',
    name: 'Providers',
    component: Providers,
  },
  {
    path: '/providers/:id',
    name: 'ProviderProfileView',
    component: ProviderProfile,
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// Navigation guards
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  const user = JSON.parse(localStorage.getItem('user') || 'null')

  // Redirect to login if route requires auth and no token
  if (to.meta.requiresAuth && !token) {
    return next({ name: 'Login' })
  }

  // Redirect to home if guest-only route and already logged in
  if (to.meta.guest && token) {
    return next({ name: 'Profile' })
  }

  // Check role-based access
  if (to.meta.role && user?.role !== to.meta.role) {
    return next({ name: 'Profile' })
  }

  next()
})

export default router
