<template>
  <footer class="footer">
    <div class="container">
      <div class="footer-grid">
        <div class="footer-brand">
          <h3>🗓 Appointly</h3>
          <p>Simplifying appointment booking for customers and service providers across India.</p>
        </div>
        <div class="footer-col">
          <h4>Quick Links</h4>
          <router-link to="/">Home</router-link>
          <router-link to="/providers">Browse Providers</router-link>
          <router-link to="/register">Sign Up</router-link>
        </div>
        <div class="footer-col">
          <h4>Support</h4>
          <a href="#">Help Center</a>
          <a href="#">Contact Us</a>
          <a href="#">Privacy Policy</a>
        </div>
      </div>
      <div class="footer-bottom">
        <p>&copy; {{ year }} Appointly. All rights reserved.</p>
      </div>
    </div>
  </footer>
</template>

<script>
export default {
  name: 'AppFooter',
  computed: {
    year() {
      return new Date().getFullYear()
    },
  },
}
</script>

<style scoped>
.footer {
  background: var(--color-dark);
  color: var(--color-gray-300);
  padding: 56px 0 0;
  margin-top: auto;
}

.footer-grid {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr;
  gap: 40px;
}

.footer-brand h3 {
  color: var(--color-white);
  font-size: 1.3rem;
  margin-bottom: 12px;
}

.footer-brand p {
  font-size: 0.9rem;
  line-height: 1.6;
  max-width: 300px;
}

.footer-col h4 {
  color: var(--color-white);
  font-size: 0.9rem;
  font-weight: 700;
  margin-bottom: 16px;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.footer-col a {
  display: block;
  font-size: 0.9rem;
  color: var(--color-gray-400);
  padding: 5px 0;
  transition: var(--transition);
}

.footer-col a:hover {
  color: var(--color-primary);
}

.footer-bottom {
  margin-top: 40px;
  padding: 20px 0;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  text-align: center;
  font-size: 0.85rem;
}

@media (max-width: 640px) {
  .footer-grid {
    grid-template-columns: 1fr;
    gap: 28px;
  }
}
</style>
