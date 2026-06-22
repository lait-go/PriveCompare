<template>
  <header class="header">
    <div class="container header-inner">
      <router-link to="/" class="logo">
        <svg width="28" height="28" viewBox="0 0 32 32" fill="none">
          <rect width="32" height="32" rx="8" fill="#2563eb"/>
          <path d="M8 12h16M8 16h12M8 20h8" stroke="white" stroke-width="2" stroke-linecap="round"/>
        </svg>
        PriceCompare
      </router-link>

      <form class="search-form" @submit.prevent="search">
        <input
          v-model="query"
          type="search"
          placeholder="Поиск товаров..."
          class="search-input"
        />
        <button type="submit" class="search-btn">Найти</button>
      </form>

      <nav class="nav">
        <router-link to="/catalog">Каталог</router-link>
        <router-link to="/cart" class="cart-link">
          Корзина
          <span v-if="cartStore.itemCount" class="cart-badge">{{ cartStore.itemCount }}</span>
        </router-link>
        <template v-if="authStore.isAuthenticated">
          <router-link to="/profile">Профиль</router-link>
          <router-link v-if="authStore.isAdmin" to="/admin" class="admin-link">Админ</router-link>
          <button class="btn btn-outline btn-sm" @click="logout">Выйти</button>
        </template>
        <template v-else>
          <router-link to="/login" class="btn btn-outline btn-sm">Войти</router-link>
        </template>
      </nav>
    </div>
  </header>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore, useCartStore } from '../stores'

const router = useRouter()
const authStore = useAuthStore()
const cartStore = useCartStore()
const query = ref('')

function search() {
  if (query.value.trim()) {
    router.push({ name: 'catalog', query: { q: query.value.trim() } })
  }
}

function logout() {
  authStore.logout()
  router.push('/')
}
</script>

<style scoped>
.header {
  background: var(--surface);
  border-bottom: 1px solid var(--border);
  position: sticky;
  top: 0;
  z-index: 100;
}
.header-inner {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  height: 64px;
}
.logo {
  display: flex;
  align-items: center;
  gap: .5rem;
  font-weight: 700;
  font-size: 1.125rem;
  color: var(--text);
  flex-shrink: 0;
}
.search-form {
  flex: 1;
  display: flex;
  max-width: 520px;
}
.search-input {
  flex: 1;
  padding: .5rem 1rem;
  border: 1.5px solid var(--border);
  border-right: none;
  border-radius: var(--radius-sm) 0 0 var(--radius-sm);
  font-size: .875rem;
}
.search-input:focus {
  outline: none;
  border-color: var(--primary);
}
.search-btn {
  padding: .5rem 1.25rem;
  background: var(--primary);
  color: white;
  border: none;
  border-radius: 0 var(--radius-sm) var(--radius-sm) 0;
  font-weight: 600;
  font-size: .875rem;
  cursor: pointer;
}
.search-btn:hover { background: var(--primary-dark); }
.nav {
  display: flex;
  align-items: center;
  gap: 1rem;
  flex-shrink: 0;
}
.nav a {
  color: var(--text-muted);
  font-size: .875rem;
  font-weight: 500;
  transition: color var(--transition);
}
.nav a:hover, .nav a.router-link-active {
  color: var(--primary);
}
.cart-link {
  position: relative;
}
.cart-badge {
  position: absolute;
  top: -8px;
  right: -12px;
  background: var(--danger);
  color: white;
  font-size: .625rem;
  font-weight: 700;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}
.admin-link { color: var(--warning) !important; }
@media (max-width: 768px) {
  .header-inner { flex-wrap: wrap; height: auto; padding: .75rem 1rem; }
  .search-form { order: 3; max-width: 100%; width: 100%; }
  .nav { gap: .5rem; }
}
</style>
