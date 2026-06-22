<template>
  <div class="container checkout-page">
    <h1 class="page-title">Оформление заказа</h1>

    <div v-if="success" class="success-card card animate-in">
      <div class="success-icon">✓</div>
      <h2>Заказ оформлен!</h2>
      <p>{{ successMessage }}</p>
      <div class="success-actions">
        <router-link to="/profile" class="btn btn-primary">Мои заказы</router-link>
        <router-link to="/" class="btn btn-outline">На главную</router-link>
      </div>
    </div>

    <form v-else class="checkout-form" @submit.prevent="submit">
      <div class="form-card card">
        <h3>Данные для доставки</h3>
        <div class="form-group">
          <label>Имя *</label>
          <input v-model="form.name" type="text" class="form-control" required />
        </div>
        <div class="form-group">
          <label>Телефон *</label>
          <input v-model="form.phone" type="tel" class="form-control" required placeholder="+7 (999) 123-45-67" />
        </div>
        <div class="form-group">
          <label>Адрес доставки *</label>
          <textarea v-model="form.address" class="form-control" rows="3" required placeholder="Город, улица, дом, квартира"></textarea>
        </div>
        <div class="form-group">
          <label>Комментарий</label>
          <textarea v-model="form.comment" class="form-control" rows="2" placeholder="Пожелания к заказу"></textarea>
        </div>
      </div>

      <div class="summary-card card">
        <h3>Ваш заказ</h3>
        <div class="summary-row">
          <span>Товаров:</span>
          <span>{{ cartStore.itemCount }}</span>
        </div>
        <div class="summary-row total">
          <span>Итого:</span>
          <strong>{{ formatPrice(cartStore.total) }}</strong>
        </div>
        <p class="payment-note">Оплата не требуется — это демо-режим</p>
        <button type="submit" class="btn btn-primary btn-lg" style="width:100%" :disabled="submitting">
          {{ submitting ? 'Оформляем...' : 'Подтвердить заказ' }}
        </button>
        <p v-if="error" class="alert alert-error" style="margin-top:1rem">{{ error }}</p>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useCartStore, useAuthStore } from '../stores'
import { ordersApi, formatPrice } from '../services/api'

const cartStore = useCartStore()
const authStore = useAuthStore()

const form = reactive({
  name: authStore.user?.name || '',
  phone: '',
  address: '',
  comment: '',
})

const submitting = ref(false)
const success = ref(false)
const successMessage = ref('')
const error = ref('')

onMounted(() => cartStore.fetchCart())

async function submit() {
  submitting.value = true
  error.value = ''
  try {
    const { data } = await ordersApi.create(form)
    success.value = true
    successMessage.value = data.message
    await cartStore.fetchCart()
  } catch (e) {
    error.value = e.response?.data?.error || 'Ошибка оформления заказа'
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.checkout-form {
  display: grid;
  grid-template-columns: 1fr 360px;
  gap: 1.5rem;
  align-items: start;
}
.form-card, .summary-card {
  padding: 1.5rem;
}
.form-card h3, .summary-card h3 {
  font-size: 1.125rem;
  margin-bottom: 1.25rem;
}
.summary-row {
  display: flex;
  justify-content: space-between;
  padding: .5rem 0;
  font-size: .875rem;
}
.summary-row.total {
  border-top: 1px solid var(--border);
  margin: .5rem 0 1rem;
  padding-top: 1rem;
  font-size: 1.125rem;
}
.payment-note {
  font-size: .8125rem;
  color: var(--text-muted);
  margin-bottom: 1rem;
}
.success-card {
  text-align: center;
  padding: 3rem 2rem;
  max-width: 480px;
  margin: 0 auto;
}
.success-icon {
  width: 64px;
  height: 64px;
  background: #dcfce7;
  color: var(--success);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2rem;
  margin: 0 auto 1rem;
}
.success-card h2 { margin-bottom: .5rem; }
.success-card p { color: var(--text-muted); margin-bottom: 1.5rem; }
.success-actions { display: flex; gap: 1rem; justify-content: center; }
@media (max-width: 768px) {
  .checkout-form { grid-template-columns: 1fr; }
}
</style>
