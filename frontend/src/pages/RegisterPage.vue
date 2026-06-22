<template>
  <div class="container auth-page">
    <div class="auth-card card animate-in">
      <h1>Регистрация</h1>
      <p class="auth-subtitle">Создайте аккаунт для оформления заказов</p>

      <form @submit.prevent="submit">
        <div class="form-group">
          <label>Имя</label>
          <input v-model="name" type="text" class="form-control" required />
        </div>
        <div class="form-group">
          <label>Email</label>
          <input v-model="email" type="email" class="form-control" required />
        </div>
        <div class="form-group">
          <label>Пароль</label>
          <input v-model="password" type="password" class="form-control" required minlength="6" />
        </div>
        <p v-if="error" class="alert alert-error">{{ error }}</p>
        <button type="submit" class="btn btn-primary btn-lg" style="width:100%" :disabled="loading">
          {{ loading ? 'Регистрация...' : 'Зарегистрироваться' }}
        </button>
      </form>

      <p class="auth-footer">
        Уже есть аккаунт? <router-link to="/login">Войти</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores'

const router = useRouter()
const authStore = useAuthStore()

const name = ref('')
const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

async function submit() {
  loading.value = true
  error.value = ''
  try {
    await authStore.register(email.value, password.value, name.value)
    router.push('/')
  } catch (e) {
    error.value = e.response?.data?.error || 'Ошибка регистрации'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-page {
  display: flex;
  justify-content: center;
  padding: 2rem 0;
}
.auth-card {
  width: 100%;
  max-width: 420px;
  padding: 2rem;
}
.auth-card h1 { font-size: 1.5rem; margin-bottom: .25rem; }
.auth-subtitle { color: var(--text-muted); margin-bottom: 1.5rem; font-size: .875rem; }
.auth-footer { text-align: center; margin-top: 1.5rem; font-size: .875rem; color: var(--text-muted); }
</style>
