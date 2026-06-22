<template>
  <div class="layout">
    <AppHeader />
    <main class="main">
      <router-view />
    </main>
    <AppFooter />
  </div>
</template>

<script setup>
import AppHeader from '../components/AppHeader.vue'
import AppFooter from '../components/AppFooter.vue'
import { useCartStore } from '../stores'
import { onMounted } from 'vue'

const cartStore = useCartStore()
onMounted(() => cartStore.fetchCart())
</script>

<style scoped>
.layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}
.main {
  flex: 1;
  padding: 2rem 0;
}
</style>
