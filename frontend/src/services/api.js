import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api/v1',
  headers: { 'Content-Type': 'application/json' },
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  const sessionId = localStorage.getItem('sessionId')
  if (sessionId) {
    config.headers['X-Session-ID'] = sessionId
  }
  return config
})

export default api

export const authApi = {
  register: (data) => api.post('/auth/register', data),
  login: (data) => api.post('/auth/login', data),
}

export const homeApi = {
  get: () => api.get('/home'),
}

export const productsApi = {
  list: (params) => api.get('/products', { params }),
  get: (id) => api.get(`/products/${id}`),
}

export const categoriesApi = {
  list: () => api.get('/categories'),
  get: (id) => api.get(`/categories/${id}`),
}

export const storesApi = {
  list: () => api.get('/stores'),
}

export const cartApi = {
  get: () => api.get('/cart'),
  add: (data) => api.post('/cart', data),
  update: (id, quantity) => api.put(`/cart/${id}?quantity=${quantity}`),
  remove: (id) => api.delete(`/cart/${id}`),
  clear: () => api.delete('/cart'),
}

export const ordersApi = {
  create: (data) => api.post('/orders', data),
  my: () => api.get('/orders/my'),
}

export const adminApi = {
  categories: {
    create: (data) => api.post('/admin/categories', data),
    update: (id, data) => api.put(`/admin/categories/${id}`, data),
    delete: (id) => api.delete(`/admin/categories/${id}`),
  },
  stores: {
    create: (data) => api.post('/admin/stores', data),
    update: (id, data) => api.put(`/admin/stores/${id}`, data),
    delete: (id) => api.delete(`/admin/stores/${id}`),
  },
  products: {
    create: (data) => api.post('/admin/products', data),
    update: (id, data) => api.put(`/admin/products/${id}`, data),
    delete: (id) => api.delete(`/admin/products/${id}`),
  },
  prices: {
    create: (data) => api.post('/admin/prices', data),
    update: (id, data) => api.put(`/admin/prices/${id}`, data),
    delete: (id) => api.delete(`/admin/prices/${id}`),
  },
  orders: {
    list: () => api.get('/admin/orders'),
    updateStatus: (id, status) => api.patch(`/admin/orders/${id}/status`, { status }),
  },
}

export function effectivePrice(price) {
  return price.discount_price ?? price.price
}

export function minPrice(prices) {
  if (!prices?.length) return 0
  return Math.min(...prices.map(effectivePrice))
}

export function formatPrice(value) {
  return new Intl.NumberFormat('ru-RU', {
    style: 'currency',
    currency: 'RUB',
    maximumFractionDigits: 0,
  }).format(value)
}

export function getSessionId() {
  let id = localStorage.getItem('sessionId')
  if (!id) {
    id = crypto.randomUUID()
    localStorage.setItem('sessionId', id)
  }
  return id
}
