<template>
  <div class="container">
    <section class="hero animate-in">
      <h1>Найдите лучшую цену</h1>
      <p>Сравнивайте цены в {{ stores.length }} магазинах и экономьте на каждой покупке</p>
      <form class="hero-search" @submit.prevent="search">
        <input v-model="query" type="search" placeholder="Что ищете? Например, молоко..." class="hero-input" />
        <button type="submit" class="btn btn-primary btn-lg">Найти</button>
      </form>
    </section>

    <section v-if="categories.length" class="section animate-in">
      <h2 class="section-title">Категории</h2>
      <div class="categories-grid">
        <router-link
          v-for="cat in categories"
          :key="cat.id"
          :to="{ name: 'catalog', query: { category_id: cat.id } }"
          class="category-card card"
        >
          <img :src="cat.image" :alt="cat.name" />
          <span>{{ cat.name }}</span>
        </router-link>
      </div>
    </section>

    <section v-if="popularProducts.length" class="section animate-in">
      <div class="section-header">
        <h2 class="section-title">Популярные товары</h2>
        <router-link to="/catalog?sort=popularity">Все →</router-link>
      </div>
      <div class="grid grid-4">
        <ProductCard v-for="p in popularProducts" :key="p.id" :product="p" />
      </div>
    </section>

    <section v-if="bestDeals.length" class="section animate-in">
      <div class="section-header">
        <h2 class="section-title">Лучшие предложения</h2>
        <router-link to="/catalog?has_discount=true">Все →</router-link>
      </div>
      <div class="grid grid-4">
        <ProductCard v-for="p in bestDeals" :key="p.id" :product="p" />
      </div>
    </section>

    <section v-if="recentChanges.length" class="section animate-in">
      <h2 class="section-title">Последние изменения цен</h2>
      <div class="recent-changes card">
        <div v-for="change in recentChanges" :key="change.id" class="change-item">
          <router-link :to="`/product/${change.product_id}`" class="change-product">
            {{ change.product?.name }}
          </router-link>
          <span class="change-store">{{ change.store?.name }}</span>
          <span class="change-price">{{ formatPrice(effectivePrice(change)) }}</span>
        </div>
      </div>
    </section>

    <section v-if="stores.length" class="section animate-in">
      <h2 class="section-title">Магазины</h2>
      <div class="stores-grid">
        <div v-for="store in stores" :key="store.id" class="store-card card">
          <img :src="store.logo" :alt="store.name" class="store-logo" />
          <div>
            <h3>{{ store.name }}</h3>
            <span class="store-rating">★ {{ store.rating }}</span>
            <p>{{ store.description }}</p>
            <span class="store-meta">Доставка от {{ formatPrice(store.delivery_cost) }}</span>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { homeApi, formatPrice, effectivePrice } from '../services/api'
import ProductCard from '../components/ProductCard.vue'

const router = useRouter()
const query = ref('')
const categories = ref([])
const stores = ref([])
const popularProducts = ref([])
const bestDeals = ref([])
const recentChanges = ref([])

onMounted(async () => {
  const { data } = await homeApi.get()
  categories.value = data.categories || []
  stores.value = data.stores || []
  popularProducts.value = data.popular_products || []
  bestDeals.value = data.best_deals || []
  recentChanges.value = data.recent_price_changes || []
})

function search() {
  if (query.value.trim()) {
    router.push({ name: 'catalog', query: { q: query.value.trim() } })
  }
}
</script>

<style scoped>
.hero {
  text-align: center;
  padding: 3rem 1rem;
  background: linear-gradient(135deg, var(--primary-light) 0%, #f0fdf4 100%);
  border-radius: var(--radius);
  margin-bottom: 2.5rem;
}
.hero h1 {
  font-size: 2.25rem;
  font-weight: 800;
  margin-bottom: .5rem;
}
.hero p {
  color: var(--text-muted);
  font-size: 1.125rem;
  margin-bottom: 1.5rem;
}
.hero-search {
  display: flex;
  max-width: 560px;
  margin: 0 auto;
  gap: .75rem;
}
.hero-input {
  flex: 1;
  padding: .875rem 1.25rem;
  border: 2px solid var(--border);
  border-radius: var(--radius-sm);
  font-size: 1rem;
}
.hero-input:focus {
  outline: none;
  border-color: var(--primary);
}
.section { margin-bottom: 2.5rem; }
.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1.25rem;
}
.section-title {
  font-size: 1.375rem;
  font-weight: 700;
}
.categories-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 1rem;
}
.category-card {
  text-align: center;
  padding: 1rem;
  color: var(--text);
  text-decoration: none;
}
.category-card img {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  object-fit: cover;
  margin: 0 auto .75rem;
}
.category-card span {
  font-size: .8125rem;
  font-weight: 600;
}
.recent-changes { padding: .5rem 0; }
.change-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: .875rem 1.25rem;
  border-bottom: 1px solid var(--border);
}
.change-item:last-child { border-bottom: none; }
.change-product {
  flex: 1;
  font-weight: 500;
  color: var(--text);
}
.change-store {
  color: var(--text-muted);
  font-size: .875rem;
}
.change-price {
  font-weight: 700;
  color: var(--primary);
}
.stores-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1rem;
}
.store-card {
  display: flex;
  gap: 1rem;
  padding: 1.25rem;
}
.store-card .store-logo {
  width: 64px;
  height: 64px;
  border-radius: var(--radius-sm);
  object-fit: cover;
  flex-shrink: 0;
}
.store-card h3 { font-size: 1rem; margin-bottom: .25rem; }
.store-rating { color: var(--warning); font-size: .8125rem; }
.store-card p {
  font-size: .8125rem;
  color: var(--text-muted);
  margin: .375rem 0;
}
.store-meta { font-size: .75rem; color: var(--primary); font-weight: 500; }
@media (max-width: 768px) {
  .hero h1 { font-size: 1.75rem; }
  .hero-search { flex-direction: column; }
}
</style>
