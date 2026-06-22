<template>
  <aside class="filters">
    <h3 class="filters-title">Фильтры</h3>

    <div class="form-group">
      <label>Категория</label>
      <select v-model="local.category_id" class="form-control" @change="emit">
        <option value="">Все категории</option>
        <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
      </select>
    </div>

    <div class="form-group">
      <label>Магазин</label>
      <select v-model="local.store_id" class="form-control" @change="emit">
        <option value="">Все магазины</option>
        <option v-for="store in stores" :key="store.id" :value="store.id">{{ store.name }}</option>
      </select>
    </div>

    <div class="form-group">
      <label>Бренд</label>
      <input v-model="local.brand" type="text" class="form-control" placeholder="Например, Nestle" @input="debouncedEmit" />
    </div>

    <div class="form-group">
      <label>Цена от</label>
      <input v-model.number="local.min_price" type="number" class="form-control" min="0" @change="emit" />
    </div>

    <div class="form-group">
      <label>Цена до</label>
      <input v-model.number="local.max_price" type="number" class="form-control" min="0" @change="emit" />
    </div>

    <div class="form-group">
      <label class="checkbox-label">
        <input v-model="local.has_discount" type="checkbox" @change="emit" />
        Только со скидкой
      </label>
    </div>

    <div class="form-group">
      <label>Сортировка</label>
      <select v-model="local.sort" class="form-control" @change="emit">
        <option value="">По умолчанию</option>
        <option value="price_asc">Цена ↑</option>
        <option value="price_desc">Цена ↓</option>
        <option value="name">По названию</option>
        <option value="popularity">По популярности</option>
        <option value="rating">По рейтингу магазина</option>
      </select>
    </div>

    <button class="btn btn-outline" style="width:100%" @click="reset">Сбросить</button>
  </aside>
</template>

<script setup>
import { reactive, watch } from 'vue'

const props = defineProps({
  categories: { type: Array, default: () => [] },
  stores: { type: Array, default: () => [] },
  modelValue: { type: Object, default: () => ({}) },
})

const emitUpdate = defineEmits(['update:modelValue'])

const local = reactive({
  category_id: '',
  store_id: '',
  brand: '',
  min_price: '',
  max_price: '',
  has_discount: false,
  sort: '',
})

watch(() => props.modelValue, (v) => Object.assign(local, v), { immediate: true })

let debounceTimer
function debouncedEmit() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(emit, 400)
}

function emit() {
  const filters = {}
  if (local.category_id) filters.category_id = local.category_id
  if (local.store_id) filters.store_id = local.store_id
  if (local.brand) filters.brand = local.brand
  if (local.min_price) filters.min_price = local.min_price
  if (local.max_price) filters.max_price = local.max_price
  if (local.has_discount) filters.has_discount = true
  if (local.sort) filters.sort = local.sort
  emitUpdate('update:modelValue', filters)
}

function reset() {
  Object.assign(local, {
    category_id: '', store_id: '', brand: '',
    min_price: '', max_price: '', has_discount: false, sort: '',
  })
  emit()
}
</script>

<style scoped>
.filters {
  background: var(--surface);
  border-radius: var(--radius);
  padding: 1.25rem;
  box-shadow: var(--shadow);
  height: fit-content;
  position: sticky;
  top: 80px;
}
.filters-title {
  font-size: 1rem;
  font-weight: 600;
  margin-bottom: 1rem;
}
.checkbox-label {
  display: flex;
  align-items: center;
  gap: .5rem;
  font-size: .875rem;
  cursor: pointer;
}
</style>
