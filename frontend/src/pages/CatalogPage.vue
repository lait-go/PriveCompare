<template>
  <div class="container catalog-page">
    <h1 class="page-title">Каталог</h1>

    <div class="catalog-layout">
      <ProductFilters
        v-model="filters"
        :categories="categories"
        :stores="stores"
        @update:model-value="loadProducts"
      />

      <div class="catalog-content">
        <div v-if="route.query.q" class="search-info">
          Результаты поиска: «{{ route.query.q }}» — {{ total }} товаров
        </div>

        <div v-if="loading" class="grid grid-4">
          <div v-for="i in 8" :key="i" class="skeleton" style="height:280px"></div>
        </div>

        <div v-else-if="products.length" class="grid grid-4">
          <ProductCard v-for="p in products" :key="p.id" :product="p" />
        </div>

        <div v-else class="empty-state">
          <p>Товары не найдены. Попробуйте изменить фильтры.</p>
        </div>

        <div v-if="totalPages > 1" class="pagination">
          <button
            class="btn btn-outline btn-sm"
            :disabled="page <= 1"
            @click="changePage(page - 1)"
          >← Назад</button>
          <span>{{ page }} / {{ totalPages }}</span>
          <button
            class="btn btn-outline btn-sm"
            :disabled="page >= totalPages"
            @click="changePage(page + 1)"
          >Вперёд →</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { productsApi, categoriesApi, storesApi } from '../services/api'
import ProductCard from '../components/ProductCard.vue'
import ProductFilters from '../components/ProductFilters.vue'

const route = useRoute()
const router = useRouter()

const products = ref([])
const categories = ref([])
const stores = ref([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const perPage = 20
const filters = ref({})

const totalPages = computed(() => Math.ceil(total.value / perPage))

async function loadProducts() {
  loading.value = true
  try {
    const params = {
      ...filters.value,
      q: route.query.q || undefined,
      page: page.value,
      per_page: perPage,
    }
    const { data } = await productsApi.list(params)
    products.value = data.data || []
    total.value = data.total || 0
  } finally {
    loading.value = false
  }
}

function changePage(p) {
  page.value = p
  loadProducts()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

watch(() => route.query, () => {
  page.value = 1
  if (route.query.category_id) filters.value.category_id = route.query.category_id
  if (route.query.has_discount) filters.value.has_discount = true
  if (route.query.sort) filters.value.sort = route.query.sort
  loadProducts()
}, { immediate: true })

onMounted(async () => {
  const [cats, strs] = await Promise.all([categoriesApi.list(), storesApi.list()])
  categories.value = cats.data
  stores.value = strs.data
})
</script>

<style scoped>
.catalog-layout {
  display: grid;
  grid-template-columns: 260px 1fr;
  gap: 1.5rem;
  align-items: start;
}
.search-info {
  background: var(--primary-light);
  padding: .75rem 1rem;
  border-radius: var(--radius-sm);
  margin-bottom: 1rem;
  font-size: .875rem;
}
.empty-state {
  text-align: center;
  padding: 4rem 2rem;
  color: var(--text-muted);
}
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  margin-top: 2rem;
}
@media (max-width: 768px) {
  .catalog-layout { grid-template-columns: 1fr; }
}
</style>
