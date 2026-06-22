<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">Категории</h1>
      <button class="btn btn-primary" @click="openForm()">+ Добавить</button>
    </div>

    <table class="admin-table card">
      <thead>
        <tr><th>ID</th><th>Название</th><th>Slug</th><th>Действия</th></tr>
      </thead>
      <tbody>
        <tr v-for="item in items" :key="item.id">
          <td>{{ item.id }}</td>
          <td>{{ item.name }}</td>
          <td>{{ item.slug }}</td>
          <td class="actions">
            <button class="btn btn-outline btn-sm" @click="openForm(item)">Изменить</button>
            <button class="btn btn-danger btn-sm" @click="remove(item.id)">Удалить</button>
          </td>
        </tr>
      </tbody>
    </table>

    <div v-if="showForm" class="modal-overlay" @click.self="showForm = false">
      <form class="modal card" @submit.prevent="save">
        <h3>{{ editing ? 'Редактировать' : 'Новая категория' }}</h3>
        <div class="form-group"><label>Название</label><input v-model="form.name" class="form-control" required /></div>
        <div class="form-group"><label>Slug</label><input v-model="form.slug" class="form-control" required /></div>
        <div class="form-group"><label>Описание</label><textarea v-model="form.description" class="form-control" rows="2"></textarea></div>
        <div class="form-group"><label>Изображение URL</label><input v-model="form.image" class="form-control" /></div>
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
import { categoriesApi, adminApi } from '../../services/api'

const items = ref([])
const showForm = ref(false)
const editing = ref(null)
const form = reactive({ name: '', slug: '', description: '', image: '' })

onMounted(load)

async function load() {
  const { data } = await categoriesApi.list()
  items.value = data
}

function openForm(item = null) {
  editing.value = item
  Object.assign(form, item || { name: '', slug: '', description: '', image: '' })
  showForm.value = true
}

async function save() {
  if (editing.value) {
    await adminApi.categories.update(editing.value.id, form)
  } else {
    await adminApi.categories.create(form)
  }
  showForm.value = false
  load()
}

async function remove(id) {
  if (confirm('Удалить категорию?')) {
    await adminApi.categories.delete(id)
    load()
  }
}
</script>

<style scoped>
@import './admin.css';
</style>
