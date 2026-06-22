<template>
  <div class="container">
    <h1 class="page-title">Личный кабинет</h1>

    <div class="profile-card card">
      <div class="profile-header">
        <div class="avatar">{{ initials }}</div>
        <div>
          <h2>{{ authStore.user?.name }}</h2>
          <p>{{ authStore.user?.email }}</p>
          <span v-if="authStore.isAdmin" class="badge badge-warning">Администратор</span>
        </div>
      </div>
    </div>

    <section class="orders-section">
      <h3>История заказов</h3>

      <div v-if="loading" class="skeleton" style="height:120px"></div>

      <div v-else-if="!orders.length" class="empty-orders card">
        <p>У вас пока нет заказов</p>
        <router-link to="/catalog" class="btn btn-primary">Перейти в каталог</router-link>
      </div>

      <div v-else class="orders-list">
        <div v-for="order in orders" :key="order.id" class="order-card card">
          <div class="order-header">
            <span class="order-id">Заказ #{{ order.id }}</span>
            <span :class="['badge', statusClass(order.status)]">{{ statusLabel(order.status) }}</span>
            <span class="order-date">{{ formatDate(order.created_at) }}</span>
          </div>
          <div class="order-items">
            <div v-for="item in order.items" :key="item.id" class="order-item">
              <span>{{ item.product?.name }} × {{ item.quantity }}</span>
              <span>{{ item.store?.name }}</span>
              <span>{{ formatPrice(item.price * item.quantity) }}</span>
            </div>
          </div>
          <div class="order-footer">
            <span>{{ order.address }}</span>
            <strong>{{ formatPrice(order.total) }}</strong>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '../stores'
import { ordersApi, formatPrice } from '../services/api'

const authStore = useAuthStore()
const orders = ref([])
const loading = ref(true)

const initials = computed(() => {
  const name = authStore.user?.name || ''
  return name.split(' ').map(w => w[0]).join('').slice(0, 2).toUpperCase()
})

onMounted(async () => {
  try {
    const { data } = await ordersApi.my()
    orders.value = data
  } finally {
    loading.value = false
  }
})

function formatDate(d) {
  return new Date(d).toLocaleDateString('ru-RU', {
    day: 'numeric', month: 'long', year: 'numeric', hour: '2-digit', minute: '2-digit',
  })
}

function statusLabel(s) {
  const map = {
    pending: 'Ожидает',
    processing: 'Обрабатывается',
    delivering: 'Доставляется',
    delivered: 'Доставлен',
    cancelled: 'Отменён',
  }
  return map[s] || s
}

function statusClass(s) {
  const map = {
    pending: 'badge-warning',
    processing: 'badge-primary',
    delivering: 'badge-primary',
    delivered: 'badge-success',
    cancelled: 'badge-danger',
  }
  return map[s] || 'badge-primary'
}
</script>

<style scoped>
.profile-card { padding: 1.5rem; margin-bottom: 2rem; }
.profile-header { display: flex; align-items: center; gap: 1rem; }
.avatar {
  width: 64px;
  height: 64px;
  background: var(--primary);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.25rem;
  font-weight: 700;
}
.profile-header h2 { font-size: 1.25rem; }
.profile-header p { color: var(--text-muted); font-size: .875rem; }
.orders-section h3 { font-size: 1.125rem; margin-bottom: 1rem; }
.empty-orders { text-align: center; padding: 3rem; }
.empty-orders p { color: var(--text-muted); margin-bottom: 1rem; }
.orders-list { display: flex; flex-direction: column; gap: 1rem; }
.order-card { padding: 1.25rem; }
.order-header {
  display: flex;
  align-items: center;
  gap: .75rem;
  margin-bottom: .75rem;
  flex-wrap: wrap;
}
.order-id { font-weight: 600; }
.order-date { color: var(--text-muted); font-size: .8125rem; margin-left: auto; }
.order-items { border-top: 1px solid var(--border); padding-top: .75rem; }
.order-item {
  display: flex;
  justify-content: space-between;
  font-size: .875rem;
  padding: .25rem 0;
  color: var(--text-muted);
}
.order-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: .75rem;
  padding-top: .75rem;
  border-top: 1px solid var(--border);
  font-size: .875rem;
}
.order-footer strong { font-size: 1.125rem; color: var(--primary); }
</style>
