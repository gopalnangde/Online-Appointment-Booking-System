<template>
  <header class="navbar-wrapper">
    <div class="container">
      <nav class="navbar-card">
        <!-- Left Brand & Navigation Pills -->
        <div class="navbar-left">
          <router-link to="/" class="navbar-brand">
            <i class="pi pi-calendar-plus brand-icon"></i>
            <span class="brand-text">Appointly</span>
          </router-link>

          <div class="nav-pills">
            <router-link to="/appointments" class="nav-pill" v-if="isLoggedIn">
              Appointments
            </router-link>
            <router-link to="/providers" class="nav-pill">
              Browse Providers
            </router-link>
            <router-link v-if="isLoggedIn && user?.role === 'ServiceProvider'" to="/provider/dashboard" class="nav-pill">
              My Practice
            </router-link>
            <router-link to="/profile" class="nav-pill" v-if="isLoggedIn">
              Profile
            </router-link>
          </div>
        </div>

        <!-- Right Side: Guest Actions or User Profile Dropdown -->
        <div class="navbar-right">
          <template v-if="!isLoggedIn">
            <router-link to="/login" class="nav-pill">Sign In</router-link>
            <router-link to="/register">
              <PButton label="Get Started" icon="pi pi-user-plus" class="p-button-sm p-button-primary" />
            </router-link>
          </template>

          <template v-else>
            <div class="user-dropdown-wrapper" ref="dropdownRef">
              <button class="user-profile-pill" @click="dropdownOpen = !dropdownOpen">
                <div class="user-avatar-sm">{{ getInitials(user?.name) }}</div>
                <span class="user-name">{{ user?.name || 'Account' }}</span>
                <span class="role-badge">{{ user?.role === 'ServiceProvider' ? 'Provider' : 'Client' }}</span>
                <i class="pi pi-chevron-down dropdown-chevron" :class="{ rotated: dropdownOpen }"></i>
              </button>

              <!-- Dropdown Menu -->
              <transition name="fade-slide">
                <div v-if="dropdownOpen" class="dropdown-menu">
                  <div class="dropdown-header">
                    <p class="dh-name">{{ user?.name }}</p>
                    <p class="dh-email">{{ user?.email }}</p>
                    <span class="dh-role">{{ user?.role === 'ServiceProvider' ? '🏥 Service Provider' : '👤 Customer' }}</span>
                  </div>

                  <div class="dropdown-divider"></div>

                  <router-link to="/profile" class="dropdown-item" @click="dropdownOpen = false">
                    <i class="pi pi-user"></i> My Profile
                  </router-link>

                  <router-link to="/appointments" class="dropdown-item" @click="dropdownOpen = false">
                    <i class="pi pi-calendar"></i> My Appointments
                  </router-link>

                  <router-link v-if="user?.role === 'ServiceProvider'" to="/provider/dashboard" class="dropdown-item" @click="dropdownOpen = false">
                    <i class="pi pi-briefcase"></i> My Practice Dashboard
                  </router-link>

                  <div class="dropdown-divider"></div>

                  <!-- EXPLICIT LOGOUT OPTION FOR ALL USERS INCLUDING SERVICE PROVIDERS -->
                  <button class="dropdown-item logout-item" @click="logout">
                    <i class="pi pi-power-off"></i> Logout Account
                  </button>
                </div>
              </transition>
            </div>
          </template>

          <!-- Mobile Hamburger Toggle -->
          <button class="mobile-toggle" @click="mobileOpen = !mobileOpen">
            <i :class="mobileOpen ? 'pi pi-times' : 'pi pi-bars'"></i>
          </button>
        </div>
      </nav>

      <!-- Mobile Dropdown Menu -->
      <transition name="fade-slide">
        <div v-if="mobileOpen" class="mobile-menu-card">
          <template v-if="isLoggedIn">
            <div class="mobile-user-info">
              <div class="user-avatar-sm">{{ getInitials(user?.name) }}</div>
              <div>
                <strong>{{ user?.name }}</strong>
                <p style="font-size: 0.8rem; color: var(--color-gray-500);">{{ user?.email }}</p>
              </div>
            </div>
            <router-link to="/appointments" class="mobile-link" @click="mobileOpen = false">
              <i class="pi pi-calendar"></i> Appointments
            </router-link>
            <router-link to="/providers" class="mobile-link" @click="mobileOpen = false">
              <i class="pi pi-search"></i> Browse Providers
            </router-link>
            <router-link v-if="user?.role === 'ServiceProvider'" to="/provider/dashboard" class="mobile-link" @click="mobileOpen = false">
              <i class="pi pi-briefcase"></i> My Practice
            </router-link>
            <router-link to="/profile" class="mobile-link" @click="mobileOpen = false">
              <i class="pi pi-user"></i> Profile
            </router-link>
            <button class="mobile-link logout-link" @click="logout">
              <i class="pi pi-power-off"></i> Logout
            </button>
          </template>
          <template v-else>
            <router-link to="/providers" class="mobile-link" @click="mobileOpen = false">
              <i class="pi pi-search"></i> Browse Providers
            </router-link>
            <router-link to="/login" class="mobile-link" @click="mobileOpen = false">
              <i class="pi pi-sign-in"></i> Sign In
            </router-link>
            <router-link to="/register" class="mobile-link" @click="mobileOpen = false">
              <i class="pi pi-user-plus"></i> Get Started
            </router-link>
          </template>
        </div>
      </transition>
    </div>
  </header>
</template>

<script>
export default {
  name: 'Navbar',
  data() {
    return {
      isLoggedIn: false,
      user: null,
      dropdownOpen: false,
      mobileOpen: false,
    }
  },
  watch: {
    $route() {
      this.checkAuth()
    },
  },
  methods: {
    checkAuth() {
      const token = localStorage.getItem('token')
      const userData = localStorage.getItem('user')
      this.isLoggedIn = !!token
      this.user = userData ? JSON.parse(userData) : null
    },
    getInitials(name) {
      if (!name) return '?'
      return name
        .split(' ')
        .map(w => w[0])
        .join('')
        .toUpperCase()
        .slice(0, 2)
    },
    logout() {
      this.dropdownOpen = false
      this.mobileOpen = false
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      this.checkAuth()
      window.dispatchEvent(new Event('auth-change'))
      this.$router.push('/login')
    },
    handleClickOutside(e) {
      if (this.$refs.dropdownRef && !this.$refs.dropdownRef.contains(e.target)) {
        this.dropdownOpen = false
      }
    },
  },
  created() {
    this.checkAuth()
  },
  mounted() {
    this.checkAuth()
    window.addEventListener('auth-change', this.checkAuth)
    window.addEventListener('storage', this.checkAuth)
  },
  beforeUnmount() {
    window.removeEventListener('auth-change', this.checkAuth)
    window.removeEventListener('storage', this.checkAuth)
  },
}
</script>

<style scoped>
.navbar-wrapper {
  padding: 16px 0;
  position: sticky;
  top: 0;
  z-index: 1000;
}

.navbar-card {
  background: var(--color-white);
  border-radius: var(--radius-full);
  padding: 10px 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 4px 20px rgba(0, 102, 255, 0.08);
  border: 1px solid var(--color-gray-200);
}

.navbar-left {
  display: flex;
  align-items: center;
  gap: 32px;
}

.navbar-brand {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 1.35rem;
  font-weight: 800;
  color: var(--color-primary);
  letter-spacing: -0.02em;
}

.brand-icon {
  font-size: 1.5rem;
  color: var(--color-primary);
}

.nav-pills {
  display: flex;
  align-items: center;
  gap: 8px;
}

.nav-pill {
  padding: 8px 18px;
  font-size: 0.9rem;
  font-weight: 600;
  color: var(--color-gray-600);
  border-radius: var(--radius-full);
  transition: var(--transition);
}

.nav-pill:hover,
.nav-pill.router-link-active {
  color: var(--color-primary);
  background: var(--color-primary-pale);
}

.navbar-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

/* User Dropdown Pill */
.user-dropdown-wrapper {
  position: relative;
}

.user-profile-pill {
  display: flex;
  align-items: center;
  gap: 10px;
  background: var(--color-gray-100);
  border: 1px solid var(--color-gray-200);
  padding: 6px 16px 6px 8px;
  border-radius: var(--radius-full);
  cursor: pointer;
  transition: var(--transition);
}

.user-profile-pill:hover {
  background: var(--color-primary-pale);
  border-color: var(--color-primary-light);
}

.user-avatar-sm {
  width: 32px;
  height: 32px;
  border-radius: var(--radius-full);
  background: var(--color-primary);
  color: var(--color-white);
  font-weight: 700;
  font-size: 0.82rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.user-name {
  font-size: 0.88rem;
  font-weight: 700;
  color: var(--color-dark);
}

.role-badge {
  font-size: 0.7rem;
  font-weight: 700;
  background: var(--color-primary-light);
  color: var(--color-primary);
  padding: 2px 8px;
  border-radius: 10px;
  text-transform: uppercase;
}

.dropdown-chevron {
  font-size: 0.75rem;
  color: var(--color-gray-500);
  transition: transform 0.2s ease;
}

.dropdown-chevron.rotated {
  transform: rotate(180deg);
}

/* Dropdown Menu Popup */
.dropdown-menu {
  position: absolute;
  right: 0;
  top: calc(100% + 12px);
  width: 250px;
  background: var(--color-white);
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-gray-200);
  box-shadow: var(--shadow-lg);
  padding: 12px;
  z-index: 1001;
}

.dropdown-header {
  padding: 8px 12px;
}

.dh-name {
  font-weight: 700;
  font-size: 0.95rem;
  color: var(--color-dark);
}

.dh-email {
  font-size: 0.8rem;
  color: var(--color-gray-500);
  margin-top: 2px;
}

.dh-role {
  display: inline-block;
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--color-primary);
  margin-top: 6px;
}

.dropdown-divider {
  height: 1px;
  background: var(--color-gray-200);
  margin: 8px 0;
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  padding: 10px 12px;
  font-size: 0.88rem;
  font-weight: 600;
  color: var(--color-dark-secondary);
  border-radius: var(--radius-sm);
  background: none;
  border: none;
  cursor: pointer;
  transition: var(--transition);
  text-align: left;
}

.dropdown-item:hover {
  background: var(--color-gray-100);
  color: var(--color-primary);
}

.dropdown-item.logout-item {
  color: var(--color-danger);
}

.dropdown-item.logout-item:hover {
  background: var(--color-danger-light);
  color: var(--color-danger);
}

/* Mobile Layout */
.mobile-toggle {
  display: none;
  background: none;
  border: none;
  font-size: 1.3rem;
  color: var(--color-dark);
  cursor: pointer;
}

.mobile-menu-card {
  margin-top: 12px;
  background: var(--color-white);
  border-radius: var(--radius-lg);
  padding: 16px;
  border: 1px solid var(--color-gray-200);
  box-shadow: var(--shadow-md);
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.mobile-user-info {
  display: flex;
  align-items: center;
  gap: 12px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--color-gray-200);
  margin-bottom: 8px;
}

.mobile-link {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 14px;
  font-size: 0.92rem;
  font-weight: 600;
  color: var(--color-dark);
  border-radius: var(--radius-sm);
  background: none;
  border: none;
  width: 100%;
}

.mobile-link.logout-link {
  color: var(--color-danger);
}

.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.2s ease;
}
.fade-slide-enter-from,
.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

@media (max-width: 900px) {
  .nav-pills {
    display: none;
  }
  .mobile-toggle {
    display: block;
  }
}
</style>
