<template>
  <div class="split-auth-container">
    <div class="auth-overlay"></div>

    <div class="container split-auth-inner">
      <!-- Left Side: Streamlined Copy -->
      <div class="auth-left-content">
        <div class="auth-brand-badge">
          <i class="pi pi-calendar-plus brand-icon"></i>
          <span>Appointly</span>
        </div>

        <h1 class="auth-hero-title">
          Create <span class="highlight">Account</span>
        </h1>

        <p class="auth-hero-desc">
          Register in under 60 seconds as a client or service provider.
        </p>

        <div class="feature-bullets">
          <div class="bullet-item">
            <i class="pi pi-check-circle bullet-icon"></i>
            <span>Book verified professionals 24/7</span>
          </div>
          <div class="bullet-item">
            <i class="pi pi-check-circle bullet-icon"></i>
            <span>Manage practice appointments & schedules</span>
          </div>
        </div>
      </div>

      <!-- Right Side: Glassmorphism Register Form -->
      <div class="auth-right-content">
        <div class="glass-card">
          <div class="glass-logo">
            <i class="pi pi-calendar-plus logo-icon"></i>
            <span class="logo-text">Appointly</span>
          </div>

          <h2 class="auth-title">Create Account</h2>

          <!-- Message Alerts -->
          <PMessage v-if="error" severity="error" :closable="false" style="margin-bottom: 18px;">
            {{ error }}
          </PMessage>
          <PMessage v-if="validationErrors.length" severity="error" :closable="false" style="margin-bottom: 18px;">
            <ul style="margin: 0; padding-left: 16px;">
              <li v-for="(err, i) in validationErrors" :key="i">{{ err }}</li>
            </ul>
          </PMessage>
          <PMessage v-if="success" severity="success" :closable="false" style="margin-bottom: 18px;">
            {{ success }}
          </PMessage>

          <form class="auth-form" @submit.prevent="handleRegister">
            <div class="form-group">
              <label class="glass-label">Full Name</label>
              <div class="input-wrapper">
                <i class="pi pi-user input-icon"></i>
                <PInputText
                  v-model="form.name"
                  type="text"
                  placeholder="John Doe"
                  class="glass-input"
                  required
                />
              </div>
            </div>

            <div class="form-group">
              <label class="glass-label">Email</label>
              <div class="input-wrapper">
                <i class="pi pi-envelope input-icon"></i>
                <PInputText
                  v-model="form.email"
                  type="email"
                  placeholder="username@gmail.com"
                  class="glass-input"
                  required
                />
              </div>
            </div>

            <div class="form-group">
              <label class="glass-label">Phone Number</label>
              <div class="input-wrapper">
                <i class="pi pi-phone input-icon"></i>
                <PInputText
                  v-model="form.phone"
                  type="text"
                  placeholder="9876543210"
                  maxlength="10"
                  class="glass-input"
                  required
                />
              </div>
            </div>

            <div class="form-group">
              <label class="glass-label">Password</label>
              <div class="input-wrapper">
                <i class="pi pi-lock input-icon"></i>
                <PInputText
                  v-model="form.password"
                  type="password"
                  placeholder="Password"
                  class="glass-input"
                  required
                />
              </div>
            </div>

            <div class="form-group">
              <label class="glass-label">Account Role</label>
              <PSelect
                v-model="form.role"
                :options="roles"
                optionLabel="label"
                optionValue="value"
                placeholder="Select account role"
                class="glass-select"
                required
              />
            </div>

            <PButton
              type="submit"
              label="Create Account"
              :loading="loading"
              class="glass-submit-btn"
            />
          </form>

          <div class="auth-footer">
            Already have an account?
            <router-link to="/login" class="auth-link">Sign in</router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { registerUser } from '../services/api.js'

export default {
  name: 'RegisterPage',
  data() {
    return {
      form: {
        name: '',
        email: '',
        phone: '',
        password: '',
        role: '',
      },
      roles: [
        { label: 'Customer (Book Appointments)', value: 'Customer' },
        { label: 'Service Provider (Offer Services)', value: 'ServiceProvider' },
      ],
      loading: false,
      error: '',
      success: '',
      validationErrors: [],
    }
  },
  methods: {
    async handleRegister() {
      this.error = ''
      this.success = ''
      this.validationErrors = []
      this.loading = true

      try {
        const { status, data } = await registerUser(this.form)

        if (data.success) {
          this.success = data.message + ' — Redirecting to login...'
          setTimeout(() => {
            this.$router.push('/login')
          }, 1500)
        } else {
          if (Array.isArray(data.data)) {
            this.validationErrors = data.data
          } else {
            this.error = data.message
          }
        }
      } catch (err) {
        this.error = 'Unable to connect to the server'
      } finally {
        this.loading = false
      }
    },
  },
}
</script>

<style scoped>
.split-auth-container {
  min-height: calc(100vh - 90px);
  display: flex;
  align-items: center;
  position: relative;
  background: linear-gradient(135deg, rgba(2, 44, 115, 0.9) 0%, rgba(0, 82, 204, 0.85) 50%, rgba(0, 102, 255, 0.85) 100%),
              url('/auth_bg.png') center/cover no-repeat;
  padding: 60px 0;
  overflow: hidden;
}

.split-auth-inner {
  display: grid;
  grid-template-columns: 1.1fr 0.9fr;
  gap: 48px;
  align-items: center;
  position: relative;
  z-index: 10;
}

@media (max-width: 960px) {
  .split-auth-inner {
    grid-template-columns: 1fr;
  }
  .auth-left-content {
    text-align: center;
  }
  .bullet-item {
    justify-content: center;
  }
}

.auth-left-content {
  color: #FFFFFF;
}

.auth-brand-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(8px);
  padding: 6px 16px;
  border-radius: var(--radius-full);
  font-size: 0.85rem;
  font-weight: 700;
  margin-bottom: 24px;
}

.brand-icon {
  color: #60A5FA;
  font-size: 1.1rem;
}

.auth-hero-title {
  font-size: 3rem;
  font-weight: 800;
  line-height: 1.15;
  margin-bottom: 16px;
  letter-spacing: -0.02em;
}

.auth-hero-title .highlight {
  color: #60A5FA;
}

.auth-hero-desc {
  font-size: 1.1rem;
  opacity: 0.9;
  line-height: 1.6;
  margin-bottom: 32px;
  max-width: 520px;
}

.feature-bullets {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.bullet-item {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 1rem;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.92);
}

.bullet-icon {
  color: #60A5FA;
  font-size: 1.25rem;
}

.social-proof-pill {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  background: rgba(0, 0, 0, 0.25);
  backdrop-filter: blur(10px);
  padding: 10px 20px;
  border-radius: var(--radius-full);
  border: 1px solid rgba(255, 255, 255, 0.15);
  font-size: 0.88rem;
}

.stars {
  color: #F59E0B;
  letter-spacing: 2px;
  font-weight: 800;
}

.auth-right-content {
  display: flex;
  justify-content: center;
}

.glass-card {
  width: 100%;
  max-width: 480px;
  background: rgba(255, 255, 255, 0.14);
  backdrop-filter: blur(24px);
  -webkit-backdrop-filter: blur(24px);
  border: 1px solid rgba(255, 255, 255, 0.28);
  border-radius: 24px;
  padding: 40px 36px;
  box-shadow: 0 30px 60px rgba(0, 0, 0, 0.3);
  color: #FFFFFF;
}

.glass-logo {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  font-size: 1.4rem;
  font-weight: 800;
  color: #FFFFFF;
  margin-bottom: 20px;
}

.logo-icon {
  font-size: 1.6rem;
  color: #60A5FA;
}

.auth-title {
  font-size: 1.75rem;
  font-weight: 700;
  color: #FFFFFF;
  margin-bottom: 24px;
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.glass-label {
  font-size: 0.85rem;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.9);
}

.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: 16px;
  color: #64748B;
  font-size: 1rem;
  z-index: 2;
}

.glass-input {
  width: 100%;
  padding-left: 44px !important;
  border-radius: 12px !important;
  background: #FFFFFF !important;
  border: 1px solid rgba(255, 255, 255, 0.4) !important;
  color: #0F172A !important;
  font-size: 0.92rem !important;
  height: 46px;
}

.glass-input::placeholder {
  color: #94A3B8 !important;
}

.glass-select {
  width: 100%;
  border-radius: 12px !important;
  background: #FFFFFF !important;
  border: 1px solid rgba(255, 255, 255, 0.4) !important;
  color: #0F172A !important;
}

.glass-submit-btn {
  width: 100%;
  height: 48px;
  background: #062863 !important;
  border-color: #062863 !important;
  color: #FFFFFF !important;
  font-weight: 700 !important;
  font-size: 0.95rem !important;
  border-radius: 12px !important;
  margin-top: 8px;
  transition: all 0.2s ease !important;
}

.glass-submit-btn:hover {
  background: #041B44 !important;
  border-color: #041B44 !important;
  transform: translateY(-1px);
}

.auth-footer {
  margin-top: 28px;
  text-align: center;
  font-size: 0.88rem;
  color: rgba(255, 255, 255, 0.8);
}

.auth-link {
  color: #FFFFFF;
  font-weight: 700;
  text-decoration: underline;
  margin-left: 4px;
}

.auth-link:hover {
  color: #60A5FA;
}
</style>
