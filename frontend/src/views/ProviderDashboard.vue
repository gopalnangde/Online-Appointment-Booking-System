<template>
  <div class="dashboard-wrapper">
    <div class="container">
      <div class="page-header" style="display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 16px;">
        <div>
          <h1 class="page-title">
            <i class="pi pi-briefcase" style="margin-right: 10px; color: var(--color-primary-dark);"></i>
            My Practice
          </h1>
          <p class="page-subtitle">Manage your service provider profile & available appointment time slots</p>
        </div>
        <router-link to="/profile">
          <PButton label="Back to Profile" icon="pi pi-arrow-left" class="p-button-outlined p-button-sm" />
        </router-link>
      </div>

      <!-- Loading -->
      <div v-if="loading" style="text-align: center; padding: 60px;">
        <PProgressSpinner style="width: 50px; height: 50px;" strokeWidth="4" />
      </div>

      <!-- No Profile Yet — Show Create Form -->
      <div v-else-if="!profile && !showCreateForm" class="empty-state">
        <i class="pi pi-building" style="font-size: 3.5rem; color: var(--color-gray-400); margin-bottom: 16px;"></i>
        <h3>No Practice Profile Yet</h3>
        <p>Set up your provider profile so customers can find and book appointments with you.</p>
        <PButton
          label="Create Provider Profile"
          icon="pi pi-plus"
          class="p-button-lg p-button-primary"
          @click="initCreateForm"
        />
      </div>

      <!-- Create Profile Form -->
      <div v-else-if="showCreateForm && !profile" class="card" style="max-width: 750px; margin: 0 auto;">
        <h2 style="font-size: 1.3rem; font-weight: 700; margin-bottom: 24px;">Create Your Provider Profile</h2>

        <PMessage v-if="formError" severity="error" :closable="false" style="margin-bottom: 16px;">
          {{ formError }}
        </PMessage>
        <PMessage v-if="formValidation.length" severity="error" :closable="false" style="margin-bottom: 16px;">
          <ul style="margin: 0; padding-left: 16px;">
            <li v-for="(err, i) in formValidation" :key="i">{{ err }}</li>
          </ul>
        </PMessage>

        <form @submit.prevent="handleCreate" style="display: flex; flex-direction: column; gap: 18px;">
          <div class="grid grid-2">
            <div class="form-group">
              <label class="form-label">Specialization <span class="req">*</span></label>
              <PInputText v-model="form.specialization" placeholder="e.g. Dentist, Hair Stylist" class="w-full" required />
            </div>
            <div class="form-group">
              <label class="form-label">Experience (Years) <span class="req">*</span></label>
              <PInputText v-model.number="form.experience" type="number" placeholder="5" min="0" class="w-full" required />
            </div>
          </div>

          <div class="form-group">
            <label class="form-label">Description</label>
            <PTextarea v-model="form.description" placeholder="Tell customers about your services..." rows="3" class="w-full" />
          </div>

          <div class="form-group">
            <label class="form-label">Address <span class="req">*</span></label>
            <PInputText v-model="form.address" placeholder="Street address" class="w-full" required />
          </div>

          <div class="grid grid-3">
            <div class="form-group">
              <label class="form-label">City <span class="req">*</span></label>
              <PInputText v-model="form.city" placeholder="Mumbai" class="w-full" required />
            </div>
            <div class="form-group">
              <label class="form-label">State <span class="req">*</span></label>
              <PInputText v-model="form.state" placeholder="Maharashtra" class="w-full" required />
            </div>
            <div class="form-group">
              <label class="form-label">PIN Code <span class="req">*</span></label>
              <PInputText v-model="form.pin_code" placeholder="400001" maxlength="6" class="w-full" required />
            </div>
          </div>

          <!-- Manage Time Slots Section -->
          <div class="slots-section">
            <div class="slots-section-header">
              <label class="form-label mb-0"><i class="pi pi-clock" style="margin-right: 6px;"></i> Configure Available Time Slots</label>
              <button type="button" class="text-btn" @click="resetFormSlots">Reset Default Slots</button>
            </div>
            <p class="section-hint">Customers will pick from these time slots when booking appointments with you.</p>

            <div class="slots-tag-list">
              <span v-for="(slot, idx) in form.available_slots" :key="idx" class="slot-tag">
                <i class="pi pi-clock" style="font-size: 0.75rem;"></i> {{ slot }}
                <button type="button" class="remove-slot-btn" @click="removeFormSlot(idx)">&times;</button>
              </span>
            </div>

            <div class="add-slot-row">
              <input v-model="newSlotInput" type="text" placeholder="e.g. 08:30 AM or 07:00 PM" class="form-input slot-input" />
              <PButton type="button" label="Add Slot" icon="pi pi-plus" class="p-button-outlined p-button-sm" @click="addFormSlot" />
            </div>
          </div>

          <div style="display: flex; gap: 12px; justify-content: flex-end; margin-top: 12px;">
            <PButton label="Cancel" icon="pi pi-times" class="p-button-outlined p-button-secondary" @click="showCreateForm = false" />
            <PButton type="submit" label="Create Profile" icon="pi pi-check" :loading="submitting" class="p-button-primary" />
          </div>
        </form>
      </div>

      <!-- Profile Exists — Show Details + Edit -->
      <div v-else-if="profile">
        <PMessage v-if="formSuccess" severity="success" :closable="true" style="margin-bottom: 20px;" @close="formSuccess = ''">
          {{ formSuccess }}
        </PMessage>

        <!-- View Mode -->
        <div v-if="!editMode" class="profile-grid">
          <div class="profile-sidebar">
            <div class="profile-avatar-lg">{{ initials }}</div>
            <h2>{{ profile.name }}</h2>
            <p class="email">{{ profile.specialization }}</p>
            <div style="margin-top: 12px;">
              <PTag :value="`${profile.experience} yrs experience`" severity="primary" />
            </div>
            <div class="divider" style="margin: 20px 0;"></div>
            <PButton label="Edit Practice Profile" icon="pi pi-user-edit" class="p-button-primary w-full p-button-sm" @click="startEdit" />
          </div>

          <div class="profile-details">
            <h3><i class="pi pi-id-card" style="margin-right: 8px;"></i> Practice Information</h3>
            <div class="detail-grid">
              <div class="detail-item">
                <label>Specialization</label>
                <p>{{ profile.specialization }}</p>
              </div>
              <div class="detail-item">
                <label>Experience</label>
                <p>{{ profile.experience }} years</p>
              </div>
              <div class="detail-item">
                <label>Email</label>
                <p>{{ profile.email }}</p>
              </div>
              <div class="detail-item">
                <label>Phone</label>
                <p>{{ profile.phone }}</p>
              </div>
              <div class="detail-item">
                <label>Address</label>
                <p>{{ profile.address }}</p>
              </div>
              <div class="detail-item">
                <label>City</label>
                <p>{{ profile.city }}</p>
              </div>
              <div class="detail-item">
                <label>State</label>
                <p>{{ profile.state }}</p>
              </div>
              <div class="detail-item">
                <label>PIN Code</label>
                <p>{{ profile.pin_code }}</p>
              </div>
            </div>

            <!-- View Available Time Slots -->
            <div style="margin-top: 24px;">
              <div class="flex-between">
                <label class="section-label">Configured Available Time Slots</label>
                <button type="button" class="text-btn" @click="startEdit">Change Time Slots</button>
              </div>
              <div class="slots-tag-list" style="margin-top: 8px;">
                <span v-for="(slot, idx) in (profile.available_slots || defaultSlots)" :key="idx" class="slot-tag readonly">
                  <i class="pi pi-clock" style="font-size: 0.75rem;"></i> {{ slot }}
                </span>
              </div>
            </div>

            <div v-if="profile.description" style="margin-top: 24px;">
              <label class="section-label">About Practice</label>
              <p style="font-size: 0.95rem; color: var(--color-gray-500); line-height: 1.7;">{{ profile.description }}</p>
            </div>
          </div>
        </div>

        <!-- My Patients & Clients Roster Section -->
        <div v-if="!editMode" style="margin-top: 24px;" class="card">
          <div class="flex-between" style="margin-bottom: 16px;">
            <h3 style="font-size: 1.15rem; font-weight: 700; color: var(--color-dark); margin: 0;">
              <i class="pi pi-users" style="margin-right: 8px; color: var(--color-primary);"></i>
              My Patients & Clients
            </h3>
            <router-link to="/appointments">
              <PButton label="View All Appointments" icon="pi pi-arrow-right" class="p-button-text p-button-sm" />
            </router-link>
          </div>

          <div v-if="uniquePatients.length === 0" style="text-align: center; padding: 24px; color: var(--color-gray-500);">
            <i class="pi pi-user-minus" style="font-size: 2rem; margin-bottom: 8px;"></i>
            <p>No patient records found yet. Bookings will populate your patient list.</p>
          </div>

          <div v-else class="grid grid-3">
            <div v-for="patient in uniquePatients" :key="patient.id" class="patient-card">
              <div class="patient-card-top">
                <div class="patient-avatar-small">{{ getPatientInitials(patient.name) }}</div>
                <div>
                  <h4 class="patient-card-name">{{ patient.name }}</h4>
                  <span class="patient-card-email">{{ patient.email }}</span>
                </div>
              </div>
              <div class="patient-card-body">
                <p><i class="pi pi-phone" style="color: var(--color-primary); margin-right: 6px;"></i> {{ patient.phone || 'N/A' }}</p>
                <p><i class="pi pi-calendar" style="color: var(--color-primary); margin-right: 6px;"></i> {{ patient.appointmentCount }} Booking(s)</p>
              </div>
              <button class="view-patient-btn" @click="openPatientModal(patient)">
                <i class="pi pi-eye"></i> View Profile
              </button>
            </div>
          </div>
        </div>

        <!-- Edit Mode -->
        <div v-else class="card" style="max-width: 750px; margin: 0 auto;">
          <h2 style="font-size: 1.3rem; font-weight: 700; margin-bottom: 24px;">Edit Provider Profile & Time Slots</h2>

          <PMessage v-if="formError" severity="error" :closable="false" style="margin-bottom: 16px;">
            {{ formError }}
          </PMessage>

          <form @submit.prevent="handleUpdate" style="display: flex; flex-direction: column; gap: 18px;">
            <div class="grid grid-2">
              <div class="form-group">
                <label class="form-label">Specialization</label>
                <PInputText v-model="editForm.specialization" class="w-full" />
              </div>
              <div class="form-group">
                <label class="form-label">Experience (Years)</label>
                <PInputText v-model.number="editForm.experience" type="number" min="0" class="w-full" />
              </div>
            </div>

            <div class="form-group">
              <label class="form-label">Description</label>
              <PTextarea v-model="editForm.description" rows="3" class="w-full" />
            </div>

            <div class="form-group">
              <label class="form-label">Address</label>
              <PInputText v-model="editForm.address" class="w-full" />
            </div>

            <div class="grid grid-3">
              <div class="form-group">
                <label class="form-label">City</label>
                <PInputText v-model="editForm.city" class="w-full" />
              </div>
              <div class="form-group">
                <label class="form-label">State</label>
                <PInputText v-model="editForm.state" class="w-full" />
              </div>
              <div class="form-group">
                <label class="form-label">PIN Code</label>
                <PInputText v-model="editForm.pin_code" maxlength="6" class="w-full" />
              </div>
            </div>

            <!-- Manage Time Slots Section in Edit Mode -->
            <div class="slots-section">
              <div class="slots-section-header">
                <label class="form-label mb-0"><i class="pi pi-clock" style="margin-right: 6px;"></i> Practice Time Slots</label>
                <button type="button" class="text-btn" @click="resetEditSlots">Reset Default Slots</button>
              </div>
              <p class="section-hint">Add or remove time slots available for customer bookings.</p>

              <div class="slots-tag-list">
                <span v-for="(slot, idx) in editForm.available_slots" :key="idx" class="slot-tag">
                  <i class="pi pi-clock" style="font-size: 0.75rem;"></i> {{ slot }}
                  <button type="button" class="remove-slot-btn" @click="removeEditSlot(idx)">&times;</button>
                </span>
              </div>

              <div class="add-slot-row">
                <input v-model="editNewSlotInput" type="text" placeholder="e.g. 08:30 AM or 07:00 PM" class="form-input slot-input" />
                <PButton type="button" label="Add Time Slot" icon="pi pi-plus" class="p-button-outlined p-button-sm" @click="addEditSlot" />
              </div>
            </div>

            <div style="display: flex; gap: 12px; justify-content: flex-end; margin-top: 12px;">
              <PButton label="Cancel" icon="pi pi-times" class="p-button-outlined p-button-secondary" @click="editMode = false" />
              <PButton type="submit" label="Save Changes" icon="pi pi-check" :loading="submitting" class="p-button-primary" />
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Patient Profile Modal -->
    <PatientProfileModal
      v-if="showPatientModal && selectedPatient"
      :customerId="selectedPatient.id"
      :customerData="selectedPatient"
      @close="showPatientModal = false"
    />
  </div>
</template>

<script>
import { getProviderProfile, createProviderProfile, updateProviderProfile, getMyAppointments } from '../services/api.js'
import PatientProfileModal from '../components/PatientProfileModal.vue'

const DEFAULT_SLOTS = [
  '09:00 AM',
  '10:00 AM',
  '11:00 AM',
  '12:00 PM',
  '02:00 PM',
  '03:00 PM',
  '04:00 PM',
  '05:00 PM',
  '06:00 PM',
]

export default {
  name: 'ProviderDashboard',
  components: { PatientProfileModal },
  data() {
    return {
      profile: null,
      loading: true,
      showCreateForm: false,
      editMode: false,
      submitting: false,
      formError: '',
      formSuccess: '',
      formValidation: [],
      newSlotInput: '',
      editNewSlotInput: '',
      defaultSlots: DEFAULT_SLOTS,
      uniquePatients: [],
      showPatientModal: false,
      selectedPatient: null,
      form: {
        specialization: '',
        description: '',
        address: '',
        city: '',
        state: '',
        pin_code: '',
        experience: 0,
        available_slots: [...DEFAULT_SLOTS],
      },
      editForm: {
        available_slots: [],
      },
    }
  },
  computed: {
    initials() {
      if (!this.profile?.name) return '?'
      return this.profile.name
        .split(' ')
        .map(w => w[0])
        .join('')
        .toUpperCase()
        .slice(0, 2)
    },
  },
  methods: {
    initCreateForm() {
      this.form.available_slots = [...DEFAULT_SLOTS]
      this.showCreateForm = true
    },

    addFormSlot() {
      if (this.newSlotInput.trim()) {
        this.form.available_slots.push(this.newSlotInput.trim())
        this.newSlotInput = ''
      }
    },

    removeFormSlot(idx) {
      this.form.available_slots.splice(idx, 1)
    },

    resetFormSlots() {
      this.form.available_slots = [...DEFAULT_SLOTS]
    },

    addEditSlot() {
      if (this.editNewSlotInput.trim()) {
        this.editForm.available_slots.push(this.editNewSlotInput.trim())
        this.editNewSlotInput = ''
      }
    },

    removeEditSlot(idx) {
      this.editForm.available_slots.splice(idx, 1)
    },

    resetEditSlots() {
      this.editForm.available_slots = [...DEFAULT_SLOTS]
    },

    async fetchProfile() {
      try {
        const { status, data } = await getProviderProfile()
        if (data.success && data.data) {
          this.profile = data.data
        }
      } catch (err) {
        // No profile yet — handled cleanly
      } finally {
        this.loading = false
      }
    },

    async handleCreate() {
      this.formError = ''
      this.formValidation = []
      this.submitting = true

      try {
        const { status, data } = await createProviderProfile(this.form)
        if (data.success) {
          this.formSuccess = data.message
          this.showCreateForm = false
          await this.fetchProfile()
        } else {
          if (Array.isArray(data.data)) {
            this.formValidation = data.data
          } else {
            this.formError = data.message
          }
        }
      } catch (err) {
        this.formError = 'Unable to connect to the server'
      } finally {
        this.submitting = false
      }
    },

    startEdit() {
      const currentSlots = this.profile.available_slots && this.profile.available_slots.length > 0
        ? [...this.profile.available_slots]
        : [...DEFAULT_SLOTS]

      this.editForm = {
        specialization: this.profile.specialization,
        description: this.profile.description,
        address: this.profile.address,
        city: this.profile.city,
        state: this.profile.state,
        pin_code: this.profile.pin_code,
        experience: this.profile.experience,
        available_slots: currentSlots,
      }
      this.editMode = true
      this.formError = ''
      this.formSuccess = ''
    },

    async handleUpdate() {
      this.formError = ''
      this.submitting = true

      try {
        const { status, data } = await updateProviderProfile(this.editForm)
        if (data.success) {
          this.formSuccess = data.message
          this.editMode = false
          await this.fetchProfile()
        } else {
          this.formError = data.message
        }
      } catch (err) {
        this.formError = 'Unable to connect to the server'
      } finally {
        this.submitting = false
      }
    },

    async fetchPatients() {
      try {
        const { data } = await getMyAppointments()
        if (data.success && Array.isArray(data.data)) {
          const map = new Map()
          data.data.forEach(apt => {
            if (apt.customer && apt.customer.id) {
              const cid = apt.customer.id
              if (!map.has(cid)) {
                map.set(cid, {
                  ...apt.customer,
                  appointmentCount: 1,
                })
              } else {
                const existing = map.get(cid)
                existing.appointmentCount++
              }
            }
          })
          this.uniquePatients = Array.from(map.values())
        }
      } catch (err) {
        console.error('Error loading provider patients:', err)
      }
    },

    getPatientInitials(name) {
      if (!name) return '?'
      return name
        .split(' ')
        .map(w => w[0])
        .join('')
        .toUpperCase()
        .slice(0, 2)
    },

    openPatientModal(patient) {
      this.selectedPatient = patient
      this.showPatientModal = true
    },
  },
  mounted() {
    this.fetchProfile()
    this.fetchPatients()
  },
}
</script>

<style scoped>
.slots-section {
  background: var(--color-gray-100);
  border: 1px solid var(--color-gray-200);
  border-radius: var(--radius-md);
  padding: 16px;
  margin: 4px 0;
}

.slots-section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.section-hint {
  font-size: 0.82rem;
  color: var(--color-gray-500);
  margin-bottom: 12px;
}

.section-label {
  display: block;
  font-size: 0.78rem;
  font-weight: 600;
  color: var(--color-gray-400);
  text-transform: uppercase;
  letter-spacing: 0.06em;
  margin-bottom: 6px;
}

.slots-tag-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 14px;
}

.slot-tag {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: var(--color-white);
  border: 1px solid var(--color-gray-200);
  color: var(--color-dark-secondary);
  padding: 6px 12px;
  border-radius: var(--radius-full);
  font-size: 0.84rem;
  font-weight: 600;
}

.slot-tag.readonly {
  background: var(--color-primary-pale);
  border-color: var(--color-primary-light);
  color: var(--color-primary-dark);
}

.remove-slot-btn {
  background: none;
  border: none;
  color: var(--color-gray-400);
  font-size: 1rem;
  font-weight: 700;
  cursor: pointer;
  padding: 0 0 0 4px;
  line-height: 1;
}

.remove-slot-btn:hover {
  color: var(--color-danger);
}

.add-slot-row {
  display: flex;
  gap: 8px;
  align-items: center;
}

.slot-input {
  max-width: 220px;
}

.text-btn {
  background: none;
  border: none;
  color: var(--color-primary);
  font-size: 0.82rem;
  font-weight: 700;
  cursor: pointer;
  padding: 0;
}

.text-btn:hover {
  text-decoration: underline;
}

.flex-between {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.form-label {
  font-size: 0.88rem;
  font-weight: 700;
  color: var(--color-dark-secondary);
}

.req {
  color: var(--color-danger);
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

.w-full {
  width: 100%;
}

.mb-0 {
  margin-bottom: 0 !important;
}

.patient-card {
  background: var(--color-white);
  border: 1px solid var(--color-gray-200);
  border-radius: var(--radius-md);
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  box-shadow: var(--shadow-xs);
  transition: var(--transition);
}

.patient-card:hover {
  border-color: var(--color-primary-light);
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
}

.patient-card-top {
  display: flex;
  align-items: center;
  gap: 12px;
}

.patient-avatar-small {
  width: 40px;
  height: 40px;
  border-radius: var(--radius-full);
  background: var(--color-primary-pale);
  color: var(--color-primary-dark);
  font-weight: 800;
  font-size: 0.95rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.patient-card-name {
  font-size: 0.98rem;
  font-weight: 700;
  color: var(--color-dark);
  margin: 0;
}

.patient-card-email {
  font-size: 0.78rem;
  color: var(--color-gray-500);
  display: block;
}

.patient-card-body {
  font-size: 0.84rem;
  color: var(--color-gray-600);
  display: flex;
  flex-direction: column;
  gap: 4px;
  background: var(--color-gray-100);
  padding: 10px;
  border-radius: 6px;
}

.patient-card-body p {
  margin: 0;
  display: flex;
  align-items: center;
}

.view-patient-btn {
  background: var(--color-primary-pale);
  color: var(--color-primary);
  border: 1px solid var(--color-primary-light);
  border-radius: var(--radius-md);
  padding: 8px 12px;
  font-size: 0.82rem;
  font-weight: 700;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  transition: var(--transition);
}

.view-patient-btn:hover {
  background: var(--color-primary);
  color: var(--color-white);
}
</style>
