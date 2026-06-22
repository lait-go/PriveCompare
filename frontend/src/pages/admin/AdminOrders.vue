<template>
  <div>
    <h1 class="page-title">Заказы</h1>

    <table class="admin-table card">
      <thead>
        <tr>
          <th>ID</th><th>Пользователь</th><th>Имя</th><th>Телефон</th>
          <th>Сумма</th><th>Статус</th><th>Дата</th><th>Действия</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="order in orders" :key="order.id">
          <td>#{{ order.id }}</td>
          <td>{{ order.user?.email }}</td>
          <td>{{ order.name }}</td>
          <td>{{ order.phone }}</td>
          <td>{{ formatPrice(order.total) }}</td>
          <td>
            <select :value="order.status" class="status-select" @change="updateStatus(order.id, $event.target.value)">
              <option value="pending">Ожидает</option>
              <option value="processing">Обрабатывается</option>
              <option value="delivering">Доставляется</option>
              <option value="delivered">Доставлен</option>
              <option value="cancelled">Отменён</option>
            </select>
          </td>
          <td>{{ formatDate(order.created_at) }}</td>
          <td>
            <button class="btn btn-outline btn-sm" @click="expanded = expanded === order.id ? null : order.id">
              {{ expanded === order.id ? 'Скрыть' : 'Детали' }}
            </button>
          </td>
        </tr>
      </tbody>
    </table>

    <div v-if="expandedOrder" class="order-detail card">
      <h3>Заказ #{{ expandedOrder.id }}</h3>
      <p><strong>Адрес:</strong> {{ expandedOrder.address }}</p>
      <p v-if="expandedOrder.comment"><strong>Комментарий:</strong> {{ expandedOrder.comment }}</p>
      <ul>
        <li v-for="item in expandedOrder.items" :key="item.id">
          {{ item.product?.name }} × {{ item.quantity }} — {{ item.store?.name }} — {{ formatPrice(item.price) }}
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { adminApi, formatPrice } from '../../services/api'

const orders = ref([])
const expanded = ref(null)

const expandedOrder = computed(() => orders.value.find(o => o.id === expanded.value))

onMounted(load)

async function load() {
  const { data } = await adminApi.orders.list()
  orders.value = data
}

async function updateStatus(id, status) {
  await adminApi.orders.updateStatus(id, status)
  load()
}

function formatDate(d) {
  return new Date(d).toLocaleDateString('ru-RU', {
    day: 'numeric', month: 'short', year: 'numeric',
  })
}
</script>

<style scoped>
@import './admin.css';
.status-select {
  padding: .25rem .5rem;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  font-size: .8125rem;
}
.order-detail {
  margin-top: 1rem;
  padding: 1.25rem;
}
.order-detail ul {
  margin-top: .75rem;
  padding-left: 1.25rem;
  font-size: .875rem;
}
</style>
