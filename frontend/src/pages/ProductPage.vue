<template>
  <div class="container product-page">
    <div v-if="loading" class="skeleton" style="height:400px"></div>

    <template v-else-if="product">
      <div class="product-layout animate-in">
        <div class="product-gallery card">
          <img :src="product.image" :alt="product.name" />
        </div>

        <div class="product-info">
          <span class="badge badge-primary">{{ product.category?.name }}</span>
          <h1 class="product-title">{{ product.name }}</h1>
          <p v-if="product.brand" class="product-brand">Бренд: {{ product.brand }}</p>
          <p class="product-desc">{{ product.description }}</p>

          <div class="specs">
            <div v-if="product.unit"><strong>Единица:</strong> {{ product.unit }}</div>
            <div v-if="product.weight_volume"><strong>Объём/вес:</strong> {{ product.weight_volume }}</div>
            <div v-if="product.barcode"><strong>Штрихкод:</strong> {{ product.barcode }}</div>
          </div>

          <div class="best-offer">
            <span>Лучшая цена:</span>
            <strong>{{ formatPrice(minPrice(product.prices)) }}</strong>
          </div>
        </div>
      </div>

      <section class="prices-section">
        <h2>Сравнение цен в магазинах</h2>
        <PriceTable :prices="product.prices" @add-to-cart="addToCart" />
      </section>

      <div class="product-actions">
        <router-link to="/checkout" class="btn btn-primary btn-lg">Перейти к оформлению</router-link>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { productsApi, formatPrice, minPrice } from '../services/api'
import { useCartStore } from '../stores'
import PriceTable from '../components/PriceTable.vue'

const route = useRoute()
const router = useRouter()
const cartStore = useCartStore()

const product = ref(null)
const loading = ref(true)

onMounted(async () => {
  try {
    const { data } = await productsApi.get(route.params.id)
    product.value = data
  } finally {
    loading.value = false
  }
})

async function addToCart(price) {
  await cartStore.addItem(product.value.id, price.store_id, 1)
  router.push('/cart')
}
</script>

<style scoped>
.product-layout {
  display: grid;
  grid-template-columns: 400px 1fr;
  gap: 2rem;
  margin-bottom: 2rem;
}
.product-gallery {
  aspect-ratio: 1;
  overflow: hidden;
}
.product-gallery img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.product-title {
  font-size: 1.75rem;
  font-weight: 700;
  margin: .75rem 0;
}
.product-brand {
  color: var(--text-muted);
  margin-bottom: .75rem;
}
.product-desc {
  color: var(--text-muted);
  margin-bottom: 1.25rem;
  line-height: 1.6;
}
.specs {
  display: flex;
  flex-direction: column;
  gap: .375rem;
  font-size: .875rem;
  margin-bottom: 1.25rem;
  padding: 1rem;
  background: var(--bg);
  border-radius: var(--radius-sm);
}
.best-offer {
  display: flex;
  align-items: baseline;
  gap: .75rem;
  font-size: 1.125rem;
}
.best-offer strong {
  font-size: 1.75rem;
  color: var(--primary);
}
.prices-section {
  margin-bottom: 2rem;
}
.prices-section h2 {
  font-size: 1.25rem;
  margin-bottom: 1rem;
}
.product-actions {
  display: flex;
  gap: 1rem;
}
@media (max-width: 768px) {
  .product-layout { grid-template-columns: 1fr; }
}
</style>
