<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">Цены</h1>
      <button class="btn btn-primary" @click="openForm()">+ Добавить</button>
    </div>

    <div class="filters-bar">
      <select v-model="filterProduct" class="form-control" style="max-width:300px" @change="load">
        <option value="">Все товары</option>
        <option v-for="p in products" :key="p.id" :value="p.id">{{ p.name }}</option>
      </select>
    </div>

    <table class="admin-table card">
      <thead>
        <tr><th>ID</th><th>Товар</th><th>Магазин</th><th>Цена</th><th>Скидка</th><th>Действия</th></tr>
      </thead>
      <tbody>
        <tr v-for="item in prices" :key="item.id">
          <td>{{ item.id }}</td>
          <td>{{ getProductName(item.product_id) }}</td>
          <td>{{ getStoreName(item.store_id) }}</td>
          <td>{{ item.price }} ₽</td>
          <td>{{ item.discount_price ? item.discount_price + ' ₽' : '—' }}</td>
          <td class="actions">
            <button class="btn btn-outline btn-sm" @click="openForm(item)">Изменить</button>
            <button class="btn btn-danger btn-sm" @click="remove(item.id)">Удалить</button>
          </td>
        </tr>
      </tbody>
    </table>

    <div v-if="showForm" class="modal-overlay" @click.self="showForm = false">
      <form class="modal card" @submit.prevent="save">
        <h3>{{ editing ? 'Редактировать' : 'Новая цена' }}</h3>
        <div class="form-group">
          <label>Товар</label>
          <select v-model="form.product_id" class="form-control" required>
            <option v-for="p in products" :key="p.id" :value="p.id">{{ p.name }}</option>
          </select>
        </div>
        <div class="form-group">
          <label>Магазин</label>
          <select v-model="form.store_id" class="form-control" required>
            <option v-for="s in stores" :key="s.id" :value="s.id">{{ s.name }}</option>
          </select>
        </div>
        <div class="form-group"><label>Цена</label><input v-model.number="form.price" type="number" step="0.01" class="form-control" required /></div>
        <div class="form-group"><label>Цена со скидкой</label><input v-model.number="form.discount_price" type="number" step="0.01" class="form-control" /></div>
        <div class="form-group"><label>% скидки</label><input v-model.number="form.discount_percent" type="number" step="0.1" class="form-control" /></div>
        <div class="form-group"><label>Популярность</label><input v-model.number="form.popularity" type="number" class="form-control" /></div>
        <div class="modal-actions">
          <button type="button" class="btn btn-outline" @click="showForm = false">Отмена</button>
          <button type="submit" class="btn btn-primary">Сохранить</button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { productsApi, storesApi, adminApi } from '../../services/api'
import api from '../../services/api'

const products = ref([])
const stores = ref([])
const prices = ref([])
const filterProduct = ref('')
const showForm = ref(false)
const editing = ref(null)
const form = reactive({
  product_id: null, store_id: null, price: 0,
  discount_price: null, discount_percent: null, popularity: 0,
})

onMounted(async () => {
  const [prods, strs] = await Promise.all([
    productsApi.list({ per_page: 100 }),
    storesApi.list(),
  ])
  products.value = prods.data.data || []
  stores.value = strs.data
  load()
})

async function load() {
  if (filterProduct.value) {
    const { data } = await api.get('/prices', { params: { product_id: filterProduct.value } })
    prices.value = data
  } else {
    const allPrices = []
    for (const p of products.value.slice(0, 30)) {
      const { data } = await api.get('/prices', { params: { product_id: p.id } })
      allPrices.push(...data)
    }
    prices.value = allPrices
  }
}

function getProductName(id) {
  return products.value.find(p => p.id === id)?.name || id
}

function getStoreName(id) {
  return stores.value.find(s => s.id === id)?.name || id
}

function openForm(item = null) {
  editing.value = item
  Object.assign(form, item || {
    product_id: products.value[0]?.id,
    store_id: stores.value[0]?.id,
    price: 0, discount_price: null, discount_percent: null, popularity: 0,
  })
  showForm.value = true
}

async function save() {
  const payload = { ...form }
  if (!payload.discount_price) payload.discount_price = null
  if (!payload.discount_percent) payload.discount_percent = null
  if (editing.value) {
    await adminApi.prices.update(editing.value.id, payload)
  } else {
    await adminApi.prices.create(payload)
  }
  showForm.value = false
  load()
}

async function remove(id) {
  if (confirm('Удалить цену?')) {
    await adminApi.prices.delete(id)
    load()
  }
}
</script>

<style scoped>
@import './admin.css';
.filters-bar { margin-bottom: 1rem; }
</style>
