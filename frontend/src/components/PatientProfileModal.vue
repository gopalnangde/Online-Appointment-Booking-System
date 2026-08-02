<template>
  <PDialog
    :visible="true"
    modal
    header="Patient / Customer Profile"
    :style="{ width: '90vw', maxWidth: '620px' }"
    @update:visible="$emit('close')"
  >
    <template #header>
      <div class="patient-modal-header">
        <h3 class="dialog-title">
          <i class="pi pi-user" style="margin-right: 8px; color: var(--color-primary);"></i>
          Patient / Client Profile
        </h3>
        <p class="dialog-subtitle">Detailed information & appointment history</p>
      </div>
    </template>

    <div v-if="loading" style="text-align: center; padding: 40px;">
      <PProgressSpinner style="width: 40px; height: 40px;" strokeWidth="4" />
    </div>

    <div v-else-if="patient" class="patient-profile-content">
      <!-- Profile Top Header -->
      <div class="patient-card-header">
        <div class="patient-avatar">
          {{ initials }}
        </div>
        <div class="patient-info-meta">
          <h2 class="patient-name">{{ patient.name }}</h2>
          <div class="patient-badges">
            <span class="role-badge"><i class="pi pi-id-card"></i> {{ patient.role || 'Customer' }}</span>
            <span class="id-badge">ID: #{{ patient.id }}</span>
          </div>
        </div>
      </div>

      <div class="divider"></div>

      <!-- Contact Information Cards -->
      <div class="info-section">
        <h4 class="section-title"><i class="pi pi-address-book"></i> Contact Details</h4>
        <div class="info-grid">
          <div class="info-item">
            <label>Email Address</label>
            <div class="info-value-row">
              <i class="pi pi-envelope"></i>
              <a :href="`mailto:${patient.email}`" class="contact-link">{{ patient.email }}</a>
            </div>
          </div>

          <div class="info-item">
            <label>Phone Number</label>
            <div class="info-value-row">
              <i class="pi pi-phone"></i>
              <a :href="`tel:${patient.phone}`" class="contact-link">{{ patient.phone || 'N/A' }}</a>
            </div>
          </div>

          <div class="info-item">
            <label>Member Since</label>
            <div class="info-value-row">
              <i class="pi pi-calendar"></i>
              <span>{{ formatDate(patient.created_at) }}</span>
            </div>
          </div>

          <div class="info-item">
            <label>Total Bookings With You</label>
            <div class="info-value-row">
              <i class="pi pi-bookmark"></i>
              <span class="booking-count-chip">{{ patientAppointments.length }} Appointment(s)</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Patient Appointment History Section -->
      <div class="info-section style-history" v-if="patientAppointments.length > 0">
        <h4 class="section-title"><i class="pi pi-history"></i> Appointment History with Patient</h4>
        <div class="history-list">
          <div v-for="apt in patientAppointments" :key="apt.id" class="history-item">
            <div class="history-header">
              <div class="history-service">{{ apt.service_title }}</div>
              <PTag :value="apt.status" :severity="getTagSeverity(apt.status)" />
            </div>
            <div class="history-meta">
              <span><i class="pi pi-calendar"></i> {{ apt.booking_date }}</span>
              <span><i class="pi pi-clock"></i> {{ apt.booking_time }}</span>
            </div>
            <p v-if="apt.notes" class="history-notes">
              <strong>Notes:</strong> {{ apt.notes }}
            </p>
          </div>
        </div>
      </div>

      <!-- Direct Quick Action Controls -->
      <div class="dialog-actions">
        <a :href="`mailto:${patient.email}`" class="action-btn email-btn">
          <i class="pi pi-envelope"></i> Send Email
        </a>
        <a :href="`tel:${patient.phone}`" class="action-btn call-btn">
          <i class="pi pi-phone"></i> Call Patient
        </a>
        <PButton
          type="button"
          label="Close"
          icon="pi pi-times"
          class="p-button-outlined p-button-secondary"
          @click="$emit('close')"
        />
      </div>
    </div>

    <!-- Fallback error -->
    <div v-else class="empty-state">
      <p>Unable to load patient profile details.</p>
    </div>
  </PDialog>
</template>

<script>
import { getUserById, getMyAppointments } from '../services/api.js'

export default {
  name: 'PatientProfileModal',
  props: {
    customerId: {
      type: [Number, String],
      required: true,
    },
    customerData: {
      type: Object,
      default: null,
    },
  },
  data() {
    return {
      patient: this.customerData,
      loading: !this.customerData,
      patientAppointments: [],
    }
  },
  computed: {
    initials() {
      if (!this.patient?.name) return '?'
      return this.patient.name
        .split(' ')
        .map(w => w[0])
        .join('')
        .toUpperCase()
        .slice(0, 2)
    },
  },
  mounted() {
    this.fetchPatientInfo()
    this.fetchPatientAppointments()
  },
  methods: {
    async fetchPatientInfo() {
      if (this.patient && this.patient.email && this.patient.phone) {
        this.loading = false
        return
      }
      this.loading = true
      try {
        const { data } = await getUserById(this.customerId)
        if (data.success && data.data) {
          this.patient = data.data
        }
      } catch (err) {
        console.error('Error fetching user profile:', err)
      } finally {
        this.loading = false
      }
    },

    async fetchPatientAppointments() {
      try {
        const { data } = await getMyAppointments()
        if (data.success && data.data) {
          // Filter appointments for this specific customer
          this.patientAppointments = data.data.filter(
            a => a.customer_id === Number(this.customerId) || a.customer?.id === Number(this.customerId)
          )
        }
      } catch (err) {
        console.error('Error loading patient appointments history:', err)
      }
    },

    formatDate(dateStr) {
      if (!dateStr) return 'N/A'
      const d = new Date(dateStr)
      if (isNaN(d)) return dateStr
      return d.toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
        year: 'numeric',
      })
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
  },
}
</script>

<style scoped>
.patient-modal-header {
  margin-bottom: 4px;
}

.dialog-title {
  font-size: 1.3rem;
  font-weight: 700;
  color: var(--color-dark);
}

.dialog-subtitle {
  font-size: 0.85rem;
  color: var(--color-gray-500);
}

.patient-profile-content {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.patient-card-header {
  display: flex;
  align-items: center;
  gap: 16px;
  background: linear-gradient(135deg, rgba(0, 102, 255, 0.06), rgba(0, 102, 255, 0.01));
  padding: 16px;
  border-radius: var(--radius-md);
  border: 1px solid rgba(0, 102, 255, 0.12);
}

.patient-avatar {
  width: 56px;
  height: 56px;
  border-radius: var(--radius-full);
  background: var(--color-primary);
  color: var(--color-white);
  font-size: 1.3rem;
  font-weight: 800;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 12px rgba(0, 102, 255, 0.2);
}

.patient-info-meta {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.patient-name {
  font-size: 1.25rem;
  font-weight: 800;
  color: var(--color-dark);
  margin: 0;
}

.patient-badges {
  display: flex;
  gap: 8px;
  align-items: center;
}

.role-badge {
  background: var(--color-primary-pale);
  color: var(--color-primary-dark);
  font-size: 0.76rem;
  font-weight: 700;
  padding: 2px 10px;
  border-radius: var(--radius-full);
}

.id-badge {
  font-size: 0.76rem;
  color: var(--color-gray-500);
  font-weight: 600;
}

.divider {
  height: 1px;
  background: var(--color-gray-200);
  margin: 0;
}

.info-section {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.section-title {
  font-size: 0.92rem;
  font-weight: 700;
  color: var(--color-dark-secondary);
  display: flex;
  align-items: center;
  gap: 8px;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

@media (max-width: 520px) {
  .info-grid {
    grid-template-columns: 1fr;
  }
}

.info-item {
  background: var(--color-gray-100);
  padding: 10px 14px;
  border-radius: var(--radius-md);
  border: 1px solid var(--color-gray-200);
}

.info-item label {
  display: block;
  font-size: 0.76rem;
  font-weight: 700;
  color: var(--color-gray-500);
  text-transform: uppercase;
  margin-bottom: 4px;
}

.info-value-row {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.92rem;
  font-weight: 600;
  color: var(--color-dark);
}

.info-value-row i {
  color: var(--color-primary);
  font-size: 0.9rem;
}

.contact-link {
  color: var(--color-primary);
  text-decoration: none;
  font-weight: 600;
  transition: var(--transition);
}

.contact-link:hover {
  text-decoration: underline;
  color: var(--color-primary-dark);
}

.booking-count-chip {
  background: var(--color-success-light);
  color: var(--color-success);
  padding: 2px 8px;
  border-radius: 6px;
  font-size: 0.8rem;
}

.history-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  max-height: 180px;
  overflow-y: auto;
  padding-right: 4px;
}

.history-item {
  background: var(--color-white);
  border: 1px solid var(--color-gray-200);
  border-radius: var(--radius-md);
  padding: 10px 14px;
}

.history-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.history-service {
  font-weight: 700;
  font-size: 0.88rem;
  color: var(--color-dark);
}

.history-meta {
  display: flex;
  gap: 14px;
  font-size: 0.78rem;
  color: var(--color-gray-500);
}

.history-notes {
  font-size: 0.8rem;
  color: var(--color-gray-600);
  margin-top: 6px;
  background: var(--color-gray-100);
  padding: 6px 10px;
  border-radius: 6px;
}

.dialog-actions {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 10px;
  margin-top: 8px;
}

.action-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  border-radius: var(--radius-md);
  font-size: 0.85rem;
  font-weight: 700;
  text-decoration: none;
  transition: var(--transition);
}

.email-btn {
  background: var(--color-primary-pale);
  color: var(--color-primary);
}

.email-btn:hover {
  background: var(--color-primary);
  color: var(--color-white);
}

.call-btn {
  background: var(--color-success-light);
  color: var(--color-success);
}

.call-btn:hover {
  background: var(--color-success);
  color: var(--color-white);
}
</style>
