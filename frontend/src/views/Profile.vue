<template>
  <div class="dashboard-wrapper">
    <div class="container">
      <div class="page-header">
        <h1 class="page-title">
          <i class="pi pi-user" style="margin-right: 10px; color: var(--color-primary-dark);"></i>
          My Account Profile
        </h1>
        <p class="page-subtitle">View and update your personal details and account settings</p>
      </div>

      <div v-if="loading" style="text-align: center; padding: 60px;">
        <PProgressSpinner style="width: 50px; height: 50px;" strokeWidth="4" />
      </div>

      <div v-else-if="user" class="profile-grid">
        <!-- Sidebar -->
        <div class="profile-sidebar">
          <div class="profile-avatar-lg">{{ initials }}</div>
          <h2>{{ user.name }}</h2>
          <p class="email">{{ user.email }}</p>
          <div style="margin-top: 12px;">
            <PTag :value="formatRole(user.role)" severity="primary" />
          </div>
          <div class="divider" style="margin: 20px 0;"></div>
          
          <div style="display: flex; flex-direction: column; gap: 12px;">
            <PButton
              v-if="!isEditing"
              label="Edit Profile"
              icon="pi pi-user-edit"
              class="p-button-primary p-button-sm w-full"
              @click="startEditing"
            />
            <PButton
              v-else
              label="Cancel Editing"
              icon="pi pi-times"
              class="p-button-outlined p-button-secondary p-button-sm w-full"
              @click="cancelEditing"
            />

            <router-link to="/appointments">
              <PButton label="My Appointments" icon="pi pi-calendar" class="p-button-outlined p-button-sm w-full" />
            </router-link>
            
            <router-link
              v-if="user.role === 'ServiceProvider'"
              to="/provider/dashboard"
            >
              <PButton label="Provider Dashboard" icon="pi pi-briefcase" class="p-button-outlined p-button-info p-button-sm w-full" />
            </router-link>
          </div>
        </div>

        <!-- Main Details / Edit Form -->
        <div class="profile-details">
          <!-- Action Banners -->
          <PMessage v-if="successMsg" severity="success" :closable="true" class="mb-4" @close="successMsg = ''">
            {{ successMsg }}
          </PMessage>
          <PMessage v-if="errorMsg" severity="error" :closable="true" class="mb-4" @close="errorMsg = ''">
            {{ errorMsg }}
          </PMessage>

          <!-- Read-Only View -->
          <div v-if="!isEditing">
            <div class="details-header-row">
              <h3><i class="pi pi-id-card" style="margin-right: 8px;"></i> Account Information</h3>
              <PButton
                label="Edit Details"
                icon="pi pi-pencil"
                class="p-button-outlined p-button-sm"
                @click="startEditing"
              />
            </div>
            
            <div class="detail-grid">
              <div class="detail-item">
                <label>Full Name</label>
                <p>{{ user.name }}</p>
              </div>
              <div class="detail-item">
                <label>Email Address</label>
                <p>{{ user.email }}</p>
              </div>
              <div class="detail-item">
                <label>Phone Number</label>
                <p>{{ user.phone }}</p>
              </div>
              <div class="detail-item">
                <label>Account Role</label>
                <p>{{ formatRole(user.role) }}</p>
              </div>
            </div>
          </div>

          <!-- Edit Profile Form -->
          <div v-else class="edit-form-card">
            <h3 class="form-title">
              <i class="pi pi-user-edit" style="margin-right: 8px; color: var(--color-primary);"></i>
              Update Personal Information
            </h3>
            <p class="form-subtitle">Make changes to your account profile details below.</p>
            
            <form @submit.prevent="saveProfile" class="profile-edit-form">
              <div class="form-group">
                <label for="name">Full Name <span class="required">*</span></label>
                <input
                  id="name"
                  v-model="form.name"
                  type="text"
                  class="form-input"
                  placeholder="Enter your full name"
                  required
                />
              </div>

              <div class="form-group">
                <label for="email">Email Address <span class="required">*</span></label>
                <input
                  id="email"
                  v-model="form.email"
                  type="email"
                  class="form-input"
                  placeholder="Enter your email address"
                  required
                />
              </div>

              <div class="form-group">
                <label for="phone">Phone Number (10 digits) <span class="required">*</span></label>
                <input
                  id="phone"
                  v-model="form.phone"
                  type="tel"
                  class="form-input"
                  placeholder="e.g. 9876543210"
                  maxlength="10"
                  required
                />
              </div>

              <div class="form-group">
                <label for="password">New Password <span class="optional">(Optional - leave blank to keep current)</span></label>
                <input
                  id="password"
                  v-model="form.password"
                  type="password"
                  class="form-input"
                  placeholder="Enter new password (min. 8 characters)"
                  minlength="8"
                />
              </div>

              <div class="form-actions">
                <PButton
                  type="button"
                  label="Cancel"
                  class="p-button-outlined p-button-secondary"
                  @click="cancelEditing"
                  :disabled="saving"
                />
                <PButton
                  type="submit"
                  label="Save Changes"
                  icon="pi pi-check"
                  class="p-button-primary"
                  :loading="saving"
                />
              </div>
            </form>
          </div>
        </div>
      </div>

      <PMessage v-else severity="error" :closable="false" style="max-width: 500px;">
        Failed to load profile. Please try logging in again.
      </PMessage>
    </div>
  </div>
</template>

<script>
import { getProfile, updateUserProfile } from '../services/api.js'

export default {
  name: 'ProfilePage',
  data() {
    return {
      user: null,
      loading: true,
      saving: false,
      isEditing: false,
      successMsg: '',
      errorMsg: '',
      form: {
        name: '',
        email: '',
        phone: '',
        password: '',
      },
    }
  },
  computed: {
    initials() {
      if (!this.user?.name) return '?'
      return this.user.name
        .split(' ')
        .map(w => w[0])
        .join('')
        .toUpperCase()
        .slice(0, 2)
    },
  },
  methods: {
    formatRole(role) {
      return role === 'ServiceProvider' ? 'Service Provider' : role
    },

    startEditing() {
      this.form = {
        name: this.user.name || '',
        email: this.user.email || '',
        phone: this.user.phone || '',
        password: '',
      }
      this.isEditing = true
      this.successMsg = ''
      this.errorMsg = ''
    },

    cancelEditing() {
      this.isEditing = false
      this.errorMsg = ''
    },

    async fetchProfile() {
      this.loading = true
      try {
        const { status, data } = await getProfile()
        if (data.success && data.data) {
          this.user = data.data
        }
      } catch (err) {
        console.error('Failed to fetch profile:', err)
      } finally {
        this.loading = false
      }
    },

    async saveProfile() {
      this.saving = true
      this.successMsg = ''
      this.errorMsg = ''

      const payload = {
        name: this.form.name,
        email: this.form.email,
        phone: this.form.phone,
      }

      if (this.form.password.trim()) {
        payload.password = this.form.password.trim()
      }

      try {
        const { status, data } = await updateUserProfile(payload)
        if (data.success && data.data) {
          this.user = data.data
          // Update localStorage user so Navbar updates name immediately
          const cachedUser = JSON.parse(localStorage.getItem('user') || '{}')
          const updatedUser = { ...cachedUser, ...data.data }
          localStorage.setItem('user', JSON.stringify(updatedUser))
          window.dispatchEvent(new Event('auth-change'))

          this.successMsg = 'Your profile has been updated successfully!'
          this.isEditing = false
          
          setTimeout(() => {
            this.successMsg = ''
          }, 4000)
        } else {
          this.errorMsg = data.message || 'Failed to update profile'
        }
      } catch (err) {
        console.error('Error updating profile:', err)
        this.errorMsg = err.response?.data?.message || err.message || 'An error occurred while updating profile'
      } finally {
        this.saving = false
      }
    },
  },
  mounted() {
    this.fetchProfile()
  },
}
</script>

<style scoped>
.details-header-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.details-header-row h3 {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--color-dark);
}

.edit-form-card {
  background: var(--color-white);
  padding: 8px 0;
}

.form-title {
  font-size: 1.25rem;
  font-weight: 800;
  color: var(--color-dark);
  margin-bottom: 4px;
}

.form-subtitle {
  font-size: 0.88rem;
  color: var(--color-gray-500);
  margin-bottom: 24px;
}

.profile-edit-form {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group label {
  font-size: 0.88rem;
  font-weight: 700;
  color: var(--color-dark-secondary);
}

.required {
  color: var(--color-danger);
}

.optional {
  font-size: 0.78rem;
  font-weight: 500;
  color: var(--color-gray-400);
}

.form-input {
  width: 100%;
  padding: 10px 14px;
  border: 1.5px solid var(--color-gray-200);
  border-radius: var(--radius-md);
  font-size: 0.92rem;
  transition: var(--transition);
  outline: none;
}

.form-input:focus {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(0, 102, 255, 0.12);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 12px;
  padding-top: 18px;
  border-top: 1px solid var(--color-gray-200);
}

.w-full {
  width: 100%;
}

.mb-4 {
  margin-bottom: 16px;
}
</style>
