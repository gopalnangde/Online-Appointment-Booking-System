<template>
  <div class="provider-profile-wrapper">
    <div class="container">
      <!-- Top Navigation / Back Button -->
      <div class="top-nav-bar">
        <button class="back-btn" @click="$router.push('/providers')">
          <i class="pi pi-arrow-left"></i>
          <span>Back to Service Providers</span>
        </button>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="loading-state">
        <PProgressSpinner style="width: 50px; height: 50px;" strokeWidth="4" />
        <p class="mt-3 text-gray">Loading provider details...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="error-state">
        <i class="pi pi-exclamation-triangle error-icon"></i>
        <h3>Unable to load provider profile</h3>
        <p>{{ error }}</p>
        <PButton label="Back to Providers" icon="pi pi-arrow-left" class="mt-3 p-button-outlined" @click="$router.push('/providers')" />
      </div>

      <!-- Provider Profile Content -->
      <div v-else-if="provider" class="profile-content">
        <!-- Hero Header Card -->
        <div class="provider-header-card">
          <div class="header-main-info">
            <div class="avatar-container">
              <div class="large-avatar">{{ getInitials(provider.name) }}</div>
              <span class="online-indicator" title="Active Provider"></span>
            </div>

            <div class="header-details">
              <div class="name-row">
                <h1 class="provider-name">{{ provider.name }}</h1>
                <PTag :value="provider.specialization" severity="primary" class="spec-badge" />
              </div>

              <div class="meta-row">
                <div class="meta-item">
                  <i class="pi pi-clock meta-icon"></i>
                  <span>{{ provider.experience }} Years Experience</span>
                </div>
                <div class="meta-item">
                  <i class="pi pi-map-marker meta-icon"></i>
                  <span>{{ provider.city }}, {{ provider.state }}</span>
                </div>
                <div v-if="provider.email" class="meta-item">
                  <i class="pi pi-envelope meta-icon"></i>
                  <span>{{ provider.email }}</span>
                </div>
              </div>

              <!-- Overall Rating Pill Header -->
              <div class="rating-header-pill">
                <div class="stars-display">
                  <i
                    v-for="i in 5"
                    :key="i"
                    :class="['pi', i <= Math.round(overallRating) ? 'pi-star-fill active' : 'pi-star']"
                  ></i>
                </div>
                <span class="rating-score">{{ overallRating > 0 ? overallRating.toFixed(1) : 'New' }}</span>
                <span class="reviews-count">({{ totalReviews }} {{ totalReviews === 1 ? 'review' : 'reviews' }})</span>
              </div>
            </div>
          </div>

          <div class="header-action-area">
            <PButton
              label="Book Appointment"
              icon="pi pi-calendar-plus"
              class="p-button-primary p-button-lg book-now-btn"
              @click="onBookClick"
            />
          </div>
        </div>

        <!-- Main Content Layout (Grid) -->
        <div class="profile-grid">
          <!-- Left Column: Bio & Reviews -->
          <div class="left-column">
            <!-- About Card -->
            <div class="content-card">
              <h2 class="card-title">
                <i class="pi pi-user card-title-icon"></i>
                About Provider
              </h2>
              <div class="divider"></div>
              <p class="bio-text">
                {{ provider.description || 'This service provider has not provided a detailed description yet.' }}
              </p>
            </div>

            <!-- Ratings & Reviews Section -->
            <div class="content-card">
              <div class="reviews-header">
                <h2 class="card-title">
                  <i class="pi pi-star-fill card-title-icon star-yellow"></i>
                  Ratings & Reviews
                </h2>
                <span class="badge-count">{{ totalReviews }} total</span>
              </div>
              <div class="divider"></div>

              <!-- Rating Summary Banner -->
              <div class="rating-summary-banner">
                <div class="rating-big-box">
                  <span class="big-score">{{ overallRating > 0 ? overallRating.toFixed(1) : '0.0' }}</span>
                  <div class="big-stars">
                    <i
                      v-for="i in 5"
                      :key="i"
                      :class="['pi', i <= Math.round(overallRating) ? 'pi-star-fill active' : 'pi-star']"
                    ></i>
                  </div>
                  <span class="sub-text">Based on {{ totalReviews }} customer {{ totalReviews === 1 ? 'review' : 'reviews' }}</span>
                </div>

                <div class="rating-breakdown">
                  <div v-for="star in [5, 4, 3, 2, 1]" :key="star" class="breakdown-row">
                    <span class="star-label">{{ star }} ★</span>
                    <div class="progress-bar-bg">
                      <div
                        class="progress-bar-fill"
                        :style="{ width: getStarPercentage(star) + '%' }"
                      ></div>
                    </div>
                    <span class="star-count">{{ getStarCount(star) }}</span>
                  </div>
                </div>
              </div>

              <!-- Reviews List -->
              <div class="reviews-list-container">
                <div v-if="reviewsLoading" class="reviews-loading">
                  <PProgressSpinner style="width: 36px; height: 36px;" strokeWidth="3" />
                  <span>Loading reviews...</span>
                </div>

                <div v-else-if="reviews.length === 0" class="no-reviews-state">
                  <i class="pi pi-comment-slash no-reviews-icon"></i>
                  <h4>No Reviews Yet</h4>
                  <p>This service provider hasn't received any customer reviews yet.</p>
                </div>

                <div v-else class="reviews-list">
                  <div v-for="review in reviews" :key="review.id" class="review-card">
                    <div class="review-card-header">
                      <div class="reviewer-info">
                        <div class="reviewer-avatar">
                          {{ getInitials(review.customer?.name || 'Customer') }}
                        </div>
                        <div>
                          <h4 class="reviewer-name">{{ review.customer?.name || 'Verified Customer' }}</h4>
                          <span class="review-date">{{ formatDate(review.created_at) }}</span>
                        </div>
                      </div>

                      <div class="review-rating-stars">
                        <i
                          v-for="i in 5"
                          :key="i"
                          :class="['pi', i <= review.rating ? 'pi-star-fill active' : 'pi-star']"
                        ></i>
                        <span class="rating-num">{{ review.rating }}/5</span>
                      </div>
                    </div>

                    <p class="review-comment" v-if="review.comment">
                      "{{ review.comment }}"
                    </p>
                    <p class="review-comment empty-comment" v-else>
                      <em>No detailed comment provided.</em>
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Right Column: Quick Details & Contact -->
          <div class="right-column">
            <div class="content-card sidebar-card">
              <h3 class="sidebar-title">Provider Overview</h3>
              <div class="divider"></div>

              <div class="info-list">
                <div class="info-item">
                  <div class="info-icon-box">
                    <i class="pi pi-briefcase"></i>
                  </div>
                  <div>
                    <label class="info-label">Specialization</label>
                    <p class="info-value">{{ provider.specialization }}</p>
                  </div>
                </div>

                <div class="info-item">
                  <div class="info-icon-box">
                    <i class="pi pi-id-card"></i>
                  </div>
                  <div>
                    <label class="info-label">Experience</label>
                    <p class="info-value">{{ provider.experience }} Years in Practice</p>
                  </div>
                </div>

                <div class="info-item">
                  <div class="info-icon-box">
                    <i class="pi pi-map-marker"></i>
                  </div>
                  <div>
                    <label class="info-label">Address & Location</label>
                    <p class="info-value">{{ provider.address }}</p>
                    <p class="info-subvalue">{{ provider.city }}, {{ provider.state }} - {{ provider.pin_code }}</p>
                  </div>
                </div>

                <div v-if="provider.email" class="info-item">
                  <div class="info-icon-box">
                    <i class="pi pi-envelope"></i>
                  </div>
                  <div>
                    <label class="info-label">Email Address</label>
                    <p class="info-value">{{ provider.email }}</p>
                  </div>
                </div>

                <div v-if="provider.phone" class="info-item">
                  <div class="info-icon-box">
                    <i class="pi pi-phone"></i>
                  </div>
                  <div>
                    <label class="info-label">Contact Phone</label>
                    <p class="info-value">{{ provider.phone }}</p>
                  </div>
                </div>
              </div>

              <div class="sidebar-booking-box">
                <h4>Ready to Schedule?</h4>
                <p>Select your date and time for an appointment.</p>
                <PButton
                  label="Book Appointment"
                  icon="pi pi-calendar-plus"
                  class="p-button-primary w-full mt-3"
                  @click="onBookClick"
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Booking Modal -->
    <BookingModal
      v-if="showBookingModal && provider"
      :provider="provider"
      @close="showBookingModal = false"
      @success="showBookingModal = false"
    />
  </div>
</template>

<script>
import { getProviderById, getProviderReviews } from '../services/api.js'
import BookingModal from '../components/BookingModal.vue'

export default {
  name: 'ProviderProfileView',
  components: { BookingModal },
  data() {
    return {
      provider: null,
      reviews: [],
      overallRating: 0,
      totalReviews: 0,
      loading: true,
      reviewsLoading: true,
      error: '',
      showBookingModal: false,
    }
  },
  computed: {
    isLoggedIn() {
      return !!localStorage.getItem('token')
    },
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

    formatDate(dateStr) {
      if (!dateStr) return ''
      const d = new Date(dateStr)
      if (isNaN(d.getTime())) return dateStr
      return d.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
      })
    },

    getStarCount(star) {
      return this.reviews.filter(r => r.rating === star).length
    },

    getStarPercentage(star) {
      if (!this.reviews.length) return 0
      const count = this.getStarCount(star)
      return Math.round((count / this.reviews.length) * 100)
    },

    async fetchProviderDetails() {
      this.loading = true
      this.error = ''
      const providerId = this.$route.params.id

      try {
        const { status, data } = await getProviderById(providerId)
        if (data.success && data.data) {
          this.provider = data.data
          if (this.provider.average_rating) {
            this.overallRating = Number(this.provider.average_rating)
          }
          if (this.provider.total_reviews) {
            this.totalReviews = Number(this.provider.total_reviews)
          }
        } else {
          this.error = data.message || 'Provider profile not found'
        }
      } catch (err) {
        console.error('Failed to load provider profile:', err)
        this.error = 'Error connecting to server'
      } finally {
        this.loading = false
      }
    },

    async fetchReviews() {
      if (!this.provider) return
      this.reviewsLoading = true
      try {
        const targetId = this.provider.user_id || this.provider.id
        const { status, data } = await getProviderReviews(targetId)
        if (data.success && data.data) {
          this.reviews = data.data.reviews || []
          if (data.data.stats) {
            this.overallRating = Number(data.data.stats.average_rating || 0)
            this.totalReviews = Number(data.data.stats.total_reviews || 0)
          }
        }
      } catch (err) {
        console.error('Failed to fetch provider reviews:', err)
      } finally {
        this.reviewsLoading = false
      }
    },

    onBookClick() {
      if (!this.isLoggedIn) {
        this.$router.push('/login')
        return
      }
      this.showBookingModal = true
    },
  },
  async mounted() {
    await this.fetchProviderDetails()
    if (this.provider) {
      await this.fetchReviews()
    }
  },
}
</script>

<style scoped>
.provider-profile-wrapper {
  padding: 10px 0 60px;
}

.top-nav-bar {
  margin-bottom: 20px;
}

.back-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: var(--color-white);
  border: 1px solid var(--color-gray-200);
  padding: 8px 18px;
  border-radius: var(--radius-full);
  font-size: 0.9rem;
  font-weight: 600;
  color: var(--color-dark-secondary);
  cursor: pointer;
  transition: var(--transition);
  box-shadow: var(--shadow-sm);
}

.back-btn:hover {
  background: var(--color-primary-pale);
  color: var(--color-primary);
  border-color: var(--color-primary-light);
}

/* Loading & Error States */
.loading-state,
.error-state {
  text-align: center;
  padding: 80px 20px;
  background: var(--color-white);
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-gray-200);
}

.error-icon {
  font-size: 3rem;
  color: var(--color-danger);
  margin-bottom: 16px;
}

/* Hero Header Card */
.provider-header-card {
  background: var(--color-white);
  border-radius: var(--radius-xl);
  padding: 32px;
  border: 1px solid var(--color-gray-200);
  box-shadow: var(--shadow-md);
  margin-bottom: 28px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 24px;
  flex-wrap: wrap;
}

.header-main-info {
  display: flex;
  align-items: center;
  gap: 24px;
  flex: 1;
  min-width: 300px;
}

.avatar-container {
  position: relative;
  flex-shrink: 0;
}

.large-avatar {
  width: 90px;
  height: 90px;
  border-radius: 50%;
  background: linear-gradient(135deg, #0052CC 0%, #0066FF 100%);
  color: #fff;
  font-size: 2.2rem;
  font-weight: 800;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8px 24px rgba(0, 102, 255, 0.25);
  border: 3px solid #fff;
}

.online-indicator {
  position: absolute;
  bottom: 4px;
  right: 4px;
  width: 16px;
  height: 16px;
  background: #10B981;
  border: 3px solid #fff;
  border-radius: 50%;
}

.header-details {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.name-row {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.provider-name {
  font-size: 1.9rem;
  font-weight: 800;
  color: var(--color-dark);
  margin: 0;
}

.spec-badge {
  background: var(--color-primary-dark) !important;
  color: #fff !important;
  font-size: 0.85rem !important;
}

.meta-row {
  display: flex;
  align-items: center;
  gap: 20px;
  flex-wrap: wrap;
  color: var(--color-gray-500);
  font-size: 0.9rem;
  font-weight: 500;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.meta-icon {
  color: var(--color-primary);
}

.rating-header-pill {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: #FEF3C7;
  padding: 6px 14px;
  border-radius: var(--radius-full);
  margin-top: 4px;
  width: fit-content;
}

.stars-display {
  display: flex;
  gap: 2px;
}

.stars-display i.pi-star-fill.active {
  color: #F59E0B;
}

.stars-display i.pi-star {
  color: #D1D5DB;
}

.rating-score {
  font-weight: 800;
  color: #92400E;
  font-size: 0.95rem;
}

.reviews-count {
  font-size: 0.85rem;
  color: #B45309;
  font-weight: 600;
}

.book-now-btn {
  padding: 12px 28px !important;
  font-size: 1.05rem !important;
  font-weight: 700 !important;
  box-shadow: 0 4px 16px rgba(0, 102, 255, 0.3) !important;
}

/* Grid Layout */
.profile-grid {
  display: grid;
  grid-template-columns: 1fr 340px;
  gap: 28px;
}

@media (max-width: 900px) {
  .profile-grid {
    grid-template-columns: 1fr;
  }
}

.content-card {
  background: var(--color-white);
  border-radius: var(--radius-lg);
  padding: 28px;
  border: 1px solid var(--color-gray-200);
  box-shadow: var(--shadow-sm);
  margin-bottom: 28px;
}

.card-title {
  font-size: 1.3rem;
  font-weight: 800;
  color: var(--color-dark);
  display: flex;
  align-items: center;
  gap: 10px;

}

.card-title-icon {
  color: var(--color-primary);
  font-size: 1.2rem;
}

.star-yellow {
  color: #F59E0B !important;
}

.divider {
  height: 1px;
  background: var(--color-gray-200);
  margin: 16px 0;
}

.bio-text {
  font-size: 0.98rem;
  line-height: 1.7;
  color: var(--color-dark-secondary);
  white-space: pre-line;
}

/* Ratings & Reviews Section */
.reviews-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.badge-count {
  background: var(--color-primary-pale);
  color: var(--color-primary-dark);
  font-weight: 700;
  font-size: 0.82rem;
  padding: 4px 12px;
  border-radius: var(--radius-full);
}

.rating-summary-banner {
  display: flex;
  gap: 32px;
  align-items: center;
  background: #FAFBFD;
  border: 1px solid var(--color-gray-200);
  border-radius: var(--radius-md);
  padding: 24px;
  margin-bottom: 28px;
  flex-wrap: wrap;
}

.rating-big-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-width: 140px;
  padding-right: 24px;
  border-right: 1px solid var(--color-gray-200);
}

@media (max-width: 600px) {
  .rating-big-box {
    border-right: none;
    padding-right: 0;
    padding-bottom: 16px;
    border-bottom: 1px solid var(--color-gray-200);
    width: 100%;
  }
}

.big-score {
  font-size: 3rem;
  font-weight: 800;
  color: var(--color-dark);
  line-height: 1;
}

.big-stars {
  display: flex;
  gap: 4px;
  margin: 8px 0;
  font-size: 1.1rem;
}

.big-stars i.active {
  color: #F59E0B;
}

.big-stars i.pi-star {
  color: #D1D5DB;
}

.sub-text {
  font-size: 0.8rem;
  color: var(--color-gray-500);
  text-align: center;
}

.rating-breakdown {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-width: 220px;
}

.breakdown-row {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 0.85rem;
}

.star-label {
  width: 32px;
  font-weight: 700;
  color: var(--color-dark-secondary);
}

.progress-bar-bg {
  flex: 1;
  height: 8px;
  background: #E5E7EB;
  border-radius: 4px;
  overflow: hidden;
}

.progress-bar-fill {
  height: 100%;
  background: #F59E0B;
  border-radius: 4px;
  transition: width 0.4s ease;
}

.star-count {
  width: 24px;
  text-align: right;
  color: var(--color-gray-500);
  font-weight: 600;
}

/* Reviews List */
.reviews-loading,
.no-reviews-state {
  text-align: center;
  padding: 40px 20px;
  color: var(--color-gray-500);
}

.no-reviews-icon {
  font-size: 2.5rem;
  color: var(--color-gray-400);
  margin-bottom: 12px;
}

.reviews-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.review-card {
  background: var(--color-white);
  border: 1px solid var(--color-gray-200);
  border-radius: var(--radius-md);
  padding: 20px;
  transition: var(--transition);
}

.review-card:hover {
  border-color: var(--color-primary-light);
  box-shadow: 0 4px 12px rgba(0, 102, 255, 0.05);
}

.review-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  flex-wrap: wrap;
  gap: 8px;
}

.reviewer-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.reviewer-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: var(--color-primary-pale);
  color: var(--color-primary-dark);
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.95rem;
}

.reviewer-name {
  font-size: 0.95rem;
  font-weight: 700;
  color: var(--color-dark);
  margin: 0;
}

.review-date {
  font-size: 0.78rem;
  color: var(--color-gray-400);
}

.review-rating-stars {
  display: flex;
  align-items: center;
  gap: 4px;
}

.review-rating-stars i.active {
  color: #F59E0B;
}

.review-rating-stars i.pi-star {
  color: #D1D5DB;
}

.rating-num {
  font-size: 0.85rem;
  font-weight: 700;
  color: var(--color-dark);
  margin-left: 6px;
}

.review-comment {
  font-size: 0.92rem;
  color: var(--color-dark-secondary);
  line-height: 1.5;

}

.empty-comment {
  color: var(--color-gray-400);
}

/* Sidebar Card */
.sidebar-title {
  font-size: 1.15rem;
  font-weight: 800;
  color: var(--color-dark);
}

.info-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
  margin-bottom: 24px;
}

.info-item {
  display: flex;
  align-items: flex-start;
  gap: 14px;
}

.info-icon-box {
  width: 38px;
  height: 38px;
  border-radius: var(--radius-md);
  background: var(--color-primary-pale);
  color: var(--color-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.1rem;
  flex-shrink: 0;
}

.info-label {
  font-size: 0.78rem;
  font-weight: 700;
  color: var(--color-gray-400);
  text-transform: uppercase;
  letter-spacing: 0.04em;
  display: block;
}

.info-value {
  font-size: 0.92rem;
  font-weight: 700;
  color: var(--color-dark);
  margin-top: 2px;
}

.info-subvalue {
  font-size: 0.85rem;
  color: var(--color-gray-500);
}

.sidebar-booking-box {
  background: linear-gradient(135deg, #F0F6FF 0%, #E6F0FF 100%);
  border-radius: var(--radius-md);
  padding: 20px;
  text-align: center;
  border: 1px solid var(--color-primary-pale);
}

.sidebar-booking-box h4 {
  font-size: 1.05rem;
  font-weight: 800;
  color: var(--color-primary-dark);
  margin-bottom: 4px;
}

.sidebar-booking-box p {
  font-size: 0.85rem;
  color: var(--color-gray-500);
}

.w-full {
  width: 100%;
}

.mt-3 {
  margin-top: 12px;
}
</style>
