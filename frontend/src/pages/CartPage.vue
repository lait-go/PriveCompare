<template>
  <div class="container">
    <h1 class="page-title">Корзина</h1>

    <div v-if="cartStore.loading" class="skeleton" style="height:200px"></div>

    <div v-else-if="!items.length" class="empty-cart card">
      <p>Корзина пуста</p>
      <router-link to="/catalog" class="btn btn-primary">Перейти в каталог</router-link>
    </div>

    <template v-else>
      <div class="cart-layout">
        <div class="cart-items card">
          <div v-for="item in items" :key="item.id" class="cart-item">
            <img :src="item.product?.image" :alt="item.product?.name" class="item-image" />
            <div class="item-info">
              <router-link :to="`/product/${item.product_id}`" class="item-name">
                {{ item.product?.name }}
              </router-link>
              <span class="item-store">{{ item.store?.name }}</span>
              <span class="item-price">{{ formatPrice(getItemPrice(item)) }}</span>
            </div>
            <div class="item-quantity">
              <button class="qty-btn" @click="updateQty(item, item.quantity - 1)">−</button>
              <span>{{ item.quantity }}</span>
              <button class="qty-btn" @click="updateQty(item, item.quantity + 1)">+</button>
            </div>
            <div class="item-total">{{ formatPrice(getItemPrice(item) * item.quantity) }}</div>
            <button class="remove-btn" @click="removeItem(item.id)">✕</button>
          </div>
        </div>

        <div class="cart-summary card">
          <h3>Итого</h3>
          <div class="summary-row">
            <span>Товаров:</span>
            <span>{{ cartStore.itemCount }}</span>
          </div>
          <div class="summary-row total">
            <span>Сумма:</span>
            <strong>{{ formatPrice(cartStore.total) }}</strong>
          </div>
          <router-link to="/checkout" class="btn btn-primary btn-lg" style="width:100%;margin-top:1rem">
            Оформить заказ
          </router-link>
          <button class="btn btn-outline" style="width:100%;margin-top:.5rem" @click="clearCart">
            Очистить корзину
          </button>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useCartStore } from '../stores'
import { formatPrice, effectivePrice } from '../services/api'

const cartStore = useCartStore()

const items = computed(() => cartStore.cart?.items || [])

onMounted(() => cartStore.fetchCart())

function getItemPrice(item) {
  const price = item.product?.prices?.find(p => p.store_id === item.store_id)
  if (price) return effectivePrice(price)
  return 0
}

async function updateQty(item, qty) {
  if (qty < 1) return
  await cartStore.updateQuantity(item.id, qty)
}

async function removeItem(id) {
  await cartStore.removeItem(id)
}

async function clearCart() {
  if (confirm('Очистить корзину?')) {
    await cartStore.clearCart()
  }
}
</script>

<style scoped>
.empty-cart {
  text-align: center;
  padding: 4rem 2rem;
}
.empty-cart p {
  color: var(--text-muted);
  margin-bottom: 1.5rem;
  font-size: 1.125rem;
}
.cart-layout {
  display: grid;
  grid-template-columns: 1fr 320px;
  gap: 1.5rem;
  align-items: start;
}
.cart-items { padding: 0; }
.cart-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem 1.25rem;
  border-bottom: 1px solid var(--border);
}
.cart-item:last-child { border-bottom: none; }
.item-image {
  width: 64px;
  height: 64px;
  border-radius: var(--radius-sm);
  object-fit: cover;
}
.item-info { flex: 1; }
.item-name {
  display: block;
  font-weight: 600;
  color: var(--text);
  margin-bottom: .25rem;
}
.item-store {
  font-size: .8125rem;
  color: var(--text-muted);
  display: block;
}
.item-price {
  font-size: .875rem;
  color: var(--primary);
  font-weight: 600;
}
.item-quantity {
  display: flex;
  align-items: center;
  gap: .5rem;
}
.qty-btn {
  width: 32px;
  height: 32px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  background: var(--surface);
  cursor: pointer;
  font-size: 1rem;
}
.qty-btn:hover { border-color: var(--primary); color: var(--primary); }
.item-total {
  font-weight: 700;
  min-width: 80px;
  text-align: right;
}
.remove-btn {
  background: none;
  border: none;
  color: var(--text-muted);
  cursor: pointer;
  font-size: 1rem;
  padding: .25rem;
}
.remove-btn:hover { color: var(--danger); }
.cart-summary {
  padding: 1.5rem;
  position: sticky;
  top: 80px;
}
.cart-summary h3 {
  font-size: 1.125rem;
  margin-bottom: 1rem;
}
.summary-row {
  display: flex;
  justify-content: space-between;
  padding: .5rem 0;
  font-size: .875rem;
}
.summary-row.total {
  border-top: 1px solid var(--border);
  margin-top: .5rem;
  padding-top: 1rem;
  font-size: 1.125rem;
}
@media (max-width: 768px) {
  .cart-layout { grid-template-columns: 1fr; }
  .cart-item { flex-wrap: wrap; }
}
</style>
