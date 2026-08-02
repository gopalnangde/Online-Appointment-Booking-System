<template>
  <PDialog
    :visible="true"
    modal
    header="Leave a Review"
    :style="{ width: '90vw', maxWidth: '500px' }"
    @update:visible="$emit('close')"
  >
    <template #header>
      <div>
        <h3 class="dialog-title">
          <i class="pi pi-star-fill" style="margin-right: 8px; color: #F1C40F;"></i>
          Leave a Review
        </h3>
        <p class="dialog-subtitle" v-if="appointment">
          For appointment on <strong>{{ appointment.service_title }}</strong>
        </p>
      </div>
    </template>

    <PMessage v-if="success" severity="success" :closable="false" class="mb-3">
      Thank you! Your review has been submitted.
    </PMessage>

    <PMessage v-if="error" severity="error" :closable="false" class="mb-3">
      {{ error }}
    </PMessage>

    <form @submit.prevent="handleReviewSubmit" class="review-form">
      <div class="form-group">
        <label class="form-label">Rating</label>
        <div class="rating-container">
          <PRating v-model="rating" :stars="5" :cancel="false" />
          <span class="rating-label" v-if="rating > 0">
            {{ rating }}/5 — {{ getRatingLabel(rating) }}
          </span>
        </div>
      </div>

      <div class="form-group">
        <label class="form-label">Your Feedback / Review</label>
        <PTextarea
          v-model="comment"
          placeholder="Share details of your experience with this provider..."
          rows="4"
          class="w-full"
          required
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
          label="Submit Review"
          icon="pi pi-send"
          :loading="loading"
          :disabled="rating === 0"
          class="p-button-primary"
        />
      </div>
    </form>
  </PDialog>
</template>

<script>
import { createReview } from '../services/api.js'

export default {
  name: 'ReviewModal',
  props: {
    appointment: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      rating: 5,
      comment: '',
      loading: false,
      error: '',
      success: false,
    }
  },
  methods: {
    getRatingLabel(r) {
      const labels = ['', 'Poor', 'Fair', 'Good', 'Very Good', 'Excellent']
      return labels[r] || ''
    },
    async handleReviewSubmit() {
      if (this.rating === 0) {
        this.error = 'Please select a rating'
        return
      }

      this.loading = true
      this.error = ''
      this.success = false

      try {
        const payload = {
          appointment_id: this.appointment.id,
          provider_id: this.appointment.provider_id,
          rating: this.rating,
          comment: this.comment,
        }

        const { status, data } = await createReview(payload)

        if (data.success) {
          this.success = true
          setTimeout(() => {
            this.$emit('success')
            this.$emit('close')
          }, 1200)
        } else {
          this.error = data.message || 'Failed to submit review'
        }
      } catch (err) {
        this.error = 'An error occurred while submitting review'
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

.review-form {
  display: flex;
  flex-direction: column;
  gap: 18px;
  margin-top: 12px;
}

.rating-container {
  display: flex;
  align-items: center;
  gap: 14px;
}

.rating-label {
  font-size: 0.95rem;
  font-weight: 600;
  color: var(--color-dark-secondary);
}

.dialog-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 12px;
}

.w-full {
  width: 100%;
}

.mb-3 {
  margin-bottom: 12px;
}
</style>
