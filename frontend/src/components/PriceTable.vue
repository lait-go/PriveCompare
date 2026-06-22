<template>
  <div class="price-table-wrapper">
    <table class="price-table">
      <thead>
        <tr>
          <th>Магазин</th>
          <th>Цена</th>
          <th>Доставка</th>
          <th>Мин. заказ</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="price in sortedPrices"
          :key="price.id"
          :class="{ 'best-price': price.id === bestPriceId }"
        >
          <td>
            <div class="store-cell">
              <img :src="price.store.logo" :alt="price.store.name" class="store-logo" />
              <div>
                <strong>{{ price.store.name }}</strong>
                <span class="store-rating">★ {{ price.store.rating }}</span>
              </div>
            </div>
          </td>
          <td>
            <div class="price-cell">
              <span v-if="price.discount_price" class="old-price">{{ formatPrice(price.price) }}</span>
              <span class="current-price">{{ formatPrice(effectivePrice(price)) }}</span>
              <span v-if="price.discount_percent" class="badge badge-danger">-{{ price.discount_percent }}%</span>
            </div>
          </td>
          <td>{{ formatPrice(price.store.delivery_cost) }}</td>
          <td>{{ formatPrice(price.store.min_order) }}</td>
          <td>
            <button
              class="btn btn-primary btn-sm"
              :disabled="!price.in_stock"
              @click="$emit('add-to-cart', price)"
            >
              {{ price.in_stock ? 'В корзину' : 'Нет в наличии' }}
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { effectivePrice, formatPrice } from '../services/api'

const props = defineProps({
  prices: { type: Array, required: true },
})

defineEmits(['add-to-cart'])

const sortedPrices = computed(() =>
  [...props.prices].sort((a, b) => effectivePrice(a) - effectivePrice(b))
)

const bestPriceId = computed(() => sortedPrices.value[0]?.id)
</script>

<style scoped>
.price-table-wrapper {
  overflow-x: auto;
  border-radius: var(--radius);
  border: 1px solid var(--border);
}
.price-table {
  width: 100%;
  border-collapse: collapse;
  background: var(--surface);
}
.price-table th {
  text-align: left;
  padding: .875rem 1rem;
  font-size: .8125rem;
  font-weight: 600;
  color: var(--text-muted);
  background: var(--bg);
  border-bottom: 1px solid var(--border);
}
.price-table td {
  padding: 1rem;
  border-bottom: 1px solid var(--border);
  vertical-align: middle;
}
.price-table tr:last-child td { border-bottom: none; }
.price-table tr.best-price {
  background: #f0fdf4;
}
.store-cell {
  display: flex;
  align-items: center;
  gap: .75rem;
}
.store-logo {
  width: 40px;
  height: 40px;
  border-radius: var(--radius-sm);
  object-fit: cover;
}
.store-rating {
  display: block;
  font-size: .75rem;
  color: var(--warning);
}
.price-cell {
  display: flex;
  align-items: center;
  gap: .5rem;
  flex-wrap: wrap;
}
.old-price {
  text-decoration: line-through;
  color: var(--text-muted);
  font-size: .875rem;
}
.current-price {
  font-size: 1.125rem;
  font-weight: 700;
  color: var(--primary);
}
</style>
