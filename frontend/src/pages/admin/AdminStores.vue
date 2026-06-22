<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">Магазины</h1>
      <button class="btn btn-primary" @click="openForm()">+ Добавить</button>
    </div>

    <table class="admin-table card">
      <thead>
        <tr><th>ID</th><th>Название</th><th>Рейтинг</th><th>Доставка</th><th>Мин. заказ</th><th>Действия</th></tr>
      </thead>
      <tbody>
        <tr v-for="item in items" :key="item.id">
          <td>{{ item.id }}</td>
          <td>{{ item.name }}</td>
          <td>★ {{ item.rating }}</td>
          <td>{{ item.delivery_cost }} ₽</td>
          <td>{{ item.min_order }} ₽</td>
          <td class="actions">
            <button class="btn btn-outline btn-sm" @click="openForm(item)">Изменить</button>
            <button class="btn btn-danger btn-sm" @click="remove(item.id)">Удалить</button>
          </td>
        </tr>
      </tbody>
    </table>

    <div v-if="showForm" class="modal-overlay" @click.self="showForm = false">
      <form class="modal card" @submit.prevent="save">
        <h3>{{ editing ? 'Редактировать' : 'Новый магазин' }}</h3>
        <div class="form-group"><label>Название</label><input v-model="form.name" class="form-control" required /></div>
        <div class="form-group"><label>Логотип URL</label><input v-model="form.logo" class="form-control" /></div>
        <div class="form-group"><label>Рейтинг</label><input v-model.number="form.rating" type="number" step="0.1" class="form-control" /></div>
        <div class="form-group"><label>Описание</label><textarea v-model="form.description" class="form-control" rows="2"></textarea></div>
        <div class="form-group"><label>Стоимость доставки</label><input v-model.number="form.delivery_cost" type="number" class="form-control" /></div>
        <div class="form-group"><label>Мин. сумма заказа</label><input v-model.number="form.min_order" type="number" class="form-control" /></div>
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
import { storesApi, adminApi } from '../../services/api'

const items = ref([])
const showForm = ref(false)
const editing = ref(null)
const form = reactive({ name: '', logo: '', rating: 4.0, description: '', delivery_cost: 0, min_order: 0 })

onMounted(load)

async function load() {
  const { data } = await storesApi.list()
  items.value = data
}

function openForm(item = null) {
  editing.value = item
  Object.assign(form, item || { name: '', logo: '', rating: 4.0, description: '', delivery_cost: 0, min_order: 0 })
  showForm.value = true
}

async function save() {
  if (editing.value) {
    await adminApi.stores.update(editing.value.id, form)
  } else {
    await adminApi.stores.create(form)
  }
  showForm.value = false
  load()
}

async function remove(id) {
  if (confirm('Удалить магазин?')) {
    await adminApi.stores.delete(id)
    load()
  }
}
</script>

<style scoped>
@import './admin.css';
</style>
