<template>
  <PDialog
    :visible="true"
    modal
    header="Book an Appointment"
    :style="{ width: '90vw', maxWidth: '560px' }"
    @update:visible="$emit('close')"
  >
    <template #header>
      <div>
        <h3 class="dialog-title">
          <i class="pi pi-calendar-plus" style="margin-right: 8px; color: var(--color-primary-dark);"></i>
          Book an Appointment
        </h3>
        <p class="dialog-subtitle" v-if="provider">
          with <strong>{{ provider.name }}</strong> ({{ provider.specialization }})
        </p>
      </div>
    </template>

    <PMessage v-if="success" severity="success" :closable="false" class="mb-3">
      Appointment booked successfully! Redirecting...
    </PMessage>

    <PMessage v-if="error" severity="error" :closable="false" class="mb-3">
      {{ error }}
    </PMessage>

    <form @submit.prevent="handleBooking" class="booking-form">
      <div class="form-group">
        <label class="form-label">Service Type <span class="req">*</span></label>
        <PInputText
          v-model="form.service_title"
          placeholder="e.g. General Consultation, Follow-up"
          class="w-full"
          required
        />
      </div>

      <div class="form-group">
        <label class="form-label">Select Date <span class="req">*</span></label>
        <input
          v-model="form.booking_date"
          type="date"
          class="form-input"
          :min="minDate"
          required
        />
      </div>

      <!-- Multiple Time Slots Selection -->
      <div class="form-group">
        <div class="slots-header">
          <label class="form-label mb-0">Select Available Time Slot <span class="req">*</span></label>
          <button
            type="button"
            class="toggle-custom-btn"
            @click="isCustomTime = !isCustomTime"
          >
            {{ isCustomTime ? 'Choose Preset Slot' : '+ Custom Time' }}
          </button>
        </div>

        <!-- Time Slot Chips -->
        <div v-if="!isCustomTime" class="slots-grid">
          <button
            v-for="slot in availableSlots"
            :key="slot"
            type="button"
            class="slot-chip"
            :class="{ active: form.booking_time === slot }"
            @click="form.booking_time = slot"
          >
            <i class="pi pi-clock slot-icon"></i>
            <span>{{ slot }}</span>
            <i v-if="form.booking_time === slot" class="pi pi-check slot-check"></i>
          </button>
        </div>

        <!-- Custom Time Input -->
        <div v-else class="custom-time-wrapper">
          <input
            v-model="customTimeInput"
            type="time"
            class="form-input"
            @change="applyCustomTime"
          />
          <p class="hint-text">Selected custom time: <strong>{{ form.booking_time || 'None' }}</strong></p>
        </div>

        <p v-if="!form.booking_time" class="validation-warning">
          <i class="pi pi-exclamation-circle"></i> Please select a time slot above
        </p>
      </div>

      <div class="form-group">
        <label class="form-label">Additional Notes (Optional)</label>
        <PTextarea
          v-model="form.notes"
          placeholder="Describe your request or any specific requirements..."
          rows="3"
          class="w-full"
        />
      </div>

      <div class="dialog-actions">
        <PButton
          type="button"
          label="Cancel"
          icon="pi pi-times"
          class="p-button-outlined p-button-secondary"
          @click="$emit('close')"
        />
        <PButton
          type="submit"
          label="Confirm Booking"
          icon="pi pi-check"
          :loading="loading"
          :disabled="!form.booking_time"
          class="p-button-primary"
        />
      </div>
    </form>
  </PDialog>
</template>

<script>
import { createAppointment } from '../services/api.js'

export default {
  name: 'BookingModal',
  props: {
    provider: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      form: {
        service_title: '',
        booking_date: '',
        booking_time: '',
        notes: '',
      },
      isCustomTime: false,
      customTimeInput: '',
      loading: false,
      error: '',
      success: false,
    }
  },
  computed: {
    minDate() {
      const today = new Date()
      return today.toISOString().split('T')[0]
    },
    availableSlots() {
      if (this.provider && Array.isArray(this.provider.available_slots) && this.provider.available_slots.length > 0) {
        return this.provider.available_slots
      }
      return [
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
    },
  },
  mounted() {
    if (this.provider && this.provider.specialization) {
      this.form.service_title = `${this.provider.specialization} Consultation`
    }
    // Set default date to today
    this.form.booking_date = this.minDate
    // Auto-select first slot
    if (this.availableSlots.length > 0) {
      this.form.booking_time = this.availableSlots[0]
    }
  },
  methods: {
    applyCustomTime() {
      if (this.customTimeInput) {
        // Format 24h input (e.g. 14:30) to 12h AM/PM format
        const [h, m] = this.customTimeInput.split(':')
        let hour = parseInt(h, 10)
        const ampm = hour >= 12 ? 'PM' : 'AM'
        hour = hour % 12 || 12
        const formattedHour = hour < 10 ? `0${hour}` : `${hour}`
        this.form.booking_time = `${formattedHour}:${m} ${ampm}`
      }
    },

    async handleBooking() {
      if (!this.form.booking_time) {
        this.error = 'Please select an appointment time slot.'
        return
      }

      this.loading = true
      this.error = ''
      this.success = false

      try {
        const rawProviderId = this.provider?.user_id || this.provider?.userId || this.provider?.id
        const providerId = Number(rawProviderId)

        if (!providerId) {
          this.error = 'Invalid provider selection. Please re-select a provider.'
          this.loading = false
          return
        }

        const payload = {
          provider_id: providerId,
          service_title: this.form.service_title,
          booking_date: this.form.booking_date,
          booking_time: this.form.booking_time,
          notes: this.form.notes || '',
        }

        const { status, data } = await createAppointment(payload)

        if (data.success) {
          this.success = true
          setTimeout(() => {
            this.$emit('success')
            this.$router.push('/appointments')
          }, 1200)
        } else {
          this.error = data.message || 'Failed to book appointment'
        }
      } catch (err) {
        this.error = err.response?.data?.message || err.message || 'An error occurred while creating appointment'
      } finally {
        this.loading = false
      }
    },
  },
}
</script>

<style scoped>
.dialog-title {
  font-size: 1.3rem;
  font-weight: 700;
  color: var(--color-dark);
}

.dialog-subtitle {
  font-size: 0.88rem;
  color: var(--color-gray-500);
  margin-top: 4px;
}

.booking-form {
  display: flex;
  flex-direction: column;
  gap: 18px;
  margin-top: 12px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-label {
  font-size: 0.88rem;
  font-weight: 700;
  color: var(--color-dark-secondary);
}

.req {
  color: var(--color-danger);
}

.slots-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.toggle-custom-btn {
  background: none;
  border: none;
  color: var(--color-primary);
  font-size: 0.82rem;
  font-weight: 700;
  cursor: pointer;
  padding: 0;
  transition: var(--transition);
}

.toggle-custom-btn:hover {
  text-decoration: underline;
  color: var(--color-primary-dark);
}

.slots-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(115px, 1fr));
  gap: 10px;
}

.slot-chip {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  background: var(--color-white);
  border: 1.5px solid var(--color-gray-200);
  border-radius: var(--radius-md);
  padding: 10px 12px;
  font-size: 0.86rem;
  font-weight: 700;
  color: var(--color-dark-secondary);
  cursor: pointer;
  transition: var(--transition);
  position: relative;
}

.slot-chip:hover {
  border-color: var(--color-primary-light);
  background: var(--color-primary-pale);
  color: var(--color-primary);
}

.slot-chip.active {
  background: var(--color-primary);
  border-color: var(--color-primary);
  color: var(--color-white);
  box-shadow: 0 4px 12px rgba(0, 102, 255, 0.25);
}

.slot-icon {
  font-size: 0.8rem;
}

.slot-check {
  font-size: 0.8rem;
  margin-left: 2px;
}

.custom-time-wrapper {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.hint-text {
  font-size: 0.8rem;
  color: var(--color-gray-500);
}

.validation-warning {
  font-size: 0.8rem;
  color: var(--color-danger);
  margin-top: 4px;
  display: flex;
  align-items: center;
  gap: 4px;
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

.dialog-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 16px;
}

.w-full {
  width: 100%;
}

.mb-0 {
  margin-bottom: 0 !important;
}

.mb-3 {
  margin-bottom: 12px;
}
</style>
