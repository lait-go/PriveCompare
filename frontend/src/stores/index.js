import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi, cartApi, getSessionId } from '../services/api'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))
  const token = ref(localStorage.getItem('token') || null)

  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  getSessionId()

  async function register(email, password, name) {
    const { data } = await authApi.register({ email, password, name })
    setAuth(data)
    return data
  }

  async function login(email, password) {
    const { data } = await authApi.login({ email, password })
    setAuth(data)
    return data
  }

  function setAuth(data) {
    token.value = data.token
    user.value = data.user
    localStorage.setItem('token', data.token)
    localStorage.setItem('user', JSON.stringify(data.user))
  }

  function logout() {
    token.value = null
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  return { user, token, isAuthenticated, isAdmin, register, login, logout }
})

export const useCartStore = defineStore('cart', () => {
  const cart = ref(null)
  const loading = ref(false)

  const itemCount = computed(() =>
    cart.value?.items?.reduce((sum, i) => sum + i.quantity, 0) || 0
  )

  const total = computed(() => {
    if (!cart.value?.items) return 0
    return cart.value.items.reduce((sum, item) => {
      const price = item.product?.prices?.find(p => p.store_id === item.store_id)
      const effective = price?.discount_price ?? price?.price ?? 0
      return sum + effective * item.quantity
    }, 0)
  })

  async function fetchCart() {
    loading.value = true
    try {
      const { data } = await cartApi.get()
      cart.value = data
    } finally {
      loading.value = false
    }
  }

  async function addItem(productId, storeId, quantity = 1) {
    const { data } = await cartApi.add({ product_id: productId, store_id: storeId, quantity })
    cart.value = data
  }

  async function updateQuantity(itemId, quantity) {
    const { data } = await cartApi.update(itemId, quantity)
    cart.value = data
  }

  async function removeItem(itemId) {
    const { data } = await cartApi.remove(itemId)
    cart.value = data
  }

  async function clearCart() {
    await cartApi.clear()
    cart.value = { items: [] }
  }

  return { cart, loading, itemCount, total, fetchCart, addItem, updateQuantity, removeItem, clearCart }
})
