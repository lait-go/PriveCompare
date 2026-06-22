<template>
  <div class="container auth-page">
    <div class="auth-card card animate-in">
      <h1>Вход</h1>
      <p class="auth-subtitle">Войдите в свой аккаунт</p>

      <form @submit.prevent="submit">
        <div class="form-group">
          <label>Email</label>
          <input v-model="email" type="email" class="form-control" required />
        </div>
        <div class="form-group">
          <label>Пароль</label>
          <input v-model="password" type="password" class="form-control" required />
        </div>
        <p v-if="error" class="alert alert-error">{{ error }}</p>
        <button type="submit" class="btn btn-primary btn-lg" style="width:100%" :disabled="loading">
          {{ loading ? 'Вход...' : 'Войти' }}
        </button>
      </form>

      <p class="auth-footer">
        Нет аккаунта? <router-link to="/register">Зарегистрироваться</router-link>
      </p>

      <div class="demo-accounts">
        <p>Демо-аккаунты:</p>
        <button class="btn btn-outline btn-sm" @click="demoLogin('user')">Пользователь</button>
        <button class="btn btn-outline btn-sm" @click="demoLogin('admin')">Администратор</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

async function submit() {
  loading.value = true
  error.value = ''
  try {
    await authStore.login(email.value, password.value)
    router.push(route.query.redirect || '/')
  } catch (e) {
    error.value = e.response?.data?.error || 'Ошибка входа'
  } finally {
    loading.value = false
  }
}

function demoLogin(type) {
  if (type === 'admin') {
    email.value = 'admin@pricecompare.ru'
    password.value = 'admin123'
  } else {
    email.value = 'user@pricecompare.ru'
    password.value = 'user123'
  }
  submit()
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
.auth-card h1 {
  font-size: 1.5rem;
  margin-bottom: .25rem;
}
.auth-subtitle {
  color: var(--text-muted);
  margin-bottom: 1.5rem;
  font-size: .875rem;
}
.auth-footer {
  text-align: center;
  margin-top: 1.5rem;
  font-size: .875rem;
  color: var(--text-muted);
}
.demo-accounts {
  margin-top: 1.5rem;
  padding-top: 1.5rem;
  border-top: 1px solid var(--border);
  text-align: center;
}
.demo-accounts p {
  font-size: .8125rem;
  color: var(--text-muted);
  margin-bottom: .75rem;
}
.demo-accounts .btn { margin: 0 .25rem; }
</style>
