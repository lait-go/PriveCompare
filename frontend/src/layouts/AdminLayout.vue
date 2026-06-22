<template>
  <div class="admin-layout">
    <aside class="sidebar">
      <div class="sidebar-brand">
        <router-link to="/">PriceCompare</router-link>
        <span class="admin-badge">Admin</span>
      </div>
      <nav class="sidebar-nav">
        <router-link to="/admin" exact-active-class="active">Дашборд</router-link>
        <router-link to="/admin/categories" active-class="active">Категории</router-link>
        <router-link to="/admin/stores" active-class="active">Магазины</router-link>
        <router-link to="/admin/products" active-class="active">Товары</router-link>
        <router-link to="/admin/prices" active-class="active">Цены</router-link>
        <router-link to="/admin/orders" active-class="active">Заказы</router-link>
      </nav>
      <router-link to="/" class="back-link">← На сайт</router-link>
    </aside>
    <main class="admin-main">
      <router-view />
    </main>
  </div>
</template>

<style scoped>
.admin-layout {
  display: flex;
  min-height: 100vh;
}
.sidebar {
  width: 240px;
  background: #1e293b;
  color: white;
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}
.sidebar-brand {
  margin-bottom: 2rem;
}
.sidebar-brand a {
  color: white;
  font-size: 1.25rem;
  font-weight: 700;
}
.admin-badge {
  display: inline-block;
  background: var(--primary);
  padding: .125rem .5rem;
  border-radius: 4px;
  font-size: .6875rem;
  margin-left: .5rem;
  vertical-align: middle;
}
.sidebar-nav {
  display: flex;
  flex-direction: column;
  gap: .25rem;
  flex: 1;
}
.sidebar-nav a {
  color: #94a3b8;
  padding: .625rem .875rem;
  border-radius: var(--radius-sm);
  font-size: .875rem;
  font-weight: 500;
  transition: all var(--transition);
}
.sidebar-nav a:hover,
.sidebar-nav a.active {
  background: rgba(255,255,255,.1);
  color: white;
}
.back-link {
  color: #94a3b8;
  font-size: .875rem;
  margin-top: auto;
}
.admin-main {
  flex: 1;
  padding: 2rem;
  background: var(--bg);
  overflow: auto;
}
@media (max-width: 768px) {
  .admin-layout { flex-direction: column; }
  .sidebar { width: 100%; }
  .sidebar-nav { flex-direction: row; flex-wrap: wrap; }
}
</style>
