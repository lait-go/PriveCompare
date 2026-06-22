<template>
  <router-link :to="`/product/${product.id}`" class="product-card card">
    <div class="product-image">
      <img :src="imageSrc" :alt="product.name" loading="lazy" @error="onImageError" />
      <span v-if="hasDiscount" class="discount-badge">-{{ maxDiscount }}%</span>
    </div>
    <div class="product-body">
      <span v-if="product.category" class="product-category">{{ product.category.name }}</span>
      <h3 class="product-name">{{ product.name }}</h3>
      <p v-if="product.brand" class="product-brand">{{ product.brand }}</p>
      <div class="product-price">
        <span class="price-from">от {{ formatPrice(min) }}</span>
        <span v-if="storeCount > 1" class="store-count">{{ storeCount }} магазинов</span>
      </div>
    </div>
  </router-link>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { minPrice, effectivePrice, formatPrice } from '../services/api'

const FALLBACK_IMAGE = '/images/products/milk.jpg'

const props = defineProps({
  product: { type: Object, required: true },
})

const imageSrc = ref(props.product.image)

watch(() => props.product.image, (url) => { imageSrc.value = url })

function onImageError() {
  imageSrc.value = FALLBACK_IMAGE
}

const min = computed(() => minPrice(props.product.prices))
const storeCount = computed(() => props.product.prices?.length || 0)

const hasDiscount = computed(() =>
  props.product.prices?.some(p => p.discount_price != null)
)

const maxDiscount = computed(() => {
  const discounts = props.product.prices
    ?.filter(p => p.discount_percent)
    .map(p => p.discount_percent) || []
  return discounts.length ? Math.max(...discounts).toFixed(0) : 0
})
</script>

<style scoped>
.product-card {
  display: block;
  color: inherit;
  text-decoration: none;
}
.product-image {
  position: relative;
  aspect-ratio: 1;
  overflow: hidden;
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
}
.product-image img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  padding: .5rem;
  transition: transform .3s ease;
}
.product-card:hover .product-image img {
  transform: scale(1.05);
}
.discount-badge {
  position: absolute;
  top: .75rem;
  left: .75rem;
  background: var(--danger);
  color: white;
  padding: .25rem .5rem;
  border-radius: var(--radius-sm);
  font-size: .75rem;
  font-weight: 700;
}
.product-body {
  padding: 1rem;
}
.product-category {
  font-size: .75rem;
  color: var(--primary);
  font-weight: 500;
}
.product-name {
  font-size: .9375rem;
  font-weight: 600;
  margin: .25rem 0;
  line-height: 1.3;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.product-brand {
  font-size: .8125rem;
  color: var(--text-muted);
}
.product-price {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: .75rem;
}
.price-from {
  font-size: 1.125rem;
  font-weight: 700;
  color: var(--primary);
}
.store-count {
  font-size: .75rem;
  color: var(--text-muted);
}
</style>
