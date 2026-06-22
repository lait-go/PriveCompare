import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores'

const routes = [
  {
    path: '/',
    component: () => import('../layouts/MainLayout.vue'),
    children: [
      { path: '', name: 'home', component: () => import('../pages/HomePage.vue') },
      { path: 'catalog', name: 'catalog', component: () => import('../pages/CatalogPage.vue') },
      { path: 'product/:id', name: 'product', component: () => import('../pages/ProductPage.vue') },
      { path: 'cart', name: 'cart', component: () => import('../pages/CartPage.vue') },
      { path: 'checkout', name: 'checkout', component: () => import('../pages/CheckoutPage.vue'), meta: { auth: true } },
      { path: 'profile', name: 'profile', component: () => import('../pages/ProfilePage.vue'), meta: { auth: true } },
      { path: 'login', name: 'login', component: () => import('../pages/LoginPage.vue') },
      { path: 'register', name: 'register', component: () => import('../pages/RegisterPage.vue') },
    ],
  },
  {
    path: '/admin',
    component: () => import('../layouts/AdminLayout.vue'),
    meta: { auth: true, admin: true },
    children: [
      { path: '', name: 'admin', component: () => import('../pages/admin/AdminDashboard.vue') },
      { path: 'categories', name: 'admin-categories', component: () => import('../pages/admin/AdminCategories.vue') },
      { path: 'stores', name: 'admin-stores', component: () => import('../pages/admin/AdminStores.vue') },
      { path: 'products', name: 'admin-products', component: () => import('../pages/admin/AdminProducts.vue') },
      { path: 'prices', name: 'admin-prices', component: () => import('../pages/admin/AdminPrices.vue') },
      { path: 'orders', name: 'admin-orders', component: () => import('../pages/admin/AdminOrders.vue') },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior: () => ({ top: 0 }),
})

router.beforeEach((to, from, next) => {
  const auth = useAuthStore()
  if (to.meta.auth && !auth.isAuthenticated) {
    next({ name: 'login', query: { redirect: to.fullPath } })
  } else if (to.meta.admin && !auth.isAdmin) {
    next({ name: 'home' })
  } else {
    next()
  }
})

export default router
