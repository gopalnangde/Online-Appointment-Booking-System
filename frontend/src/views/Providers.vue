<template>
  <div class="dashboard-wrapper">
    <div class="container">
      <div class="header-action-card">
        <div class="header-title-group">
          <h1 class="header-title">
            Service Providers
            <i class="pi pi-users title-icon"></i>
          </h1>
        </div>

        <div class="search-bar-container">
          <i class="pi pi-search search-bar-icon"></i>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search provider name, specialization, or city..."
            class="search-bar-input"
          />
          <button v-if="searchQuery" class="clear-search-btn" @click="searchQuery = ''" title="Clear search">
            <i class="pi pi-times-circle"></i>
          </button>
        </div>
      </div>

      <!-- Loading -->
      <div v-if="loading" style="text-align: center; padding: 60px;">
        <PProgressSpinner style="width: 50px; height: 50px;" strokeWidth="4" />
      </div>

      <!-- Empty State -->
      <div v-else-if="filteredProviders.length === 0" class="empty-state">
        <i class="pi pi-search-minus" style="font-size: 3rem; color: var(--color-gray-400); margin-bottom: 16px;"></i>
        <h3>No Providers Found</h3>
        <p>No service providers match your search criteria.</p>
      </div>

      <!-- Provider Cards Grid -->
      <div v-else>
        <div class="grid grid-3">
          <div
            v-for="provider in filteredProviders"
            :key="provider.id"
            class="provider-card flex-card"
          >
            <div class="provider-card-header clickable" @click="viewProfile(provider.id)">
              <div class="provider-avatar">{{ getInitials(provider.name) }}</div>
              <div>
                <h3 class="provider-name-link">{{ provider.name }}</h3>
                <div class="spec-rating-row">
                  <PTag :value="provider.specialization" severity="primary" class="spec-tag" />
                  <div class="card-rating-badge">
                    <i class="pi pi-star-fill star-icon"></i>
                    <span class="rating-num">{{ provider.average_rating ? provider.average_rating.toFixed(1) : 'New' }}</span>
                    <span v-if="provider.total_reviews" class="rating-count">({{ provider.total_reviews }})</span>
                  </div>
                </div>
              </div>
            </div>
            <div class="provider-card-body flex-body">
              <div class="provider-info-row">
                <i class="pi pi-map-marker icon"></i>
                <span>{{ provider.city }}, {{ provider.state }}</span>
              </div>
              <div class="provider-info-row">
                <i class="pi pi-home icon"></i>
                <span>{{ provider.address }}</span>
              </div>
              <div class="provider-info-row">
                <i class="pi pi-clock icon"></i>
                <span>{{ provider.experience }} years experience</span>
              </div>
              <div v-if="provider.email" class="provider-info-row">
                <i class="pi pi-envelope icon"></i>
                <span>{{ provider.email }}</span>
              </div>
              <div v-if="provider.description" style="margin-top: 4px;">
                <p style="font-size: 0.85rem; color: var(--color-gray-500); line-height: 1.5;">
                  {{ truncate(provider.description, 90) }}
                </p>
              </div>

              <div class="card-footer card-actions-grid">
                <PButton
                  label="View Profile"
                  icon="pi pi-user"
                  class="p-button-outlined p-button-secondary action-btn"
                  @click="viewProfile(provider.id)"
                />
                <PButton
                  label="Book"
                  icon="pi pi-calendar-plus"
                  class="p-button-primary action-btn"
                  @click="onBookClick(provider)"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- Pagination -->
        <div v-if="totalPages > 1" class="pagination">
          <PButton
            label="Previous"
            icon="pi pi-chevron-left"
            class="p-button-outlined p-button-sm"
            :disabled="currentPage <= 1"
            @click="goToPage(currentPage - 1)"
          />
          <span class="page-info">Page {{ currentPage }} of {{ totalPages }}</span>
          <PButton
            label="Next"
            icon="pi pi-chevron-right"
            iconPos="right"
            class="p-button-outlined p-button-sm"
            :disabled="currentPage >= totalPages"
            @click="goToPage(currentPage + 1)"
          />
        </div>
      </div>
    </div>

    <!-- Booking Modal -->
    <BookingModal
      v-if="showBookingModal && selectedProvider"
      :provider="selectedProvider"
      @close="showBookingModal = false"
      @success="showBookingModal = false"
    />
  </div>
</template>

<script>
import { getAllProviders } from '../services/api.js'
import BookingModal from '../components/BookingModal.vue'

export default {
  name: 'ProvidersPage',
  components: { BookingModal },
  data() {
    return {
      providers: [],
      searchQuery: '',
      loading: true,
      currentPage: 1,
      limit: 12,
      total: 0,
      showBookingModal: false,
      selectedProvider: null,
    }
  },
  computed: {
    totalPages() {
      return Math.ceil(this.total / this.limit)
    },
    isLoggedIn() {
      return !!localStorage.getItem('token')
    },
    filteredProviders() {
      if (!this.searchQuery.trim()) return this.providers
      const q = this.searchQuery.toLowerCase()
      return this.providers.filter(p => {
        const name = (p.name || '').toLowerCase()
        const spec = (p.specialization || '').toLowerCase()
        const city = (p.city || '').toLowerCase()
        return name.includes(q) || spec.includes(q) || city.includes(q)
      })
    },
  },
  methods: {
    viewProfile(providerId) {
      this.$router.push(`/providers/${providerId}`)
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

    truncate(text, length) {
      if (!text) return ''
      return text.length > length ? text.substring(0, length) + '...' : text
    },

    async fetchProviders() {
      this.loading = true
      try {
        const { status, data } = await getAllProviders(this.currentPage, this.limit)
        if (data.success && data.data) {
          this.providers = data.data.providers || []
          this.total = data.data.total || 0
        }
      } catch (err) {
        console.error('Failed to fetch providers:', err)
      } finally {
        this.loading = false
      }
    },

    onBookClick(provider) {
      if (!this.isLoggedIn) {
        this.$router.push('/login')
        return
      }
      this.selectedProvider = provider
      this.showBookingModal = true
    },

    goToPage(page) {
      this.currentPage = page
      this.fetchProviders()
      window.scrollTo({ top: 0, behavior: 'smooth' })
    },
  },
  mounted() {
    this.fetchProviders()
  },
}
</script>

<style scoped>
.header-action-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 20px;
  margin-bottom: 24px;
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

.flex-card {
  display: flex;
  flex-direction: column;
}

.flex-body {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.spec-tag {
  margin-top: 4px;
  background: var(--color-primary-dark);
  color: #fff;
  font-size: 0.75rem;
}

.card-footer {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid var(--color-gray-100);
}

.clickable {
  cursor: pointer;
  transition: var(--transition);
}

.clickable:hover .provider-name-link {
  color: var(--color-primary);
  text-decoration: underline;
}

.provider-name-link {
  font-size: 1.15rem;
  font-weight: 700;
  color: var(--color-dark);
  transition: var(--transition);
}

.spec-rating-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  margin-top: 4px;
}

.card-rating-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  background: #FEF3C7;
  color: #92400E;
  font-size: 0.78rem;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: var(--radius-full);
}

.card-rating-badge .star-icon {
  color: #F59E0B;
  font-size: 0.75rem;
}

.card-rating-badge .rating-count {
  color: #B45309;
  font-weight: 500;
}

.card-actions-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

.action-btn {
  width: 100%;
  font-size: 0.85rem !important;
  padding: 8px 12px !important;
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
  margin-top: 40px;
  padding: 20px 0;
}

.page-info {
  font-size: 0.9rem;
  font-weight: 500;
  color: var(--color-gray-500);
}

.w-full {
  width: 100%;
}
</style>
