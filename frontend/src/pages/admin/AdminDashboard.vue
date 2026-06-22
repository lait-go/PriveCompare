<template>
  <div>
    <h1 class="page-title">Панель администратора</h1>
    <div class="grid grid-3">
      <router-link to="/admin/categories" class="dash-card card">
        <h3>Категории</h3>
        <p>Управление категориями товаров</p>
      </router-link>
      <router-link to="/admin/stores" class="dash-card card">
        <h3>Магазины</h3>
        <p>Управление магазинами</p>
      </router-link>
      <router-link to="/admin/products" class="dash-card card">
        <h3>Товары</h3>
        <p>Управление товарами</p>
      </router-link>
      <router-link to="/admin/prices" class="dash-card card">
        <h3>Цены</h3>
        <p>Управление ценами</p>
      </router-link>
      <router-link to="/admin/orders" class="dash-card card">
        <h3>Заказы</h3>
        <p>Просмотр и управление заказами</p>
      </router-link>
      <a href="http://localhost:8080/swagger/index.html" target="_blank" class="dash-card card">
        <h3>API Docs</h3>
        <p>Swagger документация</p>
      </a>
    </div>
  </div>
</template>

<style scoped>
.dash-card {
  padding: 1.5rem;
  color: var(--text);
  text-decoration: none;
  display: block;
}
.dash-card h3 { font-size: 1.125rem; margin-bottom: .375rem; }
.dash-card p { color: var(--text-muted); font-size: .875rem; }
</style>
