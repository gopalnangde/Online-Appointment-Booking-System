<template>
  <div class="dashboard-wrapper">
    <div class="container">
      <!-- Top Action Bar (Matching Reference Medical Dashboard UI) -->
      <div class="header-action-card">
        <div class="header-title-group">
          <h1 class="header-title">
            Appointments
            <i class="pi pi-calendar title-icon"></i>
          </h1>
        </div>

        <div class="search-bar-container">
          <i class="pi pi-search search-bar-icon"></i>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search appointments, patients, providers..."
            class="search-bar-input"
          />
          <button v-if="searchQuery" class="clear-search-btn" @click="searchQuery = ''" title="Clear search">
            <i class="pi pi-times-circle"></i>
          </button>
        </div>

        <div class="header-actions-right">
          <router-link v-if="!isProvider" to="/providers">
            <PButton label="+ New Appointment" class="p-button-primary" />
          </router-link>
        </div>
      </div>

      <!-- Secondary Controls & Filters Row (Matching Screenshot Filter Row) -->
      <div class="filters-card">
        <div class="filters-left">
          <!-- Status Select Filter -->
          <select v-model="statusFilter" class="pill-select">
            <option value="All">All Statuses</option>
            <option value="Pending">Pending</option>
            <option value="Confirmed">Confirmed</option>
            <option value="Completed">Completed</option>
            <option value="Cancelled">Cancelled</option>
          </select>

          <label class="checkbox-label">
            <input type="checkbox" v-model="hideCompleted" />
            <span>Hide Completed</span>
          </label>
        </div>

        <!-- Filter Tabs Pill Group -->
        <div class="filter-tabs-pills">
          <button
            v-for="tab in tabs"
            :key="tab"
            class="tab-pill-btn"
            :class="{ active: activeTab === tab }"
            @click="activeTab = tab"
          >
            {{ tab }}
            <span class="tab-count-badge" v-if="getTabCount(tab) > 0">
              {{ getTabCount(tab) }}
            </span>
          </button>
        </div>
      </div>

      <!-- Action Banners -->
      <PMessage v-if="actionSuccess" severity="success" :closable="true" class="mb-4" @close="actionSuccess = ''">
        {{ actionSuccess }}
      </PMessage>
      <PMessage v-if="actionError" severity="error" :closable="true" class="mb-4" @close="actionError = ''">
        {{ actionError }}
      </PMessage>

      <!-- Loading State -->
      <div v-if="loading" style="text-align: center; padding: 60px;">
        <PProgressSpinner style="width: 50px; height: 50px;" strokeWidth="4" />
      </div>

      <!-- Empty State -->
      <div v-else-if="filteredAppointments.length === 0" class="empty-state">
        <i class="pi pi-calendar-times" style="font-size: 3rem; color: var(--color-gray-400); margin-bottom: 16px;"></i>
        <h3>No Appointments Found</h3>
        <p>No appointment records match your current filter or search criteria.</p>
        <router-link v-if="!isProvider" to="/providers">
          <PButton label="Browse Service Providers" icon="pi pi-search" class="p-button-sm p-button-primary" />
        </router-link>
      </div>

      <!-- Data Table View (Matching Screenshot Grid Table) -->
      <div v-else class="table-responsive">
        <table class="custom-table">
          <thead>
            <tr>
              <th style="width: 40px;"><input type="checkbox" /></th>
              <th>Time & Date</th>
              <th>{{ isProvider ? 'Patient / Customer' : 'Service Provider' }}</th>
              <th>Service Title</th>
              <th>Category</th>
              <th>Location</th>
              <th>Status</th>
              <th style="text-align: right;">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="apt in filteredAppointments" :key="apt.id">
              <td><input type="checkbox" /></td>
              <td>
                <div class="time-cell">
                  <span class="time-text">{{ apt.booking_time }}</span>
                  <span class="date-subtext">{{ formatDate(apt.booking_date) }}</span>
                </div>
              </td>
              <td>
                <div class="user-cell">
                  <div class="user-avatar-tiny">
                    {{ getInitials(isProvider ? apt.customer?.name : apt.provider?.name) }}
                  </div>
                  <div>
                    <span
                      class="user-name-text"
                      :class="{ 'clickable-provider-link': !isProvider, 'clickable-patient-link': isProvider }"
                      @click="isProvider ? openPatientModal(apt.customer, apt.customer_id) : viewProviderProfile(apt)"
                      :title="isProvider ? 'Click to view patient profile' : 'View provider profile & ratings'"
                    >
                      {{ isProvider ? (apt.customer?.name || 'Customer') : (apt.provider?.name || 'Provider') }}
                    </span>
                    <span class="user-sub-email">{{ isProvider ? apt.customer?.email : apt.provider_profile?.specialization }}</span>
                  </div>
                </div>
              </td>
              <td>
                <strong style="color: var(--color-dark);">{{ apt.service_title }}</strong>
                <p v-if="apt.notes" class="notes-truncate" :title="apt.notes">{{ apt.notes }}</p>
              </td>
              <td>
                <span class="category-chip">{{ apt.provider_profile?.specialization || 'General' }}</span>
              </td>
              <td>
                <span>{{ apt.provider_profile?.city || 'Online / Clinic' }}</span>
              </td>
              <td>
                <PTag :value="apt.status" :severity="getTagSeverity(apt.status)" />
              </td>
              <td style="text-align: right;">
                <div class="action-buttons-group">
                  <!-- Service Provider Status & Patient Profile Controls -->
                  <template v-if="isProvider">
                    <button
                      class="icon-action-btn btn-view-patient"
                      title="View Patient Profile"
                      @click="openPatientModal(apt.customer, apt.customer_id)"
                    >
                      <i class="pi pi-user"></i>
                    </button>

                    <button
                      v-if="apt.status === 'Pending'"
                      class="icon-action-btn btn-confirm"
                      title="Confirm Appointment"
                      :disabled="updatingId === apt.id"
                      @click="changeStatus(apt.id, 'Confirmed')"
                    >
                      <i class="pi pi-check"></i>
                    </button>

                    <button
                      v-if="apt.status === 'Confirmed'"
                      class="icon-action-btn btn-complete"
                      title="Mark Completed"
                      :disabled="updatingId === apt.id"
                      @click="changeStatus(apt.id, 'Completed')"
                    >
                      <i class="pi pi-check-circle"></i>
                    </button>

                    <button
                      v-if="apt.status === 'Pending' || apt.status === 'Confirmed'"
                      class="icon-action-btn btn-cancel"
                      title="Cancel Appointment"
                      :disabled="updatingId === apt.id"
                      @click="changeStatus(apt.id, 'Cancelled')"
                    >
                      <i class="pi pi-times"></i>
                    </button>
                  </template>

                  <!-- Customer Actions -->
                  <template v-else>
                    <button
                      v-if="apt.status === 'Pending' || apt.status === 'Confirmed'"
                      class="icon-action-btn btn-cancel"
                      title="Cancel Appointment"
                      :disabled="updatingId === apt.id"
                      @click="changeStatus(apt.id, 'Cancelled')"
                    >
                      <i class="pi pi-times"></i>
                    </button>

                    <button
                      v-if="apt.status === 'Completed'"
                      class="icon-action-btn btn-review"
                      title="Leave Review"
                      @click="openReviewModal(apt)"
                    >
                      <i class="pi pi-star-fill"></i>
                    </button>
                  </template>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Review Modal -->
    <ReviewModal
      v-if="showReviewModal && selectedAppointment"
      :appointment="selectedAppointment"
      @close="showReviewModal = false"
      @success="handleReviewSuccess"
    />

    <!-- Patient Profile Modal -->
    <PatientProfileModal
      v-if="showPatientModal && selectedCustomerId"
      :customerId="selectedCustomerId"
      :customerData="selectedCustomerData"
      @close="showPatientModal = false"
    />
  </div>
</template>

<script>
import { getMyAppointments, updateAppointmentStatus } from '../services/api.js'
import ReviewModal from '../components/ReviewModal.vue'
import PatientProfileModal from '../components/PatientProfileModal.vue'

export default {
  name: 'AppointmentsPage',
  components: { ReviewModal, PatientProfileModal },
  data() {
    return {
      appointments: [],
      loading: true,
      searchQuery: '',
      statusFilter: 'All',
      hideCompleted: false,
      activeTab: 'All',
      tabs: ['All', 'Pending', 'Confirmed', 'Completed', 'Cancelled'],
      updatingId: null,
      actionSuccess: '',
      actionError: '',
      showReviewModal: false,
      selectedAppointment: null,
      showPatientModal: false,
      selectedCustomerId: null,
      selectedCustomerData: null,
    }
  },
  computed: {
    user() {
      return JSON.parse(localStorage.getItem('user') || 'null')
    },
    isProvider() {
      return this.user?.role === 'ServiceProvider'
    },
    filteredAppointments() {
      return this.appointments.filter(a => {
        // Tab filter
        if (this.activeTab !== 'All' && a.status !== this.activeTab) return false
        // Select filter
        if (this.statusFilter !== 'All' && a.status !== this.statusFilter) return false
        // Hide completed checkbox
        if (this.hideCompleted && a.status === 'Completed') return false

        // Search query filter
        if (this.searchQuery.trim()) {
          const q = this.searchQuery.toLowerCase()
          const serviceTitle = (a.service_title || '').toLowerCase()
          const customerName = (a.customer?.name || '').toLowerCase()
          const providerName = (a.provider?.name || '').toLowerCase()
          const spec = (a.provider_profile?.specialization || '').toLowerCase()

          return (
            serviceTitle.includes(q) ||
            customerName.includes(q) ||
            providerName.includes(q) ||
            spec.includes(q)
          )
        }

        return true
      })
    },
  },
  mounted() {
    this.fetchAppointments()
  },
  methods: {
    getInitials(name) {
      if (!name) return '?'
      return name
        .split(' ')
        .map(w => w[0])
        .join('')
        .toUpperCase()
        .slice(0, 2)
    },

    getTabCount(tab) {
      if (tab === 'All') return this.appointments.length
      return this.appointments.filter(a => a.status === tab).length
    },

    getTagSeverity(status) {
      switch (status) {
        case 'Pending':
          return 'warn'
        case 'Confirmed':
          return 'success'
        case 'Completed':
          return 'info'
        case 'Cancelled':
          return 'danger'
        default:
          return 'secondary'
      }
    },

    formatDate(dateStr) {
      if (!dateStr) return ''
      const date = new Date(dateStr)
      if (isNaN(date)) return dateStr
      return date.toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
        year: 'numeric',
      })
    },

    async fetchAppointments() {
      this.loading = true
      try {
        const { status, data } = await getMyAppointments()
        if (data.success && data.data) {
          this.appointments = data.data
        }
      } catch (err) {
        console.error('Error fetching appointments:', err)
      } finally {
        this.loading = false
      }
    },

    async changeStatus(id, newStatus) {
      this.updatingId = id
      this.actionSuccess = ''
      this.actionError = ''

      try {
        const { status, data } = await updateAppointmentStatus(id, newStatus)
        if (data && data.success) {
          this.actionSuccess = `Appointment status updated to ${newStatus}`
          await this.fetchAppointments()
          setTimeout(() => {
            this.actionSuccess = ''
          }, 3000)
        } else {
          this.actionError = data?.message || 'Failed to update appointment status'
        }
      } catch (err) {
        this.actionError = err.response?.data?.message || err.message || 'Error updating status'
      } finally {
        this.updatingId = null
      }
    },

    openReviewModal(apt) {
      this.selectedAppointment = apt
      this.showReviewModal = true
    },

    openPatientModal(customer, customerId) {
      this.selectedCustomerData = customer || null
      this.selectedCustomerId = customerId || customer?.id
      this.showPatientModal = true
    },

    viewProviderProfile(apt) {
      const targetId = apt.provider_profile?.id || apt.provider_id
      if (targetId) {
        this.$router.push(`/providers/${targetId}`)
      }
    },

    handleReviewSuccess() {
      this.actionSuccess = 'Review submitted successfully!'
      setTimeout(() => {
        this.actionSuccess = ''
      }, 3000)
    },
  },
}
</script>

<style scoped>
/* Top Header Action Card */
.header-action-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 20px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.header-title-group {
  display: flex;
  align-items: center;
}

.header-title {
  font-size: 1.8rem;
  font-weight: 800;
  color: var(--color-primary);
  display: flex;
  align-items: center;
  gap: 8px;
}

.title-icon {
  font-size: 1.4rem;
  color: var(--color-primary);
}

.search-bar-container {
  position: relative;
  display: flex;
  align-items: center;
  flex: 1;
  max-width: 580px;
  min-width: 320px;
  background: var(--color-white);
  border: 1.5px solid var(--color-gray-200);
  border-radius: var(--radius-full);
  padding: 4px 18px;
  box-shadow: 0 2px 10px rgba(0, 102, 255, 0.04);
  transition: var(--transition);
}

.search-bar-container:focus-within,
.search-bar-container:hover {
  border-color: var(--color-primary);
  box-shadow: 0 4px 18px rgba(0, 102, 255, 0.12);
}

.search-bar-icon {
  font-size: 1.1rem;
  color: var(--color-primary);
  margin-right: 12px;
}

.search-bar-input {
  flex: 1;
  border: none !important;
  outline: none !important;
  background: transparent !important;
  padding: 10px 0 !important;
  font-size: 0.92rem !important;
  font-weight: 500;
  color: var(--color-dark) !important;
  box-shadow: none !important;
}

.search-bar-input::placeholder {
  color: var(--color-gray-400);
  font-size: 0.9rem;
}

.clear-search-btn {
  background: none;
  border: none;
  color: var(--color-gray-400);
  cursor: pointer;
  padding: 4px;
  font-size: 1rem;
  transition: var(--transition);
  display: flex;
  align-items: center;
}

.clear-search-btn:hover {
  color: var(--color-danger);
}

/* Secondary Filters Bar */
.filters-card {
  background: var(--color-white);
  border-radius: var(--radius-lg);
  padding: 12px 20px;
  border: 1px solid var(--color-gray-200);
  box-shadow: var(--shadow-xs);
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.filters-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.pill-select {
  padding: 6px 16px;
  border-radius: var(--radius-full);
  border: 1px solid var(--color-gray-200);
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--color-dark-secondary);
  background: var(--color-gray-100);
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--color-gray-600);
  cursor: pointer;
}

.filter-tabs-pills {
  display: flex;
  align-items: center;
  gap: 6px;
}

.tab-pill-btn {
  padding: 6px 14px;
  border-radius: var(--radius-full);
  font-size: 0.82rem;
  font-weight: 600;
  color: var(--color-gray-600);
  background: none;
  border: none;
  cursor: pointer;
  transition: var(--transition);
  display: flex;
  align-items: center;
  gap: 6px;
}

.tab-pill-btn:hover {
  background: var(--color-gray-100);
  color: var(--color-primary);
}

.tab-pill-btn.active {
  background: var(--color-primary);
  color: var(--color-white);
}

.tab-count-badge {
  background: rgba(255, 255, 255, 0.25);
  padding: 1px 6px;
  border-radius: 10px;
  font-size: 0.75rem;
}

/* Table Specific Cell Helpers */
.time-cell {
  display: flex;
  flex-direction: column;
}

.time-text {
  font-weight: 700;
  color: var(--color-dark);
}

.date-subtext {
  font-size: 0.78rem;
  color: var(--color-gray-500);
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-avatar-tiny {
  width: 32px;
  height: 32px;
  border-radius: var(--radius-full);
  background: var(--color-primary-light);
  color: var(--color-primary);
  font-weight: 700;
  font-size: 0.8rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.user-name-text {
  display: block;
  font-weight: 700;
  font-size: 0.88rem;
}

.clickable-provider-link {
  cursor: pointer;
  color: var(--color-primary);
  transition: var(--transition);
}

.clickable-provider-link:hover {
  text-decoration: underline;
  color: var(--color-primary-dark);
}

.clickable-patient-link {
  cursor: pointer;
  color: var(--color-primary);
  transition: var(--transition);
}

.clickable-patient-link:hover {
  text-decoration: underline;
  color: var(--color-primary-dark);
}

.btn-view-patient {
  background: var(--color-primary-pale);
  color: var(--color-primary);
}
.btn-view-patient:hover {
  background: var(--color-primary);
  color: var(--color-white);
}

.user-sub-email {
  display: block;
  font-size: 0.76rem;
  color: var(--color-gray-500);
}

.category-chip {
  background: var(--color-gray-100);
  color: var(--color-gray-600);
  padding: 4px 10px;
  border-radius: var(--radius-full);
  font-size: 0.78rem;
  font-weight: 600;
}

.notes-truncate {
  font-size: 0.78rem;
  color: var(--color-gray-500);
  margin-top: 2px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 220px;
}

/* Action Icons */
.action-buttons-group {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
}

.icon-action-btn {
  width: 32px;
  height: 32px;
  border-radius: var(--radius-full);
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.85rem;
  transition: var(--transition);
}

.btn-confirm {
  background: var(--color-success-light);
  color: var(--color-success);
}
.btn-confirm:hover {
  background: var(--color-success);
  color: var(--color-white);
}

.btn-complete {
  background: var(--color-primary-pale);
  color: var(--color-primary);
}
.btn-complete:hover {
  background: var(--color-primary);
  color: var(--color-white);
}

.btn-cancel {
  background: var(--color-danger-light);
  color: var(--color-danger);
}
.btn-cancel:hover {
  background: var(--color-danger);
  color: var(--color-white);
}

.btn-review {
  background: var(--color-warning-light);
  color: var(--color-warning);
}
.btn-review:hover {
  background: var(--color-warning);
  color: var(--color-white);
}

.mb-4 {
  margin-bottom: 16px;
}
</style>
