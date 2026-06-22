<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">Товары</h1>
      <button class="btn btn-primary" @click="openForm()">+ Добавить</button>
    </div>

    <table class="admin-table card">
      <thead>
        <tr><th>ID</th><th>Название</th><th>Категория</th><th>Бренд</th><th>Действия</th></tr>
      </thead>
      <tbody>
        <tr v-for="item in items" :key="item.id">
          <td>{{ item.id }}</td>
          <td>{{ item.name }}</td>
          <td>{{ item.category?.name }}</td>
          <td>{{ item.brand }}</td>
          <td class="actions">
            <button class="btn btn-outline btn-sm" @click="openForm(item)">Изменить</button>
            <button class="btn btn-danger btn-sm" @click="remove(item.id)">Удалить</button>
          </td>
        </tr>
      </tbody>
    </table>

    <div v-if="showForm" class="modal-overlay" @click.self="showForm = false">
      <form class="modal card" @submit.prevent="save">
        <h3>{{ editing ? 'Редактировать' : 'Новый товар' }}</h3>
        <div class="form-group"><label>Название</label><input v-model="form.name" class="form-control" required /></div>
        <div class="form-group"><label>Описание</label><textarea v-model="form.description" class="form-control" rows="2"></textarea></div>
        <div class="form-group">
          <label>Категория</label>
          <select v-model="form.category_id" class="form-control" required>
            <option v-for="c in categories" :key="c.id" :value="c.id">{{ c.name }}</option>
          </select>
        </div>
        <div class="form-group"><label>Бренд</label><input v-model="form.brand" class="form-control" /></div>
        <div class="form-group"><label>Единица</label><input v-model="form.unit" class="form-control" /></div>
        <div class="form-group"><label>Вес/объём</label><input v-model="form.weight_volume" class="form-control" /></div>
        <div class="form-group"><label>Изображение URL</label><input v-model="form.image" class="form-control" /></div>
        <div class="form-group"><label>Штрихкод</label><input v-model="form.barcode" class="form-control" /></div>
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
import { productsApi, categoriesApi, adminApi } from '../../services/api'

const items = ref([])
const categories = ref([])
const showForm = ref(false)
const editing = ref(null)
const form = reactive({
  name: '', description: '', category_id: null, brand: '', unit: 'шт',
  weight_volume: '', image: '', barcode: '',
})

onMounted(async () => {
  const [prods, cats] = await Promise.all([
    productsApi.list({ per_page: 100 }),
    categoriesApi.list(),
  ])
  items.value = prods.data.data || []
  categories.value = cats.data
})

async function load() {
  const { data } = await productsApi.list({ per_page: 100 })
  items.value = data.data || []
}

function openForm(item = null) {
  editing.value = item
  Object.assign(form, item || {
    name: '', description: '', category_id: categories.value[0]?.id,
    brand: '', unit: 'шт', weight_volume: '', image: '', barcode: '',
  })
  showForm.value = true
}

async function save() {
  if (editing.value) {
    await adminApi.products.update(editing.value.id, form)
  } else {
    await adminApi.products.create(form)
  }
  showForm.value = false
  load()
}

async function remove(id) {
  if (confirm('Удалить товар?')) {
    await adminApi.products.delete(id)
    load()
  }
}
</script>

<style scoped>
@import './admin.css';
</style>
