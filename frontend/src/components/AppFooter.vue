<template>
  <footer class="footer">
    <div class="container footer-inner">
      <div class="footer-brand">
        <strong>PriceCompare</strong>
        <p>Сравнивайте цены и экономьте на покупках</p>
      </div>
      <div class="footer-links">
        <router-link to="/catalog">Каталог</router-link>
        <router-link to="/cart">Корзина</router-link>
        <a href="http://localhost:8080/swagger/index.html" target="_blank">API Docs</a>
      </div>
      <p class="footer-copy">&copy; 2026 PriceCompare. Демо-проект.</p>
    </div>
  </footer>
</template>

<style scoped>
.footer {
  background: var(--surface);
  border-top: 1px solid var(--border);
  padding: 2rem 0;
  margin-top: auto;
}
.footer-inner {
  display: flex;
  flex-wrap: wrap;
  gap: 2rem;
  align-items: flex-start;
}
.footer-brand p {
  color: var(--text-muted);
  font-size: .875rem;
  margin-top: .25rem;
}
.footer-links {
  display: flex;
  gap: 1.5rem;
  margin-left: auto;
}
.footer-links a {
  color: var(--text-muted);
  font-size: .875rem;
}
.footer-copy {
  width: 100%;
  color: var(--text-muted);
  font-size: .8125rem;
  margin-top: 1rem;
  padding-top: 1rem;
  border-top: 1px solid var(--border);
}
</style>
